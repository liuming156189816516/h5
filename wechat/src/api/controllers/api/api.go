package api

import (
	"api/controllers"
	"api/models/api"
	info "api/webstru"
	"comm/goError"
	"github.com/json-iterator/go"
)

type DemoController struct {
	controllers.AdminController
}

func (this *DemoController) Demo() {
	req := &info.ApiReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &api.ApiServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.Api(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
