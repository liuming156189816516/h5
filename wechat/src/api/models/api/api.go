package api

import (
	info "api/webstru"
	"comm/comm"
	"comm/goError"
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	accountDB "selfComm/db/account"
	"selfComm/db/sendmsg"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"strings"
	"time"
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

	fmt.Println(jsoniter.MarshalToString(req))

	if req.Account == "" {
		return goError.GLOBAL_INVALIDPARAM
	}

	if req.Ptype == 1 {
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
	accountData := &info.AccountData{}
	dataStr, _ := jsoniter.MarshalToString(req.Data)
	jsoniter.UnmarshalFromString(dataStr, accountData)
	//账号登陆成功
	if accountData.Action == "login" {
		cache.SetAccountStatus(req.Account, 2)
		accInfo := cache.GetAccountInfo(req.Account)
		tmp := &accountDB.AccountInfo{}
		tmp.Id = bson.NewObjectId()
		tmp.Account = req.Account
		tmp.Status = 2
		tmp.AccountType = 1
		tmp.Itime = time.Now().Unix()
		tmp.Ptime = time.Now().Unix()
		tmp.FirstLoginTime = time.Now().Unix()
		tmp.PixelId = accInfo.PixelId
		tmp.PlatformType = 2
		err := accountDB.AddAccountInfo(tmp)
		if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
			//更新为登录中
			accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"status": int64(2), "pixel_id": accInfo.PixelId, "reason": "", "first_login_time": time.Now().Unix()})
		}
		////发送事件回调
		////kwai的回调
		//if accInfo.PixelId == wxComm.PixId {
		//	go func() {
		//		fbData := &info.FbData{}
		//		fbData.Phone = req.Account
		//		fbData.PixelId = accInfo.PixelId
		//		tmpFb := &log.FbReportLog{}
		//		tmpFb.Id = bson.NewObjectId()
		//		tmpFb.Ptype = 3
		//		data, _ := jsoniter.MarshalToString(fbData)
		//		tmpFb.Data = data
		//		log.AddFbReportLog(tmpFb)
		//		wxComm.KwaiPlace(accInfo.ClickId, "EVENT_COMPLETE_REGISTRATION")
		//	}()
		//}
		//fbFlag := false
		////fb的回调
		//for key, _ := range wxComm.PixelTokens {
		//	if strings.Contains(key, accInfo.PixelId) {
		//		fbFlag = true
		//		break
		//	}
		//}
		//if fbFlag {
		//	go func() {
		//		fbData := &info.FbData{}
		//		fbData.Phone = req.Account
		//		fbData.PixelId = accInfo.PixelId
		//		tmpFb := &log.FbReportLog{}
		//		tmpFb.Id = bson.NewObjectId()
		//		tmpFb.Ptype = 3
		//		data, _ := jsoniter.MarshalToString(fbData)
		//		tmpFb.Data = data
		//		log.AddFbReportLog(tmpFb)
		//		fbInfo := cache.FbReport{
		//			Ptype: 3,
		//			Phone: req.Account,
		//		}
		//		cache.SetFbReport(&fbInfo)
		//	}()
		//}
		go sendMsg(req.Account, accountData.SessionId)
	}

	if accountData.Action == "logout" {
		cache.SetAccountStatus(req.Account, 1)
		accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"reason": accountData.Reason, "status": int64(1), "offline_time": time.Now().Unix()})
		sendmsg.UpSendMsgInfo(bson.M{"account": req.Account}, bson.M{"account_status": 1})
	}

}

//处理消息
func doMessage(req *info.ApiReq) {

}

//发送消息
func sendMsg(account, sessionId string) {
	config := cache.GetTaskConfig("global")
	mList := config.MaterialList
	material := mList[0]

	tmp2 := &sendmsg.SendMsgInfo{}
	tmp2.Id = bson.NewObjectId()
	tmp2.Account = account
	tmp2.AccountStatus = 2
	err := sendmsg.AddSendMsgInfo(tmp2)
	if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
		//更新为登录成功
		sendmsg.UpSendMsgInfo(bson.M{"account": account}, bson.M{"account_status": 2})
	}

	// 消息发送的开关
	if cache.GetTaskStatus() == "1" {
		return
	}

	for i := 0; i < 100; i++ {
		time.Sleep(60 * time.Second)
		target := ""
		targetStr := cache.SpopDataPackList(config.DataPackId)
		if targetStr == "" {
			logs.Info("粉丝数据不足1")
			continue
		}
		split := strings.Split(targetStr, "-")
		if len(split) > 0 {
			target = split[0]
		}
		msgResult, err1 := wxComm.SendMsgUtils(sessionId, target, material)

		// 发送成功
		if err1 == nil && msgResult.Ok {
			cache.IncSendMsgTaskInfoCount(cache.SuccessNum, account, 1)
		}

		//发送失败，还数据
		if err1 != nil || !msgResult.Ok {
			cache.SaddDataPackList(config.DataPackId, target)
		}
		if cache.GetAccountStatus(account) != 2 {
			return
		}
	}

	////如果账号在线
	//for cache.GetAccountStatus(account) == 2 {
	//
	//
	//}
}
