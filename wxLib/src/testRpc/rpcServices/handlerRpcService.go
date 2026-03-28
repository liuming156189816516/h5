package rpcServices

import (
	"common/goError"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"testRpc/servicesApi"
)

/**
初始化RPC服务
*/
func InitRpc() error {
	var err error
	err = natsRpc.HandlerNrpcCmd("GetUser", GetUserHandler, false)
	if err != nil {
		return err
	}
	return nil
}

func GetUserHandler(req *natsRpc.NatsMsg) int32 {
	logs.Debug("GetUserHandler req:%+v", req)
	if req == nil {
		logs.Error("req is nil")
		return natsRpc.ESMR_FAILED
	}
	logs.Debug("GetUserHandler req:%+v", req)
	req_param := servicesApi.GetUserReq{}
	err := jsoniter.Unmarshal(req.GetMsgData(), &req_param)
	if err != nil {
		logs.Debug("Unmarshal err:%+v", err)
		req.Response(goError.GLOBAL_INVALIDPARAM.Ret, goError.GLOBAL_INVALIDPARAM.Msg)
		return natsRpc.ESMR_FAILED
	}
	rsp_param := servicesApi.GetUserRsp{
		Uid:  req_param.Uid,
		Name: req_param.Uid + "_test",
	}
	req.ResponeSucc(rsp_param)

	// 发送事件
	pushUserMsg := servicesApi.PushUserMsgEvent{
		Uid:  "123",
		Name: "456",
		Msg:  "789",
	}
	if err := pushUserMsg.PushEvent(); err != nil {
		logs.Debug("PushEvent err:%+v", err)
	}

	// reload config
	natsRpc.NotifyReloadConfig("reloadTestRpcConfig", nil)

	return natsRpc.ESMR_SUCCEED
}
