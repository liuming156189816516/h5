package redisDeal

import (
	"comm/goError"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

func RedisConnDo(conn redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	reply, err = conn.Do(commandName, args...)
	//logs.Debug("RedisConnDo cmd:%s, args:%+v, err:%+v", commandName, args, err)
	return
}

func RedisConnSend(conn redis.Conn, commandName string, args ...interface{}) (err error) {
	err = conn.Send(commandName, args...)
	//logs.Debug("RedisConnDo cmd:%s, args:%+v, err:%+v", commandName, args, err)
	return
}

func TmplRedisDo(commandName string, args ...interface{}) (reply interface{}, err error) {
	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool  Failed:%s 致命错误", redis_conn.Err().Error())
		return nil, goError.Err_NoRedis
	}
	defer redis_conn.Close()
	return RedisConnDo(redis_conn, commandName, args...)
}

func TmplRedisSend(commandName string, args ...interface{}) (err error) {
	redis_conn := GetRedisPool().Get()
	if redis_conn.Err() != nil {
		logs.Error("GetRedisPool Failed:%s 致命错误", redis_conn.Err().Error())
		return goError.Err_NoRedis
	}
	defer redis_conn.Close()
	return RedisConnSend(redis_conn, commandName, args...)
}

////////////////
//把redis 结果转换成[]int64
func RedisInt64s(reply interface{}, err error) ([]int64, error) {
	var int64s []int64
	values, err := redis.Values(reply, err)
	if err != nil {
		return int64s, err
	}
	if err := redis.ScanSlice(values, &int64s); err != nil {
		return int64s, err
	}
	return int64s, nil
}
