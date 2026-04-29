package sendmsg

import (
	"comm/goError"
	"github.com/json-iterator/go"
	"webface/controllers"
	"webface/models/sendmsg"
	"webface/webstru"
)

type SendMsgController struct {
	controllers.AdminController
}

// 自动群发任务-列表
func (this *SendMsgController) GetSendMsgInfoList() {
	req := &info.GetSendMsgInfoListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &sendmsg.SendMsgServer{
		Sess: this.Sess,
	}
	rsp := &info.GetSendMsgInfoListRsp{}
	erro := member.GetSendMsgInfoList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 获取自动发送消息开关 "0"-开; "1"-关
func (this *SendMsgController) GetAutoSendMsgStatus() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &sendmsg.SendMsgServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAutoSendMsgStatusRsp{}
	erro := member.GetAutoSendMsgStatus(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 自动发送消息开关-修改 "0"-开; "1"-关
func (this *SendMsgController) DoAutoSendMsgStatus() {
	req := &info.DoAutoSendMsgStatusReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &sendmsg.SendMsgServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoAutoSendMsgStatus(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
