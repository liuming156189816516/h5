package cache

// 账号信息
type AccountInfo struct {
	Account      string `json:"account"`       //账号
	AccountType  int64  `json:"account_type"`  //账号类型 1-个人号 2-商业号
	Token        string `json:"token"`         //token
	PlatformType int64  `json:"platform_type"` //平台类型 1-云控 2-APP
	Synckeys     string `json:"synckeys"`
	NickName     string `json:"nick_name"` //昵称
	LogId        string `json:"log_id"`    //修改昵称日志id
	PixelId      string `json:"pixel_id"`  //kwai pixelId
	ClickId      string `json:"click_id"`  //kwai clickid
}

type ProxyIpInfo struct {
	Type       string `json:"type"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Pwd        string `json:"pwd"`
	Country    string `json:"country"`     //国家
	IpCategory int64  `json:"ip_category"` //ip类别 1-静态ip 2-动态ip
	IpType     int64  `json:"ip_type"`     //ip类型 1-ipV4 2-ipv6
	IpId       string `json:"ip_id"`
}

type SendMsgTaskInfo struct {
	Account      string     `json:"account"`       //账号
	DataPackId   string     `json:"data_pack_id"`  //数据包id
	MaterialList []Material `json:"material_list"` //素材列表
}

type HangUpAccount struct {
	Uid     string `json:"uid"`
	Account string `json:"account"`
}

// 素材
type Material struct {
	MaterialId string `json:"material_id"` //素材id
	Content    string `json:"content"`     //内容
	Type       int64  `json:"type"`        //类型 1-文字 2-图片 3-语音 4-视频 5-语音呼叫 6-名片
	Remark     string `json:"remark"`
}

type SendMsgRecord struct {
	DataPackId   string     `json:"data_pack_id"`
	RecordId     string     `json:"record_id"`
	Account      string     `json:"account"`
	Target       string     `json:"target"`
	Lid          string     `json:"lid"`
	Content      string     `json:"content"`       //消息内容
	SendTime     int64      `json:"send_time"`     //推送时间
	IsRead       int64      `json:"is_read"`       //是否已读 1-未读 2-已读
	IsArrived    int64      `json:"is_arrived"`    //是否送达 0-未送达 2-已送达
	Status       int64      `json:"status"`        //1-发送中 2-发送失败 3-发送成功
	MsgStatus    int64      `json:"msg_status"`    //发送状态 1-发送成功 2-发送失败
	MaterialList []Material `json:"material_list"` //素材列表
	MaterialType int64      `json:"material_type"` //类型 1-文字 2-图片 3-语音 4-视频 5-语音呼叫 6-名片
}

type GlobalSetting struct {
	IsAutoTranslate int64  `json:"is_auto_translate"` //1-否 2-是
	ReceiveTarget   string `json:"receive_target"`    //接收消息翻译语言
	ChatTranslate   int64  `json:"chat_translate"`    //聊天框翻译类型 1-手动翻译 2-自动翻译
	SendTarget      string `json:"send_target"`       //聊天发送语言
}

// 账号信息
type SysConfig struct {
	AutoPullGroup int64 `json:"auto_pull_group"` //是否炸群 0-否 1-是
}

// =======================================================
// 账户信息
type AppUserInfo struct {
	Uid          string  `json:"uid"`
	InviteCode   string  `json:"invite_code"`   //当前邀请码
	FinviteCode  string  `json:"finvite_code"`  //父级邀请码
	Fuid         string  `json:"fuid"`          //父级id
	TinviteCode  string  `json:"tinvite_code"`  //顶级邀请码
	Tuid         string  `json:"tuid"`          //顶级id
	DefaultGid   string  `json:"default_gid"`   //未分组id
	TutorialList []int64 `json:"tutorial_list"` //教程,我知道了
}

// 任务配置
type TaskConfigInfo struct {
	DataPackId   string     `json:"data_pack_id"`  //数据包id
	Link         string     `json:"link"`          //链接
	MaterialList []Material `json:"material_list"` //素材列表
}

type CreateGroupSettlementInfo struct {
	Uid         string            `json:"uid"`
	AppUid      string            `json:"app_uid"`
	Fuid        string            `json:"fuid"`
	TaskInfoId  string            `json:"task_info_id"`
	InviteLink  string            `json:"invite_link"`
	InvalidTime int64             `json:"invalid_time"` //失效时间
	Targets     string            `json:"targets"`
	Status      int64             `json:"status"` //1-开始任务 2-进行中 3-结算中 4-已结束
	DataPackId  string            `json:"data_pack_id"`
	QueryMap    map[string]string `json:"query_map"` //jid==>num
	EventTime   int64             `json:"event_time"`
	TaskType    int64             `json:"task_type"` //1-拉群 2-拉粉
}

type AutoGroupSettlementInfo struct {
	Uid         string `json:"uid"`
	AppUid      string `json:"app_uid"`
	Fuid        string `json:"fuid"`
	TaskInfoId  string `json:"task_info_id"`
	InvalidTime int64  `json:"invalid_time"` //失效时间
	Targets     string `json:"targets"`
	Status      int64  `json:"status"` //1-开始任务 2-进行中 3-结算中 4-已结束
	DataPackId  string `json:"data_pack_id"`
	SerNo       string `json:"ser_no" bson:"ser_no"` //任务编号
	Account     string `json:"account"`
	Synckeys    string `json:"synckeys"`
}

type GroupLinkSettlementInfo struct {
	Uid        string `json:"uid"`
	AppUid     string `json:"app_uid"`
	Fuid       string `json:"fuid"`
	TaskInfoId string `json:"task_info_id"`
	InviteLink string `json:"invite_link"`
}

type WithdrawListInfo struct {
	Tuid    string `json:"tuid"`
	OrderId string `json:"order_id"`
	Time    int64  `json:"time"`
	BillId  string `json:"bill_id"`
	AppUid  string `json:"app_uid"`
	Amount  int64  `json:"amount"`
}

type Duration struct {
	LastCheckTime        int64 `json:"last_check_time"`        //上一次检查时间
	TotalTime            int64 `json:"total_time"`             //挂机总时长
	JsTime               int64 `json:"js_time"`                //已经计算过的时长
	DurationBonus        int64 `json:"duration_bonus"`         //挂机收益
	ContinuousOnlineTime int64 `json:"continuous_online_time"` //连续在线时长
}

type AiApprovalTaskInfo struct {
	Tuid        string `json:"tuid"`
	Uuid        string `json:"uuid"`         //任务id
	UserAccount string `json:"user_account"` //用户ws账号
	TaskAccount string `json:"task_account"` //任务ws账号
	Time        int64  `json:"time"`         //创建任务时间
}

type AiDataPackWsTaskInfo struct {
	Tuid        string `json:"tuid"`
	Uuid        string `json:"uuid"`         //任务id
	TaskAccount string `json:"task_account"` //任务ws账号
	Time        int64  `json:"time"`
}

// 账户信息
type QrcodeInfo struct {
	QrCode string `json:"qr_code"`
	IpId   string `json:"ip_id"`
	Type   string `json:"type"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pwd    string `json:"pwd"`
	Time   int64  `json:"time"`
}

