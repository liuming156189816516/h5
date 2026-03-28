package redisDeal

import (
	"comm/goError"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"
)

//get int64
func RedisDoGetInt(key string) int64 {
	ret, err := redis.Int64(TmplRedisDo("GET", key))
	if err != nil {
		return 0
	}
	return ret
}

//get string
func RedisDoGetStr(key string) string {
	ret, err := redis.String(TmplRedisDo("GET", key))
	if err != nil {
		return ""
	}
	return ret
}

//直接操作key 的
//匹配对应的key
func RedisSendPattern(pattern string) []string {
	ret, err := redis.Strings(TmplRedisDo("KEYS", pattern))
	if err != nil {
		return []string{}
	}
	return ret
}

//del
func RedisSendDel(key string) error {
	return TmplRedisSend("DEL", key)
}

//ttl
func RedisSendTtl(key string, ttl int64) error {
	return TmplRedisSend("EXPIRE", key, ttl)
}

//ttl -2 不存在 -1 永久
func RedisDoGetTtl(key string) int64 {
	ret, err := redis.Int64(TmplRedisDo("TTL", key))
	if err != nil {
		return -2
	}
	return ret
}

//set
func RedisSendSet(key string, data interface{}, ttl ...int64) error {
	str, ok := data.(string)
	if !ok { //不是字符串
		datastr, err := jsoniter.MarshalToString(data)
		if err != nil {
			logs.Error("解析失败 key:%s,value:%+v error:%+v", key, data, err)
			return err
		}
		str = datastr
	}
	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool %s Failed:%s", key, redis_conn.Err().Error())
		return goError.Err_NoRedis
	}
	defer redis_conn.Close()
	RedisConnSend(redis_conn, "SET", key, str)
	if len(ttl) > 0 {
		RedisConnSend(redis_conn, "EXPIRE", key, ttl[0])
	}
	return nil
}

//互斥写一个 数字 返回0 说明不成功 上锁用 不能 写值
func SetnxRedis(key string, value interface{}, ttl ...int64) int {
	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool %s Failed:%s", key, redis_conn.Err().Error())
		return 0
	}
	defer redis_conn.Close()
	ret, err := redis.Int(RedisConnDo(redis_conn, "SETNX", key, value))
	if len(ttl) > 0 && ret == 1 {
		ttl := ttl[0]
		RedisConnSend(redis_conn, "EXPIRE", key, ttl)
	}
	if err != nil {
		return 0
	}
	return ret
}

//IncrBy
func RedisSendInc(key string, inc int64, ttl ...int64) error {

	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool %s Failed:%s", key, redis_conn.Err().Error())
		return goError.Err_NoRedis
	}
	defer redis_conn.Close()
	RedisConnSend(redis_conn, "INCRBY", key, inc)
	if len(ttl) > 0 {
		RedisConnSend(redis_conn, "EXPIRE", key, ttl[0])
	}
	return nil
}

//IncrBy
func RedisDoInc(key string, inc int64, ttl ...int64) int64 {

	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool %s Failed:%s", key, redis_conn.Err().Error())
		return 0
	}
	defer redis_conn.Close()
	ret, _ := redis.Int64(RedisConnDo(redis_conn, "INCRBY", key, inc))
	if len(ttl) > 0 {
		RedisConnSend(redis_conn, "EXPIRE", key, ttl[0])
	}
	return ret
}
