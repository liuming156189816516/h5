package test

import (
	"comm/goError"
	"github.com/json-iterator/go"
	"script/controllers"
	"script/models/test"
	info "script/webstru"
)

type DemoController struct {
	controllers.AdminController
}

// 测试
func (this *DemoController) Demo() {
	req := &info.DemoReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &test.DemoServer{
		Sess: this.Sess,
	}
	rsp := &info.DemoRsp{}
	erro := member.Demo(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 定时任务开关
func (this *DemoController) CheckStatus() {
	req := &info.DemoReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &test.DemoServer{
		Sess: this.Sess,
	}
	rsp := &info.DemoRsp{}
	erro := member.CheckStatus(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
