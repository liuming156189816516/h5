package api

import (
	info "api/webstru"
	"comm/comm"
	"comm/goError"
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
		tmp := &accountDB.AccountInfo{}
		tmp.Id = bson.NewObjectId()
		tmp.Account = req.Account
		tmp.Status = 2
		tmp.AccountType = 1
		tmp.Itime = time.Now().Unix()
		tmp.Ptime = time.Now().Unix()
		tmp.FirstLoginTime = time.Now().Unix()
		tmp.PlatformType = 2
		err := accountDB.AddAccountInfo(tmp)
		if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
			//更新为登录中
			accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"status": int64(2), "reason": "", "first_login_time": time.Now().Unix()})
		}
		// 开关控制 "0" - 开 "1" - 关
		if cache.GetTaskStatus() == "0" {
			go wxComm.AutoSendMsg(req.Account, accountData.SessionId, accountData.Node)
		} else {
			//添加进缓存中，后续使用定时任务发送
			cacheTmp := &cache.AutoSendMsgTaskInfo{}
			cacheTmp.Account = req.Account
			cacheTmp.SessionId = accountData.SessionId
			cacheTmp.Node = accountData.Node
			cache.SetAutoSendMsgTaskInfo(cacheTmp)
		}
	}

	if accountData.Action == "logout" {
		cache.SetAccountStatus(req.Account, 1)
		accountInfo := accountDB.GetOneAccountInfo(bson.M{"account": req.Account})
		if accountInfo.Status == 2 {
			sendmsg.UpSendMsgInfo(bson.M{"account": req.Account}, bson.M{"account_status": 1})
		} else {
			sendmsg.UpSendMsgInfo(bson.M{"account": req.Account}, bson.M{"account_status": 1})
			accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"reason": accountData.Reason, "status": int64(1), "offline_time": time.Now().Unix()})
		}

	}

}

//处理消息
func doMessage(req *info.ApiReq) {

}
