package test

import (
	"comm/comm"
	"comm/goError"
	"comm/redisDeal"
	"comm/redisKeys"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"gopkg.in/mgo.v2/bson"
	info "script/webstru"
	"selfComm/db/account"
	"selfComm/db/log"
	"selfComm/db/sendmsg"
	"selfComm/wxComm/cache"
	"serApi/dllApi"
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
	if req.Type == "1" {
		cont := int64(0)
		listSendMsgInfo := sendmsg.GetListSendMsgInfo(bson.M{}, -1)
		if len(listSendMsgInfo) != 0 {
			for _, msgInfo := range listSendMsgInfo {
				cont = cont + msgInfo.SucessNum
			}
			cont1 := cont / int64(len(listSendMsgInfo))

			rsp.Message =
				"发送数据数量：" + cast.ToString(cont) +
					"；发送账号数量：" + cast.ToString(len(listSendMsgInfo)) +
					"；平均发送条数：" + cast.ToString(cont1)
		}
	}

	if req.Type == "2" {
		sendmsg.DelSendMsgInfo(bson.M{})
	}

	if req.Type == "3" {
		count()
	}

	if req.Type == "4" {
		account.DelAccountFile(bson.M{})
		account.DelAccountLog(bson.M{})
		account.DelAccountInfo(bson.M{})
		all := redisDeal.RedisDoHGetAll(redisKeys.GetAllAccountListKey())
		for s, _ := range all {
			cache.DelAllAccountList(s)
			cache.DelAccountStatus(s)
			cache.DelAccountInfo(s)
		}
		sendmsg.DelSendMsgInfo(bson.M{})

	}
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

func count() {
	kwViewMap := make(map[string]string)
	fbViewMap := make(map[string]string)
	//内容查看
	reportList := log.GetListFbReportLog(bson.M{"ptype": 1}, -1)
	for _, report := range reportList {
		tmp := &info.FbData{}
		jsoniter.UnmarshalFromString(report.Data, &tmp)
		if tmp.ClickId != "" {
			kwViewMap[tmp.ClickId] = "1"
		}
		if tmp.Fbclid != "" {
			fbViewMap[tmp.Fbclid] = "1"
		}
	}
	fmt.Println("============================================")
	fmt.Println("kwai内容查看去除重复：", len(kwViewMap))
	fmt.Println("fb内容查看去除重复：", len(fbViewMap))

	kwCodeMap := make(map[string]string)
	fbCodeMap := make(map[string]string)
	//获取验证码
	reportList1 := log.GetListFbReportLog(bson.M{"ptype": 2}, -1)
	for _, report := range reportList1 {
		tmp := &info.FbData{}
		jsoniter.UnmarshalFromString(report.Data, &tmp)
		if tmp.AreaCode != "" && tmp.Account != "" {
			if tmp.ClickId != "" {
				kwCodeMap[tmp.AreaCode+tmp.Account] = "1"
			}
			if tmp.Fbclid != "" {
				fbCodeMap[tmp.AreaCode+tmp.Account] = "1"
			}
		}
	}

	fmt.Println("kwai验证码去重复", len(kwCodeMap))
	fmt.Println("kwai验证码去重复数据", kwCodeMap)
	fmt.Println("fb验证码去重复", len(fbCodeMap))
	fmt.Println("fb验证码去重复数据", fbCodeMap)
}

func testNoProxy(phone string) {
	//获取ip
	uuid := phone
	dreq := &dllApi.VfcodeCreateReq{
		Id:   uuid,
		Code: "77777777",
		//Proxy:       proxy,
		AccountType: 1,
	}
	dRsp, err := dllApi.VfcodeCreate(dreq, -1, true, 30)
	fmt.Println(err)
	fmt.Println(jsoniter.MarshalToString(dRsp))
}
