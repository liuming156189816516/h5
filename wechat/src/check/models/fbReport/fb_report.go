package fbReport

import (
	"comm/comm"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/fb"
	"selfComm/db/log"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"

	"time"
)

func TaskRun() {
	go func() {
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(10 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				val := GetStaticData("fbReport")
				if val != "" {
					//正在执行
					continue
				}
				go fbReport(t)
			}
		}
	}()
}

func fbReport(t time.Time) {
	defer DelStaticData("fbReport")
	SetStaticData("fbReport", "fbReport")

	for i := 0; i < 30; i++ {
		// 从缓存中取出一条数据
		report := cache.GetFbReport()
		if report.Ptype == 0 {
			break
		}
		toString, _ := jsoniter.MarshalToString(report)
		fmt.Println("fb任务上报=======>", toString)
		var eventName string
		switch report.Ptype {
		case 1:
			eventName = "ViewContent" // 内容查看
		case 2:
			eventName = "SubmitVerification" // 获取验证码
		case 3:
			eventName = "SessionSuccess" // 挂机成功
		default:
			continue
		}

		if report.Ptype == 1 {
			wxComm.FbReport("", eventName, report.Fbclid, report.Fbp, report.PixelId, "")
		}

		if report.Ptype == 2 {
			tmp := &fb.Fb{}
			tmp.Id = bson.NewObjectId()
			tmp.Phone = report.Phone
			tmp.Fbclid = report.Fbclid
			tmp.Fbp = report.Fbp
			tmp.PixelId = report.PixelId
			err := fb.AddFb(tmp)
			if err == nil {
				eventID := comm.Md5(report.Phone + report.PixelId)
				wxComm.FbReport(eventID, eventName, report.Fbclid, report.Fbp, report.PixelId, "")

				tmp1 := &log.FbLog{}
				tmp1.Id = bson.NewObjectId()
				tmp1.EventID = eventID
				tmp1.EventName = eventName
				tmp1.Phone = report.Phone
				tmp1.Fbclid = report.Fbclid
				tmp1.Fbp = report.Fbp
				tmp1.PixelId = report.PixelId
				log.AddFbLog(tmp1)

			}
		}

		if report.Ptype == 3 {
			oneFb := fb.GetOneFb(bson.M{"phone": report.Phone})
			if oneFb.Id.Hex() != "" {
				eventID := comm.Md5(oneFb.Phone + oneFb.PixelId)
				wxComm.FbReport(eventID, eventName, oneFb.Fbclid, oneFb.Fbp, oneFb.PixelId, oneFb.Id.Hex())
			}
		}
	}
}
