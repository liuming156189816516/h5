package dllApi

type CallHangupReq struct {
	Account string
	Target  string
	Id      string
}

type CallHangupRsp struct {
	Code int64
	Msg  string
}

type MessageSendReq struct {
	Account string
	Target  string
	Type    int64 // 1-文字 2-图片 3-音频 4-视频 5-语音呼叫 6-名片 7-删除
	Content string
	Remark  string //图片文字
	MsgId   string
}

type MessageSendRsp struct {
	Code   int64
	Msg    string
	DataId string
}

type MessageDecryptReq struct {
	Uid     string
	Account string
	Type    string
	Key     string
	Url     string
}

type MessageDecryptRsp struct {
	Content string
}

type UrlDecryptRsp struct {
	Status int    `json:"Status"`
	Body   string `json:"Body"`
}

type ContactHeadimgReq struct {
	Account string
	Target  string
}

type ContactHeadimgRsp struct {
	Content string
}

type AccountLoginReq struct {
	Uid     string
	Account string
}

type AccountLoginRsp struct {
	Code int64
	Msg  string
}

type GroupInviteReq struct {
	Account   string
	AdAccount string
	Group     string
	Targets   []string
}

type GroupInviteRsp struct {
}

type GroupJoinReq struct {
	Account string
	Code    string
}

type GroupJoinRsp struct {
	Code    int64
	Qid     string
	Message string
}

type GroupHeadimgReq struct {
	Account string
	Qid     string
	Image   string
}

type GroupHeadimgRsp struct {
	Code int64
}

type GroupNameReq struct {
	Account string
	Qid     string
	Name    string
}

type GroupNameRsp struct {
	Code int64
}

type GroupDescribeReq struct {
	Account string
	Qid     string
	Content string
}

type GroupDescribeRsp struct {
	Code int64
}

type AccountAddParamSocks5 struct {
	Type string `json:"type"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}


type GroupInfoReq struct {
	Account string `json:"account"`
	Qid     string `json:"qid"`
}

type GroupInfoRsp struct {
	Creator      string   `json:"creator"` //创建者
	Members      []string `json:"members"` //群成员ws号
	Lids         []string `json:"lids"`    //群成员的lid
	Name         string   `json:"name"`         //群名称
	Announcement bool     `json:"announcement"` //为true代表只能由管理员发言
	Id           string   `json:"id"`           //群id
}

type PhoneQueryReq struct {
	Account string   `json:"account"`
	Numbers []string `json:"numbers"`
}

type PhoneQueryRsp struct {
	Code     int64
	QueryMap map[string]string
	UpList   []string
}

type ContactListInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

type VfcodeCreateReq struct {
	Id          string
	Proxy       AccountAddParamSocks5
	AccountType int64
}

type VfcodeCreateRsp struct {
	QrCode string `json:"qr_code"`
}

type VfcodeCheckReq struct {
	Id string
}

type VfcodeCheckRsp struct {
	Code int64
	Data interface{}
}


type AccountCheckReq struct {
	Account string
}

type AccountCheckRsp struct {
	Code   int64
	ErrMsg string
}

type MessageSendAsynReq struct {
	Account string
	Target  string
	Type    int64 // 1-文字 2-图片 3-音频 4-视频 5-语音呼叫 6-名片 7-删除
	Content string
	Remark  string //图片文字
	MsgId   string
}

type MessageSendAsynRsp struct {
}

