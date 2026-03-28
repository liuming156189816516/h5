package main

import (
	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"natsRpc"
	_ "testRpc/routers"
	"testRpc/rpcServices"
)

func main() {

	// 加载nats配置文件
	natsRpc.LoadConfig("./conf/NatsConfig.toml")
	// 启动 nats
	srvName := "testRpc"
	natsRpc.StartNats(srvName)

	// 初始化rpc请求 initRpc
	if err := rpcServices.InitRpc(); err != nil {
		log.Fatalf("InitRpc err:%+v", err)
	}
	// 初始化事件 initEvent
	if err := rpcServices.InitEvent(); err != nil {
		log.Fatalf("InitEvent err:%+v", err)
	}

	beego.Run()
}
