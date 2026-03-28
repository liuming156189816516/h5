package redispool_test

import (
	"fmt"
	"qqgame/baselib/redispool"
)

func test() {
	if pool, err := redispool.NewPool("/etc/goservice/pc/social/redis.toml"); err != nil {
		fmt.Println(err)
		return
	}
	conn := pool.Get()
	defer conn.Close()
	conn.Send("SET", "XxxKey", "YyyValue")
}
