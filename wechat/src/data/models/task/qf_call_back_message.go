package task

import (
	"comm/event"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"natsRpc"
	"selfComm/wxComm/cache"
	"serApi/dataApi"
	"sync"
)

var lock = sync.Mutex{}

func AccountMessageCallBackEventHandler(msg *natsRpc.NatsMsg) int32 {
	//错误
	req := &event.ReceiveMessagesReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	defer msg.ResponeSucc(&dataApi.AccountSendMsgResultRsp{})

	//删除任务
	cache.DelDoTask(req.TaskId)

	target := req.From
	if req.FromType == 1 {
		target = cache.GetSendMsgPhoneLid(target)
	}
	lock.Lock()
	defer lock.Unlock()
	record := cache.GetSendMsgRecordInfo(req.Account + "_" + target)
	if req.Type == "notify" && req.Ctype == "receipt" && req.Read == false {
		//已送达
		if record.IsArrived == 0 {
			record.IsArrived = 1
			cache.IncSendMsgTaskInfoCount(cache.ArrivedNum, record.Account, record.Account, 1)
		}
	}
	cache.SetSendMsgRecord(record)
	return natsRpc.ESMR_SUCCEED
}
