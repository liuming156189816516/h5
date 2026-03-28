package dataApi

type AccountLoginResultData struct {
	Errno   int64  `json:"errno"`
	Account string `json:"account"`
	Errmsg  string `json:"errmsg"`
	Data   struct{
		Token string
	} `json:"data"`
}

type AccountLoginResultRsp struct {
}

type AccountSendMsgResultData struct {
	Errno  int64  `json:"errno"`
	Data   string `json:"data"`
	Errmsg string `json:"errmsg"`
	To     string `json:"to"`
}

type AccountSendMsgResultRsp struct {
}
