package info

type GetSendMsgInfoListReq struct {
	Account       string `json:"account"`        //账号查询
	AccountStatus int64  `json:"account_status"` //账号状态 1-离线 2-在线
	StartTime     int64  `json:"start_time"`     //开始时间
	EndTime       int64  `json:"end_time"`       //结束时间
	AccountGroup  string `json:"account_group"`  //账号分组
	Page          int64  `json:"page"`
	Limit         int64  `json:"limit"`
}

type GetSendMsgInfoListRsp struct {
	Total        int64                     `json:"total"`
	List         []*GetSendMsgInfoListInfo `json:"list"`
	SuccessCount int64                     `json:"success_count"` //发送完成总数
	ArrivedCount int64                     `json:"arrived_count"` //送达完成总数
	ReadCount    int64                     `json:"read_count"`    //已读完成总数
	Average      int64                     `json:"average"`       //平均发送数
	ArrivedRate  int64                     `json:"arrived_rate"`  //送达率
	ReadRate     int64                     `json:"read_rate"`     //已读率
}
type GetSendMsgInfoListInfo struct {
	Id            string `json:"id"`
	Account       string `json:"account"`         //账号
	AccountStatus int64  `json:"account_status" ` //账号状态 1-离线 2-在线
	SucessNum     int64  `json:"sucess_num"`      //已完成数量
	ArrivedNum    int64  `json:"arrived_num"`     //已送达
	ReadNum       int64  `json:"read_num"`        //已读
	Reason        string `json:"reason"`          //原因
	Itime         int64  `json:"itime"`           //创建时间
	Ptime         int64  `json:"ptime"`           //更新时间
}

type DoAutoSendMsgStatusReq struct {
	AutoSendMsgStatus string `json:"auto_send_msg_status"` //自动发送消息开关 "0" - 开; "1" - 关
}

type GetAutoSendMsgStatusRsp struct {
	AutoSendMsgStatus string `json:"auto_send_msg_status"` //自动发送消息开关 "0" - 开; "1" - 关
}
