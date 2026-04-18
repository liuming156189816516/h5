package api

import (
	info "api/webstru"
	"comm/comm"
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/log"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"strings"
)

// 群发
type ApiServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *ApiServer) getUid() string {
	return this.Sess.Uid
}

func (this *ApiServer) Api(req *info.ApiReq, rsp *info.NullRsp) *goError.ErrRsp {

	if req.Ptype == 1 {

		//账号登陆成功
		//发送事件回调
		accInfo := cache.GetAccountInfo(req.Account)
		//kwai的回调
		if accInfo.PixelId == wxComm.PixId {
			go func() {
				fbData := &info.FbData{}
				fbData.Phone = req.Account
				fbData.PixelId = accInfo.PixelId
				tmp := &log.FbReportLog{}
				tmp.Id = bson.NewObjectId()
				tmp.Ptype = 3
				data, _ := jsoniter.MarshalToString(fbData)
				tmp.Data = data
				log.AddFbReportLog(tmp)
				wxComm.KwaiPlace(accInfo.ClickId, "EVENT_COMPLETE_REGISTRATION")
			}()
		}
		//fb的回调
		for key, _ := range wxComm.PixelTokens {
			if strings.Contains(key, accInfo.PixelId) {
				go func() {
					fbData := &info.FbData{}
					fbData.Phone = req.Account
					fbData.PixelId = accInfo.PixelId
					tmp := &log.FbReportLog{}
					tmp.Id = bson.NewObjectId()
					tmp.Ptype = 3
					data, _ := jsoniter.MarshalToString(fbData)
					tmp.Data = data
					log.AddFbReportLog(tmp)
					fbInfo := cache.FbReport{
						Ptype: 3,
						Phone: req.Account,
					}
					cache.SetFbReport(&fbInfo)
				}()
			}
		}

		//处理账号
		go doAccount(req)
	}
	if req.Ptype == 2 {
		//处理消息
		go doMessage(req)
	}
	return nil
}

//处理账号
func doAccount(req *info.ApiReq) {

}

//处理消息
func doMessage(req *info.ApiReq) {

}