type AiApprovalTaskInfoTg struct {
	Tuid        string `json:"tuid"`
	UserId      int64  `json:"user_id"`      //tg用户Id
	UserAccount string `json:"user_account"` //用户ws账号
	Time        int64  `json:"time"`         //创建任务时间
}

type AiUploadTaskInfoWs struct {
	Tuid     string `json:"tuid"`
	TaskId   string `json:"task_id"` //任务id
	Phone    string `json:"phone"`   //手机号
	Link     string `json:"link"`    //链接
	ImgUrl   string `json:"img_url"` //图片
	Uuid     string `json:"uuid"`
	Time     int64  `json:"time"`
	DataType int64  `json:"data_type"` //1-数据号 2-监控号
}

type AiApprovalTaskInfoWs struct {
	Tuid   string `json:"tuid"`
	TaskId string `json:"task_id"` //任务id
	AppUid string `json:"app_uid"`
	Time   int64  `json:"time"` //创建任务时间
}

type AiTaskBatchResultInfoWs struct {
	Tuid                string         `json:"tuid"`
	BatchId             string         `json:"batch_id"`               //任务id
	TaskId              string         `json:"task_id"`                //任务id
	IsUploadMonitorFlag bool           `json:"is_upload_monitor_flag"` //是否上传监控号true-上传；false-未上传
	Time                int64          `json:"time"`                   //提交审批任务时间
	AiMsgInfos          []*AiMsgInfoWs `json:"ai_msg_infos"`
}

type AiMsgInfoWs struct {
	Id       string `json:"id"`
	Phone    string `json:"phone"`
	AppUid   string `json:"app_uid"`
	DataType int64  `json:"data_type"` //1-数据号 2-监控号
}

type AccountSocks5Info struct {
	Type          string `json:"type"`
	ClientIp      string `json:"client_ip"`      //客户端ip
	ClientCountry string `json:"client_country"` // 客户端国家代码
	Host          string `json:"host"`
	Port          string `json:"port"`
	User          string `json:"user"`
	Pwd           string `json:"pwd"`
}

type AccountSocks5Map struct {
	UserMap map[string]*AccountSocks5Info `json:"user_map"`
}
