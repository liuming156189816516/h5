package info

type DemoReq struct {
	Account   []string `json:"account"`
	Numbers   []string `json:"numbers"`
	Group     string   `json:"group"`
	AdAccount string   `json:"ad_account"`
	Ip        string   `json:"ip"`
	Ua        string   `json:"ua"`
	Type      string   `json:"type"`
}

type DemoRsp struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ExportUserContactReq struct {
	Key       string `json:"key"`
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
}

type ExportUserContactRsp struct {
	ContactUrl string `json:"contact_url"`
}

type GetHistoryStatisInfo struct {
	RegisterNum             int64 `json:"register_num"`                //新增注册
	TodayNewActiveUserNum   int64 `json:"today_new_active_user_num"`   //新增活跃用户
	TodayActiveUserNum      int64 `json:"today_active_user_num"`       //活跃用户
	TodayCreateGroupTaskNum int64 `json:"today_create_group_task_num"` //拉群任务数
	DataNum                 int64 `json:"data_num"`                    //推广资源数
	BountyAmount            int64 `json:"bounty_amount"`               //任务收益
	CommissionAmount        int64 `json:"commission_amount"`           //返佣收益
	WithdrawUserNum         int64 `json:"withdraw_user_num"`           //提现人数
	WithdrawAmount          int64 `json:"withdraw_amount"`             //提现扣款
	AdjustAmount            int64 `json:"adjust_amount"`               //人工调整
	PullFanTaskNum          int64 `json:"pull_fan_task_num"`           //拉粉任务数
	PullFanDataNum          int64 `json:"pull_fan_data_num"`           //拉粉资源数
	PullFanBountyAmount     int64 `json:"pull_fan_bounty_amount"`      //拉粉赏金
	PullFanCommissionAmount int64 `json:"pull_fan_commission_amount"`  //拉粉返佣
}

type RefundReq struct {
	Id      string `json:"id"`
	Account string `json:"account"`
}

type ShuntReq struct {
	LiveCode string `json:"live_code"` //1-访问量 2-点击量
}

type ShuntRsp struct {
	VisitNum int64 `json:"visit_num"` //访问量
	ClickNum int64 `json:"click_num"` //点击量
}
type ShuntData struct {
	Qudao string `json:"qudao"`
}
