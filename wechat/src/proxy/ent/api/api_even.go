package api

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"natsRpc"
	"serApi/dllApi"
)

// 挂断通话
func CallHangupHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.CallHangupReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.CallHangupRsp{}
	err = CallHangup(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 自主进群
func GroupJoinHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.GroupJoinReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.GroupJoinRsp{}
	err = GroupJoin(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 自主进群
func GroupInfoHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.GroupInfoReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.GroupInfoRsp{}
	err = GroupInfo(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 查询注册
func PhoneQueryHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.PhoneQueryReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.PhoneQueryRsp{}
	err = PhoneQuery(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 消息发送
func MessageSendHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.MessageSendReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.MessageSendRsp{}
	err = MessageSend(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 创建关联的验证码
func VfcodeCreateHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.VfcodeCreateReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.VfcodeCreateRsp{}
	err = VfcodeCreate(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 检测关联的验证码
func VfcodeCheckHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.VfcodeCheckReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.VfcodeCheckRsp{}
	err = VfcodeCheck(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

// 检测账号
func AccountCheckHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.AccountCheckReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.AccountCheckRsp{}
	err = AccountCheck(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}

func MessageSendAsynHandler(msg *natsRpc.NatsMsg) int32 {
	req := &dllApi.MessageSendAsynReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	rsp := dllApi.MessageSendAsynRsp{}
	err = MessageSendAsyn(req, &rsp)
	if err != nil {
		msg.Response(-1000, fmt.Sprintf("请求错误:%s", err.Error()))
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc(rsp)
	return natsRpc.ESMR_SUCCEED
}
