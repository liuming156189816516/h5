package main

import (
	"comm/redisDeal"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"mframe"
	"mlog"
	"natsRpc"
	"os"
	"proxy/ent"
	"strconv"
	"strings"
	"time"
	"utils"
)

/*
*
初始化事件
*/
func InitEvent(os string, ver int) error {
	err := error(nil)

	if err = ent.InitTaskEvent(os, ver); err != nil {
		return err
	}

	if err = ent.InitApiHandler(os, ver); err != nil {
		return err
	}

	return err
}

func main() {
	//日志系统
	l := logs.GetLogger()
	logs.EnableFuncCallDepth(true)
	beego.SetLogFuncCall(true)
	logs.Async()
	logs.Async(1e3)
	//beego.BeeLogger.DelLogger("console") //去掉 控制台日志
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/proxy.log","separate":["error"]}`)
	//正式上线 需要屏蔽 debug 日志。。不然会被日志搞半死
	//logs.SetLevel(logs.LevelInfo) //日志级别

	//随机种子
	rand.Seed(time.Now().UnixNano())

	serverName := beego.AppConfig.String("appname")
	serverId := utils.StrToInt64(beego.AppConfig.String("serverid"))
	if serverName == "" || serverId == 0 {
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
		println("StartRedis failed: ", err.Error())
		return
	}

	// 加载nats配置文件
	natsCfg := &natsRpc.NatsConfig{}
	natsCfg.Timeout, _ = beego.AppConfig.Int("Nats::Timeout")
	natsCfgServers := beego.AppConfig.String("Nats::Servers")
	natsCfg.User = beego.AppConfig.String("Nats::User")
	natsCfg.Password = beego.AppConfig.String("Nats::Password")
	natsCfg.Secure, _ = beego.AppConfig.Bool("Nats::Secure")
	natsCfg.Servers = strings.Split(natsCfgServers, ",")
	l.Println("natsCfg cfg ", natsCfg)
	natsRpc.LoadBeegoConfig(natsCfg)
	if err := natsRpc.StartNats(serverName, int32(serverId)); err != nil {
		println("InitEvent", err.Error())
		return
	}

	ver := 0
	if len(os.Args) > 1 {
		sid, err := strconv.Atoi(os.Args[1])
		if err == nil {
			ver = int(sid)
		}
	}
	// 初始化事件 initEvent
	if err := InitEvent("ios", ver); err != nil {
		println("InitEvent", err.Error())
		return
	}

	//启动日志
	mframe.Start(serverName, &mframe.FrameOption{Svrid: -1, OnLoadConfig: func() {
		logs.SetLevel(int(mlog.GetLogLevel()))
	}})
	defer mframe.Stop()

	//启动服务器
	l.Println("start proxy")
	beego.Run()
}
