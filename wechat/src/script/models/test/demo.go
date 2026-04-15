package test

import (
	"comm/comm"
	"comm/goError"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	info "script/webstru"
	"selfComm/db/log"
	"selfComm/wxComm/cache"
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
	count1 := 0
	//内容查看
	reportList := log.GetListFbReportLog(bson.M{"ptype": 1}, -1)
	for _, report := range reportList {
		tmp := &info.FbData{}
		jsoniter.UnmarshalFromString(report.Data, &tmp)
		if tmp.ClickId != "" {
			count1++
		}
	}
	fmt.Println("内容查看：", count1)
	count2 := 0

	qrCodeMap := make(map[string]string)
	//获取验证码
	reportList1 := log.GetListFbReportLog(bson.M{"ptype": 2}, -1)
	for _, report := range reportList1 {
		tmp := &info.QrCode{}
		jsoniter.UnmarshalFromString(report.Data, &tmp)
		if tmp.ClickId != "" {
			qrCodeMap[tmp.AreaCode+tmp.Account] = "1"
			count2++
		}
	}

	fmt.Println("验证码", count2)
	fmt.Println("验证码去重复", len(qrCodeMap))
	return nil
}

func (this *DemoServer) CheckStatus(req *info.DemoReq, rsp *info.DemoRsp) *goError.ErrRsp {
	if req.Status == "0" || req.Status == "1" {
		cache.SetTaskStatus(req.Status)
	} else {
		rsp.Message = "参数有误"
	}
	if req.Status == "0" {
		rsp.Message = "定时任务已开启"
	}
	if req.Status == "1" {
		rsp.Message = "定时任务已关闭"
	}
	return nil
}
