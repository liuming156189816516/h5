package models

import (
	"comm/event"
	"data/models/task"
	"natsRpc"
)

// task事件
func InitTaskHandler() error {
	//登录
	err := natsRpc.HandlerNrpc(event.TaskTypeUserLoginEventBack, task.TaskUserLoginEventBackHandler)
	if err != nil {
		return err
	}

	//设置商店名称
	err = natsRpc.HandlerNrpc(event.TaskTypeShopNameEventBack, task.TaskTypeShopNameEventBackHandler)
	if err != nil {
		return err
	}

	//离线
	err = natsRpc.HandlerNrpc(event.TaskTypeUserLogoutEventBack, task.TaskTypeUserLogoutEventBackHandler)
	if err != nil {
		return err
	}

	//移除
	err = natsRpc.HandlerNrpc(event.TaskTypeUserRemoveEventBack, task.TaskTypeUserRemoveEventBackHandler)
	if err != nil {
		return err
	}

	//处理消息发送结果
	err = natsRpc.HandlerNrpc(event.AccountMessageResultEvent, task.AccountMessageResultEventHandler)
	if err != nil {
		return err
	}

	//处理已读和已回复消息回调
	err = natsRpc.HandlerNrpc(event.AccountMessageCallBackEvent, task.AccountMessageCallBackEventHandler)
	if err != nil {
		return err
	}

	//处理接听电话回调
	err = natsRpc.HandlerNrpc(event.AccountAcceptCallCallBackEvent, task.AccountAcceptCallCallBackEventHandler)
	if err != nil {
		return err
	}

	//消息发送
	err = natsRpc.HandlerNrpc(event.TaskTypeSendMsgEventBack, task.TaskTypeSendMsgEventBackHandler)
	if err != nil {
		return err
	}

	return err
}
