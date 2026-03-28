package httpQury

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"proxy/ent/info"
)

// 请求
func NrpcDllCallDll(method, mod string, body interface{}, timeouts ...int32) *info.ResponseResult {
	toString, _ := jsoniter.MarshalToString(body)
	msg, err := natsRpc.NrpcDllCall(mod, 100, "", method+"|"+toString, timeouts...)
	if err != nil {
		return &info.ResponseResult{
			Code:    -1001,
			Message: "请求 错误 ",
		}
	}
	newRet := &info.ResponseResultNew{}
	err = jsoniter.Unmarshal(msg.Data, newRet)
	//logs.Info("读取到的数据 %+v", string(msg.Data))
	if err != nil {
		logs.Error("NatsCall %s.%s Unmarshal rsp %s failed:%s", mod, "", string(msg.Data), err.Error())
		return nil
	}
	ret := &info.ResponseResult{}
	ret.Code = newRet.Ret
	ret.Message = newRet.Msg
	ret.Data = newRet.Data
	ret.Success = true
	if newRet.Ret != 0 {
		ret.Success = false
	}
	return ret
}
