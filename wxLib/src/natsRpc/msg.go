package natsRpc

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/go-nats"
	"time"
)

var RpcTimeout = errors.New(`{"errno":401, "err":"Timeout"}`)
var RpcNoService = errors.New(`{"errno":404, "err":"No Service"}`)
var RpcNoRequest = errors.New(`{"errno":405, "err":"No Request"}`)

const (
	ESMR_SUCCEED int32 = 0
	ESMR_FAILED  int32 = 1
	ESMR_TIMEOUT int32 = 2
)

type Session struct {
	SvrFE string `json:"SvrFE,omitempty"`
	SvrID int32  `json:"SvrID,omitempty"`
	Mod   string `json:"Mod,omitempty"`
	Cmd   string `json:"Cmd,omitempty"`
	Time  int64  `json:"Time,omitempty"`
}

func NewSession(mod string, cmd string) *Session {
	sess := &Session{}
	//sess.SvrFE = GetServerName()
	//sess.SvrID = GetServerID()
	sess.Mod = mod
	sess.Cmd = cmd
	sess.Time = time.Now().Unix()

	return sess
}
func (s *Session) GetServerID() int32 {
	if s == nil {
		return 0
	}
	return s.SvrID
}
func (s *Session) GetServerFE() string {
	if s == nil {
		return ""
	}
	return s.SvrFE
}
func (s *Session) GetMod() string {
	if s == nil {
		return ""
	}
	return s.Mod
}

func (s *Session) GetCmd() string {
	if s == nil {
		return ""
	}
	return s.Cmd
}

type NatsMsg struct {
	Sess   Session `json:"Sess"`
	ErrNo  int32   `json:"err_no"`
	ErrStr string  `json:"err_str"`
	MsgData []byte    `json:"Data,omitempty"`
	//下面的没用
	NatsMsg *nats.Msg `json:"-"`
	//NatsMsgStr *string   `json:"-"`
}

func (m *NatsMsg) GetMsgData() []byte {
	if m == nil {
		return nil
	}
	return m.MsgData
}

func (m *NatsMsg) GetSession() *Session {
	if m == nil {
		return nil
	}
	return &m.Sess
}
func (m *NatsMsg) GetServerID() int32 {
	if m == nil {
		return 0
	}
	return m.Sess.GetServerID()
}

func (m *NatsMsg) GetServerFE() string {
	if m == nil {
		return ""
	}
	return m.Sess.GetServerFE()
}
func (m *NatsMsg) GetMod() string {
	if m == nil {
		return ""
	}
	return m.Sess.GetMod()
}

func (m *NatsMsg) GetCmd() string {
	if m == nil {
		return ""
	}
	return m.Sess.GetCmd()
}

func (m *NatsMsg) GetMsgErrNo() int32 {
	if m == nil {
		return 0
	}
	return m.ErrNo
}
func (m *NatsMsg) GetMsgErrStr() string {
	if m == nil {
		return ""
	}
	return m.ErrStr
}
func (m *NatsMsg) GetReply() string {
	if m == nil || m.NatsMsg == nil {
		return ""
	}
	return m.NatsMsg.Reply
}

func (m *NatsMsg) GetConn() *nats.Conn {
	if m == nil || m.NatsMsg == nil {
		return nil
	}
	return g_natsConn
}

func (m *NatsMsg) Marshal(para interface{}) error {
	if para == nil {
		m.MsgData = nil
		return nil
	}
	var err error
	m.MsgData, err = jsoniter.Marshal(para)
	return err
}

func (m *NatsMsg) Unmarshal(para interface{}) error {
	if para == nil {
		return nil
	}
	var err error
	err = jsoniter.Unmarshal(m.MsgData, para)
	return err
}

func genErrNatsMsg(errno int32, errstr string) *NatsMsg {
	rspmsg := &NatsMsg{Sess: *NewSession("", "")}
	//rspmsg.MsgBody.Mod = GetServerName()
	rspmsg.ErrNo = errno
	rspmsg.ErrStr = errstr

	return rspmsg
}

func (m *NatsMsg) Response(errno int32, errstr string, rspparam ...interface{}) (err error) {
	if m == nil || m.GetReply() == ""{
		return nil
	}
	rpcRsp := &NatsMsg{}
	rpcRsp.ErrNo = errno
	rpcRsp.ErrStr = errstr
	if len(rspparam) > 0 && rspparam[0] != nil{
		rpcRsp.Marshal(rspparam[0])
	}
	return NrpcReply(m, rpcRsp)
}
func (m *NatsMsg) ResponeSucc(rspparam interface{}) (err error) {
	errstr :=""
	if rspparam != nil{
		if s, ok := rspparam.(string);ok{
			errstr = s
			rspparam = nil
		}
	}
	return m.Response(0, errstr, rspparam)
}
