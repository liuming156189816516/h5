package api

import (
	info "api/webstru"
	"comm/comm"
	"comm/goError"
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
		accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"status": int64(2), "reason": "", "first_login_time": time.Now().Unix()})
		go sendMsg(req.Account, accountData.SessionId, accountData.Node)
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

//发送消息
func sendMsg(account, sessionId string, node string) {
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
		//time.Sleep(10 * time.Second)
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
		msgResult, err1 := wxComm.SendMsgUtils(sessionId, target, material, node)

		// 发送成功
		if err1 == nil && msgResult.Ok {
			cache.IncSendMsgTaskInfoCount(cache.SuccessNum, account, 1)
		}

		//发送失败，还异常数据
		if err1 != nil || !msgResult.Ok {
			cache.SaddDataPackListErr(config.DataPackId, target)
		}

		if cache.GetAccountStatus(account) != 2 {
			return
		}
	}
}
