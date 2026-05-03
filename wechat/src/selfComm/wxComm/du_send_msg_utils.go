package wxComm

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"selfComm/wxComm/cache"
	"selfComm/wxComm/wxHttp"
)

type SendMsgUtilsReq struct {
	To                    string      `json:"to"`
	Text                  string      `json:"text"`
	PreviewType           string      `json:"previewType"`
	ContextInfo           ContextInfo `json:"contextInfo"`
	InviteLinkGroupTypeV2 string      `json:"inviteLinkGroupTypeV2"`
}

type ContextInfo struct {
	ForwardingScore int             `json:"forwardingScore"`
	IsForwarded     bool            `json:"isForwarded"`
	ExternalAdReply ExternalAdReply `json:"externalAdReply"`
	ForwardOrigin   string          `json:"forwardOrigin"`
}

type ExternalAdReply struct {
	Title                         string `json:"title"`
	Body                          string `json:"body"`
	MediaType                     string `json:"mediaType"`
	MediaURL                      string `json:"mediaUrl"`
	Thumbnail                     string `json:"thumbnail"`
	ContainsAutoReply             bool   `json:"containsAutoReply"`
	RenderLargerThumbnail         bool   `json:"renderLargerThumbnail"`
	ShowAdAttribution             bool   `json:"showAdAttribution"`
	ClickToWhatsappCall           bool   `json:"clickToWhatsappCall"`
	AdContextPreviewDismissed     bool   `json:"adContextPreviewDismissed"`
	AutomatedGreetingMessageShown bool   `json:"automatedGreetingMessageShown"`
	DisableNudge                  bool   `json:"disableNudge"`
	WtwaAdFormat                  bool   `json:"wtwaAdFormat"`
}

type SendMsgUtilsRsp struct {
	Ok        bool   `json:"ok"`
	MessageId string `json:"messageId"`
	SessionId string `json:"sessionId"`
	Node      string `json:"node"`
}

// 杜哥发送消息接口
func SendMsgUtils(sessionId, target string, material cache.Material, node string) (*SendMsgUtilsRsp, error) {
	//api := "http://47.251.15.67:8787/api/session/" + sessionId + "/send/business-extended"
	api := "https://tls.v168.vip/api/session/" + sessionId + "/send/business-extended"
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"
	headerMap["X-Backend-Node"] = node
	advertise := Advertise{}
	jsoniter.UnmarshalFromString(material.Content, &advertise)
	image, _, _ := GetBase64ByUrl(advertise.Img)
	param := map[string]interface{}{
		"to":          target + "@s.whatsapp.net",
		"text":        advertise.Remark,
		"previewType": "NONE",
		"contextInfo": map[string]interface{}{
			"forwardingScore": 2,
			"isForwarded":     false,
			"externalAdReply": map[string]interface{}{
				"title":                         advertise.Title,
				"body":                          "",
				"mediaType":                     "VIDEO",
				"mediaUrl":                      advertise.Content,
				"thumbnail":                     image, //图片 base64，
				"containsAutoReply":             false,
				"renderLargerThumbnail":         false,
				"showAdAttribution":             false,
				"clickToWhatsappCall":           false,
				"adContextPreviewDismissed":     false,
				"automatedGreetingMessageShown": false,
				"disableNudge":                  false,
				"wtwaAdFormat":                  false,
			},
			"forwardOrigin": "UNKNOWN",
		},
		"inviteLinkGroupTypeV2": "DEFAULT",
	}
	/*paramStr, _ := jsoniter.MarshalToString(param)
	logs.Info("SendMsgUtils param: " + paramStr)*/
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: param,
		Timeout: 30,
	})
	ret := &SendMsgUtilsRsp{}
	if rsp.Err == nil {
		//logs.Info("SendMsgUtils result: " + string(rsp.Body))
		jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	} else {
		logs.Info("SendMsgUtils err  sessionId: " + sessionId + "to: " + target + rsp.Err.Error())
		return ret, rsp.Err
	}

	return ret, nil
}
