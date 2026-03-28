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
	"time"
)

// 下线
func TaskTypeUserLogoutEventBackHandler(msg *natsRpc.NatsMsg) int32 {
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
	req := &event.TaskUserLogoutEventReq{}
	reqStr, _ := jsoniter.MarshalToString(rsp.Req)
	err = jsoniter.UnmarshalFromString(reqStr, &req)

	//删除任务
	cache.DelDoTask(req.TaskId)

	if rsp.Code == 0 {
		for _, accStr := range req.Account {
			//下线成功
			account.UpAccountInfo(bson.M{"account": accStr}, bson.M{"reason": "账号退出", "status": int64(1), "offline_time": time.Now().Unix()}, req.Uid)
			cache.SetAccountStatus(req.Uid, accStr, 1)
		}
	}
	return natsRpc.ESMR_SUCCEED
}
