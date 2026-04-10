package info

type DemoReq struct {
	Status string `json:"status"` //"0" 开启，"1" 关闭
}

type DemoRsp struct {
	Message string `json:"message"` //定时任务已开启/定时任务已关闭
}
