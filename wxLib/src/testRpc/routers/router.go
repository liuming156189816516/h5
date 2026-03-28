package routers

import (
	"github.com/astaxie/beego"
	"testRpc/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.AutoRouter(&controllers.MainController{})
}
