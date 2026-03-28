// @APIVersion 1.0.0
// @Title 云控后台API
// @Description urlPrefix /v1/
package routers

import (
	"app/controllers"
	"app/controllers/account"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//跨域设置
	var FilterGateWay = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		//允许访问源
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		//允许post访问
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,ContentType,Authorization,accept,accept-encoding, authorization, content-type") //header的类型
		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	beego.InsertFilter("*", beego.BeforeRouter, FilterGateWay)
	beego.Router("*", &controllers.BaseController{}, "OPTIONS:Options")
	//beego.Router("/Shunt/?:uid/?:workId", &work.WorkShuntController{})
	InitRouter()
}
func InitRouter() {
	beego.AutoRouter(&account.AccountController{})

	beego.InsertFilter("/ae/*", beego.BeforeStatic, func(ctx *context.Context) {
		ctx.Output.Header("Cache-control", "max-age=5")
		ctx.Output.Header("Content-Type", "application/download")
	})
	beego.SetStaticPath("/ae", "out")

	beego.InsertFilter("/upload/*", beego.BeforeStatic, func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		//允许访问源
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		ctx.Output.Header("Cache-control", "max-age=5")
		ctx.Output.Header("Content-Type", "application/download")
	})
	beego.SetStaticPath("/upload", "material")
}
