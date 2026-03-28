package dllApi

import (
	"errors"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"serApi"
)

// 挂断电话
func CallHangup(reqparam *CallHangupReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *CallHangupRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "CallHangup", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->CallHangup failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(CallHangupRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { // 不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "CallHangup", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->CallHangup failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 通过链接进群
func GroupJoin(reqparam *GroupJoinReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *GroupJoinRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "GroupJoin", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->GroupJoin failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(GroupJoinRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { //不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "GroupJoin", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->GroupJoin failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 获取群的详细信息
func GroupInfo(reqparam *GroupInfoReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *GroupInfoRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "GroupInfo", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->GroupInfo failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(GroupInfoRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { //不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "GroupInfo", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->GroupInfo failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 查询注册
func PhoneQuery(reqparam *PhoneQueryReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *PhoneQueryRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "PhoneQuery", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->PhoneQuery failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(PhoneQueryRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { //不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "PhoneQuery", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->PhoneQuery failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 消息发送
func MessageSend(reqparam *MessageSendReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *MessageSendRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "MessageSend", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->MessageSend failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(MessageSendRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { // 不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "MessageSend", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->MessageSend failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 消息发送
func MessageSendAsyn(reqparam *MessageSendAsynReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *MessageSendAsynRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "MessageSendAsyn", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->MessageSendAsyn failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(MessageSendAsynRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { // 不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "MessageSendAsyn", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->MessageSendAsyn failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 创建关联的验证码
func VfcodeCreate(reqparam *VfcodeCreateReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *VfcodeCreateRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "VfcodeCreate", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->VfcodeCreate failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(VfcodeCreateRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { //不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "VfcodeCreate", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->VfcodeCreate failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 检测关联的验证码
func VfcodeCheck(reqparam *VfcodeCheckReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *VfcodeCheckRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "VfcodeCheck", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->VfcodeCheck failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(VfcodeCheckRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { //不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "VfcodeCheck", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->VfcodeCheck failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

// 检测账号
func AccountCheck(reqparam *AccountCheckReq, svrid int32, needRsp bool, timeout ...int32) (rsppara *AccountCheckRsp, err error) {
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall(serApi.ServerWechatDll, svrid, "AccountCheck", reqparam, timeout...)
		if err != nil {
			logs.Error("ServerWechatDll->AccountCheck failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, errors.New("请求错误")
		}
		rsppara = new(AccountCheckRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { //不需要回包
		rpcReq := &natsRpc.NatsMsg{}
		rpcReq.Marshal(reqparam)
		err := natsRpc.NrpcSend(serApi.ServerWechatDll, svrid, "AccountCheck", rpcReq)
		if err != nil {
			logs.Error("ServerWechatDll->AccountCheck failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}

