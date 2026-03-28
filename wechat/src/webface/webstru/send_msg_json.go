package info

type GetSendMsgTaskListReq struct {
	Name   string `json:"name"`
	Status int64  `json:"status"` //任务状态 1-创建成功 2-执行中 3-关闭任务 3-停止群发 4-已完成
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
}

type GetSendMsgTaskListRsp struct {
	Total int64                     `json:"total"`
	List  []*GetSendMsgTaskListInfo `json:"list"`
}
type GetSendMsgTaskListInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"` //名称
	//OnlineNum       int64      `json:"online_num"`       //在线数量
	AccountNum      int64      `json:"account_num"`      //账号总数
	SucessNum       int64      `json:"sucess_num"`       //已完成数量
	SendNum         int64      `json:"send_num"`         //群发总数
	CompleteMessage int64      `json:"complete_message"` //话术完整数量
	Status          int64      `json:"status"`           //任务状态 1-创建成功 2-执行中 3-关闭任务 3-停止群发 4-已完成
	Reason          string     `json:"reason"`           //原因
	ConfigStr       string     `json:"config_str"`       //配置信息
	MaterialList    []Material `json:"material_list"`    //素材列表
	DataPackName    string     `json:"data_pack_name"`   //数据包
	Itime           int64      `json:"itime"`            //创建时间
	ArrivedNum      int64      `json:"arrived_num"`      //已送达
	ArrivedRate     string     `json:"arrived_rate"`     //送达率
	InProNum        int64      `json:"in_pro_num"`       //执行中数量
}

type GetSendMsgInfoListReq struct {
	TaskId        string `json:"task_id"`        //必传
	Account       string `json:"account"`        //账号查询
	AccountStatus int64  `json:"account_status"` //账号状态 1-离线 2-在线
	Status        int64  `json:"status"`         //任务状态 1-创建成功 2-执行中 3-关闭任务 3-停止群发 4-已完成
	Sort          string `json:"sort"`           //已完成/群发数排序正序==>sucess_num   已完成/群发数排序倒序==>-sucess_num 已读率正序===>read_rate 已读率倒序===>-read_rate 回复率正序===>reply_rate  回复率倒序===>-reply_rate 接收消息总数正序====>receive_msg_num  接收消息总数倒序====>-receive_msg_num
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

type AddSendMsgTaskReq struct {
	Name            string     `json:"name"`              //活动名称
	GroupIds        []string   `json:"group_ids"`         //账号分组
	DataType        int64      `json:"data_type"`         //ws数据类型 1-ws粉丝数据
	DataPackId      string     `json:"data_pack_id"`      //数据包id
	SendType        int64      `json:"send_type"`         //1-系统默认模式
	SendNum         int64      `json:"send_num"`          //账号每群发条数
	MinTime         int64      `json:"min_time"`          //间隔最小时间
	MaxTime         int64      `json:"max_time"`          //间隔最大时间
	SpeechSkillType int64      `json:"speech_skill_type"` //群发话术类型 1-自主添加
	MaterialList    []Material `json:"material_list"`     //素材列表
	Replenish       int64      `json:"replenish"`         //设置补发 1-不设置补发
	ConfigStr       string     `json:"config_str"`        //配置信息
}

// 素材
type Material struct {
	MaterialId string `json:"material_id"` //素材id
	Content    string `json:"content"`     //内容
	Type       int64  `json:"type"`        //类型 1-文字 2-图片 3-语音 4-视频 5-语音呼叫 6-名片
	Remark     string `json:"remark"`      //描述
}

type GetSendMsgGroupRsp struct {
	Total int64                      `json:"total"`
	List  []*GetSendMsgGroupListInfo `json:"list"`
}
type GetSendMsgGroupListInfo struct {
	GroupId   string `json:"group_id"`
	Name      string `json:"name"`       //名称
	Count     int64  `json:"count"`      //数量
	OnlineNum int64  `json:"online_num"` //在线数量
}

type DoBatchDelSendMsgTaskReq struct {
	Ids []string `json:"ids"` //ids
}

type DoBatchStopSendMsgTaskReq struct {
	Ids []string `json:"ids"` //ids
}

type DoBatchCloseSendMsgTaskReq struct {
	Ids []string `json:"ids"` //ids
}

type DoOutTaskExcelReq struct {
	Id string `json:"id"` //id
}

type DoOutTaskExcelRsp struct {
	Url string `json:"url"` //路径
}

type DoOutTaskRecordExcelReq struct {
	Id string `json:"id"` //id
}

type DoOutTaskRecordExcelRsp struct {
	Url string `json:"url"` //路径
}

type GetDataSumReq struct {
	Id string `json:"id"` //id
}

type GetDataSumRsp struct {
	AccountNum       int64   `json:"account_num"`        //账号总数
	AccountSucessNum int64   `json:"account_sucess_num"` //已完成账号
	TargetNum        int64   `json:"target_num"`         //目标粉丝数
	PushNum          int64   `json:"push_num"`           //成功推送粉丝
	ReadRate         float64 `json:"read_rate"`          //已读率
	ReplyRate        float64 `json:"reply_rate"`         //回复率
	ReceiveMsgNum    int64   `json:"receive_msg_num"`    //接受消息总数
	ArrivedNum       int64   `json:"arrived_num"`        //已送达
}
