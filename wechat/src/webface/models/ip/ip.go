package ip

import (
	"comm/comm"
	"comm/cos"
	"comm/goError"
	"comm/mgoDeal"
	"comm/tableName"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"selfComm/db/account"
	"selfComm/db/ip"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"strings"
	"time"
	"utils"
	info "webface/webstru"
)

// ip
type IpServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *IpServer) getUid() string {
	return this.Sess.Uid
}

// ip分组-列表
func (this *IpServer) GetIpGroupList(req *info.GetIpGroupListReq, rsp *info.GetIpGroupListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	where := bson.M{}
	rsp.List = []*info.GetIpGroupListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "-itime", start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetIpGroupListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
			count := ip.GetCountIp(bson.M{"group_id": tmp.Id})
			tmp.Count = count
		}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// ip分组-操作
func (this *IpServer) DoIpGroup(req *info.DoIpGroupReq, rsp *info.NullRsp) *goError.ErrRsp {
	if req.Ptype == 1 {
		//新增
		tmp := &ip.IpGroup{}
		tmp.Id = bson.NewObjectId()
		tmp.Name = req.Name
		ip.AddIpGroup(tmp)
	}
	if req.Ptype == 2 {
		//编辑
		where := bson.M{}
		where["_id"] = bson.ObjectIdHex(req.Id)
		update := bson.M{}
		update["name"] = req.Name
		ip.UpIpGroup(where, update)
	}
	if req.Ptype == 3 {
		//删除
		for _, delId := range req.DelId {
			count := ip.GetCountIp(bson.M{"group_id": delId})
			if count > 0 {
				return goError.DelIpGroupErr
			}
			where := bson.M{}
			where["_id"] = bson.ObjectIdHex(delId)
			ip.DelIpGroup(where)
		}
	}
	return nil
}

