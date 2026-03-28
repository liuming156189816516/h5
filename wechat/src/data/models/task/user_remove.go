package task

import (
	"comm/event"
	"data/models/stru"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"selfComm/wxComm/cache"
)

//账号移除
func TaskTypeUserRemoveEventBackHandler(msg *natsRpc.NatsMsg) int32 {
	//解析
	rsp := &stru.ResponseResult{}
	err := json.Unmarshal(msg.MsgData, &rsp)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc("")
	if rsp.TaskId == "" {
		logs.Error("数据错误", rsp)
		return natsRpc.ESMR_SUCCEED
	}

	//解析req
	req := &event.TaskUserRemoveEventReq{}
	reqStr, _ := jsoniter.MarshalToString(rsp.Req)
	err = jsoniter.UnmarshalFromString(reqStr, &req)

	//删除任务
	cache.DelDoTask(req.TaskId)

	return natsRpc.ESMR_SUCCEED
}
