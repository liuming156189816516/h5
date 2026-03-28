package event

// task 事件
const (
	//登录
	TaskTypeUserLoginEvent     = "TaskTypeUserLoginEvent"
	TaskTypeUserLoginEventBack = "TaskTypeUserLoginEventBack"

	//退出
	TaskTypeUserLogoutEvent     = "TaskTypeUserLogoutEvent"
	TaskTypeUserLogoutEventBack = "TaskTypeUserLogoutEventBack"

	//移除
	TaskTypeUserRemoveEvent     = "TaskTypeUserRemoveEvent"
	TaskTypeUserRemoveEventBack = "TaskTypeUserRemoveEventBack"

	//设置商店名称
	TaskTypeShopNameEvent     = "TaskTypeShopNameEvent"
	TaskTypeShopNameEventBack = "TaskTypeShopNameEventBack"

	//修改头像
	TaskTypeAccountHeadimgEvent     = "TaskTypeAccountHeadimgEvent"
	TaskTypeAccountHeadimgEventBack = "TaskTypeAccountHeadimgEventBack"

	//获取头像
	TaskTypeContactHeadimgEvent     = "TaskTypeContactHeadimgEvent"
	TaskTypeContactHeadimgEventBack = "TaskTypeContactHeadimgEventBack"

	//两步验证
	TaskTypeVerifyTwostepEvent     = "TaskTypeVerifyTwostepEvent"
	TaskTypeVerifyTwostepEventBack = "TaskTypeVerifyTwostepEventBack"

	//发送消息
	TaskTypeSendMsgEvent     = "TaskTypeSendMsgEvent"
	TaskTypeSendMsgEventBack = "TaskTypeSendMsgEventBack"

	//炸群
	TaskTypeZhaGroupEvent     = "TaskTypeZhaGroupEvent"
	TaskTypeZhaGroupEventBack = "TaskTypeZhaGroupEventBack"

	//处理消息发送结果
	AccountMessageResultEvent = "AccountMessageResultEvent"
	//处理消息回调
	AccountMessageCallBackEvent = "AccountMessageCallBackEvent"
	//处理接听电话回调
	AccountAcceptCallCallBackEvent = "AccountAcceptCallCallBackEvent"


)

// 都有的字段
type TaskEventTaskIdInfo struct {
	Phone    string
	TaskType int64
	TaskId   string
	TaskTime int64
}

// 登录
type TaskUserLoginEventReq struct {
	Uid     string
	Account []string
	TaskEventTaskIdInfo
}

// 下线
type TaskUserLogoutEventReq struct {
	Uid     string
	Account []string
	TaskEventTaskIdInfo
}

// 移除
type TaskUserRemoveEventReq struct {
	Uid     string
	Account []string
	TaskEventTaskIdInfo
}

// 设置商店名称
type TaskTypeShopNameEventReq struct {
	Uid      string
	Account  string
	NickName string
	LogId    string //修改昵称日志id
	IsCheck  bool   //是否检测 false 否 true 是
	TaskEventTaskIdInfo
}

// 发送消息
type TaskTypeSendMsgEventReq struct {
	Uid           string
	SendMsgTaskId string
	Account       string
	Target        string
	Lid           string
	Type          int64 // 1-文字 2-图片 3-音频 4-视频 5-语音呼叫 6-名片
	Content       string
	Remark        string //图片对应的描述
	DataPackId    string

	//================
	IsUp              bool //是否发送过消息 false-否 true-是
	SendMsgTaskInfoId string
	MaterialListStr   string
	MaterialType      int64
	UpErrNum          int64 `json:"up_err_num"` //上传通讯录失败次数
	SendNum           int64
	TaskEventTaskIdInfo
}

type ReceiveMessagesReq struct {
	Id      string `json:"id"`
	Account string `json:"account"`
	From    string `json:"from"`
	Type    string `json:"type"`
	Ctype   string `json:"ctype"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
	To      string `json:"to"`
	Read    bool   `json:"read"`
	FromType int64 `json:"from_type"` //from是手机号还是lid 0-lid 1-手机号
	TaskEventTaskIdInfo
}

type TaskTypePullGroupEventReq struct {
	Uid         string
	Step        int64 //当前步骤 0-通讯录 1-建群 2-发送消息
	Qid         string
	Account     string
	AdAccount   string
	SendAccount string
	YmAccount   string
	Qname       string
	Ad          string
	Targets     []string
	UpList      []string          //上传通讯录
	UpMap       map[string]string //jid-num对应关系
	//Phones          []string          //数据列表
	Members         []string //群成员
	ErrMembers      []string //异常的群成员
	PullGroupInfoId string
	PullNum         int64
	DataPackId      string
	ErrNum          int64
	SendMsgErr      int64
	PullGroupTaskId string
	IsAgain         int    //是否重发 0-否 1-是
	IsAnnouncement  int64  //是否禁言 0-否 1-是
	Mindex          int    //素材索引
	AutoAd          string //炒群话术
	InviteLink      string //邀请链接
	TaskEventTaskIdInfo
}

type TaskTypeBigGroupEventReq struct {
	Uid            string
	BigGroupTaskId string
	Step           int64
	//AdmingAccount  string            //管理员
	PullAccountList []string //拉手
	AdAccount       string   //营销号
	Ad              string   //话术
	Qurl            string   //邀请链接
	Qid             string   //群id
	DataPackId      string   //数据包id
	UsePullAccount  []string //已使用过的拉手
	PullNum         int64
	EorrList        []string //异常的数据
	MemberList      []string //群成员
	TargetNum       int64    //目标人数
	TarList         []string //全部的数据
	Mindex          int      //素材索引
	AutoAd          string   //炒群话术
	TaskEventTaskIdInfo
}

// 炸群
type TaskTypeZhaGroupEventReq struct {
	Ad          string
	SendAccount string
	Qid         string
	AutoAd      string //炒群话术
	Mindex      int    //素材索引
	TaskEventTaskIdInfo
}


