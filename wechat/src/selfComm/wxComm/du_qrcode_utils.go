package wxComm

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"selfComm/wxComm/wxHttp"
)

type QrcodeUtilsReq struct {
	Phone string `json:"phone"`
}

type QrcodeUtilsRsp struct {
	SessionID                 string `json:"sessionId"`
	Phone                     string `json:"phone"`
	PairingCode               string `json:"pairingCode"`
	LastPairingCodeFromServer string `json:"lastPairingCodeFromServer"`
	AlreadyConnected          bool   `json:"alreadyConnected"`
	Status                    struct {
		ID          string      `json:"id"`
		Connection  string      `json:"connection"`
		QRAvailable bool        `json:"qrAvailable"`
		User        interface{} `json:"user"` // 这里可能为 null，先用 interface{}
	} `json:"status"`
	Hint string `json:"hint"`
}

// 杜哥获取验证码接口
func QrcodeUtils(phone string) *QrcodeUtilsRsp {
	api := "http://10.0.2.50:8787/api/login"
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"

	param := &QrcodeUtilsReq{}
	param.Phone = phone

	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: param,
		Timeout: 30,
	})
	ret := &QrcodeUtilsRsp{}
	if rsp.Err == nil {
		logs.Info("QrcodeUtils result: " + string(rsp.Body))
		jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	} else {
		logs.Info("QrcodeUtils err  phone: " + phone + rsp.Err.Error())
	}

	return ret
}
