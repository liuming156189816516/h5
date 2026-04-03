package qrcode

import (
	"comm/event"
	"comm/redisDeal"
	"comm/redisKeys"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/account"
	"selfComm/wxComm/cache"
	"serApi/dllApi"
	"strings"
	"time"
)

func TaskRun() {
	go func() {
		//二维码登录
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(10 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				go checkQrcodeTask(t)
			}
		}
	}()
}

// 二维码登录
func checkQrcodeTask(t time.Time) {
	list := redisDeal.RedisSendPattern(redisKeys.GetCheckQrcodeTaskKey("*"))
	for _, key := range list {
		val := GetStaticData(key)
		if val != "" {
			//正在执行
			continue
		}
		SetStaticData(key, key)
		go checkQrcode(key)
	}
}

// 检测二维码
func checkQrcode(key string) {
	defer DelStaticData(key)
	uuid := strings.ReplaceAll(key, redisKeys.GetCheckQrcodeTaskKey(""), "")

	dreq := &dllApi.VfcodeCheckReq{
		Id: uuid,
	}
	dRsp, _ := dllApi.VfcodeCheck(dreq, -1, true, 10)
	if dRsp == nil {
		return
	}
	if dRsp.Code == -1001 {
		return
	}
	if dRsp.Code == 1000 {
		return
	}
	if dRsp.Code == 0 {
		checkData := &QrcodeCheckData{}
		cStr, _ := jsoniter.MarshalToString(dRsp.Data)
		jsoniter.UnmarshalFromString(cStr, &checkData)
		if checkData.Token != "" {
			proxyStr := cache.GetCheckQrcodeTask(uuid)
			data := CheckQrcodeTaskData{}
			jsoniter.UnmarshalFromString(proxyStr, &data)

			split := strings.Split(checkData.Token, ",")
			tmp := &account.AccountInfo{}
			tmp.Id = bson.NewObjectId()
			tmp.Account = strings.ReplaceAll(split[0], "_", "")
			tmp.GroupId = data.DefaultGid
			tmp.Status = 3
			tmp.AccountType = data.AccountType
			tmp.Token = checkData.Token
			tmp.Synckeys = checkData.Synckeys
			tmp.Itime = time.Now().Unix()
			tmp.Ptime = time.Now().Unix()
			tmp.NickName = checkData.Nickname
			tmp.PixelId = data.PixelId
			tmp.PlatformType = 2
			tmp.AreaCode = data.AreaCode
			tmp.IsProxyUser = data.IsProxyUser
			err := account.AddAccountInfo(tmp)
			if err != nil && strings.Contains(err.Error(), "E11000 duplicate key") {
				//更新为登录中
				account.UpAccountInfo(bson.M{"account": tmp.Account}, bson.M{"status": int64(3)})
			}

			cache.SetAllAccountList(tmp.Account)
			ipTmp := cache.ProxyIpInfo{}
			ipTmp.IpId = data.ProxyId
			ipTmp.User = data.User
			ipTmp.Pwd = data.Pwd
			ipTmp.Port = data.Port
			ipTmp.Host = data.Host
			ipTmp.Type = data.Type
			cache.SetProxyIp(tmp.Account, &ipTmp)
			cache.IncIpUserNum(ipTmp.IpId, 1)
			cache.SetAccountStatus(tmp.Account, tmp.Status)
			accInfo := &cache.AccountInfo{
				Account:      tmp.Account,
				AccountType:  tmp.AccountType,
				Token:        tmp.Token,
				PlatformType: tmp.PlatformType,
				PixelId:      data.PixelId,
				ClickId:      data.ClickId,
				Synckeys:     tmp.Synckeys,
			}
			cache.SetAccountInfo(tmp.Account, accInfo)
			//登录
			e := &event.TaskUserLoginEventReq{
				Account: []string{tmp.Account},
			}
			event.AddLoginTask(event.TaskTypeUserLogin, tmp.Account, e)
		}
	}
	cache.DelCheckQrcodeTask(uuid)
}
