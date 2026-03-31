package test

import (
	"comm/comm"
	"comm/goError"
	"gopkg.in/mgo.v2/bson"
	info "script/webstru"
	"selfComm/db/ip"
)

// 群发
type DemoServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *DemoServer) getUid() string {
	return this.Sess.Uid
}

func (this *DemoServer) Demo(req1 *info.DemoReq, rsp *info.DemoRsp) *goError.ErrRsp {
	tmpGroup := &ip.IpGroup{}
	tmpGroup.Id = bson.NewObjectId()
	tmpGroup.Name = "测试"
	ip.AddIpGroup(tmpGroup)

	tmpIp := &ip.Ip{}
	tmpIp.Id = bson.NewObjectId()
	tmpIp.GroupId = tmpGroup.Id.Hex()
	tmpIp.ProxyIp = "173.211.69.97:6690"
	tmpIp.ProxyUser = "scoohztt"
	tmpIp.ProxyPwd = "epbs3o70v4iv"
	tmpIp.ProxyType = "socks5"
	tmpIp.Status = 1
	tmpIp.AllotNum = 5
	tmpIp.IpType = 1
	tmpIp.IpCategory = 1
	tmpIp.ExpireStatus = 1
	tmpIp.Country = "美国"
	tmpIp.DisableStatus = 1
	ip.AddIp(tmpIp)
	return nil
}
