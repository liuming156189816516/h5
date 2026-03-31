package test

import (
	"comm/comm"
	"comm/goError"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	info "script/webstru"
	"selfComm/db/sendmsg"
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
	msgInfo := sendmsg.GetCountSendMsgInfo(bson.M{})
	fmt.Println(msgInfo)
	return nil
}
