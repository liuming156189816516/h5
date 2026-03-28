package task

import (
	"comm/event"
	"data/models/stru"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"natsRpc"
	"selfComm/db/account"
	"selfComm/wxComm/cache"
)

//登录
func TaskUserLoginEventBackHandler(msg *natsRpc.NatsMsg) int32 {
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
	req := &event.TaskUserLoginEventReq{}
	reqStr, _ := jsoniter.MarshalToString(rsp.Req)
	err = jsoniter.UnmarshalFromString(reqStr, &req)

	//删除任务
	cache.DelDoTask(req.TaskId)

	if rsp.Code != 0 {
		//失败
		for _, accStr := range req.Account {
			account.UpAccountInfo(bson.M{"account": accStr}, bson.M{"status": int64(4), "reason": rsp.Message}, req.Uid)
			cache.SetAccountStatus(req.Uid, accStr, 4)
		}
	}
	return natsRpc.ESMR_SUCCEED
}
