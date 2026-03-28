package servicesApi

import (
	"common/goError"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
)

type ServiceApi struct {
	session *natsRpc.Session
}

func NewServiceApi(s *natsRpc.Session) *ServiceApi {
	return &ServiceApi{session: s}
}

func (this *ServiceApi) GetUserInfo(reqparam *GetUserReq, svrid int32, needRsp bool) (rsppara *GetUserRsp, err error) {
	// SendAndRecv
	req := natsRpc.NatsMsg{Sess: *this.session}
	req.MsgData, err = jsoniter.Marshal(&reqparam)
	if err != nil {
		logs.Debug("testRpc->GetUser failed:%s", err)
		return nil, err
	}
	// 需要回包
	if needRsp {
		rsp, err := natsRpc.NrpcCall("testRpc", svrid, "GetUser", &req)
		if err != nil {
			logs.Debug("testRpc->GetUser failed:%s", err)
			return nil, err
		}
		if rsp.GetMsgErrNo() != 0 {
			return nil, goError.NewGoError(rsp.GetMsgErrNo(), rsp.GetMsgErrStr())
		}
		rsppara = new(GetUserRsp)
		err = jsoniter.Unmarshal(rsp.GetMsgData(), rsppara)
		return rsppara, err
	} else { // 不需要回包
		err := natsRpc.NrpcSend("testRpc", svrid, "GetUser", &req)
		if err != nil {
			logs.Debug("testRpc->GetUser failed:%s", err)
			return nil, err
		}
		return nil, err
	}
}
