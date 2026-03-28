package redisDeal

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisCfg struct {
	IdleCount int //空闲连接数上线
	Servers   string
	Password  string
}

var redis_config = &RedisCfg{}
var G_redisPool *redis.Pool = nil

func StartRedis(cfg *RedisCfg, idleCount ...int) error {
	c := 1024
	if len(idleCount) > 0 {
		c = idleCount[0]
	}
	cfg.IdleCount = c
	redis_config = cfg
	return openRedis(c)
}

//获取pool
func GetRedisPool() *redis.Pool {
	if G_redisPool == nil {
		logs.Error("找不到redis:%+v", redis_config)
		return &redis.Pool{}
	}
	return G_redisPool
}
func openRedis(idleCount int) error {
	idleCount = redis_config.IdleCount
	if idleCount == 0 {
		idleCount = 1024
	}
	G_redisPool = &redis.Pool{
		MaxIdle:     idleCount,
		MaxActive:   0, //when zero,there's no limit. https://godoc.org/github.com/garyburd/redigo/redis#Pool
		IdleTimeout: time.Minute * 10,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redisDial(redis_config.Servers, redis_config.Password)
		},
		TestOnBorrow: redisOnBorrow,
	}
	//测试一下是否链接上了
	redis_conn := GetRedisPool().Get()
	if err := redis_conn.Err(); err != nil {
		return err
	}
	defer redis_conn.Close()
	logs.Trace("Conn to Redis: %+v", redis_config.Servers)
	return nil
}

func redisDial(servers string, pass string) (redis.Conn, error) {
	c, err := redis.Dial("tcp", servers)
	if err != nil {
		logs.Error("redis.Dial(%s) failed:%s", servers, err.Error())
		return nil, err
	}
	if pass != "" {
		if _, err := c.Do("AUTH", pass); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}
