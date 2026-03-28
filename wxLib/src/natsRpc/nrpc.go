package natsRpc

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/go-nats"
	"mframe"
	"time"
)

/**
发送，不等回包
*/
func NrpcSend(mod string, svrid int32, cmd string, req *NatsMsg) (err error) {
	subj := GenReqSubject(mod, cmd, svrid)
	return NatsPublish(g_natsConn, subj, req, false, nil)
}
func NrpcCallNoReply(mod string, svrid int32, cmd string, req interface{}) (err error) {

	_, err = NrpcCall(mod, svrid, cmd, req, 0)
	return err
}

//远程调用并且回包
func NrpcCall(mod string, svrid int32, cmd string, req interface{} /*rsp interface{},*/, args ...int32) (*NatsMsg, error) {

	rpcReq := &NatsMsg{}
	rpcReq.Sess.SvrFE = GetServerName()
	rpcReq.Sess.SvrID = GetServerId()
	rpcReq.Sess.Mod = mod
	rpcReq.Sess.Cmd = cmd
	rpcReq.Sess.Time = time.Now().Unix()
	//logs.Info("rpcReq.Sess.SvrFE======>",rpcReq.Sess.SvrFE)

	rpcReq.Marshal(req)

	timeout := time.Duration(mframe.GetCallTimeout()) * time.Second

	if len(args) > 0 {
		xtime := args[0]
		timeout = time.Duration(xtime) * time.Second
	}
	if timeout <= 0 { //不等回包
		err := NatsSendMsg(g_natsConn, mod, svrid, cmd, rpcReq, nil)
		return nil, err
	}

	msg, err := NatsCallMsg(g_natsConn, mod, svrid, cmd, rpcReq, timeout, nil)
	if err != nil {
		return nil, err
	}
	rpcRsp := &NatsMsg{}
	err = jsoniter.Unmarshal(msg.Data, rpcRsp)
	//logs.Debug("读取到的数据 %+v",string(msg.Data))
	if err != nil {
		logs.Error("NatsCall %s.%s Unmarshal rsp %s failed:%s", mod, cmd, string(msg.Data), err.Error())
		return nil, err
	}
	return rpcRsp, nil
}

//回包,由于c++过来的请求不会带有replay,因此这里有两种发送模式
func NrpcReply(req *NatsMsg, rpcRsp *NatsMsg) (err error) {

	rpcRsp.Sess.SvrFE = GetServerName()
	rpcRsp.Sess.SvrID = GetServerId()
	rpcRsp.Sess.Mod = req.GetMod()
	rpcRsp.Sess.Cmd = req.GetCmd()
	rpcRsp.Sess.Time = time.Now().Unix()

	//req.Sess.Route = fmt.Sprint("%s->%s.%d", req.Sess.Route, GetServerName(), GetServerID())

	subj := req.GetReply()
	if subj == "" {
		return fmt.Errorf("No Reply")
	}

	conn := req.GetConn()
	if conn == nil {
		return fmt.Errorf("No Conn")
	}
	return NatsPublish(conn, subj, rpcRsp, true, nil)
}

/**
监听事件
*/
func NrpcHandlerEvent(server string, cmd string, handler func(msg *NatsMsg) int32) error {
	subj := fmt.Sprintf("event.%s.%s", server, cmd)
	err := DoRegistNatsHandler(g_natsConn, subj, GetServerName(), func(msg *nats.Msg) int32 {
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
	}, 0, nil)
	return err
}

/**
发布事件
*/
func NatsPulishEvent(ename string, eobj interface{}) error {
	subj := fmt.Sprintf("event.%s.%s", GetServerName(), ename)
	sess := Session{SvrFE: GetServerName(), SvrID: GetServerId(), Cmd: ename}
	msg := NatsMsg{Sess: sess}
	switch eobj.(type) {
	case []byte, jsoniter.RawMessage, string:
		msg.MsgData = eobj.([]byte)
	default:
		bdata, err := jsoniter.Marshal(eobj)
		if err != nil {
			logs.Error("Failed to jsoniter.Marshal: %+v", eobj)
			return err
		}
		msg.MsgData = bdata
	}
	NatsPublish(g_natsConn, subj, msg, false, nil)
	return nil
}

// config

//监听配置变更通知
func ListenConfig(config_key string, _func func([]byte)) {
	subj := "config." + config_key
	queue := fmt.Sprintf("%s-%d", GetServerName(), GetServerId())

	err := DoRegistNatsHandler(g_natsConn, subj, queue, func(msg *nats.Msg) int32 {
		logs.Trace("Recv_Config %s to reload config", msg.Subject)
		_func(msg.Data)
		return 0
	}, 1, nil)

	logs.Trace("ListenConfig key: %s @ %s %+v", subj, queue, err)

	subj = "config." + config_key + "." + GetServerName()
	err = DoRegistNatsHandler(g_natsConn, subj, queue, func(msg *nats.Msg) int32 {
		logs.Trace("Recv_Config %s to reload config", msg.Subject)
		_func(msg.Data)
		return 0
	}, 1, nil)
	logs.Trace("ListenConfig key: %s @ %s %+v", subj, queue, err)

	subj = "config." + config_key + "." + queue
	err = DoRegistNatsHandler(g_natsConn, subj, queue, func(msg *nats.Msg) int32 {
		logs.Trace("Recv_Config %s to reload config", msg.Subject)
		_func(msg.Data)
		return 0
	}, 1, nil)
	logs.Trace("ListenConfig key: %s @ %s %+v", subj, queue, err)
}

//NotifyReloadConfig("online", "push")
func NotifyReloadConfig(cfg string, obj interface{}) {
	subj := "config"
	subj = subj + "." + cfg
	logs.Trace("NotifyReloadConfig subj:%s,obj:%+v", subj, obj)
	NatsPublish(g_natsConn, subj, obj, false, nil)
}
