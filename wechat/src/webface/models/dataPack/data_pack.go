package dataPack

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
	"selfComm/db/dataPack"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"strings"
	"sync"
	"utils"
	info "webface/webstru"
)

// 数据包
type DataPackServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *DataPackServer) getUid() string {
	return this.Sess.Uid
}

// 数据包-列表
func (this *DataPackServer) GetDataPackList(req *info.GetDataPackListReq, rsp *info.GetDataPackListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	where := bson.M{}
	if req.Name != "" {
		where["name"] = req.Name
	}
	rsp.List = []*info.GetDataPackListInfo{}
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
		tmp := &info.GetDataPackListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
			tmp.ResidueNum = cache.ScardDataPackListCount(tmp.Id)
			tmp.ErrNum = cache.ScardDataPackListErrCount(tmp.Id)
		}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
		}
		if p, ok := data["up_num"]; ok {
			tmp.UpNum = utils.GetInt64(p)
		}
		if p, ok := data["repeat_num"]; ok {
			tmp.RepeatNum = utils.GetInt64(p)
		}
		if p, ok := data["invalid_num"]; ok {
			tmp.InvalidNum = utils.GetInt64(p)
		}
		if p, ok := data["into_num"]; ok {
			tmp.IntoNum = utils.GetInt64(p)
		}
		if p, ok := data["up_status"]; ok {
			tmp.UpStatus = utils.GetInt64(p)
		}
		if p, ok := data["itime"]; ok {
			tmp.Itime = utils.GetInt64(p)
		}
		tmp.SourceRepeatNum = tmp.UpNum - tmp.InvalidNum - tmp.IntoNum - tmp.RepeatNum
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// 上传数据包
func (this *DataPackServer) UpLoadFile(fileContent string, req *info.UpLoadFileReq, rsp *info.UpLoadFileRsp) *goError.ErrRsp {
	uid := this.getUid()
	fileContent = strings.ReplaceAll(fileContent, "\r\n", "\n")
	fileContent = strings.ReplaceAll(fileContent, "\n\n", "\n")
	phoneArr := strings.Split(fileContent, "\n")
	ioutil.WriteFile("./dataPack/"+uid+"_"+comm.Md5(fileContent)+".txt", []byte(fileContent), 0644)
	if req.Ptype == 1 {
		//新增
		tmp := &dataPack.DataPack{}
		tmp.Id = bson.NewObjectId()
		tmp.Name = req.Name
		tmp.UpNum = int64(len(phoneArr))
		tmp.UpStatus = 1
		dataPack.AddDataPack(tmp)
		go func(phoneList []string, tmpId, uid string, intoType int64) {
			memberList := cache.GetAllDataPackList()
			map2 := map[string]int64{}
			if req.IntoType == 1 {
				for _, s := range memberList {
					map2[s] = 1
				}
			}
			invalidNum := int64(0)
			repeatNum := int64(0)
			var args []interface{}
			for _, phoneStr := range phoneList {
				if !strings.Contains(phoneStr, "-") {
					if strings.Contains(phoneStr, ",") {
						phoneStrNew := strings.ReplaceAll(phoneStr, ",", "")
						toInt64 := utils.StrToInt64(phoneStrNew)
						if toInt64 <= 0 {
							invalidNum++
							continue
						}
					} else {
						toInt64 := utils.StrToInt64(phoneStr)
						if toInt64 <= 0 {
							invalidNum++
							continue
						}
					}

					if wxComm.IsChain(phoneStr) {
						invalidNum++
						continue
					}
					if _, ok := map2[phoneStr]; ok {
						repeatNum++
						continue
					}
				}
				args = append(args, phoneStr)
			}
			cache.SaddListDataPackList(tmpId, args)
			cache.SaddListAllDataPackList(args)
			cache.SaddListDataPackList2(tmpId, args)
			wh := bson.M{"_id": bson.ObjectIdHex(tmpId)}
			up := bson.M{}
			intoNum := cache.ScardDataPackListCount2(tmpId)
			up["into_num"] = intoNum
			up["up_status"] = int64(2)
			up["invalid_num"] = invalidNum
			up["repeat_num"] = repeatNum
			dataPack.UpDataPack(wh, up)
			cache.SetSchedule(tmpId, utils.IntToStr(intoNum)+"_"+utils.IntToStr(int64(len(phoneList))-intoNum))
		}(phoneArr, tmp.Id.Hex(), uid, req.IntoType)
		rsp.Id = tmp.Id.Hex()
	}
	if req.Ptype == 2 {
		//补充
		tmp := dataPack.GetByIdDataPack(req.Id)
		wh := bson.M{"_id": bson.ObjectIdHex(req.Id)}
		up := bson.M{"up_status": int64(1)}
		up["up_num"] = tmp.UpNum + int64(len(phoneArr))
		dataPack.UpDataPack(wh, up)
		go func(phoneList []string, tmpId, uid string) {
			memberList := cache.GetAllDataPackList()
			map2 := map[string]int64{}
			for _, s := range memberList {
				map2[s] = 1
			}
			invalidNum := int64(0)
			repeatNum := int64(0)
			var args []interface{}
			for _, phoneStr := range phoneList {
				if !strings.Contains(phoneStr, "-") {
					if strings.Contains(phoneStr, ",") {
						phoneStrNew := strings.ReplaceAll(phoneStr, ",", "")
						toInt64 := utils.StrToInt64(phoneStrNew)
						if toInt64 <= 0 {
							invalidNum++
							continue
						}
					} else {
						toInt64 := utils.StrToInt64(phoneStr)
						if toInt64 <= 0 {
							invalidNum++
							continue
						}
					}
					if wxComm.IsChain(phoneStr) {
						invalidNum++
						continue
					}
					if _, ok := map2[phoneStr]; ok {
						repeatNum++
						continue
					}
				}
				args = append(args, phoneStr)
			}
			cache.SaddListDataPackList(tmpId, args)
			cache.SaddListAllDataPackList(args)
			cache.SaddListDataPackList2(tmpId, args)
			wh := bson.M{"_id": bson.ObjectIdHex(tmpId)}
			up := bson.M{}
			intoNum := cache.ScardDataPackListCount2(tmpId)

			up["into_num"] = intoNum
			up["up_status"] = int64(2)
			up["invalid_num"] = invalidNum + tmp.InvalidNum
			up["repeat_num"] = repeatNum + tmp.RepeatNum
			dataPack.UpDataPack(wh, up)
			inNum := intoNum - tmp.IntoNum
			cache.SetSchedule(tmpId, utils.IntToStr(inNum)+"_"+utils.IntToStr(int64(len(phoneList))-inNum))
		}(phoneArr, tmp.Id.Hex(), uid)
		rsp.Id = req.Id
	}
	return nil
}