// ip-列表
func (this *IpServer) GetIpList(req *info.GetIpListReq, rsp *info.GetIpListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	where := bson.M{}
	if req.ProxyIp != "" {
		where["proxy_ip"] = req.ProxyIp
	}
	if req.StartTime > 0 && req.EndTime > 0 {
		where["expire_time"] = bson.M{"$gte": req.StartTime, "$lte": req.EndTime}
	}
	if req.Status > 0 {
		where["status"] = req.Status
	}
	if req.IpCategory > 0 {
		where["ip_category"] = req.IpCategory
	}
	if req.ExpireStatus > 0 {
		where["expire_status"] = req.ExpireStatus
	}
	if req.DisableStatus > 0 {
		where["disable_status"] = req.DisableStatus
	}
	if req.GroupId != "" {
		where["group_id"] = req.GroupId
	}
	rsp.List = []*info.GetIpListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, req.Sort, start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetIpListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
		}
		if p, ok := data["proxy_ip"]; ok {
			tmp.ProxyIp = utils.GetString(p)
		}
		if p, ok := data["user_num"]; ok {
			tmp.UserNum = utils.GetInt64(p)
		}
		if p, ok := data["proxy_user"]; ok {
			tmp.ProxyUser = utils.GetString(p)
		}
		if p, ok := data["status"]; ok {
			tmp.Status = utils.GetInt64(p)
		}
		if p, ok := data["allot_num"]; ok {
			tmp.AllotNum = utils.GetInt64(p)
		}
		if p, ok := data["ip_type"]; ok {
			tmp.IpType = utils.GetInt64(p)
		}
		if p, ok := data["country"]; ok {
			tmp.Country = utils.GetString(p)
		}
		if p, ok := data["disable_status"]; ok {
			tmp.DisableStatus = utils.GetInt64(p)
		}
		if p, ok := data["expire_time"]; ok {
			tmp.ExpireTime = utils.GetInt64(p)
		}
		if p, ok := data["expire_status"]; ok {
			tmp.ExpireStatus = utils.GetInt64(p)
		}
		if p, ok := data["remark"]; ok {
			tmp.Remark = utils.GetString(p)
		}
		if p, ok := data["reason"]; ok {
			tmp.Reason = utils.GetString(p)
		}
		if p, ok := data["ip_category"]; ok {
			tmp.IpCategory = utils.GetInt64(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	go func() {
		//检测ip失效状态
		listIp := ip.GetListIp(bson.M{}, -1)
		for _, ipInfo := range listIp {
			isUp := false
			expireStatus := int64(1)
			if ipInfo.ExpireTime-time.Now().Unix() <= 0 {
				expireStatus = 2
			} else if ipInfo.ExpireTime-time.Now().Unix() <= 60*60*24*3 && ipInfo.ExpireTime-time.Now().Unix() > 0 {
				expireStatus = 3
			}
			//userNum := redisDeal.RedisDoHGetInt(redisKeys.GetIpUserNumKey(uid), ipInfo.Id.Hex())
			userNum := cache.GetIpUserNum(ipInfo.Id.Hex())
			if expireStatus != ipInfo.ExpireStatus || ipInfo.UserNum != userNum {
				isUp = true
			}
			if isUp {
				ip.UpIp(bson.M{"_id": ipInfo.Id}, bson.M{"expire_status": expireStatus, "user_num": userNum})
			}
		}
	}()
	return nil
}

// 检查文件
func (this *IpServer) CheckFile(fileContent string, req *info.CheckFileReq, rsp *info.CheckFileRsp) *goError.ErrRsp {
	fileContent = strings.ReplaceAll(fileContent, "\r\n", "\n")
	fileContent = strings.ReplaceAll(fileContent, "\n\n", "\n")
	ipStrArr := strings.Split(fileContent, "\n")
	failNumber := 0
	successList := []string{}
	successNumber := 0
	exportStr := ""

	//获取所有的ip
	if req.Ptype == 1 {
		//静态ip
		for _, ipStr := range ipStrArr {
			split := strings.Split(ipStr, ",")
			if len(split) != 5 {
				failNumber++
				exportStr = exportStr + ipStr + "\n"
				continue
			} else {
				if split[4] != "socks5" && split[4] != "http" {
					failNumber++
					exportStr = exportStr + ipStr + "\n"
					continue
				}
				successNumber++
				successList = append(successList, ipStr)
			}
		}
	}
	if req.Ptype == 2 {
		//动态ip
		for _, ipStr := range ipStrArr {
			split := strings.Split(ipStr, ",")
			if len(split) != 6 {
				failNumber++
				exportStr = exportStr + ipStr + "\n"
				continue
			} else {
				if split[4] != "socks5" && split[4] != "http" {
					failNumber++
					exportStr = exportStr + ipStr + "\n"
					continue
				}
				successNumber++
				successList = append(successList, ipStr)
			}
		}
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
	rsp.FailNumber = failNumber
	rsp.SuccessNumber = successNumber
	rsp.SuccessList = successList
	return nil
}

// 添加ip
func (this *IpServer) AddIp(req *info.AddIpReq, rsp *info.AddIpRsp) *goError.ErrRsp {

	successNumber := 0
	failNumber := 0
	for _, ipStr := range req.SuccessList {
		split := strings.Split(ipStr, ",")
		if req.IpCategory == 1 {
			tmp := &ip.Ip{}
			tmp.Id = bson.NewObjectId()
			tmp.GroupId = req.GroupId
			tmp.ProxyIp = split[0] + ":" + split[1]
			tmp.ProxyUser = split[2]
			tmp.ProxyPwd = split[3]
			tmp.ProxyType = split[4]
			tmp.Status = 1
			tmp.AllotNum = req.AllotNum
			tmp.UserNum = int64(0)
			tmp.IpType = req.IpType
			tmp.IpCategory = req.IpCategory
			tmp.ExpireStatus = 1
			tmp.DisableStatus = 1
			tmp.Country = req.Country
			tmp.ExpireTime = req.ExpireTime
			if tmp.ExpireTime-time.Now().Unix() <= 0 {
				tmp.ExpireStatus = 2
			} else if tmp.ExpireTime-time.Now().Unix() <= 60*60*24*3 && tmp.ExpireTime-time.Now().Unix() > 0 {
				tmp.ExpireStatus = 3
			}
			err := ip.AddIp(tmp)
			if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
				failNumber++
				continue
			}
		}
		if req.IpCategory == 2 {
			tmp := &ip.Ip{}
			tmp.Id = bson.NewObjectId()
			tmp.GroupId = req.GroupId
			tmp.ProxyIp = split[0] + ":" + split[1]
			tmp.ProxyUser = split[2]
			tmp.ProxyPwd = split[3]
			tmp.ProxyType = split[4]
			tmp.Status = 1
			tmp.AllotNum = req.AllotNum
			tmp.UserNum = int64(0)
			tmp.IpType = 1
			tmp.IpCategory = req.IpCategory
			tmp.ExpireStatus = 1
			tmp.Country = split[5]
			tmp.ExpireTime = req.ExpireTime
			tmp.DisableStatus = 1
			if tmp.ExpireTime-time.Now().Unix() <= 0 {
				tmp.ExpireStatus = 2
			} else if tmp.ExpireTime-time.Now().Unix() <= 60*60*24*3 && tmp.ExpireTime-time.Now().Unix() > 0 {
				tmp.ExpireStatus = 3
			}
			err := ip.AddIp(tmp)
			if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
				failNumber++
				continue
			}
		}
		successNumber++
	}
	rsp.SuccessNumber = successNumber
	rsp.FailNumber = failNumber
	return nil
}

