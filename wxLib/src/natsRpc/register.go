package natsRpc

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/go-nats"
	"time"
)

//监听内部服务的请求, 完整协议解析
func HandlerNrpc(cmd string, handler func(*NatsMsg) int32) error {
	subj := GenReqSubject(GetServerName(), cmd, -1)
	if err := registNatsHandler(cmd, subj, handler, 0); err != nil {
		return err
	}
	return nil
}

//监听内部服务的请求, 完整协议解析
func HandlerNrpcNum(cmd string, handler func(*NatsMsg) int32, num int64) error {
	subj := GenReqSubject(GetServerName(), cmd, -1)
	if err := registNatsHandler(cmd, subj, handler, int(num)); err != nil {
		return err
	}
	return nil
}

//监听内部服务的请求, 完整协议解析
func HandlerNrpcByName(serverName, cmd string, handler func(*NatsMsg) int32) error {
	subj := GenReqSubject(serverName, cmd, -1)
	if err := registNatsHandler(cmd, subj, handler, 0); err != nil {
		return err
	}
	return nil
}

//监听内部服务的请求, 完整协议解析
func HandlerNrpcTmp(cmd string, handler func(*NatsMsg) int32) (*nats.Subscription, error) {
	subj := GenReqSubject(GetServerName(), cmd, -1)
	return registNatsHandlerTmp(cmd, subj, handler, 0)
}

//监听内部服务的请求, 完整协议解析
func HandlerNrpcCmd(cmd string, handler func(*NatsMsg) int32, needP2p bool) error {
	subj := GenReqSubject(GetServerName(), cmd, -1)
	if err := registNatsHandler(cmd, subj, handler, 0); err != nil {
		return err
	}
	if needP2p {
		subj = GenReqSubject(GetServerName(), cmd, GetServerId())
		if err := registNatsHandler(cmd, subj, handler, 0); err != nil {
			return err
		}
	}
	return nil
}

func HandlerNrpcEvent(server string, cmd string, handler func(msg *NatsMsg) int32, task string, tnums ...int) error {
	// TODO Add options logic
	b := 0
	if len(tnums) > 0 {
		b = tnums[0]
	}
	subj := fmt.Sprintf("event.%s.%s", server, cmd)
	if task == "" {
		task = GetServerName()
	} else {
		task = GetServerName() + "." + task
	}
	err := DoRegistNatsHandler(g_natsConn, subj, task, func(msg *nats.Msg) int32 {
		start := time.Now()
		req := NatsMsg{}
		err := jsoniter.Unmarshal(msg.Data, &req)
		if err != nil {
			logs.Error("Handle Event Unmarshal req %s failed:%s", string(msg.Data[:]), err.Error())
			t := time.Now().Sub(start)
			ReportDoRpcStat("Unmarshal", subj, ESMR_FAILED, t)
			return ESMR_FAILED
		}
		req.NatsMsg = msg
		ret := handler(&req)
		t := time.Now().Sub(start)
		s := fmt.Sprintf("%s.%d", req.GetSession().SvrFE, req.GetSession().SvrID)
		ReportDoRpcStat(s, subj, ret, t)
		return ret
	}, b, nil)

	return err
}
