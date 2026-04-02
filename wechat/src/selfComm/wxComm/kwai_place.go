package wxComm

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"selfComm/wxComm/wxHttp"
)

var PixId = "307096390156304"
var AccessToken = "ini0kAtt8me-1K8i6LW39rp5TLtBkKsd8-7LHBjxN5A"

type KwaiPlaceReq struct {
	ClickId         string `json:"clickid"`
	EventName       string `json:"event_name"`
	PixelId         string `json:"pixelId"`
	AccessToken     string `json:"access_token"`
	TestFlag        bool   `json:"testFlag"`
	TrackFlag       bool   `json:"trackFlag"`
	IsAttributed    int64  `json:"is_attributed"`
	MmpCode         string `json:"mmpcode"`
	PixelSdkVersion string `json:"pixelSdkVersion"`
}

type KwaiPlaceRsp struct {
	Result int64  `json:"result"`
	Msg    string `json:"error_msg"`
}

// kwai投放
func KwaiPlace(clickId, eventName string) *KwaiPlaceRsp {
	api := "https://www.adsnebula.com/log/common/api"
	headerMap := make(map[string]string)
	headerMap["accept"] = "application/json;charset=utf-8"
	headerMap["Content-Type"] = "application/json"

	param := &KwaiPlaceReq{}
	param.ClickId = clickId
	param.EventName = eventName
	param.PixelId = PixId
	param.AccessToken = AccessToken
	param.TestFlag = false
	param.TrackFlag = false
	param.IsAttributed = 1
	param.MmpCode = "PL"
	param.PixelSdkVersion = "9.9.9"

	paramStr, _ := jsoniter.MarshalToString(param)
	logs.Info("KwaiPlace param: " + paramStr)

	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: param,
		Timeout: 30,
	})
	ret := &KwaiPlaceRsp{}
	if rsp.Err == nil {
		logs.Info("KwaiPlace result: " + string(rsp.Body))
		jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	} else {
		logs.Info("KwaiPlace err  clickId: " + clickId + "eventName: " + eventName + "err： " + rsp.Err.Error())
	}

	return ret
}
