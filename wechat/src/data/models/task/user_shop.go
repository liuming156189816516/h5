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

//设置商店名称
func TaskTypeShopNameEventBackHandler(msg *natsRpc.NatsMsg) int32 {
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
		logs.Error("TaskTypeShopNameEventBackHandler数据错误", rsp)
		return natsRpc.ESMR_SUCCEED
	}

	//解析req
	req := &event.TaskTypeShopNameEventReq{}
	reqStr, _ := jsoniter.MarshalToString(rsp.Req)
	err = jsoniter.UnmarshalFromString(reqStr, &req)

	//删除任务
	cache.DelDoTask(req.TaskId)

	if rsp.Code == 0 {
		//成功
		account.UpAccountInfo(bson.M{"account": req.Account}, bson.M{"nick_name": req.NickName}, req.Uid)
		accInfo := cache.GetAccountInfo(req.Uid, req.Account)
		accInfo.NickName = req.NickName
		cache.SetAccountInfo(req.Uid, req.Account, accInfo)
	}
	return natsRpc.ESMR_SUCCEED
}
