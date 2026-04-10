package test

import (
	"comm/comm"
	"comm/goError"
	info "script/webstru"
	"selfComm/wxComm/cache"
)

// 群发
type DemoServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *DemoServer) getUid() string {
	return this.Sess.Uid
}

func (this *DemoServer) Demo(req *info.DemoReq, rsp *info.DemoRsp) *goError.ErrRsp {
	return nil
}

func (this *DemoServer) CheckStatus(req *info.DemoReq, rsp *info.DemoRsp) *goError.ErrRsp {
	if req.Status == "0" || req.Status == "1" {
		cache.SetTaskStatus(req.Status)
	} else {
		rsp.Message = "参数有误"
	}
	if req.Status == "0" {
		rsp.Message = "定时任务已开启"
	}
	if req.Status == "1" {
		rsp.Message = "定时任务已关闭"
	}
	return nil
}
