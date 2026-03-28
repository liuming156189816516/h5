package redisDeal

import (
	"comm/goError"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"
)

//HINCRBY do 返回结果值得
func RedisDoHincrby(key string, key2 interface{}, inc int64) (int64, error) {
	//logs.Debug("HINCRBY key%s key2:%s inc=%d", key, key2, inc)
	return redis.Int64(TmplRedisDo("HINCRBY", key, key2, inc))
}

//HINCRBY do 返回结果值得
func RedisSendHincrby(key string, key2 interface{}, inc int64, ttl ...int64) error {
	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool %s Failed:%s", key, redis_conn.Err().Error())
		return goError.Err_NoRedis
	}
	defer redis_conn.Close()
	TmplRedisSend("HINCRBY", key, key2, inc)
	if len(ttl) > 0 {
		RedisConnSend(redis_conn, "EXPIRE", key, ttl[0])
	}
	return nil
}

//hlen int64
func RedisDoHLen(key string) int64 {
	ret, err := redis.Int64(TmplRedisDo("HLEN", key))
	if err != nil {
		return 0
	}
	return ret
}

//Hkeys int64
func RedisDoHKeys(key string) []string {
	ret, err := redis.Strings(TmplRedisDo("HKEYS", key))
	if err != nil {
		return []string{}
	}
	return ret
}

//hget int64
func RedisDoHGetInt(key string, key2 interface{}) int64 {
	ret, err := redis.Int64(TmplRedisDo("HGET", key, key2))
	if err != nil {
		return 0
	}
	return ret
}

//指定得 field 是否存在 1 存在 0不存在
func RedisDoHexists(key string, key2 interface{}) int64 {
	ret, err := redis.Int64(TmplRedisDo("HEXISTS", key, key2))
	if err != nil {
		return 0
	}
	return ret
}

//hget str
func RedisDoHGetSrt(key string, key2 interface{}) string {
	ret, err := redis.String(TmplRedisDo("HGET", key, key2))
	if err != nil {
		return ""
	}
	return ret
}

//hset
func RedisSendHSet(key string, key2, data interface{}, ttl ...int64) error {
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
	RedisConnSend(redis_conn, "HSET", key, key2, str)
	if len(ttl) > 0 {
		RedisConnSend(redis_conn, "EXPIRE", key, ttl[0])
	}
	return nil
}

//hset
func RedisDoHSet(key string, key2, data interface{}) error {
	str, ok := data.(string)
	if !ok { //不是字符串
		datastr, err := jsoniter.MarshalToString(data)
		if err != nil {
			logs.Error("解析失败 key:%s,value:%+v error:%+v", key, data, err)
			return err
		}
		str = datastr
	}
	_, err := TmplRedisDo("HSET", key, key2, str)
	return err
}

//hGetAll
func RedisDoHGetAll(key string) map[string]string {
	s, err := redis.Strings(TmplRedisDo("HGETALL", key))
	if err != nil {
		return map[string]string{}
	}
	l := len(s)
	if l%2 != 0 {
		return map[string]string{}
	}
	ret := map[string]string{}
	for i := 0; i < l; i += 2 {
		ret[s[i]] = s[i+1]
	}
	return ret
}

//hscan
func RedisDoHGetAllByHscan(key string) map[string]string {
	s, err := redis.Strings(TmplRedisDo("HGETALL", key))
	if err != nil {
		return map[string]string{}
	}
	l := len(s)
	if l%2 != 0 {
		return map[string]string{}
	}
	ret := map[string]string{}
	for i := 0; i < l; i += 2 {
		ret[s[i]] = s[i+1]
	}
	return ret
}

//hdel
func RedisSendHDel(key string, key2 interface{}) error {
	return TmplRedisSend("HDEL", key, key2)
}
