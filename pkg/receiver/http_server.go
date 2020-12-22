package receiver

import (
	"net/http"
	"os"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/push"
	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/libs/mongo"
)

type HTTPService struct {
	database *mongo.Database
	server   *http.Server
	receiver *Service
}

func NewHTTPService(database *mongo.Database) *HTTPService {
	srv := &HTTPService{
		database: database,
		receiver: NewService(storage.New(database)),
	}
	return srv
}

func (srv *HTTPService) Run(addr string) error {
	handler := gin.Default()
	if os.Getenv("ENV") == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		handler = gin.New()
	}

	srv.route(handler)

	srv.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	return srv.server.ListenAndServe()
}

func (srv *HTTPService) Close() error {
	if srv.server != nil {
		_ = srv.server.Close()
	}
	return nil
}

func (srv *HTTPService) route(handler *gin.Engine, middleware ...gin.HandlerFunc) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })

	handler.Use(middleware...)
	handler.POST("/v1/receive/packet", srv.ReceivePacket)
}

func (srv *HTTPService) ReceivePacket(c *gin.Context) {
	var arg push.Packet
	if err := c.ShouldBind(&arg); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid)
		return
	}

	if err := srv.receiver.InsertPacket(c.Request.Context(), c.ClientIP(), &arg); err != nil {
		httputil.RespError(c, httputil.ErrClaimsNotFound)
		return
	}
	httputil.RespSuccess(c)
}