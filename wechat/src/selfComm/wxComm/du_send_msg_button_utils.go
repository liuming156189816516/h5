package wxComm

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"selfComm/wxComm/cache"
	"selfComm/wxComm/wxHttp"
)

type SendMsgButtonUtilsRsp struct {
	Ok        bool   `json:"ok"`
	Error     string `json:"error"`
	MessageId string `json:"messageId"`
	SessionId string `json:"sessionId"`
	Node      string `json:"node"`
}

// 杜哥发送按钮消息接口
func SendMsgbButtonUtils(sessionId, target string, material cache.Material, node string) (*SendMsgButtonUtilsRsp, error) {
	api := "https://tls.v168.vip/api/session/" + sessionId + "/send/interactive"
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"
	headerMap["X-Backend-Node"] = node
	advertise := Advertise{}
	jsoniter.UnmarshalFromString(material.Content, &advertise)
	param := map[string]interface{}{
		"to":     target,
		"title":  advertise.Title,
		"image":  advertise.Img,
		"text":   advertise.Remark,
		"footer": "",
		"url":    advertise.Content,
	}
	paramStr, _ := jsoniter.MarshalToString(param)
	logs.Info("SendMsgbButtonUtils req target:" + target + " sessionId: " + sessionId + " node: " + node + " param :" + paramStr)
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: param,
		Timeout: 60,
	})
	ret := &SendMsgButtonUtilsRsp{}
	logs.Info("SendMsgbButtonUtils rsp target:" + target + " result: " + string(rsp.Body))
	if rsp.Err == nil {
		jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	} else {
		//logs.Info("SendMsgbButtonUtils err target:" + target + " sessionId: " + sessionId + "to: " + target + rsp.Err.Error())
		return ret, rsp.Err
	}

	return ret, nil
}
