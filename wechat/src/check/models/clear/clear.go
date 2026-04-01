package clear

import (
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/sendmsg"
	"selfComm/wxComm/cache"
	"time"
	"utils"
)

func TaskRun() {
	go func() {
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(60 * time.Minute)
		for {
			select {
			case t := <-ticker.C:
				go clearCache(t)
			}
		}
	}()
}

//清除发送两天前的发送数据的缓存
func clearCache(t time.Time) {
	begin := utils.GetTimeBegin(t.Unix()) + 60*60*2 + 30*60
	end := utils.GetTimeBegin(t.Unix()) + 60*60*4
	if t.Unix() < begin || t.Unix() > end {
		return
	}
	infoList := sendmsg.GetListSendMsgInfo(bson.M{"account_status": 1, "ptime": bson.M{"$lt": t.Unix() - 60*60*24*2}}, -1)
	for _, info := range infoList {
		cache.DelSendMsgRecordInfo(info.Account)
		cache.DelSendMsgPhoneLid(info.Account)
	}
}
