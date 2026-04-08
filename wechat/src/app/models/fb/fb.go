package fb

import (
	info "app/webstru"
	"comm/comm"
	"comm/goError"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type FbService struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *FbService) getUid() string {
	return this.Sess.Uid
}

func (this *FbService) FbReport(req *info.FbReportReq, rsp *info.NullRsp) *goError.ErrRsp {

	toString, _ := jsoniter.MarshalToString(req)
	fmt.Println("toString1=================>", toString)

	//data := &info.FbData{}
	//jsoniter.UnmarshalFromString(req.Data, data)

	//if data.Fbclid != "" {
	//	//写入日志
	//	tmp := &log.FbReportLog{}
	//	tmp.Id = bson.NewObjectId()
	//	tmp.Ptype = req.Ptype
	//	tmp.Data = req.Data
	//	log.AddFbReportLog(tmp)
	//
	//	fbInfo := cache.FbReport{
	//		Ptype:   req.Ptype,
	//		Fbclid:  data.Fbclid,
	//		Fbp:     data.Fbp,
	//		PixelId: data.PixelId,
	//		Phone:   data.AreaCode + data.Account,
	//	}
	//	cache.SetFbReport(&fbInfo)
	//}
	return nil
}
