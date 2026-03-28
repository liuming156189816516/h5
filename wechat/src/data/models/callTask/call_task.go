package callTask

import (
	"comm/comm"
	"comm/event"
	"encoding/json"
	"fmt"
	"mlog"
	"natsRpc"
	"selfComm/wxComm/cache"
	"time"
)

//消息回报 总
type ResponseResult struct {
	Code     int64
	Success  bool
	Message  string
	Data     interface{}
	Debug    string
	Wxid     string
	PlanUuid string
	TaskId   string
	Req      interface{}
	SaveData interface{}
}

//事件发送处理
func CallUserOneTask(phone string, task_type int64, data map[string]interface{}, ts ...int32) *ResponseResult {

	gettaskId := ""
	if id, ok := data["TaskId"]; ok { //有task id
		t, ok := id.(string)
		if ok && t != "" {
			gettaskId = t
		}
	}
	if gettaskId == "" {
		taskId := event.GetTaskId(task_type)
		if taskId == "" {
			mlog.LogWxError(phone, "解析错误:%+v", data)
			return &ResponseResult{Code: -100, Message: fmt.Sprintf("解析错误:%+v", data)}
		}
		data["TaskId"] = taskId
		gettaskId = taskId
	}
	data["TaskId"] = gettaskId
	data["Phone"] = phone
	data["TaskType"] = task_type
	data["TaskTime"] = time.Now().Unix()

	//保存任务
	cache.SetDoTask(gettaskId, data)

	estr := ""
	switch task_type {
	case event.TaskTypeUserSendMsg:
		estr = event.TaskTypeSendMsgEvent
	}
	if estr == "" {
		mlog.LogWxError(phone, "任务类型错误", task_type)
		return &ResponseResult{Code: -100, Message: fmt.Sprintf("任务类型错误:%d", task_type)}
	}
	ttl := int32(0)
	if len(ts) > 0 {
		ttl = ts[0]
	}
	//直接转发
	ret, err := natsRpc.NrpcCall(comm.ServerWechatDll, -1, estr, data, ttl)
	if err != nil {
		mlog.LogWxError(phone, "请求错误%+v,ret=%+v", err, ret)
		return &ResponseResult{Code: -100, Message: fmt.Sprintf("请求错误%+v,ret=%+v", err, ret)}
	}
	rsp := &ResponseResult{}
	if ret != nil && len(ret.MsgData) != 0 {
		err = json.Unmarshal(ret.MsgData, &rsp)
		if err != nil {
			mlog.LogWxError(phone, "解析错误%+v", err)
			return &ResponseResult{Code: -100, Message: fmt.Sprintf("解析错误%+v", err)}
		}
	}
	return rsp
}
