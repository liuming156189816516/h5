package account

import (
	"archive/zip"
	"comm/comm"
	"comm/cos"
	"comm/goError"
	"comm/mgoDeal"
	"comm/tableName"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"selfComm/db/account"
	"selfComm/wxComm"
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
	if len(req.Accounts) > 0 {
		err := account.UpAccountInfo(bson.M{"account": bson.M{"$in": req.Accounts}}, bson.M{"group_id": req.GroupId})
		if err == nil {
			for _, acc := range req.Accounts {
				accountInfo := cache.GetAccountInfo(acc)
				accountInfo.AccountGroup = req.GroupId
				cache.SetAccountInfo(acc, accountInfo)
			}
		}
	}

	return nil
}

func (this *AccountServer) DoOutPutAccount(req *info.DoOutPutAccountReq, rsp *info.DoOutPutAccountRsp) *goError.ErrRsp {

	tmpPath := beego.AppConfig.String("tmpPath")

	// 👉 创建临时目录
	dirName := fmt.Sprintf("export_%d", time.Now().UnixNano())
	baseDir := filepath.Join(tmpPath, dirName)

	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	defer os.RemoveAll(baseDir)

	// 👉 1. 生成每个账号的 JSON 文件
	for _, acc := range req.Accounts {

		accInfo := account.GetOneAccountInfo(bson.M{"account": acc})
		if accInfo == nil || accInfo.Token == "" {
			continue
		}

		// ✅ 关键：把 Token（JSON字符串）解析成对象
		var obj interface{}
		if err := json.Unmarshal([]byte(accInfo.Token), &obj); err != nil {
			fmt.Println("非法JSON:", acc, err)
			continue
		}

		data, err := json.Marshal(obj)
		if err != nil {
			continue
		}

		// 👉 文件名：账号.json（防路径攻击）
		fileName := filepath.Base(acc) + ".json"
		filePath := filepath.Join(baseDir, fileName)

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			continue
		}
	}

	// 👉 2. 创建 zip 文件
	zipName := fmt.Sprintf("%d.zip", time.Now().UnixNano())
	zipPath := filepath.Join(tmpPath, zipName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	defer zipFile.Close()
	defer os.Remove(zipPath)

	zipWriter := zip.NewWriter(zipFile)

	// 👉 3. 写入 zip
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		// zip 内文件名（只保留文件名）
		zipEntryName := filepath.Base(path)

		writer, err := zipWriter.Create(zipEntryName)
		if err != nil {
			return nil
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}

	if err := zipWriter.Close(); err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}

	// 👉 4. 上传 zip
	fileUrl := cos.UploadAwsFile(zipPath, zipName)
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
	// ✅ 1️⃣ 随机打乱
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(accinfoList), func(i, j int) {
		accinfoList[i], accinfoList[j] = accinfoList[j], accinfoList[i]
	})

	go func(list []*account.AccountInfo) {

		// 👉 提取账号字符串（修复 bug）
		var accNames []string
		for _, v := range list {
			accNames = append(accNames, v.Account)
		}

		// 👉 批量更新状态 = 3（登录中）
		account.UpAccountInfo(
			bson.M{"account": bson.M{"$in": accNames}},
			bson.M{"status": int64(3)},
		)

		for _, accountInfo := range list {

			importJson, err := wxComm.ImportJson(
				accountInfo.Account,
				json.RawMessage(accountInfo.Token),
			)

			if err != nil || !importJson.Ok {
				account.UpAccountInfo(
					bson.M{"account": accountInfo.Account},
					bson.M{"status": int64(4)},
				)
			}

			// ✅ 2️⃣ 每50个暂停20秒
			/*if (i+1)%50 == 0 {
				fmt.Println("已处理:", i+1, "暂停20秒...")
				time.Sleep(20 * time.Second)
			}*/
		}

	}(accinfoList)
	return nil
}

// 批量下线
func (this *AccountServer) DoBatchLogout(req *info.DoBatchLogoutReq, rsp *info.NullRsp) *goError.ErrRsp {
	accList := req.Accounts
	account.UpAccountInfo(bson.M{"account": bson.M{"$in": accList}}, bson.M{"status": int64(1)})
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

	accList := []*account.AccountInfo{}
	w := bson.M{}
	w["status"] = bson.M{"$ne": int64(3)}
	w["account"] = bson.M{"$in": req.Accounts}

	accinfoList := account.GetListAccountInfo(w, -1)

	for _, accountInfo := range accinfoList {
		accList = append(accList, accountInfo)
	}

	if len(accList) == 0 {
		return goError.DoLoginErr
	}

	// ✅ 1️⃣ 随机打乱
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(accinfoList), func(i, j int) {
		accinfoList[i], accinfoList[j] = accinfoList[j], accinfoList[i]
	})

	go func(list []*account.AccountInfo) {

		// 👉 提取账号字符串（修复 bug）
		var accNames []string
		for _, v := range list {
			accNames = append(accNames, v.Account)
		}

		// 👉 批量更新状态 = 3（登录中）
		account.UpAccountInfo(
			bson.M{"account": bson.M{"$in": accNames}},
			bson.M{"status": int64(3)},
		)

		for _, accountInfo := range list {

			importJson, err := wxComm.ImportJson(
				accountInfo.Account,
				json.RawMessage(accountInfo.Token),
			)

			if err != nil || !importJson.Ok {
				account.UpAccountInfo(
					bson.M{"account": accountInfo.Account},
					bson.M{"status": int64(4)},
				)
			}

			// ✅ 2️⃣ 每50个暂停20秒
			/*if (i+1)%50 == 0 {
				fmt.Println("已处理:", i+1, "暂停20秒...")
				time.Sleep(20 * time.Second)
			}*/
		}

	}(accinfoList)

	return nil
}

