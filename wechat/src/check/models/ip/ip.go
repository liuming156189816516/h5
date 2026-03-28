package ip

import (
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/account"
	"selfComm/db/ip"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"time"
)

func TaskRun() {
	go func() {
		//校验ip
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(1 * time.Hour)
		for {
			select {
			case t := <-ticker.C:
				go checkIpTask(t)
			}
		}
	}()
}

// 校验ip
func checkIpTask(t time.Time) {
	go func() {
		listIp := ip.GetListIp(bson.M{}, -1)
		for _, ipInfo := range listIp {
			go CheckIp(ipInfo.Id.Hex())
			useNum := account.GetCountAccountInfo(bson.M{"proxy_ip": ipInfo.Id.Hex()})
			cache.SetIpUserNum(ipInfo.Id.Hex(), useNum)
		}
	}()
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
