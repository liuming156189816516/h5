package redispool

import (
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Pool struct {
	*redis.Pool
	tomlPath string
}

var poolMap = map[string]*Pool{}
var poolLock = &sync.RWMutex{}

type RedisCfg struct {
	Host string `toml:"Host"`
	//Port               int32  `toml:"Port"`
	Passwd             string `toml:"Passwd"`
	MaxIdle            int
	MaxActive          int
	Wait               bool
	IdleTimeoutSeconds int
}

//type PoolConfig struct {
//}
//
//// 配置设计
//type RedisConfig struct {
//	Redis RedisCfg
//	Pool  PoolConfig
//}

//func NewPool(tomlPath string) (*Pool, error) {
//	poolLock.RLock()
//	pool := poolMap[tomlPath]
//	poolLock.RUnlock()
//	if pool != nil {
//		return pool, nil
//	}
//	poolLock.Lock()
//	defer poolLock.Unlock()
//	if pool = poolMap[tomlPath]; pool != nil {
//		return pool, nil
//	}
//	pool = &Pool{tomlPath: tomlPath}
//	if err := baselib.InitConfig(pool.loadPool); err != nil {
//		return nil, err
//	}
//	return pool, nil
//}
//
//
//func (pool *Pool) loadPool() error {
//	tomlPath := pool.tomlPath
//	redisConf := &RedisConfig{Pool: PoolConfig{MaxIdle: 100, MaxActive: 10000, IdleTimeoutSeconds: 240, Wait: false}}
//	// 将toml格式文件，解析到redisConf对象中
//	if err := baselib.DecodeToml(tomlPath, &redisConf); err != nil {
//		logs.TRACESVR("toml.DecodeFile %s Failed, err:", tomlPath, err)
//		return err
//	}
//
//	if redisConf.Redis.Host == "" || redisConf.Redis.Port <= 0 {
//		err := fmt.Errorf("LoadRedisConfig failed, path: %s, conf: %+v", tomlPath, redisConf)
//		logs.TRACESVR("LoadRedisConfig failed, path: %s, conf: %+v", tomlPath, redisConf)
//		return err
//	}
//	logs.TRACESVR("Toml %s Read Redis Config: %+v", tomlPath, redisConf)
//
//	if pool.Pool != nil {
//		pool.Pool.Close()
//	}
//	// redisServer := fmt.Sprintf("%s:%d", redisConf.Redis.Host, redisConf.Redis.Port)
//	pool.Pool = newPool(redisConf)
//	poolMap[tomlPath] = pool
//	return nil
//}
func (pool *Pool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func NewPool(redisConf *RedisCfg) *redis.Pool {
	//server := fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
	return &redis.Pool{
		MaxIdle:     redisConf.MaxIdle,
		MaxActive:   redisConf.MaxActive, //when zero,there's no limit. https://godoc.org/github.com/garyburd/redigo/redis#Pool
		IdleTimeout: time.Duration(redisConf.IdleTimeoutSeconds) * time.Second,
		Wait:        redisConf.Wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConf.Host)
			if err != nil {
				return nil, err
			}
			if redisConf.Passwd != "" {
				if _, err = c.Do("AUTH", redisConf.Passwd); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
