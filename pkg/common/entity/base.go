package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageReq struct {
	Page  int64 `json:"page" form:"page" default:"1" `
	Limit int64 `json:"limit" from:"limit" default:"20"`
}

func (v PageReq) SetPage(opt *options.FindOptions) {
	opt.SetSkip((v.Page - 1) * v.Limit)
	opt.SetLimit(v.Limit)
}

type ObjectIDReq struct {
	ID string `json:"id" form:"id" binding:"required,len=24"`
}

type OmitObjectIDReq struct {
	ID string `json:"id" form:"id" binding:"omitempty,len=24"`
}

func (v OmitObjectIDReq) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(v.ID)
}

func (v ObjectIDReq) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(v.ID)
}

type TimeReq struct {
	BeginTsSec int64 `json:"beginTsSec" form:"beginTsSec"`
	EndTsSec   int64 `json:"endTsSec" form:"endTsSec"`
}

func (v TimeReq) BeginTime() time.Time {
	return time.Unix(v.BeginTsSec, 0)
}
func (v TimeReq) EndTime() time.Time {
	return time.Unix(v.EndTsSec, 0)
}

func (v TimeReq) InitTimeSection(t time.Duration) (b, e time.Time) {
	// 没有查询时间
	if v.BeginTsSec <= 0 && v.EndTsSec <= 0 {
		e = time.Now()
		b = e.Add(-t)
		return b, e
	}
	// 只有结束时间
	if v.BeginTsSec <= 0 && v.EndTsSec > 0 {
		e = time.Unix(v.EndTsSec, 0)
		b = e.Add(-t)
		return b, e
	}
	// 只有开始时间
	if v.BeginTsSec > 0 && v.EndTsSec <= 0 {
		b = time.Unix(v.BeginTsSec, 0)
		e = time.Now()
		return b, e
	}

	// 时间都存在
	return time.Unix(v.BeginTsSec, 0), time.Unix(v.EndTsSec, 0)
}

type ListResp struct {
	Count int64       `json:"count"`
	List  interface{} `json:"list"`
}
