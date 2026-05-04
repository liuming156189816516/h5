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
		where["account"] = bson.M{
			"$regex": "^" + req.Account,
		}
	}
	if req.AccountStatus > 0 {
		where["account_status"] = req.AccountStatus
	}
	if req.AccountGroup != "" {
		where["account_group"] = req.AccountGroup
	} else {
		if req.StartTime > 0 && req.EndTime > 0 {
			where["itime"] = bson.M{"$gte": req.StartTime, "$lte": req.EndTime}
		}
	}
	sort := "-itime"
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
	SuccessCount := int64(0)
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
			SuccessCount = SuccessCount + tmp.SucessNum
		}

		if p, ok := data["arrived_num"]; ok {
			tmp.ArrivedNum = utils.GetInt64(p)
		}
		if p, ok := data["reason"]; ok {
			tmp.Reason = utils.GetString(p)
		}
		if p, ok := data["itime"]; ok {
			tmp.Itime = utils.GetInt64(p)
		}
		if p, ok := data["ptime"]; ok {
			tmp.Ptime = utils.GetInt64(p)
		}
		rsp.List = append(rsp.List, tmp)
	}

	sumwhere := []bson.M{}
	if req.AccountGroup != "" {
		sumwhere = append(sumwhere, bson.M{"$match": bson.M{"account_group": req.AccountGroup}})
	} else {
		if req.StartTime > 0 && req.EndTime > 0 {
			sumwhere = append(sumwhere, bson.M{"$match": bson.M{"itime": bson.M{"$gte": req.StartTime, "$lte": req.EndTime}}})
		}
	}

	sumwhere = append(sumwhere, bson.M{
		"$group": bson.M{
			"_id": "null",
			// 发送成功总数
			"total_sucess": bson.M{
				"$sum": "$sucess_num",
			},
			//送达成功总数
			"total_arrived": bson.M{
				"$sum": "$arrived_num",
			},
			// 账号数量
			"account_count": bson.M{
				"$sum": 1,
			},
		},
	})

	sumRet, err1 := mgoDeal.QueryMongoSum(db, tb, sumwhere)

	var totalSucess, totalArrived, accountCount int64

	if err1 == nil {
		if p, ok := sumRet["total_sucess"]; ok {
			totalSucess = utils.GetInt64(p)
		}

		if p, ok := sumRet["total_arrived"]; ok {
			totalArrived = utils.GetInt64(p)
		}

		if p, ok := sumRet["account_count"]; ok {
			accountCount = utils.GetInt64(p)
		}
	}

	if totalSucess > 0 && accountCount > 0 {
		rsp.SuccessCount = totalSucess
		rsp.ArrivedCount = totalArrived
		rsp.Average = totalSucess / accountCount
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
			count := cache.GetSendMsgTaskInfoCount(cache.SuccessNum, msgInfo.Account)
			if msgInfo.SucessNum != count {
				up["sucess_num"] = count
			}

			count10 := cache.GetSendMsgTaskInfoCount(cache.ArrivedNum, msgInfo.Account)
			up["arrived_num"] = count10

			if len(up) > 0 {
				sendmsg.UpSendMsgInfo(bson.M{"_id": msgInfo.Id}, up)
			}
		}
	}()
	return nil
}

// 获取自动发送消息开关 "0" - 开; "1" - 关
func (this *SendMsgServer) GetAutoSendMsgStatus(req *info.NullReq, rsp *info.GetAutoSendMsgStatusRsp) *goError.ErrRsp {
	rsp.AutoSendMsgStatus = cache.GetTaskStatus()
	return nil
}

// 自动发送消息开关-修改 "0" - 开; "1" - 关
func (this *SendMsgServer) DoAutoSendMsgStatus(req *info.DoAutoSendMsgStatusReq, rsp *info.NullRsp) *goError.ErrRsp {
	cache.SetTaskStatus(req.AutoSendMsgStatus)
	return nil
}
