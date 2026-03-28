package task

import (
	"comm/event"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"selfComm/wxComm/cache"
	"serApi/dataApi"
	"strings"
)

func AccountMessageResultEventHandler(msg *natsRpc.NatsMsg) int32 {
	//错误
	req := &event.ReceiveMessagesReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Info("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	defer msg.ResponeSucc(&dataApi.AccountSendMsgResultRsp{})

	//删除任务
	cache.DelDoTask(req.TaskId)
	//解析data
	tmpData := &dataApi.AccountSendMsgResultData{}
	jsoniter.UnmarshalFromString(req.Content, &tmpData)
	target := strings.ReplaceAll(tmpData.To, "@s.whatsapp.net", "")
	if strings.Contains(target, "lid") {
		split := strings.Split(target, ".")
		target = split[0]
	}

	record := cache.GetSendMsgRecordInfo(req.Account + "_" + target)
	if tmpData.Errno == 0 && tmpData.Data != "" {
		//成功
		cache.IncSendMsgTaskInfoCount(cache.SuccessNum, record.Account, record.Account, 1)
	} else {
		//把数据还回去
		if record.DataPackId != "" {
			cache.SaddDataPackList(record.DataPackId, record.Target+"-"+target)
		}
	}
	cache.SetSendMsgRecord(record)
	return natsRpc.ESMR_SUCCEED
}
