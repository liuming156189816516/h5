package info

type FbReportReq struct {
	Ptype int64  `json:"ptype"`
	Data  string `json:"data"` // JSON字符串
}

type FbData struct {
	Fbclid   string `json:"fbclid"`
	Fbp      string `json:"fbp"`
	PixelId  string `json:"pixel_id"`
	Account  string `json:"account"`   //账号（获取验证码）
	AreaCode string `json:"area_code"` //区号（获取验证码）
}
