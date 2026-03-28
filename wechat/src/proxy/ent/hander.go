package ent

import (
	"comm/comm"
	"comm/event"
	"fmt"
	"natsRpc"
	"proxy/ent/api"
	"proxy/ent/task"
)

// 任务请求需要对应回复
func InitTaskEvent(os string, ver int) error {
	err := error(nil)

	//登录
	err = DllHanderNrpc(os, ver, event.TaskTypeUserLoginEvent, task.TaskUserLoginEventHandler)
	if err != nil {
		return err
	}

	//下线
	err = DllHanderNrpc(os, ver, event.TaskTypeUserLogoutEvent, task.TaskTypeUserLogoutEventHandler)
	if err != nil {
		return err
	}

	//移除
	err = DllHanderNrpc(os, ver, event.TaskTypeUserRemoveEvent, task.TaskTypeUserRemveEventHandler)
	if err != nil {
		return err
	}

	//设置商店名称
	err = DllHanderNrpc(os, ver, event.TaskTypeShopNameEvent, task.TaskTypeShopNameEventHandler)
	if err != nil {
		return err
	}

	//发送消息
	err = DllHanderNrpc(os, ver, event.TaskTypeSendMsgEvent, task.TaskTypeSendMsgEventHandler)
	if err != nil {
		return err
	}

	return err
}

// api接口
func InitApiHandler(os string, ver int) error {
	err := error(nil)

	//挂断通话
	err = DllHanderNrpc(os, ver, "CallHangup", api.CallHangupHandler)
	if err != nil {
		return err
	}

	//自主进群
	err = DllHanderNrpc(os, ver, "GroupJoin", api.GroupJoinHandler)
	if err != nil {
		return err
	}

	//群的详细信息
	err = DllHanderNrpc(os, ver, "GroupInfo", api.GroupInfoHandler)
	if err != nil {
		return err
	}

	//查询注册
	err = DllHanderNrpc(os, ver, "PhoneQuery", api.PhoneQueryHandler)
	if err != nil {
		return err
	}

	//消息发送
	err = DllHanderNrpc(os, ver, "MessageSend", api.MessageSendHandler)
	if err != nil {
		return err
	}

	//异步发送消息
	err = DllHanderNrpc(os, ver, "MessageSendAsyn", api.MessageSendAsynHandler)
	if err != nil {
		return err
	}

	//创建关联的验证码
	err = DllHanderNrpc(os, ver, "VfcodeCreate", api.VfcodeCreateHandler)
	if err != nil {
		return err
	}

	//检测关联的验证码
	err = DllHanderNrpc(os, ver, "VfcodeCheck", api.VfcodeCheckHandler)
	if err != nil {
		return err
	}

	//检测账号
	err = DllHanderNrpc(os, ver, "AccountCheck", api.AccountCheckHandler)
	if err != nil {
		return err
	}

	return err
}

func DllHanderNrpc(os string, ver int, cmd string, handler func(*natsRpc.NatsMsg) int32) error {
	if err := natsRpc.HandlerNrpcByName(comm.ServerWechatDll, cmd, handler); err != nil {
		return err
	}
	bRe := false
	if os != "" {
		cmd = "os_" + cmd
		bRe = true
	}
	if ver > 0 {
		cmd = fmt.Sprintf("%s_%d", cmd, ver)
		bRe = true
	}

	if bRe {
		return natsRpc.HandlerNrpcByName(comm.ServerWechatDll, cmd, handler)
	}
	return nil
}
