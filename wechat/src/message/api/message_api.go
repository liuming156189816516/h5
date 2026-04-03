package api

import (
	"comm/event"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"natsRpc"
	accountDB "selfComm/db/account"
	"selfComm/db/sendmsg"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"strings"
	"time"
	"utils"
)

// 接收消息处理
func ReceiveMessagesEventHandler(reqList []*natsRpc.ReceiveMessagesReq) int32 {
	for _, para := range reqList {

		toString, _ := jsoniter.MarshalToString(para)
		logs.Info("群发协议接收到消息===============>" + toString)

		rcontentStr, _ := jsoniter.MarshalToString(para.Content)
		rcontent := natsRpc.ReceiveMessagesContent{}
		jsoniter.UnmarshalFromString(rcontentStr, &rcontent)

		if rcontent.Group != "" {
			continue
		}

		if strings.Contains(rcontent.To, "@g") {
			continue
		}

		account := para.Account
		if para.Type == "error" {
			if para.Ctype == "stream:error" {
				decodeString, _ := jsoniter.MarshalToString(para.Content)
				if strings.Contains(string(decodeString), "replaced") {
					//账号抢登
					cache.SetReplacedAccount(account)
				}
			}

			if para.Ctype == "offline" {
				account = strings.ReplaceAll(account, "_", "")
				//账号掉线,需要重新登录
				e := &event.TaskUserLoginEventReq{
					Account: []string{account},
				}
				event.AddLoginTask(event.TaskTypeUserLogin, account, e)
			}
		}
		if para.Type == "result" {
			if para.Ctype == "/account/login" {
				//处理登录结果回调
				resultLogin(para)
			}
			if para.Ctype == "/message/send" || para.Ctype == "/call/phone" || para.Ctype == "/call/video" {
				//处理消息发送结果
				req1 := &event.ReceiveMessagesReq{}
				req1.Account = para.Account
				req1.Content = rcontentStr
				req1.Time = para.Time
				event.AddMessageTaskResult(event.TaskTypeUserMessageResult, account, req1)
			}
		}

		if para.Type == "notify" && para.Ctype == "receipt" {

			//已送达
			req1 := &event.ReceiveMessagesReq{}

			notifyContentStr, _ := jsoniter.MarshalToString(para.Content)
			notifyContent := &natsRpc.NotifyContent{}
			jsoniter.UnmarshalFromString(notifyContentStr, notifyContent)

			if strings.Contains(notifyContent.From, "@s.whatsapp.net") {
				notifyContent.From = strings.ReplaceAll(notifyContent.From, "@s.whatsapp.net", "")
				req1.FromType = 1
			}

			if strings.Contains(notifyContent.From, "lid") {
				split := strings.Split(notifyContent.From, ".")
				notifyContent.From = split[0]
				req1.FromType = 0
			}

			if strings.Contains(notifyContent.From, ".") {
				continue
			}

			req1.Type = para.Type
			req1.Ctype = para.Ctype
			req1.From = notifyContent.From
			req1.Read = notifyContent.Read
			req1.Account = para.Account
			event.AddMessageTaskResult(event.TaskTypeUserMessageCallBack, account, req1)
		}
	}
	return natsRpc.ESMR_SUCCEED
}

func resultLogin(req *natsRpc.ReceiveMessagesReq) int32 {
	accInfo := cache.GetAccountInfo(req.Account)
	accountLoginStr, _ := jsoniter.MarshalToString(req.Content)
	accountLoginContent := &natsRpc.AccountLoginContent{}
	jsoniter.UnmarshalFromString(accountLoginStr, accountLoginContent)

	if accountLoginContent.Errno == 0 {
		//登录成功
		accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"reason": "", "status": int64(2), "first_login_time": time.Now().Unix(), "offline_time": int64(0), "nick_name": accInfo.NickName})
		cache.SetAccountStatus(req.Account, 2)

		if accInfo.PixelId == wxComm.PixId && accInfo.ClickId != "" {
			//kwai发送成功的回调
			go func() {
				wxComm.KwaiPlace(accInfo.ClickId, "EVENT_COMPLETE_REGISTRATION")
			}()
		}

		//添加登录成功的任务
		tmp2 := &sendmsg.SendMsgInfo{}
		tmp2.Id = bson.NewObjectId()
		tmp2.Account = req.Account
		tmp2.AccountStatus = 2
		err := sendmsg.AddSendMsgInfo(tmp2)
		if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
			//更新为登录成功
			sendmsg.UpSendMsgInfo(bson.M{"account": req.Account}, bson.M{"account_status": 2})
		}
		//商业号需要设置商店名称
		if accInfo.PlatformType == 1 && accInfo.AccountType == 2 && accInfo.NickName == "" {
			e := &event.TaskTypeShopNameEventReq{
				Account:  req.Account,
				NickName: utils.RandStr(7),
				IsCheck:  true,
			}
			event.AddTask(event.TaskTypeShopName, req.Account, e)
		}
	} else {
		reason := accountLoginContent.Errmsg
		if strings.Contains(accountLoginContent.Errmsg, "KEY可能已失效") {
			reason = "网络错误"
		}
		if strings.Contains(accountLoginContent.Errmsg, "403") || strings.Contains(accountLoginContent.Errmsg, "402") || strings.Contains(accountLoginContent.Errmsg, "401") {
			reason = "封号"
			if strings.Contains(accountLoginContent.Errmsg, "401") {
				acInfo := cache.GetAccountInfo(req.Account)
				if acInfo.PlatformType == 2 {
					//app账号
					reason = "账号退出"
				}
			}
		}
		accountDB.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"reason": reason, "status": int64(1), "offline_time": time.Now().Unix()})
		cache.SetAccountStatus(req.Account, 1)
		ip := cache.GetProxyIp(req.Account)
		cache.IncIpUserNum(ip.IpId, -1)
		sendmsg.UpSendMsgInfo(bson.M{"account": req.Account}, bson.M{"account_status": 1})
	}
	return natsRpc.ESMR_SUCCEED
}
