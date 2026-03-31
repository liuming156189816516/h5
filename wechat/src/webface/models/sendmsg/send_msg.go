package sendmsg

import (
	"comm/comm"
	"comm/goError"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/sendmsg"
	"selfComm/wxComm/cache"
	"utils"
	info "webface/webstru"
)

// 群发
type SendMsgServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *SendMsgServer) getUid() string {
	return this.Sess.Uid
}

// 自动群发任务-列表
func (this *SendMsgServer) GetSendMsgInfoList(req *info.GetSendMsgInfoListReq, rsp *info.GetSendMsgInfoListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	where := bson.M{}
	if req.Account != "" {
		where["account"] = req.Account
	}
	if req.AccountStatus > 0 {
		where["account_status"] = req.AccountStatus
	}
	sort := "itime"
	rsp.List = []*info.GetSendMsgInfoListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, sort, start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetSendMsgInfoListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
		}
		if p, ok := data["account"]; ok {
			tmp.Account = utils.GetString(p)
		}
		if p, ok := data["account_status"]; ok {
			tmp.AccountStatus = utils.GetInt64(p)
		}
		if p, ok := data["sucess_num"]; ok {
			tmp.SucessNum = utils.GetInt64(p)
		}

		if p, ok := data["arrived_num"]; ok {
			tmp.ArrivedNum = utils.GetInt64(p)
		}
		if p, ok := data["reason"]; ok {
			tmp.Reason = utils.GetString(p)
		}
		rsp.List = append(rsp.List, tmp)
	}

	go func() {
		msgInfoList := sendmsg.GetListSendMsgInfo(bson.M{}, -1)
		for _, msgInfo := range msgInfoList {
			accountStatus := cache.GetAccountStatus(msgInfo.Account)
			if accountStatus != 2 {
				accountStatus = 1
			}
			up := bson.M{}
			if msgInfo.AccountStatus != accountStatus {
				up["account_status"] = accountStatus
			}
			count := cache.GetSendMsgTaskInfoCount(cache.SuccessNum, msgInfo.Account, msgInfo.Account)
			if msgInfo.SucessNum != count {
				up["sucess_num"] = count
			}

			count10 := cache.GetSendMsgTaskInfoCount(cache.ArrivedNum, msgInfo.Account, msgInfo.Account)
			up["arrived_num"] = count10

			if len(up) > 0 {
				sendmsg.UpSendMsgInfo(bson.M{"_id": msgInfo.Id}, up)
			}
		}
	}()
	return nil
}
