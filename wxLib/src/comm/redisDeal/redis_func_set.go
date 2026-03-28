package redisDeal

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"
)

//sadd
func RedisSendSadd(key string, data interface{}) error {
	str, ok := data.(string)
	if !ok { //不是字符串
		datastr, err := jsoniter.MarshalToString(data)
		if err != nil {
			logs.Error("解析失败 key:%s,value:%+v error:%+v", key, data, err)
			return err
		}
		str = datastr
	}
	err := TmplRedisSend("SADD", key, str)
	if err != nil {
		return err
	}
	return nil
}

//sadd 多个
func RedisSendSaddList(key string, data ...interface{}) error {
	var args []interface{}
	args = append(args, key)
	args = append(args, data...)
	err := TmplRedisSend("SADD", args...)
	if err != nil {
		return err
	}
	return nil
}

//sadd sum
func RedisSendSaddRepeat(key string, datas ...interface{}) error {
	info := []interface{}{}
	info = append(info, key)
	for id, d := range datas {
		str, ok := d.(string)
		if !ok { //不是字符串
			datastr, err := jsoniter.MarshalToString(d)
			if err != nil {
				logs.Error("解析失败 key:%s,value:%+v error:%+v", key, d, err)
				continue
			}
			str = datastr
		}
		info = append(info, str)
		if id != 0 && id%200 == 0 {
			err := TmplRedisSend("SADD", info...)
			if err != nil {
				continue
			}
			info = []interface{}{}
			info = append(info, key)
		}
	}
	if len(info) > 1 {
		err := TmplRedisSend("SADD", info...)
		if err != nil {
			return err
		}
	}

	return nil
}

//sadd 1 表示存在
func RedisDoSisMember(key string, data interface{}) int64 {
	str, ok := data.(string)
	if !ok { //不是字符串
		datastr, err := jsoniter.MarshalToString(data)
		if err != nil {
			logs.Error("解析失败 key:%s,value:%+v error:%+v", key, data, err)
			return 0
		}
		str = datastr
	}
	ret, err := redis.Int64(TmplRedisDo("SISMEMBER", key, str))
	if err != nil {
		return 0
	}
	return ret
}

//spop
func RedisDoSpop(key string) string {
	ret, err := redis.String(TmplRedisDo("SPOP", key))
	if err != nil {
		return ""
	}
	return ret
}

//SRANDMEMBER 随机娶一个值
func RedisDoSrandMember(key string) string {
	ret, err := redis.String(TmplRedisDo("SRANDMEMBER", key))
	if err != nil {
		return ""
	}
	return ret
}

//SCARD 集合数量
func RedisDoScard(key string) int64 {
	ret, err := redis.Int64(TmplRedisDo("SCARD", key))
	if err != nil {
		return 0
	}
	return ret
}

//	SMEMBERS 返回set 里面得全部数据
func RedisSmembers(key string) []string {

	ret, err := redis.Strings(TmplRedisDo("SMEMBERS", key))
	if err != nil {
		return []string{}
	}
	return ret
}

//移除一个
func RedisSrem(key string, data string) error {
	err := TmplRedisSend("SREM", key, data)
	if err != nil {
		return err
	}
	return err
}

//移除多个
func RedisSremList(key string, data ...interface{}) error {
	var args []interface{}
	args = append(args, key)
	args = append(args, data...)
	err := TmplRedisSend("SREM", args...)
	if err != nil {
		return err
	}
	return err
}