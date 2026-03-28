package account

import (
	"app/controllers"
	"app/models/account"
	info "app/webstru"
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
)

type AccountController struct {
	controllers.AdminController
}

// 获取二维码
func (this *AccountController) GetQrCode() {
	req := &info.GetQrCodeReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.GetQrCodeRsp{}
	erro := member.GetQrCode(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 获取ws列表
func (this *AccountController) GetAccountList() {
	req := &info.GetAccountResultReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAccountResultRsp{}
	erro := member.GetAccountResult(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
