package redispool

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"

	"time"
)

//redis Do 封装
//func Do(redisConfPath, commandName string, args ...interface{}) (reply interface{}, err error) {
//	//1. 获取redis连接
//	redisPool, err := NewPool(redisConfPath)
//	if err != nil {
//		return
//	}
//
//	//2. 根据参数，执行redis命令
//	reply, err = redisPool.Do(commandName, args...)
//
//	return
//}

func (pool *Pool) RedisLockUin(key string, uin uint64) error {

	strKey := fmt.Sprintf("redislock_%s_%d", key, uin)
	gap := time.Millisecond * 100
	for i := 0; i < 4; i++ {
		value, err := redis.Int64(pool.Do("INCR", strKey))
		if err != nil {
			//logs.ERRORLOG("conn.Do(INCR %s) failed, uin:%d error:%s", strKey, uin, err.Error())
			return errors.New("redis_conn.Do Failed!")
		}
		if value != 1 {
			time.Sleep(gap)
			gap *= 2
		} else {
			redis.Int64(pool.Do("EXPIRE", strKey, 60))
			//logs.TRACEPLAYER(uin, "RedisLockUin Success")
			return nil
		}
	}

	return errors.New("RedisLock Failed!")
}

func (pool *Pool) RedisUnLockUin(key string, uin uint64) {

	strKey := fmt.Sprintf("redislock_%s_%d", key, uin)
	redis.Int64(pool.Do("DEL", strKey))

	//logs.TRACEPLAYER(uin, "RedisUnLockUin Success")
	return
}
