package controllers

import (
	"comm/comm"
	"comm/goError"
	"comm/token"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/url"
	"strings"
	"time"
)

//登录之后得 路由
type AdminController struct {
	BaseController
	Sess    *comm.SessInfo // 当前的用户
	ModName string
	FunName string
}

//检查登录
func (this *AdminController) Prepare() {

	//ip 检查不需要
	//data := &ipInfo.AddrInfo{}
	//ipInfo.GetAwdbAddr(this.GetRealIp(), data)
	//logs.Debug("infodata ", data)
	//if data.AreaCode == "CN" {
	//	this.returnNoQuery()
	//	return
	//}

	//mod fun 权限用
	this.ModName, this.FunName = this.GetControllerAndAction()
	//避免空指针
	sess := &comm.SessInfo{}
	sess.Info = &comm.LoginInfo{}
	sess.Model = this.ModName
	sess.Func = this.FunName
	this.Sess = sess

	//检查登陆状态
	fakeAuth, _ := beego.AppConfig.Bool("fakeAuth")
	if _, ok := noLoginMap[this.FunName]; !ok && fakeAuth == false && this.checkLogin() == false { // 未登录
		this.returnNoAuth()
		return
	}
	if this.Ctx.Input.IsOptions() {
		logs.Debug("options")
	}
	runmode := this.Ctx.Input.RunMethod
	logs.Debug(this.ModName, this.FunName, runmode)
	//权限检查
	key := fmt.Sprintf("%s-%s", this.ModName, this.FunName)
	if _, ok := exceptionList[key]; !ok {
		//限制一下请求频率
		if !tryLock(this.Sess.Uid, key) {
			this.returnTooMach()
			return
		}
	}
}

//检查登录
func (this *AdminController) checkLogin() bool {

	tk := this.GetString("token")
	if tk == "" {
		return false
	}
	tktmp, err := url.QueryUnescape(tk)

	if err == nil {
		tk = tktmp
	}
	tk = strings.Replace(tk, " ", "+", -1)
	tokenInfo := token.CheckToken(tk)
	if tokenInfo.Token != tk { //
		logs.Debug("Get tk:%s,", tk)
		return false
	}
	if tokenInfo.Src != comm.AdminSrc {
		return false
	}
	if tokenInfo.Extime < time.Now().Unix() {
		return false
	}
	this.Sess.Uid = tokenInfo.Uid
	this.Sess.DB = tokenInfo.Db
	this.Sess.ExtTime = tokenInfo.Extime
	this.Sess.Pack = tokenInfo.Pack
	this.Sess.Src = tokenInfo.Src
	this.Sess.Info.Ip = this.GetRealIp()
	this.Sess.Info.Max = ""
	return true
}

/**
没有登录
*/
func (this *AdminController) returnNoAuth() {
	this.Data["json"] = goError.CheckLoginFailed
	this.ServeJSON()
	this.StopRun()
}

/**
没有登录
*/
func (this *AdminController) returnNoQuery() {
	this.Data["json"] = goError.QueryError
	this.ServeJSON()
	this.StopRun()
}

//请求太频繁
func (this *AdminController) returnTooMach() {
	this.Data["json"] = goError.GLOBAL_TOOMACH
	this.ServeJSON()
	this.StopRun()
}

// run after finished
func (this *AdminController) Finish() {

}