// 获取上传结果
func (this *DataPackServer) GetSchedule(req *info.GetScheduleReq, rsp *info.GetScheduleRsp) *goError.ErrRsp {
	packInfo := dataPack.GetByIdDataPack(req.Id)
	if packInfo.UpStatus == 2 {
		srt := cache.GetSchedule(req.Id)
		split := strings.Split(srt, "_")
		rsp.UpStatus = 2
		if len(split) > 0 {
			rsp.Success = utils.StrToInt64(split[0])
			rsp.Fail = utils.StrToInt64(split[1])
		}
	}
	return nil
}

// 获取剩余数量
func (this *DataPackServer) GetResidueNum(req *info.GetResidueNumReq, rsp *info.GetResidueNumRsp) *goError.ErrRsp {
	start := (req.Page - 1) * req.Limit
	smembers := cache.SmembersDataPackList(req.Id)
	if len(smembers) != 0 {
		selectedKeys := []string{}
		if start+req.Limit > int64(len(smembers)) {
			selectedKeys = smembers[start:]
		} else {
			selectedKeys = smembers[start : start+req.Limit]
		}
		for _, selectedKey := range selectedKeys {
			rsp.List = append(rsp.List, selectedKey)
		}
	}
	rsp.Total = int64(len(smembers))
	return nil
}

// 批量删除
func (this *DataPackServer) BathDel(req *info.BathDelReq, rsp *info.NullRsp) *goError.ErrRsp {
	uid := this.getUid()
	for _, id := range req.Ids {
		config := cache.GetTaskConfig("global")
		if config.DataPackId == id {
			continue
		}
		err := dataPack.DelDataPack(bson.M{"_id": bson.ObjectIdHex(id)})
		if err == nil {
			go func(uid, id string) {
				smembers := cache.SmembersDataPackList2(id)
				var args []interface{}
				for _, smember := range smembers {
					args = append(args, smember)
				}
				cache.SremListAllDataPackList(args)
				cache.DelDataPackList(id)
				cache.DelDataPackList2(id)
				cache.DelDataPackListErr(id)
			}(uid, id)
		}
	}
	return nil
}

// 批量导出
func (this *DataPackServer) DoOutPutData(req *info.DoOutPutDataReq, rsp *info.DoOutPutDataRsp) *goError.ErrRsp {
	exportStr := ""
	phoneList := []string{}
	if req.Type == 1 {
		//导出全部数据
		phoneList = cache.SmembersDataPackList2(req.Id)
	} else if req.Type == 2 {
		//导出剩余数据
		phoneList = cache.SmembersDataPackList(req.Id)
	} else if req.Type == 3 {
		//导出异常数据
		phoneList = cache.SmembersDataPackListErr(req.Id)
		go func(phoneList []string) {
			var args []interface{}
			for _, smember := range phoneList {
				args = append(args, smember)
			}
			cache.SremListAllDataPackList(args)
		}(phoneList)
	}
	array := splitArray(phoneList, 10000)
	wg := sync.WaitGroup{}
	cStrList := []string{}
	for _, arr := range array {
		wg.Add(1)
		go func(arr []string) {
			defer wg.Done()
			cStr := ""
			for _, phone := range arr {
				cStr = cStr + phone + "\n"
			}
			cStrList = append(cStrList, cStr)
		}(arr)
	}
	wg.Wait()
	for _, ctr := range cStrList {
		exportStr = exportStr + ctr
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

func splitArray(arr []string, size int) [][]string {
	var res [][]string
	for i, j := 0, size; i < len(arr); i, j = i+size, j+size {
		if j > len(arr) {
			j = len(arr)
		}
		res = append(res, arr[i:j])
	}
	return res
}
