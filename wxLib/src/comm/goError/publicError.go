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
	GLOBAL_INVALIDPARAM  = &ErrRsp{100001, "param err"}
	UserPwdErr           = &ErrRsp{100002, "用户不存在或密码错误"}
	UserDisableErr       = &ErrRsp{100003, "用户被禁用"}
	MaterailGroupNameErr = &ErrRsp{100004, "分组名称不能为未分组"}
	DelIpGroupErr        = &ErrRsp{100005, "检测到当前分组下还存在IP，请先移动IP到其余的分组或者删除IP后再进行操作"}
	DelAccountGroupErr   = &ErrRsp{100006, "检测到当前分组下还存在ws账号，请先移动ws账号到其余的分组或者删除ws账号后再进行操作"}
	DoLoginErr           = &ErrRsp{100007, "没有可登录的账号"}
	DefaultGroupErr      = &ErrRsp{100008, "不能以未分组命名"}
	GLOBAL_TOOMACH       = &ErrRsp{100009, "操作频繁"}
	UserExitErr          = &ErrRsp{100010, "用户已存在"}
	IpOperationErr       = &ErrRsp{100011, "获取ip异常"}
	AccountCodeLoginErr       = &ErrRsp{100012, "账号验证码登录异常"}

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
