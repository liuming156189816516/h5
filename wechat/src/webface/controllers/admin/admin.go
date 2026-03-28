package admin

import (
	"comm/goError"
	"github.com/json-iterator/go"
	"webface/controllers"
	"webface/models/admin"
	"webface/webstru"
)

type AdminMemberController struct {
	controllers.AdminController
}

// 登录
func (this *AdminMemberController) Login() {
	req := &info.AdminLoginReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.AdminLoginRsp{}
	erro := member.Login(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 退出登录
func (this *AdminMemberController) LoginOut() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.LoginOut(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 用户菜单
func (this *AdminMemberController) Menu() {
	req := &info.MenuReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.MenuRsp{}
	erro := member.Menu(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 菜单列表
func (this *AdminMemberController) AllMenu() {
	req := &info.MenuReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.MenuRsp{}
	erro := member.AllMenu(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 编辑菜单
func (this *AdminMemberController) DoMenu() {
	req := &info.DoMenuReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoMenu(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 角色列表
func (this *AdminMemberController) RoleList() {
	req := &info.RoleListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.RoleListRsp{}
	erro := member.RoleList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 编辑角色
func (this *AdminMemberController) DoRole() {
	req := &info.DoRoleReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoRole(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 管理员
func (this *AdminMemberController) AdminUser() {
	req := &info.GetAdminUserListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAdminUserListRsp{}
	erro := member.GetAdminUserList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 用户-操作
func (this *AdminMemberController) DoAdminUser() {
	req := &info.DoAdminUserReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &admin.AdminServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoAdminUser(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