// 设置到期时间
func (this *IpServer) DoExpireTime(req *info.DoExpireTimeReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		expireStatus := 1
		if req.ExpireTime-time.Now().Unix() <= 0 {
			expireStatus = 2
		} else if req.ExpireTime-time.Now().Unix() <= 60*60*24*3 && req.ExpireTime-time.Now().Unix() > 0 {
			expireStatus = 3
		}
		ip.UpIp(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"expire_time": req.ExpireTime, "expire_status": expireStatus})
	}
	return nil
}

// 分配规则
func (this *IpServer) DoAllotNum(req *info.DoAllotNumReq, rsp *info.NullRsp) *goError.ErrRsp {

	if req.AllotNum <= 0 {
		return goError.GLOBAL_INVALIDPARAM
	}
	for _, id := range req.Ids {
		//useNum := redisDeal.RedisDoHGetInt(redisKeys.GetIpUserNumKey(uid), id)
		useNum := cache.GetIpUserNum(id)
		if req.AllotNum > useNum {
			ip.UpIp(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"allot_num": req.AllotNum})
		}
	}
	return nil
}

// 移动分组
func (this *IpServer) DoMoveIpGroup(req *info.DoMoveIpGroupReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		ip.UpIp(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"group_id": req.GroupId})
	}
	return nil
}

// 网络检测
func (this *IpServer) DoCheckStatus(req *info.DoCheckStatusReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		go CheckIp(id)
	}
	return nil
}

// 检测ip
func CheckIp(id string) {
	//先更新ip为检测中
	w := bson.M{"_id": bson.ObjectIdHex(id)}
	up := bson.M{"status": int64(3)}
	ip.UpIp(w, up)
	ipInfo := ip.GetByIdIp(id)
	boo := wxComm.GetRealIp(ipInfo.ProxyType, ipInfo.ProxyUser, ipInfo.ProxyPwd, ipInfo.ProxyIp)
	if boo {
		w := bson.M{"_id": bson.ObjectIdHex(id)}
		up := bson.M{"status": int64(1)}
		up["reason"] = ""
		ip.UpIp(w, up)
	} else {
		w := bson.M{"_id": bson.ObjectIdHex(id)}
		up := bson.M{"status": int64(2)}
		up["reason"] = "代理连接失败，请重试"
		ip.UpIp(w, up)
	}
}

// ip启动分配
func (this *IpServer) DoStartDistribution(req *info.DoStartDistributionReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		w := bson.M{"_id": bson.ObjectIdHex(id)}
		up := bson.M{"disable_status": int64(1)}
		ip.UpIp(w, up)
	}
	return nil
}

// ip禁用分配
func (this *IpServer) DoDisableAllocation(req *info.DoDisableAllocationReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		w := bson.M{"_id": bson.ObjectIdHex(id)}
		up := bson.M{"disable_status": int64(2)}
		ip.UpIp(w, up)
	}
	return nil
}

