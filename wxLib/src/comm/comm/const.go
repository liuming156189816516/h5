package comm

// serverName
const (
	ServerWechatDll = "ws_task_wechatdll"
	ServerData      = "ws_task_data"
	ServerTask      = "ws_task_task"
	ServerWebFace   = "ws_task_webface"
	ServerCheck     = "ws_task_check"
)

// cfg 配置
const (
	RunMode = "runmode" //运行模式
)

// inc 值类型
const (
	MenuId      = "menuId"
	RoleId      = "roleId"
	ChatMessage = "chatMessage"
	AiUuid      = "aiUuid"
	AiUuidWs    = "aiUuidWs"
	AiUuidTg    = "aiUuidTg"
)

const (
	//一个粉的积分
	ONE_FRIEND_AMOUNT = "ONE_FRIEND_AMOUNT"
	//返佣比列
	ONE_FRIEND_RATE = "ONE_FRIEND_RATE"
	//提现限制积分
	WITHDRAW_LIMIT_AMOUNT = "WITHDRAW_LIMIT_AMOUNT"
	//提现次数
	WITHDRAW_LIMIT_COUNT = "WITHDRAW_LIMIT_COUNT"
	//审批限制积分
	APPROVAL_LIMIT_COUNT = "APPROVAL_LIMIT_COUNT"
)

// 投放事件
const (
	CompleteRegistration = "CompleteRegistration" //注册成功
	SubmitApplication    = "SubmitApplication"    // 邀友入群提交
	Subscribe            = "Subscribe"            // 邀友入群通过
)

const (
	//当任务到期时系统自动关闭
	Timeout = "O sistema desliga automaticamente quando uma tarefa expira"
	//任务已完成,奖励已发放
	Success = "A tarefa foi concluída e a recompensa foi distribuída"
	//一个wsApp账号拉群限制
	WsAppLimit = "A conta do WhatsApp [@] concluiu 3 tarefas manuais. Troque de conta e tente novamente"
	//无效粉丝
	InvalidFans = "Após verificação, descobri que o número do celular da tarefa não foi encontrado no grupo"
	//禁言限制
	MuteLimit = "Por favor, não configure para que apenas administradores possam enviar mensagens"
	//无效链接
	InvalidLink = "Link inválido"
	//链接被撤销
	LinkRevoked = "O link de convite foi revogado"
	//链接失效
	LinkInvalidation = "O link de convite expirou"
	//无法加入群组，群组已满
	GroupIsFull = "Não é possível entrar no grupo, o grupo está cheio"
	//进群需要管理员同意
	AdminLimit = "Não defina privilégios de administrador para entrar neste grupo"
	//无效粉丝(拉粉)
	InvalidLink2 = "Após verificação, nenhum número de celular brasileiro foi encontrado no grupo. Ou o link foi enviado várias vezes"
	//驳回
	TurnDown = "Se você recusar, entre em contato com o atendimento ao cliente"
	//提现超时
	WithdrawalTimeout = "Tempo esgotado, entre em contato com o atendimento ao cliente"
	//区号限制
	AreaCodeLimit = "55"
	//汇率
	ExchangeRate = 0.005
	//一次拉粉
	FansNum = "FANSNUM"
	//注册赠送
	RegisterGiveaway = "REGISTERGIVEAWAY"
	//下载赠送
	DownloadGiveaway = "DOWNLOADGIVEAWAY"
	//拉群任务一个粉多少金币
	LqGiveaway = "LQGIVEAWAY"
	//拉粉任务一个粉多少金币
	LfGiveaway = "LFGIVEAWAY"
	//拉群任务限制每天多少次
	LqNumLimit = "LQNUMLIMIT"
	//拉粉任务单群验资最大数量
	LfMax = "LFMAX"
	//拉粉任务单群验资最小数量
	LfMin = "LFMIN"
	//提现限制次数
	WithdrawalLimitCount = "WITHDRAWALLIMITCOUNT"
	//提现限制积分
	WithdrawalLimitAmount = "WITHDRAWALLIMITAMOUNT"
	//审批限制积分
	ApprovalLimit = "APPROVALLIMIT"
	//返佣比例 %
	RebateRate = "REBATERATE"
	//挂机一个小时的奖励
	GJGiveaway = "GJGIVEAWAY"
	//挂机一个粉多少金币
	GJGiveawayIntegral = "GJGIVEAWAYINTEGRAL"
	//ai私发积分
	AIGiveaway = "AIGIVEAWAY"
	//ai任务发放数量
	AINum = "AINUM"
	//ai任务开关 0-关 1-开
	AITaskSwitch = "AITASKSWITCH"
	//ws私发任务开关 0-关 1-开
	WSPrivateTaskSwitch = "WSPRIVATETASKSWITCH"

)
