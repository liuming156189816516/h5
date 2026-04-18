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
	//log.DelFbReportLog(bson.M{})
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

	/*tmpProxy := &cache.AccountSocks5Info{}
	lockIp := ip.GetOneLockIp()
	if lockIp.ProxyIp == "" {
		return goError.IpOperationErr
	}
	split := strings.Split(lockIp.ProxyIp, ":")
	tmpProxy.User = lockIp.ProxyUser
	tmpProxy.Pwd = lockIp.ProxyPwd
	tmpProxy.Type = lockIp.ProxyType
	tmpProxy.Host = split[0]
	tmpProxy.Port = split[1]
	proxy := dllApi.AccountAddParamSocks5{}
	proxy.User = tmpProxy.User
	proxy.Pwd = tmpProxy.Pwd
	proxy.Host = tmpProxy.Host
	proxy.Port = tmpProxy.Port
	proxy.Type = tmpProxy.Type
	uuid := "355692051682"
	dreq := &dllApi.QrcodeCreateReq{
		Id:          uuid,
		Proxy:       proxy,
		AccountType: 1,
	}
	dRsp, err := dllApi.QrcodeCreate(dreq, -1, true, 30)
	fmt.Println(err)
	fmt.Println(jsoniter.MarshalToString(dRsp))*/
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
