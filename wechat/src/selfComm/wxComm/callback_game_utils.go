package wxComm

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"selfComm/wxComm/wxHttp"
)

// 杜哥获取验证码接口
func CallbackGameUtils(param interface{}) {
	api := "https://d1bhdk5en1ht2p.cloudfront.net/api/activity/whatsapp/postDllApi"
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"
	paramStr, _ := jsoniter.MarshalToString(param)
	logs.Info("CallbackGameUtils param: ", paramStr)
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: param,
		Timeout: 60,
	})
	logs.Info("CallbackGameUtils result: ", string(rsp.Body))
}
