package rpcServices

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"testRpc/servicesApi"
)

/**
初始化事件
*/
func InitEvent() error {
	var err error
	err = natsRpc.NrpcHandlerEvent("testRpc", "PushUserMsg", PushUserMsgHandler)
	if err != nil {
		return err
	}
	natsRpc.ListenConfig("reloadTestRpcConfig", ReloadTestRpc)
	return nil
}

func ReloadTestRpc([]byte) {
	logs.Debug("ReloadTestRpc 收到更新事件")
}

// 事件处理
func PushUserMsgHandler(req *natsRpc.NatsMsg) int32 {
	logs.Debug("PushUserMsgHandler req:%+v", req)
	req_param := servicesApi.PushUserMsgEvent{}
	err := jsoniter.Unmarshal(req.GetMsgData(), &req_param)
	if err != nil {
		logs.Error("event is invalid:%s", err.Error())
		return natsRpc.ESMR_FAILED
	}
	logs.Debug("PushUserMsgHandler req_param:%+v", req_param)
	return natsRpc.ESMR_SUCCEED
}
