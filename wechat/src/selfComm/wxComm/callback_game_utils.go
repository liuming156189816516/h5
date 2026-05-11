package wxComm

import (
	"github.com/astaxie/beego/logs"
	"selfComm/wxComm/wxHttp"
)

// 杜哥获取验证码接口
func CallbackGameUtils(param interface{}) {
	api := "https://front-api.sancatalyst.com/api/activity/whatsapp/postDllApi"
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: param,
		Timeout: 30,
	})
	logs.Info(string(rsp.Body))
}
