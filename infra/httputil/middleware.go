package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/huzhongqing/qelog/infra/logs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandlerRegisterTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		WithTraceID(c)
		c.Next()
	}
}

type loggingRespWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (lrw *loggingRespWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func HandlerLogging(enable bool, skipPrefixPaths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enable {
			c.Next()
			return
		}

		// 记录请求
		path := c.Request.URL.Path
		for _, v := range skipPrefixPaths {
			if strings.HasPrefix(path, v) {
				c.Next()
				return
			}
		}
		ip := c.ClientIP()
		uri := c.Request.URL.RequestURI()
		method := c.Request.Method
		headers := c.Request.Header.Clone()
		request := ""
		if c.Request.Body != nil {
			b, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				RespError(c, err)
				c.Abort()
				return
			}
			// 这里已经读取第一次，就关掉
			_ = c.Request.Body.Close()
			body := bytes.NewBuffer(b)
			request = body.String()
			c.Request.Body = ioutil.NopCloser(body)
		}
		logs.Qezap.InfoWithCtx(c.Request.Context(), "Request", zap.Any("headers", headers), zap.String("reqBody", request),
			zap.String("uri", uri),
			logs.Qezap.ConditionOne(method),
			logs.Qezap.ConditionTwo(path),
			logs.Qezap.ConditionThree(ip))

		lrw := &loggingRespWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = lrw

		c.Next()
		baseResp := BaseResp{}
		if err := json.Unmarshal(lrw.body.Bytes(), &baseResp); err != nil || baseResp.Code != 0 {
			logs.Qezap.ErrorWithCtx(c.Request.Context(), "Response", zap.String("respBody", lrw.body.String()),
				logs.Qezap.ConditionOne(path), logs.Qezap.ConditionTwo(uri), logs.Qezap.ConditionThree(ip))
		}
	}
}

func GinLogger(skipPaths []string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(skipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			rawPath := path
			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			if param.Latency > time.Minute {
				// Truncate in a golang < 1.8 safe way
				param.Latency = param.Latency - param.Latency%time.Second
			}

			logs.Qezap.Debug("GIN", zap.String("latency", param.Latency.String()),
				zap.String("method", param.Method),
				zap.String("path", param.Path),
				zap.String("error", param.ErrorMessage),
				logs.Qezap.ConditionOne(strconv.Itoa(param.StatusCode)),
				logs.Qezap.ConditionTwo(param.Method+"_"+rawPath),
				logs.Qezap.ConditionThree(param.ClientIP))
		}
	}
}

func GinRecoveryWithLogger() gin.HandlerFunc {
	if logs.Qezap != nil {
		w := logs.Qezap.Clone().SetWriteLevel(zap.ErrorLevel).SetWritePrefix("GIN-ERROR")
		return gin.RecoveryWithWriter(w)
	}
	return gin.Recovery()
}