// 批量删除
func (this *IpServer) DoBatchDel(req *info.DoBatchDelReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		w := bson.M{"_id": bson.ObjectIdHex(id)}
		ip.DelIp(w)
		cache.DelIpUserNum(id)
	}
	return nil
}

// 修改国家
func (this *IpServer) DoUpCountry(req *info.DoUpCountryReq, rsp *info.NullRsp) *goError.ErrRsp {

	for _, id := range req.Ids {
		w := bson.M{"_id": bson.ObjectIdHex(id)}
		up := bson.M{"country": req.Country}
		ip.UpIp(w, up)
	}
	return nil
}

// 批量导出
func (this *IpServer) DoOutPutIp(req *info.DoOutPutIpReq, rsp *info.DoOutPutIpRsp) *goError.ErrRsp {

	exportStr := ""
	for _, id := range req.Ids {
		ipinfo := ip.GetByIdIp(id)
		split := strings.Split(ipinfo.ProxyIp, ":")
		ipStr := split[0] + "," + split[1] + "," + ipinfo.ProxyUser + "," + ipinfo.ProxyPwd + "," + ipinfo.ProxyType + "," + ipinfo.Country
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

// 获取ipv4分配
func (this *IpServer) GetIpV4Allot(req *info.NullReq, rsp *info.GetIpV4AllotRsp) *goError.ErrRsp {
	uid := this.Sess.Uid
	totalCount := int64(0)
	UseNum := int64(0)
	sumwhere := []bson.M{}
	sumwhere = append(sumwhere, bson.M{"$match": bson.M{"ip_type": int64(1)}})
	sumwhere = append(sumwhere, bson.M{"$match": bson.M{"ip_category": int64(1)}})
	sumwhere = append(sumwhere, bson.M{"$group": bson.M{"_id": "null", "num_sum": bson.M{"$sum": "$user_num"}}})
	sumRet, err := mgoDeal.QueryMongoSum(comm.GetUserMgoDBName(uid), tableName.GetTableIpListInfo(), sumwhere)
	if err == nil {
		if p, ok := sumRet["num_sum"]; ok {
			UseNum = utils.GetInt64(p)
		}
	}
	sumwhere2 := []bson.M{}
	sumwhere2 = append(sumwhere2, bson.M{"$match": bson.M{"ip_type": int64(1)}})
	sumwhere2 = append(sumwhere2, bson.M{"$match": bson.M{"ip_category": int64(1)}})
	sumwhere2 = append(sumwhere2, bson.M{"$group": bson.M{"_id": "null", "num_sum": bson.M{"$sum": "$allot_num"}}})
	sumRet2, err := mgoDeal.QueryMongoSum(comm.GetUserMgoDBName(uid), tableName.GetTableIpListInfo(), sumwhere2)
	if err == nil {
		if p, ok := sumRet2["num_sum"]; ok {
			totalCount = utils.GetInt64(p)
		}
	}
	rsp.TotalCount = totalCount
	rsp.UseNum = UseNum
	rsp.NoUserNum = rsp.TotalCount - rsp.UseNum
	return nil
}

// 获取ipv6分配
func (this *IpServer) GetIpV6Allot(req *info.NullReq, rsp *info.GetIpV6AllotRsp) *goError.ErrRsp {
	uid := this.Sess.Uid
	totalCount := int64(0)
	UseNum := int64(0)
	sumwhere := []bson.M{}
	sumwhere = append(sumwhere, bson.M{"$match": bson.M{"ip_type": int64(2)}})
	sumwhere = append(sumwhere, bson.M{"$match": bson.M{"ip_category": int64(1)}})
	sumwhere = append(sumwhere, bson.M{"$group": bson.M{"_id": "null", "num_sum": bson.M{"$sum": "$user_num"}}})
	sumRet, err := mgoDeal.QueryMongoSum(comm.GetUserMgoDBName(uid), tableName.GetTableIpListInfo(), sumwhere)
	if err == nil {
		if p, ok := sumRet["num_sum"]; ok {
			UseNum = utils.GetInt64(p)
		}
	}
	sumwhere2 := []bson.M{}
	sumwhere2 = append(sumwhere2, bson.M{"$match": bson.M{"ip_type": int64(2)}})
	sumwhere2 = append(sumwhere2, bson.M{"$match": bson.M{"ip_category": int64(1)}})
	sumwhere2 = append(sumwhere2, bson.M{"$group": bson.M{"_id": "null", "num_sum": bson.M{"$sum": "$allot_num"}}})
	sumRet2, err := mgoDeal.QueryMongoSum(comm.GetUserMgoDBName(uid), tableName.GetTableIpListInfo(), sumwhere2)
	if err == nil {
		if p, ok := sumRet2["num_sum"]; ok {
			totalCount = utils.GetInt64(p)
		}
	}
	rsp.TotalCount = totalCount
	rsp.UseNum = UseNum
	rsp.NoUserNum = rsp.TotalCount - rsp.UseNum
	return nil
}

// 获取动态ip分配
func (this *IpServer) GetIpDynamicAllot(req *info.NullReq, rsp *info.GetIpDynamicAllotRsp) *goError.ErrRsp {
	count1 := ip.GetCountIp(bson.M{"ip_category": int64(2)})
	rsp.TotalCount = count1
	return nil
}

// 获取国家
func (this *IpServer) GetCountryList(req *info.NullReq, rsp *info.GetCountryListRsp) *goError.ErrRsp {
	countryStr := "不丹,东帝汶,中非共和国,丹麦,乌克兰,乌兹别克斯坦,乌干达,乌拉圭,乍得,乔治亚,也门,亚美尼亚,以色列,伊拉克,伊朗,伯利兹,佛得角,俄罗斯,保加利亚,克罗地亚,关岛,冈比亚,冰岛,几内亚,几内亚比绍,列支敦士登,刚果,刚果共和国,刚果民主共和国,利比亚,利比里亚,加拿大,加纳,加蓬,匈牙利,北马里亚纳群岛,南斯拉夫,南非,博茨瓦纳,卡塔尔,卢旺达,卢森堡,印度,印度尼西亚,危地马拉,厄瓜多尔,厄立特里亚,叙利亚,古巴,台湾,吉尔吉斯斯坦,吉布提,哈萨克斯坦,哥伦比亚,哥斯达黎加,喀麦隆,图瓦卢,土库曼斯坦,土耳其,圣卢西亚,圣基茨和尼维斯,圣多美和普林西比,圣文森特和格林纳丁斯,圣文森特岛,圣皮埃尔和米克隆群岛,圣诞岛,圣赫勒拿,圣马力诺,圭亚那,坦桑尼亚,埃及,埃塞俄比亚,基里巴斯,塔吉克斯坦,塞内加尔,塞尔维亚共和国,塞拉利昂,塞浦路斯,塞舌尔,墨西哥,多哥,多米尼加,多米尼加共和国,奥地利,委内瑞拉,孟加拉,孟加拉国,安哥拉,安圭拉,安提瓜和巴布达,安提瓜岛和巴布达,安道尔,尼加拉瓜,尼日利亚,尼日尔,尼泊尔,巴勒斯坦,巴哈马,巴基斯坦,巴巴多斯,巴巴多斯岛,巴布亚新几内亚,巴拉圭,巴拿马,巴林,巴西,布基纳法索,布隆迪,希腊,帕劳群岛,库克群岛,开曼群岛,弗兰克群岛,德国,意大利,所罗门群岛,扎伊尔,托克劳,拉脱维亚,挪威,捷克共和国,摩尔多瓦,摩洛哥,摩纳哥,文莱,斐济,斯威士兰,斯洛伐克,斯洛文尼亚,斯里兰卡,新加坡,新喀里多尼亚,新西兰,日本,智利,朝鲜,柬埔寨,格林纳达,格陵兰,格鲁吉亚,梵蒂冈,比利时,毛里塔尼亚,毛里求斯,汤加,沙特阿拉伯,法国,法属圭亚那,法属波利尼西亚,法罗群岛,波兰,波多黎各,波斯尼亚和黑塞哥维那,泰国,津巴布韦,洪都拉斯,海地,澳大利亚,爱尔兰,爱沙尼亚,牙买加,特克斯和凯克特斯群岛,特立尼达和多巴哥,玻利维亚,瑙鲁,瑞典,瑞士,瓜地马拉,瓜德罗普,瓦利斯和福图纳,瓦努阿图,留尼旺岛,白俄罗斯,百慕大,直布罗陀,科威特,科摩罗,科特迪瓦,科科斯群岛,秘鲁,突尼斯,立陶宛,索马里,约旦,纳米比亚,纽埃,维尔京群岛，美属,维尔京群岛，英属,缅甸,罗马尼亚,美国,老挝,肯尼亚,芬兰,苏丹,苏里南,英国,荷兰,莫桑比克,莱索托,菲律宾,萨尔瓦多,萨摩亚,葡萄牙,蒙古,蒙特塞拉特,西班牙,诺福克,贝宁,赞比亚,越南,阿塞拜疆,阿富汗,阿尔及利亚,阿尔巴尼亚,阿拉伯联合酋长国,阿曼,阿根廷,阿鲁巴,韩国,香港,马其顿,马尔代夫,马拉维,马提尼克,马来西亚,马约特岛,马绍尔群岛,马耳他,马达加斯加,马里,黎巴嫩"
	split := strings.Split(countryStr, ",")
	rsp.CountryList = split
	return nil
}

// 编辑备注
func (this *IpServer) DoIpRemark(req *info.DoIpRemarkReq, rsp *info.DoIpRemarkRsp) *goError.ErrRsp {

	ip.UpIp(bson.M{"_id": bson.ObjectIdHex(req.Id)}, bson.M{"remark": req.Remark})
	rsp.Remark = req.Remark
	return nil
}

// ip-动态
func (this *IpServer) GetDynamicIp(req *info.NullReq, rsp *info.GetDynamicIpRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	where := bson.M{}
	where["ip_category"] = int64(2)
	where["status"] = int64(1)
	where["disable_status"] = int64(1)
	rsp.List = []*info.GetDynamicIpInfo{}
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "", -1, -1)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetDynamicIpInfo{}
		if p, ok := data["_id"]; ok {
			tmp.IpId = utils.GetString(p)
		}
		if p, ok := data["country"]; ok {
			tmp.Country = utils.GetString(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// ip-静态
func (this *IpServer) GetStaticIp(req *info.GetStaticIpReq, rsp *info.GetStaticIpRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	where := bson.M{}
	where["ip_category"] = int64(1)
	where["ip_type"] = req.IpType
	where["status"] = int64(1)
	where["disable_status"] = int64(1)
	rsp.List = []*info.GetStaticIpInfo{}
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "", -1, -1)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	m := make(map[string]string)
	for _, data := range all {
		tmp := &info.GetStaticIpInfo{}
		if p, ok := data["country"]; ok {
			tmp.Country = utils.GetString(p)
		}
		if _, ok := m[tmp.Country]; ok {
			continue
		} else {
			m[tmp.Country] = tmp.Country
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// ip校正工具
func (this *IpServer) DoResetIp(req *info.NullReq, rsp *info.NullRsp) *goError.ErrRsp {
	go func() {

		listIp := ip.GetListIp(bson.M{}, -1)
		for _, ipInfo := range listIp {
			useNum := account.GetCountAccountInfo(bson.M{"proxy_ip": ipInfo.Id.Hex()})
			cache.SetIpUserNum(ipInfo.Id.Hex(), useNum)
		}
	}()
	return nil
}

// 获取分配的ip
func (this *IpServer) GetUseList(req *info.GetUseListReq, rsp *info.GetUseListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	where := bson.M{}
	if req.StartTime > 0 && req.EndTime > 0 {
		where["ip_time"] = bson.M{"$gte": req.StartTime, "$lte": req.EndTime}
	}
	where["proxy_ip"] = req.Id

	rsp.List = []*info.GetUseListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}

	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "ip_time", start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetUseListInfo{}
		if p, ok := data["account"]; ok {
			tmp.Account = utils.GetString(p)
		}

		if p, ok := data["ip_time"]; ok {
			tmp.IpTime = utils.GetInt64(p)
		}
		tmp.ProxyIp = req.ProxyIp
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}
