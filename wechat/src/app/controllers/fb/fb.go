package fb

import (
	"app/controllers"
	"app/models/fb"
	info "app/webstru"
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
)

type FbController struct {
	controllers.AdminController
}

func (this *FbController) FbReport() {
	req := &info.FbReportReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}

	member := &fb.FbService{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.FbReport(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
