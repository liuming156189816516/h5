package main

import (
	"comm/redisDeal"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"mframe"
	"mlog"
	"natsRpc"
	"strings"
	"task/models/server"
	"time"
)

/*
*
初始化事件
*/
func InitEvent() error {
	err := error(nil)
	return err
}

func main() {

	//日志系统
	l := logs.GetLogger()
	logs.EnableFuncCallDepth(true)
	beego.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(2)
	logs.Async(1e3)

	//beego.BeeLogger.DelLogger("console")
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/task.log","separate":["debug","error"]}`)

	//正式上线 需要屏蔽 debug 日志。。不然会被日志搞半死
	//logs.SetLevel(logs.LevelInfo)
	//随机种子
	rand.Seed(time.Now().UnixNano())

	serverName := beego.AppConfig.String("appname")
	if serverName == "" {
		println("服务名称 和服务id 不对 ")
		return
	}

	// 初始化REDIS
	rediscfg := &redisDeal.RedisCfg{}
	rediscfg.IdleCount, _ = beego.AppConfig.Int("Redis::IdleCount")
	rediscfg.Servers = beego.AppConfig.String("Redis::Servers")
	rediscfg.Password = beego.AppConfig.String("Redis::Password")
	l.Println("redis cfg ", rediscfg)
	if err := redisDeal.StartRedis(rediscfg); err != nil {
		println("StartRedis failed: %s", err.Error())
		return
	}
	// 加载nats配置文件
	natsCfg := &natsRpc.NatsConfig{}
	natsCfg.Timeout, _ = beego.AppConfig.Int("Nats::Timeout")
	natsCfgServers := beego.AppConfig.String("Nats::Servers")
	natsCfg.User = beego.AppConfig.String("Nats::User")
	natsCfg.Password = beego.AppConfig.String("Nats::Password")
	natsCfg.Secure, _ = beego.AppConfig.Bool("Nats::Secure")
	l.Println("natsCfg cfg ", natsCfg)
	natsCfg.Servers = strings.Split(natsCfgServers, ",")
	natsRpc.LoadBeegoConfig(natsCfg)

	if err := natsRpc.StartNats(serverName); err != nil {
		println("InitEvent", err)
		return
	}
	// 初始化事件 initEvent
	if err := InitEvent(); err != nil {
		println("InitEvent", err)
		return
	}

	//启动日志
	mframe.Start(serverName, &mframe.FrameOption{Svrid: -1, OnLoadConfig: func() {
		logs.SetLevel(int(mlog.GetLogLevel()))
	}})
	defer mframe.Stop()
	mlog.Trace("start :%s", serverName)

	//任务 读取
	go server.RunTask()

	//启动服务器
	l.Println("start task")
	beego.Run()
}
