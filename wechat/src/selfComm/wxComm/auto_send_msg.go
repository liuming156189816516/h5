package wxComm

import (
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/sendmsg"
	"selfComm/wxComm/cache"
	"strings"
)

func AutoSendMsg(account, sessionId string, node string) {
	config := cache.GetTaskConfig("global")
	mList := config.MaterialList
	material := mList[0]
	accountInfo := cache.GetAccountInfo(account)

	tmp2 := &sendmsg.SendMsgInfo{}
	tmp2.Id = bson.NewObjectId()
	tmp2.Account = account
	tmp2.AccountStatus = 2
	tmp2.AccountGroup = accountInfo.AccountGroup

	err := sendmsg.AddSendMsgInfo(tmp2)
	if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
		sendmsg.UpSendMsgInfo(
			bson.M{"account": account},
			bson.M{"account_status": 2, "account_group": accountInfo.AccountGroup},
		)
	}

	// ❗新增：连续错误计数
	errCount := 0

	for i := 0; i < 100; i++ {

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

		msgResult, err1 := SendMsgUtils(sessionId, target, material, node)

		// 发送成功
		if err1 == nil && msgResult.Ok {
			cacheTmp := &cache.AutoSendMsgRecord{}
			cacheTmp.MessageId = msgResult.MessageId
			cacheTmp.Account = account
			cacheTmp.Target = target
			cacheTmp.IsRead = 1
			cacheTmp.IsArrived = 0
			cache.SetAutoSendMsgRecord(cacheTmp)
			cache.IncSendMsgTaskInfoCount(cache.SuccessNum, account, 1)
		}

		// 失败回收数据
		if err1 != nil || !msgResult.Ok {
			cache.SaddDataPackList(config.DataPackId, target)
		}

		// ❗错误处理
		if err1 != nil {
			errCount++
			logs.Error("发送失败 errCount=%d, err=%v", errCount, err1)

			// 连续5次错误直接退出
			if errCount >= 5 {
				logs.Error("连续5次发送失败，停止任务 account=%s", account)
				return
			}
		} else {
			// 成功则清零
			errCount = 0
		}

		// 账号状态检测
		if cache.GetAccountStatus(account) != 2 {
			return
		}
	}
}
