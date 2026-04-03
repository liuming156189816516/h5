package account

import (
	"comm/comm"
	"comm/cos"
	"comm/event"
	"comm/goError"
	"comm/mgoDeal"
	"comm/tableName"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math/rand"
	"os"
	"selfComm/db/account"
	"selfComm/db/ip"
	"selfComm/wxComm/cache"
	"strings"
	"time"
	"utils"
	info "webface/webstru"
)

// 数据包
type AccountServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *AccountServer) getUid() string {
	return this.Sess.Uid
}

// 账号分组-列表
func (this *AccountServer) GetAccountGroupList(req *info.GetAccountGroupListReq, rsp *info.GetAccountGroupListRsp) *goError.ErrRsp {

	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	where := bson.M{}
	if req.Name != "" {
		where["name"] = req.Name
	}
	rsp.List = []*info.GetAccountGroupListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "sort", start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetAccountGroupListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
			count := account.GetCountAccountInfo(bson.M{"group_id": tmp.Id})
			tmp.Count = count
		}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
			if tmp.Name == "未分组" {
				tmp.IsDefault = 1
			}
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// 账号分组-操作
func (this *AccountServer) DoAccountGroup(req *info.DoAccountGroupReq, rsp *info.NullRsp) *goError.ErrRsp {

	if req.Ptype == 1 {
		//新增
		if req.Name == "未分组" {
			return goError.DefaultGroupErr
		}
		tmp := &account.AccountGroup{}
		tmp.Id = bson.NewObjectId()
		tmp.Name = req.Name
		tmp.Sort = 0 - time.Now().Unix()
		account.AddAccountGroup(tmp)
	}
	if req.Ptype == 2 {
		//编辑
		where := bson.M{}
		where["_id"] = bson.ObjectIdHex(req.Id)
		update := bson.M{}
		update["name"] = req.Name
		account.UpAccountGroup(where, update)
	}
	if req.Ptype == 3 {
		//删除
		for _, delId := range req.DelId {
			count := account.GetCountAccountInfo(bson.M{"group_id": delId})
			if count > 0 {
				return goError.DelAccountGroupErr
			}
			where := bson.M{}
			where["_id"] = bson.ObjectIdHex(delId)
			account.DelAccountGroup(where)
		}
	}
	return nil
}

