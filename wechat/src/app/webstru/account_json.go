package info

type NullReq struct {
}

type NullRsp struct {
}

type GetQrCodeRsp struct {
	Code string `json:"code"` //验证码
}

type LoginAccountReq struct {
	AccountType int64  `json:"account_type"` //账号类型 1-个人号 2-商业号
	Code        string `json:"code"`         //验证码
	Account     string `json:"account"`      //账号
	AreaCode    string `json:"area_code"`    //区号
}

type LoginAccountRsp struct {
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
	AccountType int64  `json:"account_type"` //账号类型 1-个人号 2-商业号
	Account     string `json:"account"`      //账号
	AreaCode    string `json:"area_code"`    //区号
}

type GetAccountResultRsp struct {
	Status int64 `json:"status"` // 状态： 1-登陆中，2-成功，3-失败
}
