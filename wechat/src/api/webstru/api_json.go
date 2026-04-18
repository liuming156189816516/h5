package info

type NullRsp struct {
}

type ApiReq struct {
	Ptype   int64       `json:"ptype"`   //1-账号 2-消息
	Account string      `json:"account"` //手机号
	Data    interface{} `json:"data"`
}

type FbData struct {
	PixelId string `json:"pixel_id"`
	Phone   string `json:"phone"` //手机号
}
