package goError

import (
	"errors"
	"fmt"
)

var (
	Err_NoMysql     = errors.New("No Mysql DB")
	Err_NoRedis     = errors.New("No Redis")
	Err_NoApp       = errors.New("No App")
	Err_NoRedisTmpl = errors.New("NotRedisTmpl")
)

const (
	Http_check_ip_limit   = 601 //单个IP在30s内访问超过限制
	Http_check_uri_repeat = 602 //同一个uri重复范文
	Http_check_uri_sum    = 603 //uri checksum 失败
)

func NewError(format string, v ...interface{}) error {
	msg := fmt.Sprintf(format, v...)
	return errors.New(msg)
}

func NewGoError(n int32, format string) *ErrRsp {
	te := &ErrRsp{Ret: n}
	te.Msg = format
	return te
}

type ErrRsp struct {
	Ret int32  `json:"code"`
	Msg string `json:"msg"`
}

func GetErrorMsg(rsp ErrRsp) (int32, string) {
	return rsp.Ret, rsp.Msg
}

// error for conn
var (
	SuccRsp            = &ErrRsp{Ret: 0, Msg: "success"}
	WebSocketHeadError = &ErrRsp{400, "no net"}
	CheckLoginFailed   = &ErrRsp{401, "Re-login"}
	OpenidLoginFailed  = &ErrRsp{402, "Re-login"} //openid过期
	QueryError         = &ErrRsp{403, "No access allowed"}
)

var (
	GLOBAL_SYSTEMERROR   = &ErrRsp{100000, "sys err"}
	GLOBAL_INVALIDPARAM  = &ErrRsp{100002, "param err"}
	UserPwdErr           = &ErrRsp{100003, "用户不存在或密码错误"}
	UserDisableErr       = &ErrRsp{100004, "用户被禁用"}
	MaterailGroupNameErr = &ErrRsp{100006, "分组名称不能为未分组"}
	DelIpGroupErr        = &ErrRsp{100007, "检测到当前分组下还存在IP，请先移动IP到其余的分组或者删除IP后再进行操作"}
	DelRoleErr           = &ErrRsp{100008, "admin或主管角色,不可删除"}
	UpRoleErr            = &ErrRsp{100009, "分组名称不能为admin或者主管"}
	DelAccountGroupErr   = &ErrRsp{100010, "检测到当前分组下还存在ws账号，请先移动ws账号到其余的分组或者删除ws账号后再进行操作"}
	DoUseStatusErr       = &ErrRsp{100011, "后台端口不足，请联系客服增加端口！"}
	DoLoginErr           = &ErrRsp{100012, "没有可登录的账号"}
	NoAccountErr         = &ErrRsp{100013, "没有可用的账号"}
	TimeErr              = &ErrRsp{100014, "开始时间不能大于结束时间"}
	RecordLogErr         = &ErrRsp{100015, "暂无日志"}
	DoChatGroupErr       = &ErrRsp{100016, "修改失败,分组名称不能为默认值"}
	TranslateTextErr     = &ErrRsp{100017, "翻译失败"}
	ShareErr             = &ErrRsp{100018, "该工单分享链接已关闭或已到期，请联系对应的工单创建者。"}
	SharePwdErr          = &ErrRsp{100019, "密码错误"}

	DoCreateWorkTaskErr1 = &ErrRsp{100020, "工单历史中存在账号计数工单!!!"}
	DoCreateWorkTaskErr2 = &ErrRsp{100021, "工单历史中存在席位计数工单!!!"}
	DoPullGroupErr       = &ErrRsp{100022, "拉群分组和营销分组不可以是同一个分组"}
	DefaultGroupErr      = &ErrRsp{100025, "不能以未分组命名"}
	SysConfigKeyExitErr  = &ErrRsp{100027, "配置项key已存在,请重新设置"}
	BalanceErr           = &ErrRsp{100029, "账户余额不足"}
	DelAccountGroupErr2  = &ErrRsp{100031, "检测到当前分组存在于任务配置中,请先修改任务配置中的营销分组后再进行操作"}
	WithdrawReviseErr    = &ErrRsp{100039, "用户不存在"}

	GLOBAL_TOOMACH          = &ErrRsp{100000, "Operações frequentes"}
	RevisePwdErr            = &ErrRsp{100019, "Senha errada"}
	UserPwdEnErr            = &ErrRsp{100003, "O usuário não existe ou a senha está incorreta"}
	UserExitErr             = &ErrRsp{100005, "O usuário já existe"}
	UserRegisterCodeErr     = &ErrRsp{100023, "Erro de código de verificação"}
	InviteCodeErr           = &ErrRsp{100024, "O código de convite expirou"}
	UserNoExitErr           = &ErrRsp{100026, "O usuário não existe"}
	AmountLimitEnErr        = &ErrRsp{100030, "O valor do seu saque não pode ser menor que o valor mínimo de saque"}
	SystemEnError           = &ErrRsp{100032, "Erro no sistema, entre em contato com o atendimento ao cliente"}
	SubmitCreateTaskEnError = &ErrRsp{100034, "Erro de formato de link"}
	VerificationErr         = &ErrRsp{100035, "Erro de autenticação"}
	WithdrawCountEnErr      = &ErrRsp{100036, "Muitas retiradas hoje"}
	OperationENErr          = &ErrRsp{100037, "Operação frequente"}
	UserRuletaErr           = &ErrRsp{100040, "Por favor, conclua a tarefa para ter a chance de girar a roleta"}
	DoSignErr               = &ErrRsp{100041, "Login concluído hoje"}
	TaskSignErr             = &ErrRsp{100043, "Você deve concluir a tarefa antes de poder fazer login"}
	RegisterErr             = &ErrRsp{100044, "Restrições de registro"}
	BlackErr                = &ErrRsp{100045, "A conta foi desativada"}
	UserRuletaBalanceErr    = &ErrRsp{100047, "Saldo de conta insuficiente"}
	GetQrcodeErr            = &ErrRsp{100048, "A conta do WhatsApp já existe"}
	BindingErr              = &ErrRsp{100049, "A vinculação falhou, a conta de retirada foi vinculada"}
	SubmitAiMessageTaskErr  = &ErrRsp{100050, "Por favor, carregue as fotos primeiro"}
	SubmitTgCodeErr         = &ErrRsp{100051, "Erro no código de verificação"}
	GetQrcodeErr1           = &ErrRsp{100052, "Por favor, tente novamente em 60 segundos"}
	SubmitTgCodeErr2        = &ErrRsp{100053, "Operações frequentes, tente novamente"}
	AiDataErr        		= &ErrRsp{100054, "Dados do ventilador insuficientes, entre em contato com o atendimento ao cliente"}
	RegisterAccountPhoneErr = &ErrRsp{100055, "Usuários não brasileiros não podem se cadastrar"}
	QrcodePhoneErr          = &ErrRsp{100056, "Por favor, use um número de telefone celular começando com 55"}
	WithdrawBankErr         = &ErrRsp{100057, "Por favor, utilize o CPF para sacar dinheiro"}
)

var globalErrMap = map[int32]string{
	400: "no net",
}
var (
	FILE_NOT_FOUND = &ErrRsp{50001, "文件未找到"}
)

var (
	DATABASE_READ_ERROR   = &ErrRsp{40001, "DATABASE_READ_ERROR"}
	DATABASE_INSERT_ERROR = &ErrRsp{40002, "DATABASE_INSERT_ERROR"}
	DATABASE_UPDATE_ERROR = &ErrRsp{40003, "DATABASE_UPDATE_ERROR"}
)
