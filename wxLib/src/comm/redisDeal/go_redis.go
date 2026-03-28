package redisDeal

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)
//第二个 redis 实列
var RedisClient *redis.Client
func RedisInitialize() *redis.Client {
	DB, err := beego.AppConfig.Int("redisdb")
	if err != nil {
		DB = 0
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.String("redislink"), // redis地址
		Password: beego.AppConfig.String("redispwd"),  // redis密码，没有则留空
		DB:       DB,                                  // 默认数据库，默认是0
	})
	return RedisClient
}


