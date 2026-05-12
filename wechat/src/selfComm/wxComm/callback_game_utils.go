package wxComm

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"selfComm/wxComm/wxHttp"
)

type ApiReq struct {
	Callback string `json:"callback"` //游戏台子回调地址
}

// 杜哥获取验证码接口
func CallbackGameUtils(param interface{}) {
	paramStr, _ := jsoniter.MarshalToString(param)
	logs.Info("CallbackGameUtils param: ", paramStr)
	apiReq := &ApiReq{}
	apiReqStr, _ := jsoniter.MarshalToString(param)
	jsoniter.UnmarshalFromString(apiReqStr, apiReq)
	if apiReq.Callback != "" {
		api := apiReq.Callback
		headerMap := make(map[string]string)
		headerMap["Content-Type"] = "application/json"
		rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
			Url:     api,
			Headers: headerMap,
			Method:  "post",
			Content: param,
			Timeout: 60,
		})
		logs.Info("CallbackGameUtils result: ", string(rsp.Body))
	}

}
