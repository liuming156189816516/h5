package redisDeal

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	jsoniter "github.com/json-iterator/go"
)

//lpush
func RedisSendLpush(key string, data interface{}) error {
	str, ok := data.(string)
	if !ok { //不是字符串
		datastr, err := jsoniter.MarshalToString(data)
		if err != nil {
			logs.Error("解析失败 key:%s,value:%+v error:%+v", key, data, err)
			return err
		}
		str = datastr
	}
	err := TmplRedisSend("LPUSH", key, str)
	if err != nil {
		return err
	}
	return nil
}

//rpop
func RedisDoRpop(key string) string {
	ret, err := redis.String(TmplRedisDo("RPOP", key))
	if err != nil {
		return ""
	}
	return ret
}

//llen
func RedisDoLLen(key string) int64 {
	ret, err := redis.Int64(TmplRedisDo("LLen", key))
	if err != nil {
		return 0
	}
	return ret
}
