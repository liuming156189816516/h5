package info

import (
	"encoding/json"
	"mlog"
)

// 事件 不需要回复
func SaveLogs(wxid, name string, data interface{}) {
	str, ok := data.(string)
	if !ok {
		b, _ := json.Marshal(data)
		str = string(b)
	}
	if mlog.IsNeedTraceUid(wxid) {
		mlog.GetPlayerLoggerMng().WriteMsg(wxid,
			"%s:----------------------------------------------------\n%s\n----------------------------------------------------", name, string(str))
	}
}
