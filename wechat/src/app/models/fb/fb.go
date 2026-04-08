package fb

import (
	info "app/webstru"
	"comm/comm"
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/log"
	"selfComm/wxComm/cache"
)

type FbService struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *FbService) getUid() string {
	return this.Sess.Uid
}

func (this *FbService) FbReport(req *info.FbReportReq, rsp *info.NullRsp) *goError.ErrRsp {

	data := &info.FbData{}
	jsoniter.UnmarshalFromString(req.Data, data)

	//写入日志
	tmp := &log.FbReportLog{}
	tmp.Id = bson.NewObjectId()
	tmp.Ptype = req.Ptype
	tmp.Data = req.Data
	log.AddFbReportLog(tmp)

	fbInfo := cache.FbReport{
		Ptype:   req.Ptype,
		Fbclid:  data.Fbclid,
		Fbp:     data.Fbp,
		PixelId: data.PixelId,
		Phone:   data.AreaCode + data.Account,
	}
	cache.SetFbReport(&fbInfo)

	return nil
}
