package natsRpc

import (
	"crypto/tls"
	"fmt"
	"mlog"
	"github.com/nats-io/go-nats"
	"golang.org/x/exp/errors"
	"os"
	"time"
	"mframe"
)

var g_natsConn *nats.Conn
var serverNameToNats string
var _mod string
var _id int32

func genNameToNats() string {
	if serverNameToNats != "" {
		return serverNameToNats
	}
	hostName, _ := os.Hostname()
	serverNameToNats = fmt.Sprintf("goserver.%s.%d-%s.%d", hostName, os.Getpid(), _mod, _id)
	return serverNameToNats
}
func GetServerName() string {
	return _mod
}
func GetServerId() int32 {
	return _id
}
func ConnToNats(cfg *NatsConfig) (*nats.Conn, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 2
	}

	var tlscfg *tls.Config
	if cfg.Secure {
		tlscfg = &tls.Config{}
		tlscfg.InsecureSkipVerify = true
	}

	var NatsOpts = nats.Options{
		//	Url:            natsConfig.Url,
		User:     cfg.User,
		Password: cfg.Password,
		//	Token:          natsConfig.Token,
		Servers:        cfg.Servers,
		Name:           genNameToNats(),
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  100 * time.Millisecond,
		Timeout:        time.Duration(cfg.Timeout) * time.Second,
		Secure:         cfg.Secure,
		TLSConfig:      tlscfg,
	}

	var err = error(nil)
	newNatsConn, err := NatsOpts.Connect()
	// 初始化Nats
	if err != nil {
		fmt.Printf("nats.Connect error:%s", err)
		mlog.Error("nats.Connect error:%s", err)
		return nil, err
	}
	return newNatsConn, nil
}

func StartNats(mod string, ids ...int32) error {
	_mod = mod
	if len(ids) > 0 {
		_id = ids[0]
	}
	if _id == 0 {
		_id = 1000 //1000 开始
	}
	oldId := _id
	conn, err := ConnToNats(GetNatsConfig())
	if err != nil {
		mlog.Error("ConnToNats config:%+v, err:%+v", GetNatsConfig(), err)
		return err
	}
	if ok := askServer(conn); !ok { //
		return errors.New("服务已满")
	}
	if oldId != _id { //变了 需要重新 获取新的
		conn.Close()
		//重新连接
		serverNameToNats = ""
		conn, err = ConnToNats(GetNatsConfig())
		if err != nil {
			return err
		}
	}
	mframe.SetServerID(_id)
	g_natsConn = conn
	//注册 服务
	g_natsConn.Subscribe(genAskSubjToNats(), onServerAsk)
	mlog.Debug("ConnToNats config:%+v succ,serveri= %d", GetNatsConfig(), _id)
	return nil
}

func onServerAsk(natsMsg *nats.Msg) {
	g_natsConn.Publish(natsMsg.Reply, []byte(genNameToNats()))
}

func askServer(natsConn *nats.Conn) bool {
	for i := _id; i < 1100; i++ {
		_id = int32(i)
		reply, err := natsConn.Request(genAskSubjToNats(), []byte("are you here?"), time.Millisecond*50)
		mlog.Trace("Ask Server %s, replay: %+v, err: %+v", genAskSubjToNats(), reply, err)
		if err == nil {
			mlog.Trace("serverid here %d", _id)
			continue
		}
		return true
	}
	return false
}

func genAskSubjToNats() string {
	return fmt.Sprintf("system.service.ask.%s.%d", GetServerName(), GetServerId())
}
