package stru

//消息回报 总
type ResponseResult struct {
	Code    int64
	Success bool
	Message string
	Data    interface{}

	Account string
	TaskId  string
	Req     interface{}
}

type SendMsgData struct {
	Id      string              `json:"id"`
	Submits map[string][]string `json:"submits"`
}
