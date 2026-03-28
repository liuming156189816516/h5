package natsRpc

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/go-nats"
	"mlog"
	"strings"
	"time"
)

func GenReqSubject(mod string, cmd string, svrid int32) string {
	if svrid == 100 {
		return mod
	}
	if svrid <= 0 {
		return fmt.Sprintf("nrpc.%s.%s", mod, cmd)
	} else {
		return fmt.Sprintf("nrpc.%s.%d.%s", mod, svrid, cmd)
	}
}

type serviceSubscription struct {
	subj      string
	queue     string
	ch        chan *nats.Msg
	handler   func(*nats.Msg) int32
	threadNum int
	oSubs     []*nats.Subscription
}

func DoRegistNatsHandler(natsConn *nats.Conn, subj string, queue string, handler func(*nats.Msg) int32, threadNum int,
	onReg func(sSub *serviceSubscription)) error {

	// TODO Add options logic
	if natsConn == nil {
		mlog.Error("no natsConn")
		return errors.New("no natsConn")
	}
	//subj := fmt.Sprintf("rpc.%s.%s", server_config.ServiceName, cmd)
	//queue := server_config.ServerName
	var oSubjs []*nats.Subscription = nil
	//var err error = nil

	if threadNum <= 0 {
		oSubj, err := natsConn.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
			//起一个协程去做这个事
			if strings.Index(subj, "HeartBeat") <= 0 {
				mlog.Debug("natsConn.Recive At: %s GetMsg:%s[%s]", subj, queue, string(msg.Data))
			}
			go handler(msg)

		})
		if err != nil {
			mlog.Error("HandleRpc QueueSubscribe %s queue:%s failed:%s", subj, queue, err.Error())
			return err
		}
		mlog.Trace("natsConn.QueueSubscribe %s queue:%s, threadNum: %d succ", subj, queue, threadNum)

		oSubjs = append(oSubjs, oSubj)
	} else {
		for i := 0; i < threadNum; i++ {
			oSubj, err := natsConn.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
				if strings.Index(subj, "HeartBeat") <= 0 {
					mlog.Debug("natsConn.Recive At: %s - GetMsg:%s[%s]", subj, queue, string(msg.Data))
				}
				handler(msg)
				//go handler(msg)
			})
			if err != nil {
				mlog.Error("HandleRpc  QueueSubscribe %s,Index:%d,threadNum:%d queue:%s failed:%s", subj, i, threadNum, queue, err.Error())
				return err
			}
			mlog.Trace("natsConn.QueueSubscribe %s queue:%s,Index:%d,threadNum:%d succ", subj, queue, i, threadNum)
			oSubjs = append(oSubjs, oSubj)
		}

	}

	natsConn.Flush()
	if onReg != nil {
		sSub := &serviceSubscription{subj: subj, queue: queue, handler: handler, threadNum: threadNum, oSubs: oSubjs}
		onReg(sSub)
	}

	return nil
}

func DoRegistNatsHandlerTmp(natsConn *nats.Conn, subj string, queue string, handler func(*nats.Msg) int32, threadNum int) (*nats.Subscription, error) {

	// TODO Add options logic
	if natsConn == nil {
		mlog.Error("no natsConn")
		return nil, errors.New("no natsConn")
	}
	//subj := fmt.Sprintf("rpc.%s.%s", server_config.ServiceName, cmd)
	//queue := server_config.ServerName
	var oSubjs *nats.Subscription = nil
	//var err error = nil

	if threadNum <= 0 {
		oSubj, err := natsConn.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
			//起一个协程去做这个事
			if strings.Index(subj, "HeartBeat") <= 0 {
				mlog.Debug("natsConn.Recive At: %s GetMsg:%s[%s]", subj, queue, string(msg.Data))
			}
			go handler(msg)

		})
		if err != nil {
			mlog.Error("HandleRpc QueueSubscribe %s queue:%s failed:%s", subj, queue, err.Error())
			return nil, err
		}
		mlog.Trace("natsConn.QueueSubscribe %s queue:%s, threadNum: %d succ", subj, queue, threadNum)

		oSubjs = oSubj
	} else {
		for i := 0; i < threadNum; i++ {
			oSubj, err := natsConn.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
				if strings.Index(subj, "HeartBeat") <= 0 {
					mlog.Debug("natsConn.Recive At: %s - GetMsg:%s[%s]", subj, queue, string(msg.Data))
				}
				handler(msg)
				//go handler(msg)
			})
			if err != nil {
				mlog.Error("HandleRpc  QueueSubscribe %s,Index:%d,threadNum:%d queue:%s failed:%s", subj, i, threadNum, queue, err.Error())
				return nil, err
			}
			mlog.Trace("natsConn.QueueSubscribe %s queue:%s,Index:%d,threadNum:%d succ", subj, queue, i, threadNum)
			oSubjs = oSubj
		}

	}
	natsConn.Flush()
	return oSubjs, nil
}

