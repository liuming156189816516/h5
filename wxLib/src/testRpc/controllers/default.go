package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"testRpc/servicesApi"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

/**
示例一个rpc请求
beego url 192.168.1.24:8080/main/getuserinfo
请求参数
{
    "uid": "1003"
}
返回
{
    "uid": "1003",
    "name": "1003_test"
}
*/
func (c *MainController) GetUserInfo() {
	req := &servicesApi.GetUserReq{}
	err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Debug("err:%+v, body:%s", err, string(c.Ctx.Input.RequestBody))
		c.Ctx.WriteString("Unmarshal err")
		return
	}
	logs.Debug("MainController GetUserInfo")
	srvApi := servicesApi.NewServiceApi(&natsRpc.Session{})
	rsp, err := srvApi.GetUserInfo(req, -1, true)
	if err != nil {
		logs.Debug("GetUserInfo err:%+v", err)
	}
	logs.Debug("rsp:%+v, err:%+v", rsp, err)
	str, _ := jsoniter.MarshalToString(rsp)
	c.Ctx.WriteString(str)
	return
}
