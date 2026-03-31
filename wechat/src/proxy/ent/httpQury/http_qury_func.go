package httpQury

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"proxy/ent/info"
)

// 帐号登录
func AccountLogin(param []*info.AccountLoginParam, account string) *info.ResponseResult {
	info.SaveLogs(account, "AccountLogin-入参", param)
	ret := NrpcDllCallDll("/account/login?async=1", beego.AppConfig.String("dll_mod"), param, 30)
	info.SaveLogs(account, "AccountLogin-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 帐号离线
func AccountLogout(param []string, account string) *info.ResponseResult {
	info.SaveLogs(account, "AccountLogout-入参", param)
	ret := NrpcDllCallDll("/account/remove", beego.AppConfig.String("dll_mod"), param, 30)
	info.SaveLogs(account, "AccountLogout-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 帐号移除
func AccountRemove(param []string, account string) *info.ResponseResult {
	info.SaveLogs(account, "AccountRemove-入参", param)
	ret := NrpcDllCallDll("/account/remove", beego.AppConfig.String("dll_mod"), param, 30)
	info.SaveLogs(account, "AccountRemove-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 设置商店名称
func ShopName(param *info.ShopNameParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "ShopSetname-入参", param)
	ret := NrpcDllCallDll("/shop/setname", beego.AppConfig.String("dll_mod"), param, 30)
	info.SaveLogs(param.Account, "ShopSetname-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误",
		}
	}
	return ret
}

// 帐号的一些信息检测
func AccountCheck(param *info.AccountCheckParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "AccountCheck-入参", param)
	ret := NrpcDllCallDll("/account/check", beego.AppConfig.String("dll_mod"), param, 30)
	info.SaveLogs(param.Account, "AccountCheck-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 消息发送
func MessageSend(param *info.MessageSendParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "MessageSend-入参", param)
	ret := NrpcDllCallDll("/message/send?async=1", beego.AppConfig.String("dll_mod"), param, 10)
	info.SaveLogs(param.Account, "MessageSend-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 消息发送(异步)
func MessageSendAsyn(param *info.MessageSendParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "MessageSend3-入参", param)
	ret := NrpcDllCallDll("/message/send?async=3", beego.AppConfig.String("dll_mod"), param, 30)
	info.SaveLogs(param.Account, "MessageSend3-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误",
		}
	}
	return ret
}

// 拨打电话
func CallPhone(param *info.CallPhoneParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "CallPhone-入参", param)
	ret := NrpcDllCallDll("/call/phone", beego.AppConfig.String("dll_mod"), param, 10)
	info.SaveLogs(param.Account, "CallPhone-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 挂断通话
func CallHangup(param *info.CallHangupParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "CallHangup-入参", param)
	ret := NrpcDllCallDll("/call/hangup", beego.AppConfig.String("dll_mod"), param, 10)
	info.SaveLogs(param.Account, "CallHangup-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 查询注册
func PhoneQuery(param *info.PhoneQueryParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "PhoneQuery-入参", param)
	ret := NrpcDllCallDll("/phone/query", beego.AppConfig.String("dll_mod"), param, 10)
	info.SaveLogs(param.Account, "PhoneQuery-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 自主进群
func GroupJoin(param *info.GroupJoinParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "GroupJoin-入参", param)
	ret := NrpcDllCallDll("/group/join", beego.AppConfig.String("dll_mod"), param, 60)
	info.SaveLogs(param.Account, "GroupJoin-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 获取群的详细信息
func GroupInfo(param *info.GroupInfoParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "GroupGetinfo-入参", param)
	ret := NrpcDllCallDll("/group/getinfo", beego.AppConfig.String("dll_mod"), param, 60)
	info.SaveLogs(param.Account, "GroupGetinfo-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 创建关联的验证码
func VfcodeCreate(param *info.VfcodeCreateParam) *info.ResponseResult {
	logs.Info("VfcodeCreate param:", param)
	//info.SaveLogs(param.Account, "VfcodeCreate-入参", param)
	ret := NrpcDllCallDll("/vfcode/create", beego.AppConfig.String("dll_mod"), param, 10)
	toString, _ := jsoniter.MarshalToString(ret)
	logs.Info("VfcodeCreate ret:", toString)
	//info.SaveLogs(param.Account, "VfcodeCreate-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}

// 检测关联的验证码
func VfcodeCheck(param *info.VfcodeCheckParam) *info.ResponseResult {
	info.SaveLogs(param.Account, "VfcodeCheck-入参", param)
	ret := NrpcDllCallDll("/vfcode/check", beego.AppConfig.String("dll_mod"), param, 10)
	info.SaveLogs(param.Account, "VfcodeCheck-出参", ret)
	if ret == nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	return ret
}
