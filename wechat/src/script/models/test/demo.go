package test

import (
	"comm/comm"
	"comm/goError"
	info "script/webstru"
	"selfComm/wxComm"
)

// 群发
type DemoServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *DemoServer) getUid() string {
	return this.Sess.Uid
}

func (this *DemoServer) Demo(req1 *info.DemoReq, rsp *info.DemoRsp) *goError.ErrRsp {
	wxComm.KwaiPlace("iy7RRvunG1KahufyqVEuAg", "EVENT_BUTTON_CLICK")
	return nil
}
