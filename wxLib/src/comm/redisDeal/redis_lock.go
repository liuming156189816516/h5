package redisDeal

import (
	"fmt"
	"comm/comm"
	"time"
)

const redisLockKeyTtl = 1500 //毫秒

//redis 锁
type RedisLock struct {
	Key string
}

//上锁
func RedisDoLock(key string) *RedisLock {
	lockKey := fmt.Sprintf("%s.user.redis.lock.%s", comm.GetMgoDBName(), key)
	lock := &RedisLock{Key: lockKey}
	qpsTimer := time.NewTicker(redisLockKeyTtl * time.Millisecond) //每1.5秒
	defer func() {
		qpsTimer.Stop()
	}()
	for {
		select {
		case <-qpsTimer.C: //2秒 2秒还有人锁住 那就强制修改 锁时间 为2秒然后 返回 认为我上的锁
			TmplRedisSend("PEXPIRE", key, 2)
			return lock
		default:
			now := time.Now().UnixNano() / 1000 / 1000 //毫秒
			ret := SetnxRedis(lockKey, now, 2)         //
			if ret == 1 {                              //成功
				return lock
			}
		}
		time.Sleep(10 * time.Millisecond) //休息10毫秒
	}
	return lock
}

//释放锁
func (l *RedisLock) RedisUnLock() {
	TmplRedisSend("DEL", l.Key)
}
