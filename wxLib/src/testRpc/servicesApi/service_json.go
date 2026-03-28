package servicesApi

import "natsRpc"

type GetUserReq struct {
	Uid string `json:"uid"`
}

type GetUserRsp struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

type PushUserMsgEvent struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func (event *PushUserMsgEvent) PushEvent() error {
	return natsRpc.NatsPulishEvent("PushUserMsg", event)
}
