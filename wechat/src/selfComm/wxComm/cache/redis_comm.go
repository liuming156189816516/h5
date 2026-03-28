package cache

import (
	"comm/redisDeal"
	"comm/redisKeys"
	jsoniter "github.com/json-iterator/go"
)

// ip使用次数
func SetIpUserNum(ipId string, inc int64) {
	redisDeal.RedisDoHSet(redisKeys.GetIpUserNumKey(), ipId, inc)
}

func IncIpUserNum(ipId string, inc int64) {
	redisDeal.RedisDoHincrby(redisKeys.GetIpUserNumKey(), ipId, inc)
}

func GetIpUserNum(ipId string) int64 {
	return redisDeal.RedisDoHGetInt(redisKeys.GetIpUserNumKey(), ipId)
}
func DelIpUserNum(ipId string) {
	redisDeal.RedisSendHDel(redisKeys.GetIpUserNumKey(), ipId)
}

// ======================================================================================================================
// 正在工作中的任务列表
func SetDoTask(key2 string, data interface{}) {
	//redisDeal.RedisSendHSet(redisKeys.GetDoTaskList(), key2, data, 60*60*24*3)
}

func DelDoTask(key2 string) {
	//redisDeal.RedisSendHDel(redisKeys.GetDoTaskList(), key2)
}

// ======================================================================================================================
// 获取数据包上传状态
func SetSchedule(key2 string, data string) {
	redisDeal.RedisSendHSet(redisKeys.GetScheduleKey(), key2, data, 60*5)
}

func GetSchedule(key2 string) string {
	return redisDeal.RedisDoHGetSrt(redisKeys.GetScheduleKey(), key2)
}

// ======================================================================================================================
// 缓存数据包内容
func ScardDataPackListCount(dataPackId string) int64 {
	return redisDeal.RedisDoScard(redisKeys.GetDataPackListKey(dataPackId))
}

func SpopDataPackList(dataPackId string) string {
	return redisDeal.RedisDoSpop(redisKeys.GetDataPackListKey(dataPackId))
}

func SaddDataPackList(dataPackId, target string) {
	redisDeal.RedisSendSadd(redisKeys.GetDataPackListKey(dataPackId), target)
}

func DelDataPackList(dataPackId string) {
	redisDeal.RedisSendDel(redisKeys.GetDataPackListKey(dataPackId))
}

func SmembersDataPackList(dataPackId string) []string {
	return redisDeal.RedisSmembers(redisKeys.GetDataPackListKey(dataPackId))
}

func SaddListDataPackList(dataPackId string, args []interface{}) {
	redisDeal.RedisSendSaddList(redisKeys.GetDataPackListKey(dataPackId), args...)
}

// ======================================================================================================================
// 缓存数据包异常内容
func ScardDataPackListErrCount(dataPackId string) int64 {
	return redisDeal.RedisDoScard(redisKeys.GetDataPackListErrKey(dataPackId))
}

func SpopDataPackListErr(dataPackId string) string {
	return redisDeal.RedisDoSpop(redisKeys.GetDataPackListErrKey(dataPackId))
}

func SaddDataPackListErr(dataPackId, target string) {
	redisDeal.RedisSendSadd(redisKeys.GetDataPackListErrKey(dataPackId), target)
}

func DelDataPackListErr(dataPackId string) {
	redisDeal.RedisSendDel(redisKeys.GetDataPackListErrKey(dataPackId))
}

func SmembersDataPackListErr(dataPackId string) []string {
	return redisDeal.RedisSmembers(redisKeys.GetDataPackListErrKey(dataPackId))
}

func SaddListDataPackListErr(dataPackId string, args []interface{}) {
	redisDeal.RedisSendSaddList(redisKeys.GetDataPackListErrKey(dataPackId), args...)
}

// ======================================================================================================================
// 缓存数据包内容-所有
func SmembersDataPackList2(dataPackId string) []string {
	return redisDeal.RedisSmembers(redisKeys.GetDataPackListKey2(dataPackId))
}
func DelDataPackList2(dataPackId string) {
	redisDeal.RedisSendDel(redisKeys.GetDataPackListKey2(dataPackId))
}

func SaddListDataPackList2(dataPackId string, args []interface{}) {
	redisDeal.RedisSendSaddList(redisKeys.GetDataPackListKey2(dataPackId), args...)
}
func ScardDataPackListCount2(dataPackId string) int64 {
	return redisDeal.RedisDoScard(redisKeys.GetDataPackListKey2(dataPackId))
}

// ======================================================================================================================
// 一个账号下数据内容去重
func SremListAllDataPackList(args []interface{}) {
	redisDeal.RedisSremList(redisKeys.GetAllDataPackListKey(), args...)
}

func MemberListAllDataPackList(phone string) int64 {
	return redisDeal.RedisDoSisMember(redisKeys.GetAllDataPackListKey(), phone)
}

func SaddListAllDataPackList(args []interface{}) {
	redisDeal.RedisSendSaddList(redisKeys.GetAllDataPackListKey(), args...)
}

func GetAllDataPackList() []string {
	return redisDeal.RedisSmembers(redisKeys.GetAllDataPackListKey())
}

// ======================================================================================================================
// 任务配置
func SetTaskConfig(key string, taskInfo *TaskConfigInfo) {
	redisDeal.RedisDoHSet(redisKeys.GetTaskConfigKey(), key, taskInfo)
}

func GetTaskConfig(key string) *TaskConfigInfo {
	srt := redisDeal.RedisDoHGetSrt(redisKeys.GetTaskConfigKey(), key)
	taskInfo := &TaskConfigInfo{}
	jsoniter.UnmarshalFromString(srt, &taskInfo)
	return taskInfo
}