//发送消息到nats
func NatsPublish(natsConn *nats.Conn, subj string, req interface{}, isReply bool, CheckNatsService func(string) bool) (err error) {

	if !isReply && CheckNatsService != nil {
		if ok := CheckNatsService(subj); ok == false {
			mlog.Debug("not find service for %s", subj)
			return RpcNoService
		}
	}
	data, err := jsoniter.Marshal(req)
	if err != nil {
		mlog.Error("NatsCall Marshal req (%+v) failed:%s", req, err.Error())
		return err
	}
	start := time.Now()
	err = natsConn.Publish(subj, data)
	t := time.Now().Sub(start)
	ret := 0
	if err != nil {
		ret = 1
		mlog.Error("natsConn.Publish[%s] (%s) failed：%s", subj, string(data), err.Error())
		ReportStat("Send", subj, int(ret), t)
		return err
	}
	if strings.Index(subj, "_INBOX.") != 0 {
		ReportStat("Send", subj, int(ret), t)
	}

	if strings.Index(subj, "HeartBeat") <= 0 {
		mlog.Debug("natsConn.Publish[%s] (%+v) ", subj, string(data))
	}
	return nil
}
func NatsSendMsg(natsConn *nats.Conn, mod string, svrid int32, cmd string, req interface{}, CheckNatsService func(string) bool) error {
	if natsConn == nil {
		mlog.Error("no conn")
		return errors.New("no natsConn")
	}
	subj := GenReqSubject(mod, cmd, svrid)
	if CheckNatsService != nil {
		if ok := CheckNatsService(subj); ok == false {
			return RpcNoService
		}
	}
	data, err := jsoniter.Marshal(req)
	if err != nil {
		mlog.Error("NatsCall %s Marshal req (%+v) failed:%s", subj, req, err.Error())
		return err
	}
	start := time.Now()
	err = natsConn.Publish(subj, data)
	t := time.Now().Sub(start)
	ret := ESMR_SUCCEED
	if err != nil {
		if err == nats.ErrTimeout {
			err = errors.New(fmt.Sprintf("%s, %d", err.Error(), t/time.Millisecond))
			ret = ESMR_TIMEOUT
		} else {
			ret = ESMR_FAILED
		}
		mlog.Error("natsConn.Publish %s , failed:%s", subj, err.Error())

	}
	ReportCallRpcStat(mod, cmd, ret, 0)
	return err
}
func NatsCallMsg(natsConn *nats.Conn, mod string, svrid int32, cmd string, req interface{}, timeout time.Duration, CheckNatsService func(string) bool) (rsp *nats.Msg, err error) {

	if natsConn == nil {
		mlog.Error("no conn")
		return nil, errors.New("no natsConn")
	}
	if timeout <= 0 || timeout > 30*time.Minute {
		timeout = 30 * time.Second
	}
	subj := GenReqSubject(mod, cmd, svrid)
	if CheckNatsService != nil {
		if ok := CheckNatsService(subj); ok == false {
			return nil, RpcNoService
		}
	}
	data, err := jsoniter.Marshal(req)
	if err != nil {
		mlog.Error("NatsCall %s Marshal req (%+v) failed:%s", subj, req, err.Error())
		return nil, err
	}
	start := time.Now()
	msg, err := natsConn.Request(subj, data, timeout)

	t := time.Now().Sub(start)
	ret := ESMR_SUCCEED

	if err != nil {
		if err == nats.ErrTimeout {
			err = errors.New(fmt.Sprintf("%s, %d", err.Error(), t/time.Millisecond))
			ret = ESMR_TIMEOUT
		} else {
			ret = ESMR_FAILED
		}
		mlog.Error("natsConn.Request %s , timeout:%d failed:%s", subj, timeout/time.Millisecond, err.Error())

	} else {
		mlog.Debug("natsConn.Request %s (cost:%v)\n: %s , Response  %s", subj, t, string(data), string(msg.Data))
	}
	ReportCallRpcStat(mod, cmd, ret, t)
	return msg, err
}

//注册handler
func registNatsHandler(func_name string, subj string, handler func(*NatsMsg) int32, tNum int) error {
	err := DoRegistNatsHandler(g_natsConn, subj, GetServerName(), func(msg *nats.Msg) int32 {
		start := time.Now()
		req := NatsMsg{}
		err := jsoniter.Unmarshal(msg.Data, &req)
		if err != nil {
			if msg.Reply != "" {
				rsp := genErrNatsMsg(-1, "invalid request")
				NatsPublish(g_natsConn, msg.Reply, rsp, false, nil)
			}
			mlog.Error("HandleRpc Unmarshal req %s failed:%s", string(msg.Data[:]), err.Error())
			t := time.Now().Sub(start)
			ReportDoRpcStat("Unmarshal", subj, ESMR_FAILED, t)
			return ESMR_FAILED
		}
		if req.Sess.Cmd == "" {
			req.Sess.Cmd = func_name
		}
		req.NatsMsg = msg
		ret := handler(&req)
		t := time.Now().Sub(start)
		s := fmt.Sprintf("%s.%d", req.GetSession().SvrFE, req.GetSession().SvrID)
		ReportDoRpcStat(s, subj, ret, t)
		return ret
	}, tNum, nil)
	return err
}

//注册handler
func registNatsHandlerTmp(func_name string, subj string, handler func(*NatsMsg) int32, tNum int) (*nats.Subscription, error) {
	sub, err := DoRegistNatsHandlerTmp(g_natsConn, subj, GetServerName(), func(msg *nats.Msg) int32 {
		start := time.Now()
		req := NatsMsg{}
		err := jsoniter.Unmarshal(msg.Data, &req)
		if err != nil {
			if msg.Reply != "" {
				rsp := genErrNatsMsg(-1, "invalid request")
				NatsPublish(g_natsConn, msg.Reply, rsp, false, nil)
			}
			mlog.Error("HandleRpc Unmarshal req %s failed:%s", string(msg.Data[:]), err.Error())
			t := time.Now().Sub(start)
			ReportDoRpcStat("Unmarshal", subj, ESMR_FAILED, t)
			return ESMR_FAILED
		}
		if req.Sess.Cmd == "" {
			req.Sess.Cmd = func_name
		}
		req.NatsMsg = msg
		ret := handler(&req)
		t := time.Now().Sub(start)
		s := fmt.Sprintf("%s.%d", req.GetSession().SvrFE, req.GetSession().SvrID)
		ReportDoRpcStat(s, subj, ret, t)
		return ret
	}, tNum)
	return sub, err
}
