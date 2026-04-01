package task

import (
	"comm/event"
	"data/models/stru"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"natsRpc"
	"selfComm/db/sendmsg"
	"selfComm/wxComm/cache"
)

// 消息发送
func TaskTypeSendMsgEventBackHandler(msg *natsRpc.NatsMsg) int32 {
	//解析
	rsp := &stru.ResponseResult{}
	err := json.Unmarshal(msg.MsgData, &rsp)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	msg.ResponeSucc("")
	if rsp.TaskId == "" {
		logs.Error("数据错误", rsp)
		return natsRpc.ESMR_SUCCEED
	}

	//解析req
	req := &event.TaskTypeSendMsgEventReq{}
	reqStr, _ := jsoniter.MarshalToString(rsp.Req)
	err = jsoniter.UnmarshalFromString(reqStr, &req)
	//删除任务
	cache.DelDoTask(req.TaskId)

	if rsp.Code == -2 {
		//粉丝数据不足
		sendmsg.UpSendMsgInfo(bson.M{"account": req.Account}, bson.M{"reason": "粉丝数据不足"})
		return natsRpc.ESMR_SUCCEED
	}

	if rsp.Code != 0 && req.Account != "" {
		//失败
		record := cache.GetSendMsgRecordInfo(req.Account, req.Account+"_"+req.Target)
		//把数据还回去
		if record.DataPackId != "" {
			target := req.Target
			if req.Lid != "" {
				target = req.Target + "-" + req.Lid
			}
			cache.SaddDataPackList(req.DataPackId, target)
		}
	}
	return natsRpc.ESMR_SUCCEED
}
