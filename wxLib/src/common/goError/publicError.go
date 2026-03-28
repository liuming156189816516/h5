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

func NewError(format string, v ...interface{}) error {
	msg := fmt.Sprintf(format, v...)
	return errors.New(msg)
}

type GError struct {
	n int32
	s string
}

func (e *GError) Error() string {
	return e.s
}
func (e *GError) ErrNo() int32 {
	return e.n
}

func NewGoError(n int32, format string, v ...interface{}) *GError {
	te := &GError{n: n}
	xv := make([]interface{}, 0, len(v)+1)
	xv = append(xv, n)
	xv = append(xv, v...)
	te.s = fmt.Sprintf("[%d:"+format+"]", xv)
	return te
}

type ErrRsp struct {
	Ret int32  `json:"ret"`
	Msg string `json:"msg"`
}

func GetErrorMsg(rsp ErrRsp) (int32, string) {
	return rsp.Ret, rsp.Msg
}

//error for conn
var (
	WebSocketHeadError = &ErrRsp{400, "网络连接失败"}
	CheckLoginFailed   = &ErrRsp{401, "请重新登录"}
	ReqError           = &ErrRsp{404, "请求数据格式错误"}
	NoService          = &ErrRsp{404, "请求服务无效"}
	NoAuth             = &ErrRsp{405, "请求无权限"}
	RenewLogin         = &ErrRsp{406, "请重新登录"}
	NotAllow           = &ErrRsp{407, "NotAllow"}
	QueryKuick         = &ErrRsp{410, "操作太频繁"}
	RspError           = &ErrRsp{420, "处理失败"}
	RspAppClose        = &ErrRsp{421, "系统升级中,请耐心等待..."}
	RspTimeout         = &ErrRsp{501, "请求超时"}
	IpError            = &ErrRsp{502, "登录IP限制"}
	RegionError        = &ErrRsp{503, "登录地区限制"}
)

var GLOBAL_INVALIDPARAM = &ErrRsp{1000002, "参数错误"} // global 100
var GLOBAL_SYSTEMERROR = &ErrRsp{100000, "系统错误"}   // global 100 = &ErrRsp{1000002, "参数错误"} // global 100
