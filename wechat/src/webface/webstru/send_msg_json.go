package info

type GetSendMsgInfoListReq struct {
	Account       string `json:"account"`        //账号查询
	AccountStatus int64  `json:"account_status"` //账号状态 1-离线 2-在线
	Page          int64  `json:"page"`
	Limit         int64  `json:"limit"`
}

type GetSendMsgInfoListRsp struct {
	Total int64                     `json:"total"`
	List  []*GetSendMsgInfoListInfo `json:"list"`
}
type GetSendMsgInfoListInfo struct {
	Id            string `json:"id"`
	Account       string `json:"account"`         //账号
	AccountStatus int64  `json:"account_status" ` //账号状态 1-离线 2-在线
	SucessNum     int64  `json:"sucess_num"`      //已完成数量
	ArrivedNum    int64  `json:"arrived_num"`     //已送达
	Reason        string `json:"reason"`          //原因
}
