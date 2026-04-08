package test

import (
	"comm/comm"
	"comm/goError"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	info "script/webstru"
	"selfComm/db/fb"
	"selfComm/db/log"
)

// 群发
type DemoServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *DemoServer) getUid() string {
	return this.Sess.Uid
}

func (this *DemoServer) Demo(req *info.DemoReq, rsp *info.DemoRsp) *goError.ErrRsp {
	//查询上报表，根据手机号去重
	if req.Type == "1" {
		listFb := fb.GetListFb(bson.M{}, -1)
		fmt.Println("执行完成1=========================>", len(listFb))
		printed := make(map[string]bool)
		cont := int64(0)
		for _, fb := range listFb {
			if fb.Fbclid != "" && !printed[fb.Phone] {
				printed[fb.Phone] = true
				toString, _ := jsoniter.MarshalToString(fb)
				fmt.Println("执行完成1=========================>", toString)
				cont++
			}
		}
		fmt.Println("执行完成1=========================>", cont)
	}
	//小莫上报日志查询
	if req.Type == "2" {
		listFb := log.GetListFbReportLog(bson.M{}, -1)
		fmt.Println("执行完成2=========================>", len(listFb))
		cont := int64(0)
		for _, fb := range listFb {
			toString, _ := jsoniter.MarshalToString(fb)
			fmt.Println("执行完成2=========================>", toString)
			cont++
		}
		fmt.Println("执行完成2=========================>", cont)
	}

	//查询给fb上报日志
	if req.Type == "3" {
		listFb := log.GetListFbLog(bson.M{}, -1)
		fmt.Println("执行完成3=========================>", len(listFb))
		for _, fb := range listFb {
			toString, _ := jsoniter.MarshalToString(fb)
			fmt.Println("执行完成3=========================>", toString)
		}
	}

	return nil
}
