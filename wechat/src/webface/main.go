package main

import (
	"comm/mgoDeal"
	"comm/redisDeal"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"natsRpc"
	"strings"
	"time"
	"utils"
	_ "webface/routers"
)

func main() {

	//日志系统
	l := logs.GetLogger()
	logs.EnableFuncCallDepth(true)
	beego.SetLogFuncCall(true)
	logs.Async()
	logs.Async(1e3)
	//beego.BeeLogger.DelLogger("console")
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/webface.log","separate":["debug","error"]}`)
	//正式上线 需要屏蔽 debug 日志。。不然会被日志搞半死
	logs.SetLevel(logs.LevelInfo)
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
	// 初始化数据库
	mgocfg := &mgoDeal.MgoCfg{}
	mgocfg.Url = beego.AppConfig.Strings("Mgo::Url")
	mgocfg.Query = beego.AppConfig.String("Mgo::Query")
	mgocfg.User = beego.AppConfig.String("Mgo::User")
	mgocfg.Password = beego.AppConfig.String("Mgo::Password")
	mgocfg.PoolNum, _ = beego.AppConfig.Int("Mgo::PoolNum")
	mgocfg.Ssl, _ = beego.AppConfig.Bool("Mgo::Ssl")
	mgocfg.SyncTask, _ = beego.AppConfig.Bool("Mgo::SyncTask")
	l.Println("mgo cfg ", mgocfg)
	if err := mgoDeal.StartMgoDb(mgocfg); err != nil {
		println("StartMgoDb %+v failed: %s", mgoDeal.GetMgoCoinfig(), err.Error())
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

	if err := natsRpc.StartNats(serverName, int32(serverId)); err != nil {
		println("InitEvent", err)
		return
	}
	//启动服务器
	l.Println("start webface")
	beego.Run()

}
