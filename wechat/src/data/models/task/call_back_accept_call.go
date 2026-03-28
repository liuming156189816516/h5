package task

import (
	"comm/event"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"natsRpc"
	"selfComm/wxComm/cache"
	"serApi/dataApi"
	"serApi/dllApi"
)

func AccountAcceptCallCallBackEventHandler(msg *natsRpc.NatsMsg) int32 {
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

	target := req.To

	//挂断电话
	req1 := &dllApi.CallHangupReq{
		Account: req.Account,
		Target:  target,
		Id:      req.Id,
	}
	dllApi.CallHangup(req1, -1, false)
	return natsRpc.ESMR_SUCCEED
}
