package info

type FbReportReq struct {
	Ptype int64  `json:"ptype"` //1-内容查看 2-获取验证码 3-注册成功
	Data  string `json:"data"`  // JSON字符串
}

type FbData struct {
	Fbclid   string `json:"fbclid"`
	ClickId  string `json:"click_id"` //kwai clickid
	Fbp      string `json:"fbp"`
	PixelId  string `json:"pixel_id"`
	Account  string `json:"account"`   //账号（获取验证码）
	AreaCode string `json:"area_code"` //区号（获取验证码）
}
