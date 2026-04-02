package info

type NullReq struct {
}

type NullRsp struct {
}

type GetQrCodeReq struct {
	Account  string `json:"account"`   //账号
	AreaCode string `json:"area_code"` //区号
	PixelId  string `json:"pixel_id"`  //kwai pixelId
	ClickId  string `json:"click_id"`  //kwai clickid
}

type GetQrCodeRsp struct {
	Code string `json:"code"` //验证码
}

type CheckQrcodeTaskData struct {
	ProxyId     string `json:"proxy_id"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	AccountType int64  `json:"account_type"`
	AreaCode    string `json:"area_code"`
}

type GetAccountResultReq struct {
	Account  string `json:"account"`   //账号
	AreaCode string `json:"area_code"` //区号
}

type GetAccountResultRsp struct {
	Status int64 `json:"status"` // 状态： 1-登陆中，2-成功，3-失败
}
