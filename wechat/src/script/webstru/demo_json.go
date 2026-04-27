package info

type DemoReq struct {
	Status string `json:"status"` //"0" 开启，"1" 关闭
	Phone  string `json:"phone"`  //手机号
	Type   string `json:"type"`
}

type DemoRsp struct {
	Message string `json:"message"` //定时任务已开启/定时任务已关闭
}

type FbData struct {
	Fbclid   string `json:"fbclid"`
	ClickId  string `json:"click_id"` //kwai clickid
	Fbp      string `json:"fbp"`
	PixelId  string `json:"pixel_id"`
	Account  string `json:"account"`   //账号（获取验证码）
	AreaCode string `json:"area_code"` //区号（获取验证码）
}

type ApiReq struct {
	Ptype   int64       `json:"ptype"`   //1-账号 2-消息
	Account string      `json:"account"` //手机号
	Data    interface{} `json:"data"`
}