// 账号-列表
func (this *AccountServer) GetAccountInfoList(req *info.GetAccountInfoListReq, rsp *info.GetAccountInfoListRsp) *goError.ErrRsp {

	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	where := bson.M{}
	if req.GroupId != "" {
		where["group_id"] = req.GroupId
	}

	if req.Account != "" {
		where["account"] = req.Account
	}

	if req.Status > 0 {
		where["status"] = req.Status
	}
	if req.AccountType > 0 {
		where["account_type"] = req.AccountType
	}
	if req.PlatformType > 0 {
		where["platform_type"] = req.PlatformType
	}

	if req.IsProxyUser >= 0 {
		where["is_proxy_user"] = req.IsProxyUser
	}
	sort := "-itime"
	if req.Reason != "" {
		where["reason"] = req.Reason
	}

	if req.PixelId != "" {
		where["pixel_id"] = req.PixelId
	}

	if req.ItimeStartTime > 0 && req.ItimeEndTime > 0 {
		where["itime"] = bson.M{"$gte": req.ItimeStartTime, "$lte": req.ItimeEndTime}
	}

	if req.FirstLoginStartTime > 0 && req.FirstLoginEndTime > 0 {
		where["first_login_time"] = bson.M{"$gte": req.FirstLoginStartTime, "$lte": req.FirstLoginEndTime}
	}

	if req.OfflineStartTime > 0 && req.OfflineEndTime > 0 {
		where["offline_time"] = bson.M{"$gte": req.OfflineStartTime, "$lte": req.OfflineEndTime}
	}

	rsp.List = []*info.GetAccountInfoListInfo{}
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
		tmp := &info.GetAccountInfoListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
		}
		if p, ok := data["head"]; ok {
			tmp.Head = utils.GetString(p)
		}
		if p, ok := data["account"]; ok {
			tmp.Account = utils.GetString(p)

		}
		if p, ok := data["nick_name"]; ok {
			tmp.NickName = utils.GetString(p)
		}
		if p, ok := data["status"]; ok {
			tmp.Status = utils.GetInt64(p)
		}
		if p, ok := data["reason"]; ok {
			tmp.Reason = utils.GetString(p)
		}
		if p, ok := data["account_type"]; ok {
			tmp.AccountType = utils.GetInt64(p)
		}
		if p, ok := data["offline_time"]; ok {
			tmp.OfflineTime = utils.GetInt64(p)
		}
		if p, ok := data["first_login_time"]; ok {
			tmp.FirstLoginTime = utils.GetInt64(p)
		}
		if p, ok := data["remark"]; ok {
			tmp.Remark = utils.GetString(p)
		}
		if p, ok := data["pixel_id"]; ok {
			tmp.PixelId = utils.GetString(p)
		}
		if p, ok := data["platform_type"]; ok {
			tmp.PlatformType = utils.GetInt64(p)
		}
		if p, ok := data["itime"]; ok {
			tmp.Itime = utils.GetInt64(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// 移动至其他分组
func (this *AccountServer) DoUpGroup(req *info.DoUpGroupReq, rsp *info.NullRsp) *goError.ErrRsp {
	account.UpAccountInfo(bson.M{"account": bson.M{"$in": req.Accounts}}, bson.M{"group_id": req.GroupId})
	return nil
}

// 批量导出
func (this *AccountServer) DoOutPutAccount(req *info.DoOutPutAccountReq, rsp *info.DoOutPutAccountRsp) *goError.ErrRsp {

	exportStr := ""
	for _, acc := range req.Accounts {
		accInfo := account.GetOneAccountInfo(bson.M{"account": acc})
		ipStr := accInfo.Token
		exportStr = exportStr + ipStr + "\n"
	}
	fileName := comm.Md5(exportStr) + ".txt"
	tmpPath := beego.AppConfig.String("tmpPath")
	filePath := tmpPath + fileName
	err := ioutil.WriteFile(filePath, []byte(exportStr), 0777)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	defer os.Remove(filePath)
	fileUrl := cos.UploadAwsFile(filePath, fileName)
	rsp.Url = fileUrl
	return nil
}

// 批量删除
func (this *AccountServer) DoBatchDelAccount(req *info.DoBatchDelAccountReq, rsp *info.NullRsp) *goError.ErrRsp {

	w := bson.M{}
	//w["platform_type"] = int64(1)
	w["account"] = bson.M{"$in": req.Accounts}
	w["status"] = bson.M{"$ne": int64(2)}
	accInfoList := account.GetListAccountInfo(w, -1)
	err := account.DelAccountInfo(w)
	if err == nil {
		go func(accounts []*account.AccountInfo) {
			for _, accInfo := range accounts {
				go cache.DelAccountStatus(accInfo.Account)
				proxyIp := cache.GetProxyIp(accInfo.Account)
				go cache.IncIpUserNum(proxyIp.IpId, -1)
				go cache.DelProxyIp(accInfo.Account)
				go cache.DelAllAccountList(accInfo.Account)
				go cache.DelAccountInfo(accInfo.Account)
				//批量移除
				e := &event.TaskUserRemoveEventReq{
					Account: []string{accInfo.Account},
				}
				event.AddTask(event.TaskTypeUserRemove, "", e)
			}
		}(accInfoList)
	}
	return nil
}

// 快速上线
func (this *AccountServer) DoBatchFastLogin(req *info.DoBatchFastLoginReq, rsp *info.NullRsp) *goError.ErrRsp {

	accList := []string{}
	w := bson.M{}
	w["status"] = bson.M{"$ne": int64(3)}
	w["account"] = bson.M{"$in": req.Accounts}
	accinfoList := account.GetListAccountInfo(w, -1)
	for _, accountInfo := range accinfoList {
		accList = append(accList, accountInfo.Account)
	}
	if len(accList) == 0 {
		//没有可登录的账号
		return goError.DoLoginErr
	}
	account.UpAccountInfo(bson.M{"account": bson.M{"$in": accList}}, bson.M{"status": int64(3)})
	groupAccList := [][]string{}
	aList := []string{}
	for i, acc := range accList {
		aList = append(aList, acc)
		if len(aList) >= 50 {
			groupAccList = append(groupAccList, aList)
			aList = []string{}
		}
		if i == len(accList)-1 {
			groupAccList = append(groupAccList, aList)
		}
	}
	for _, accList2 := range groupAccList {
		if len(accList2) > 0 {
			e := &event.TaskUserLoginEventReq{
				Account: accList2,
			}
			event.AddLoginTask(event.TaskTypeUserLogin, "", e)
		}
	}
	return nil
}

// 批量下线
func (this *AccountServer) DoBatchLogout(req *info.DoBatchLogoutReq, rsp *info.NullRsp) *goError.ErrRsp {

	accList := req.Accounts
	account.UpAccountInfo(bson.M{"account": bson.M{"$in": accList}}, bson.M{"status": int64(5)})
	groupAccList := [][]string{}
	aList := []string{}
	for i, acc := range accList {
		aList = append(aList, acc)
		if len(aList) >= 50 {
			groupAccList = append(groupAccList, aList)
			aList = []string{}
		}
		if i == len(accList)-1 {
			groupAccList = append(groupAccList, aList)
		}
	}
	for _, accList2 := range groupAccList {
		if len(accList2) > 0 {
			e := &event.TaskUserLogoutEventReq{
				Account: accList2,
			}
			event.AddTask(event.TaskTypeUserLogout, "", e)
		}
	}
	return nil
}

// 分组排序
func (this *AccountServer) SortGroup(req *info.SortGroupReq, rsp *info.NullRsp) *goError.ErrRsp {

	for i, v := range req.List {
		account.UpAccountGroup(bson.M{"_id": bson.ObjectIdHex(v)}, bson.M{"sort": i})
	}
	return nil
}

// 释放ip
func (this *AccountServer) DoFreedIp(req *info.DoFreedIpReq, rsp *info.NullRsp) *goError.ErrRsp {

	accList := []string{}
	accinfoList := account.GetListAccountInfo(bson.M{"account": bson.M{"$in": req.Accounts}, "status": int64(1), "proxy_ip": bson.M{"$ne": ""}}, -1)
	for _, accountInfo := range accinfoList {
		accList = append(accList, accountInfo.Account)
	}
	account.UpAccountInfo(bson.M{"account": bson.M{"$in": accList}}, bson.M{"proxy_ip": ""})
	go func(accList2 []string) {
		for _, acc2 := range accList2 {
			proxyIp := cache.GetProxyIp(acc2)
			cache.IncIpUserNum(proxyIp.IpId, -1)
			cache.DelProxyIp(acc2)
		}
	}(accList)
	return nil
}

// 批量上线
func (this *AccountServer) DoBatchLogin(req *info.DoBatchLoginReq, rsp *info.NullRsp) *goError.ErrRsp {

	accList := []string{}
	w := bson.M{}
	w["status"] = bson.M{"$ne": int64(3)}
	w["account"] = bson.M{"$in": req.Accounts}
	accinfoList := account.GetListAccountInfo(w, -1)
	for _, accountInfo := range accinfoList {
		accList = append(accList, accountInfo.Account)
	}

	if len(accList) == 0 {
		//没有可登录的账号
		return goError.DoLoginErr
	}
	if len(accList) > 0 {
		account.UpAccountInfo(bson.M{"account": bson.M{"$in": accList}}, bson.M{"status": int64(3), "is_proxy_user": int64(0)})
		if req.IpCategory == 1 {
			for _, acc := range accList {
				proxyIp := cache.GetProxyIp(acc)
				cache.IncIpUserNum(proxyIp.IpId, -1)
				cache.DelProxyIp(acc)
			}
			//静态ip登录
			where := bson.M{}
			where["ip_category"] = int64(1)
			where["ip_type"] = req.IpType
			where["status"] = int64(1)
			where["disable_status"] = int64(1)
			where["country"] = req.Country
			listIp := ip.GetListIp(where, -1)
			ShuffleSlice(listIp)
			j := 0
			for _, ipInfo := range listIp {
				userNum := cache.GetIpUserNum(ipInfo.Id.Hex())
				if userNum >= ipInfo.AllotNum {
					continue
				}
				lnum := ipInfo.AllotNum - userNum
				for i := 0; i < int(lnum); i++ {
					ipTmp := cache.ProxyIpInfo{}
					ipTmp.IpId = ipInfo.Id.Hex()
					ipTmp.User = ipInfo.ProxyUser
					ipTmp.Pwd = ipInfo.ProxyPwd
					split := strings.Split(ipInfo.ProxyIp, ":")
					if len(split) > 0 {
						ipTmp.Port = split[len(split)-1]
						ipTmp.Host = strings.ReplaceAll(ipInfo.ProxyIp, ":"+ipTmp.Port, "")
					}
					ipTmp.Type = ipInfo.ProxyType
					ipTmp.Country = ipInfo.Country
					ipTmp.IpCategory = ipInfo.IpCategory
					ipTmp.IpType = ipInfo.IpType
					ipTmp.IpId = ipInfo.Id.Hex()
					cache.SetProxyIp(accList[j], &ipTmp)
					cache.IncIpUserNum(ipInfo.Id.Hex(), 1)
					j++
					if len(accList) <= j {
						break
					}
				}
				if len(accList) <= j {
					break
				}
			}
		}

		if req.IpCategory == 2 && req.IpType == 3 {
			//动态ip登录
			for _, acc := range accList {
				ipInfo := ip.GetByIdIp(req.IpId)
				ipTmp := cache.ProxyIpInfo{}
				ipTmp.IpId = ipInfo.Id.Hex()
				ipTmp.User = ipInfo.ProxyUser
				ipTmp.Pwd = ipInfo.ProxyPwd
				split := strings.Split(ipInfo.ProxyIp, ":")
				if len(split) > 0 {
					ipTmp.Port = split[len(split)-1]
					ipTmp.Host = strings.ReplaceAll(ipInfo.ProxyIp, ":"+ipTmp.Port, "")
				}
				ipTmp.Type = ipInfo.ProxyType
				ipTmp.Country = ipInfo.Country
				ipTmp.IpCategory = ipInfo.IpCategory
				ipTmp.IpType = ipInfo.IpType
				ipTmp.IpId = ipInfo.Id.Hex()
				cache.SetProxyIp(acc, &ipTmp)
				cache.IncIpUserNum(ipInfo.Id.Hex(), 1)
			}
		}
	}

	go func(accList []string) {
		bList := []string{}
		for _, acc := range accList {
			proxyIp := cache.GetProxyIp(acc)
			if proxyIp.Host == "" {
				cache.SetAccountStatus(acc, 1)
				account.UpAccountInfo(bson.M{"account": acc}, bson.M{"status": int64(1), "reason": "当前用户ip资源不足，请自行补充入库"})
				continue
			} else {
				go account.UpAccountInfo(bson.M{"account": acc}, bson.M{"proxy_ip": proxyIp.IpId, "ip_time": time.Now().Unix()})
			}
			bList = append(bList, acc)
		}
		groupAccList := [][]string{}
		aList := []string{}
		for i, acc := range bList {
			aList = append(aList, acc)
			if len(aList) >= 50 {
				groupAccList = append(groupAccList, aList)
				aList = []string{}
			}
			if i == len(bList)-1 {
				groupAccList = append(groupAccList, aList)
			}
		}
		for _, accList2 := range groupAccList {
			if len(accList2) > 0 {
				e := &event.TaskUserLoginEventReq{
					Account: accList2,
				}
				event.AddLoginTask(event.TaskTypeUserLogin, "", e)
			}
		}
	}(accList)
	return nil
}

func ShuffleSlice(s []*ip.Ip) {
	rand.Seed(time.Now().UnixNano())
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}
