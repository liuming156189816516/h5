package natsRpc

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/go-nats"
	"mframe"
	"mlog"
	"strings"
	"time"
	"errors"
	"fmt"
)

const (
	Ok   = "ok"
	Fail = "Fail"
)

// 远程调用协议并且回包
func NrpcDllCall(mod string, svrid int32, cmd string, req string /*rsp interface{},*/, args ...int32) (*nats.Msg, error) {
	timeout := time.Duration(mframe.GetCallTimeout()) * time.Second
	if len(args) > 0 {
		xtime := args[0]
		timeout = time.Duration(xtime) * time.Second
	}
	if timeout <= 0 { //不等回包
		err := NatsSendMsgWs(g_natsConn, mod, svrid, cmd, req, nil)
		return nil, err
	}
	msg, err := NatsCallMsgWs(g_natsConn, mod, svrid, cmd, req, timeout, nil)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// 监听协议
func HandlerDllNrpc(cmd string, handler func([]*ReceiveMessagesReq) int32) error {
	subj := GenReqSubject(GetServerName(), cmd, 100)
	if err := registNatsDllHandler(cmd, subj, handler, 0); err != nil {
		return err
	}
	return nil
}

// 注册handler
func registNatsDllHandler(func_name string, subj string, handler func([]*ReceiveMessagesReq) int32, tNum int) error {
	err := DoRegistNatsHandler(g_natsConn, subj, GetServerName(), func(msg *nats.Msg) int32 {
		req := []*ReceiveMessagesReq{}
		err := jsoniter.Unmarshal(msg.Data, &req)
		if err != nil {
			if msg.Reply != "" {
				NatsPublishDll(g_natsConn, msg.Reply, Fail, false, nil)
				return ESMR_FAILED
			}
			return ESMR_FAILED
		}
		ret := handler(req)
		if ret != ESMR_SUCCEED {
			NatsPublishDll(g_natsConn, msg.Reply, Fail, false, nil)
		} else {
			NatsPublishDll(g_natsConn, msg.Reply, Ok, false, nil)
		}
		return ret
	}, tNum, nil)
	return err
}

// 发送消息到nats
func NatsPublishDll(natsConn *nats.Conn, subj string, req string, isReply bool, CheckNatsService func(string) bool) (err error) {
	if !isReply && CheckNatsService != nil {
		if ok := CheckNatsService(subj); ok == false {
			mlog.Debug("not find service for %s", subj)
			return RpcNoService
		}
	}
	data := []byte(req)
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

type ReceiveMessagesReq struct {
	Account string      `json:"account"`
	Time    int64       `json:"time"`
	Type    string      `json:"type"`
	Ctype   string      `json:"ctype"`
	Content interface{} `json:"content"`
}

type ReceiveMessagesContent struct {
	Group string `json:"group"`
	To    string `json:"to"`
}

type AccountLoginContent struct {
	Id     int64  `json:"id"`
	Errmsg string `json:"errmsg"`
	Errno  int64  `json:"errno"`
}

type ChatContent struct {
	Id    string `json:"id"`
	Group string `json:"group"`
	From  string `json:"from"`
	To    string `json:"to"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Key   string `json:"key"`
	Url   string `json:"url"`
}

type NotifyContent struct {
	Group  string   `json:"group"`
	From   string   `json:"from"`
	Read   bool     `json:"read"`
	Msgids []string `json:"msgids"`
}

type MessageSendContent struct {
	Errno  int64  `json:"errno"`
	Data   string `json:"data"`
	Errmsg string `json:"errmsg"`
	To     string `json:"to"`
	Id     int64  `json:"id"`
}

func NatsSendMsgWs(natsConn *nats.Conn, mod string, svrid int32, cmd string, req string, CheckNatsService func(string) bool) error {
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
	//data, err := jsoniter.Marshal(req)
	//if err != nil {
	//	mlog.Error("NatsCall %s Marshal req (%+v) failed:%s", subj, req, err.Error())
	//	return err
	//}
	start := time.Now()
	err := natsConn.Publish(subj, []byte(req))
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

func NatsCallMsgWs(natsConn *nats.Conn, mod string, svrid int32, cmd string, req string, timeout time.Duration, CheckNatsService func(string) bool) (rsp *nats.Msg, err error) {

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
	//data, err := jsoniter.Marshal(req)
	//if err != nil {
	//	mlog.Error("NatsCall %s Marshal req (%+v) failed:%s", subj, req, err.Error())
	//	return nil, err
	//}

	start := time.Now()
	msg, err := natsConn.Request(subj, []byte(req), timeout)

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
		mlog.Debug("natsConn.Request %s (cost:%v)\n: %s , Response  %s", subj, t, req, string(msg.Data))
	}
	ReportCallRpcStat(mod, cmd, ret, t)
	return msg, err
}