// 入库文件-列表
func (this *AccountServer) GetAccountFileList(req *info.GetAccountFileListReq, rsp *info.GetAccountFileListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	where := bson.M{}
	if req.Name != "" {
		where["name"] = req.Name
	}
	if req.StartTime > 0 && req.EndTime > 0 {
		where["itime"] = bson.M{"$gte": req.StartTime, "$lte": req.EndTime}
	}
	rsp.List = []*info.GetAccountFileListInfo{}
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
		tmp := &info.GetAccountFileListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
		}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
		}
		if p, ok := data["account_type"]; ok {
			tmp.AccountType = utils.GetInt64(p)
		}
		if p, ok := data["status"]; ok {
			tmp.Status = utils.GetInt64(p)
		}
		if p, ok := data["success_num"]; ok {
			tmp.SuccessNum = utils.GetInt64(p)
		}
		if p, ok := data["fail_num"]; ok {
			tmp.FailNum = utils.GetInt64(p)
		}
		if p, ok := data["remark"]; ok {
			tmp.Remark = utils.GetString(p)
		}
		if p, ok := data["itime"]; ok {
			tmp.Itime = utils.GetInt64(p)
		}
		tmp.AccType = "WS"
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// 入库文件-批量删除
func (this *AccountServer) DoBathDelAccountFile(req *info.DoBathDelAccountFileReq, rsp *info.NullRsp) *goError.ErrRsp {
	for _, id := range req.Ids {
		account.DelAccountFile(bson.M{"_id": bson.ObjectIdHex(id)})
		account.DelAccountLog(bson.M{"file_id": id})
	}
	return nil
}

// 检查文件
func (this *AccountServer) CheckAccountFile(fileContent string, req *info.NullReq, rsp *info.CheckAccountFileRsp) *goError.ErrRsp {
	fileContent = strings.ReplaceAll(fileContent, "\r\n", "\n")
	fileContent = strings.ReplaceAll(fileContent, "\n\n", "\n")
	fileContent = strings.ReplaceAll(fileContent, "\r", "\n")
	ipStrArr := strings.Split(fileContent, "\n")
	failNumber := 0
	successList := []string{}
	successNumber := 0
	exportStr := ""
	for _, ipStr := range ipStrArr {
		upAccount := info.UpJsonAccount{}
		jsoniter.UnmarshalFromString(ipStr, &upAccount)
		if upAccount.Phone == "" {
			failNumber++
			exportStr = exportStr + ipStr + "\n"
			continue
		}
		successNumber++
		successList = append(successList, ipStr)
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

// 添加账号
func (this *AccountServer) AddAccount(req *info.AddAccountReq, rsp *info.AddAccountRsp) *goError.ErrRsp {
	uid := this.Sess.Uid
	tmp1 := &account.AccountFile{}
	tmp1.Id = bson.ObjectIdHex(req.FileId)
	tmp1.Name = req.Name
	tmp1.AccountType = req.AccountType
	tmp1.Remark = req.Remark
	tmp1.Status = 1
	account.AddAccountFile(tmp1)
	rsp.Id = tmp1.Id.Hex()

	go func(uid, fileId string, accountType int64) {
		successNumber := 0
		failNumber := 0
		logList := make([]interface{}, 0)
		accList := make([]interface{}, 0)
		for _, accountStr := range req.SuccessList {
			acc := ""
			token := accountStr
			upAccount := info.UpJsonAccount{}
			jsoniter.UnmarshalFromString(accountStr, &upAccount)
			acc = upAccount.Phone

			tmp2 := &account.AccountLog{}
			tmp2.Id = bson.NewObjectId()
			tmp2.FileId = fileId
			tmp2.Account = acc
			tmp2.Status = 2
			tmp2.Token = token
			tmp := &account.AccountInfo{}
			upTime := cache.GetAllAccountList(tmp2.Account)
			if upTime != "" {
				failNumber++
				tmp2.Status = 1
				tmp2.Reason = "账号重复: " + utils.TimeIntToString(utils.StrToInt64(upTime))
			} else {
				tmp.Id = bson.NewObjectId()
				tmp.Account = acc
				tmp.GroupId = req.GroupId
				tmp.Status = 1
				tmp.AccountType = accountType
				tmp.Token = token
				tmp.Itime = time.Now().Unix()
				tmp.Ptime = time.Now().Unix()
				tmp.Creator = uid
				successNumber++
				accList = append(accList, tmp)
				cache.SetAllAccountList(tmp.Account)
				cache.SetAccountStatus(tmp.Account, tmp.Status)
				accInfo := &cache.AccountInfo{
					Account:      tmp.Account,
					AccountType:  tmp.AccountType,
					Token:        tmp.Token,
					AccountGroup: req.GroupId,
				}
				cache.SetAccountInfo(tmp.Account, accInfo)
			}
			logList = append(logList, tmp2)
			time.Sleep(4 * time.Millisecond)
		}
		account.AddBatchAccountInfo(accList)
		go account.AddBatchAccountLog(logList)
		account.UpAccountFile(bson.M{"_id": bson.ObjectIdHex(fileId)}, bson.M{"status": int64(2), "success_num": successNumber, "fail_num": failNumber})
	}(uid, tmp1.Id.Hex(), req.AccountType)
	return nil
}

// 获取上传结果
func (this *AccountServer) GetAccountSchedule(req *info.GetAccountScheduleReq, rsp *info.GetAccountScheduleRsp) *goError.ErrRsp {
	fileInfo := account.GetByIdAccountFile(req.Id)
	rsp.UpStatus = fileInfo.Status
	rsp.Success = fileInfo.SuccessNum
	rsp.Fail = fileInfo.FailNum
	return nil
}
