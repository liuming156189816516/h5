package account

import (
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"webface/controllers"
	"webface/models/account"
	info "webface/webstru"
)

type AccountController struct {
	controllers.AdminController
}

// 账号分组-列表
func (this *AccountController) GetAccountGroupList() {
	req := &info.GetAccountGroupListReq{}
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
	rsp := &info.GetAccountGroupListRsp{}
	erro := member.GetAccountGroupList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 账号分组-操作
func (this *AccountController) DoAccountGroup() {
	req := &info.DoAccountGroupReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoAccountGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 账号列表
func (this *AccountController) GetAccountInfoList() {
	req := &info.GetAccountInfoListReq{}
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
	rsp := &info.GetAccountInfoListRsp{}
	erro := member.GetAccountInfoList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 移动至其他分组
func (this *AccountController) DoUpGroup() {
	req := &info.DoUpGroupReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoUpGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量导出
func (this *AccountController) DoOutPutAccount() {
	req := &info.DoOutPutAccountReq{}
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
	rsp := &info.DoOutPutAccountRsp{}
	erro := member.DoOutPutAccount(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量删除
func (this *AccountController) DoBatchDelAccount() {
	req := &info.DoBatchDelAccountReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoBatchDelAccount(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 快速上线
func (this *AccountController) DoBatchFastLogin() {
	req := &info.DoBatchFastLoginReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoBatchFastLogin(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量下线
func (this *AccountController) DoBatchLogout() {
	req := &info.DoBatchLogoutReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoBatchLogout(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 分组排序
func (this *AccountController) SortGroup() {
	req := &info.SortGroupReq{}
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
	rsp := &info.NullRsp{}
	erro := member.SortGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量上线
func (this *AccountController) DoBatchLogin() {
	req := &info.DoBatchLoginReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoBatchLogin(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 释放ip
func (this *AccountController) DoFreedIp() {
	req := &info.DoFreedIpReq{}
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
	rsp := &info.NullRsp{}
	erro := member.DoFreedIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
