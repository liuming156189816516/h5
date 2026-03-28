package redisDeal

import (
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"
)

//设置Key ZADD
func ZaddRedis(key string, score int64, data interface{}) error {
	str, ok := data.(string)
	if !ok { //不是字符串
		datastr, err := jsoniter.MarshalToString(data)
		if err != nil {
			return err
		}
		str = datastr
	}
	err := TmplRedisSend("ZADD", key, score, str)
	if err != nil {
		return err
	}
	return nil
}

//该成员的 的分数增加一个数
func ZincrbyRedis(key string, inc int64, value interface{}) error {
	err := TmplRedisSend("ZINCRBY", key, inc, value)
	if err != nil {
		return err
	}
	return nil
}

//获取成员数
func ZcardRedis(key string) int64 {
	num, err := redis.Int64(TmplRedisDo("ZCARD", key))
	if err != nil {

		return num
	}
	return num
}

//获取 分数区间里面的 成员数
func ZCountRedis(key string, min, max int64) int64 {
	num, err := redis.Int64(TmplRedisDo("ZCOUNT", key, min, max))
	if err != nil {

		return num
	}
	return num
}

//获取该成员的 分数
func ZscoreRedis(key string, value interface{}) int64 {
	point, err := redis.Int64(TmplRedisDo("ZSCORE", key, value))
	if err != nil {

		return point
	}
	return point
}

//删除该成员
func ZremRedis(key string, value interface{}) error {
	return TmplRedisSend("ZREM", key, value)
}

//获取Key 排名从大到小 排列的 的成员 start (排名) 到end
func ZrevrangeRedis(key string, start int64, end int64, score string) []string {
	args := []interface{}{}
	args = append(args, key)
	args = append(args, start)
	args = append(args, end)
	if score != "" {
		args = append(args, score)
	}
	list, err := redis.Strings(TmplRedisDo("ZREVRANGE", args...))
	if err != nil {
		return []string{}
	}
	return list
}

//获取Key 分数从大到小 排列的 的成员 start (分数) 到end
func ZrevrangebyscoreRedis(key string, start int64, end int64, score string) []string {
	args := []interface{}{}
	args = append(args, key)
	args = append(args, end)
	args = append(args, start)
	if score != "" {
		args = append(args, score)
	}
	list, err := redis.Strings(TmplRedisDo("ZREVRANGEBYSCORE", args...))
	if err != nil {
		return []string{}
	}
	return list
}
