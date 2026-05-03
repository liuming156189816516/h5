package cache

import (
	"comm/redisDeal"
	"comm/redisKeys"
	jsoniter "github.com/json-iterator/go"
	"time"
)

// ======================================================================================================================
// 账号信息
func SetAccountInfo(acc string, accountInfo *AccountInfo) {
	redisDeal.RedisDoHSet(redisKeys.GetAccountListInfoKey(), acc, accountInfo)
}

func GetAccountInfo(acc string) *AccountInfo {
	srt := redisDeal.RedisDoHGetSrt(redisKeys.GetAccountListInfoKey(), acc)
	accinfo := &AccountInfo{}
	jsoniter.UnmarshalFromString(srt, &accinfo)
	return accinfo
}

func DelAccountInfo(acc string) {
	redisDeal.RedisSendHDel(redisKeys.GetAccountListInfoKey(), acc)
}

// ======================================================================================================================
// 账号的状态 账号状态 1-离线 2-在线 3-登录中 4-登录失败 5-离线中
func SetAccountStatus(acc string, status int64) {
	redisDeal.RedisDoHSet(redisKeys.GetAccountStatusKey(), acc, status)
}

func GetAccountStatus(acc string) int64 {
	return redisDeal.RedisDoHGetInt(redisKeys.GetAccountStatusKey(), acc)
}
func DelAccountStatus(acc string) {
	redisDeal.RedisSendHDel(redisKeys.GetAccountStatusKey(), acc)
}

// ======================================================================================================================
// 所有推广号,避免推广号重复
func SetAllAccountList(account string) {
	redisDeal.RedisDoHSet(redisKeys.GetAllAccountListKey(), account, time.Now().Unix())
}

func GetAllAccountList(account string) string {
	return redisDeal.RedisDoHGetSrt(redisKeys.GetAllAccountListKey(), account)
}

func DelAllAccountList(account string) {
	redisDeal.RedisSendHDel(redisKeys.GetAllAccountListKey(), account)
}

// ======================================================================================================================
// 账号抢登列表
func SetReplacedAccount(account string) {
	redisDeal.RedisSendSet(account, 1, 60*1)
}

func GetReplacedAccount(account string) int64 {
	return redisDeal.RedisDoGetInt(account)
}

// ======================================================================================================================
// 挂机二维码检测任务
func SetCheckQrcodeTask(uuid string, data interface{}) {
	redisDeal.RedisSendSet(redisKeys.GetCheckQrcodeTaskKey(uuid), data, 180)
}

func GetCheckQrcodeTask(uuid string) string {
	return redisDeal.RedisDoGetStr(redisKeys.GetCheckQrcodeTaskKey(uuid))
}

func DelCheckQrcodeTask(uuid string) {
	redisDeal.RedisSendDel(redisKeys.GetCheckQrcodeTaskKey(uuid))
}

// ======================================================================================================================
// 账号的代理ip
func SetProxyIp(acc string, proxyInfo *ProxyIpInfo) {
	redisDeal.RedisDoHSet(redisKeys.GetAccountProxyIpKey(), acc, proxyInfo)
}

func GetProxyIp(acc string) *ProxyIpInfo {
	prox := &ProxyIpInfo{}
	srt := redisDeal.RedisDoHGetSrt(redisKeys.GetAccountProxyIpKey(), acc)
	jsoniter.UnmarshalFromString(srt, &prox)
	return prox
}

func DelProxyIp(acc string) {
	redisDeal.RedisSendHDel(redisKeys.GetAccountProxyIpKey(), acc)
}

// =================================================
// 群发任务明细
func SetSendMsgRecord(tmp *SendMsgRecord) {
	redisDeal.RedisDoHSet(redisKeys.GetSendMsgTaskInfoKey(tmp.Account), tmp.RecordId, tmp)
}

// 群发任务明细详情
func GetSendMsgRecordInfo(account, recordId string) *SendMsgRecord {
	tmp := &SendMsgRecord{}
	v := redisDeal.RedisDoHGetSrt(redisKeys.GetSendMsgTaskInfoKey(account), recordId)
	jsoniter.UnmarshalFromString(v, &tmp)
	return tmp
}

// 群发任务详情
func DelSendMsgRecordInfo(account string) {
	redisDeal.RedisSendDel(redisKeys.GetSendMsgTaskInfoKey(account))
}

// =================================================
// 群发任务phone-lid
func SetSendMsgPhoneLid(account, phone, lid string) {
	redisDeal.RedisDoHSet(redisKeys.GetSendMsgPhoneLidKey(account), phone, lid)
}

func GetSendMsgPhoneLid(account, phone string) string {
	srt := redisDeal.RedisDoHGetSrt(redisKeys.GetSendMsgPhoneLidKey(account), phone)
	return srt
}

func DelSendMsgPhoneLid(account string) {
	redisDeal.RedisSendDel(redisKeys.GetSendMsgPhoneLidKey(account))
}

// ================================================================
// 群发任务统计
const (
	SuccessNum = "SuccessNum" //已完成
	ArrivedNum = "ArrivedNum" //已送达
	ReadNum    = "ReadNum"    //已读
)

// ================================================================
// 群发任务明细统计
func IncSendMsgTaskInfoCount(countType, account string, num int64) int64 {
	hincrby, _ := redisDeal.RedisDoHincrby(redisKeys.GetSendMsgTaskInfoCountKey(), countType+"_"+account, num)
	return hincrby
}

func GetSendMsgTaskInfoCount(countType, account string) int64 {
	getInt := redisDeal.RedisDoHGetInt(redisKeys.GetSendMsgTaskInfoCountKey(), countType+"_"+account)
	return getInt
}

// ======================================================================================================================
// 发送任务列表
func SetSendMsgTaskInfo(tmp *SendMsgTaskInfo) {
	redisDeal.RedisSendLpush(redisKeys.GetAllSendMsgTaskList(), tmp)
}

func GetSendMsgTaskInfo() *SendMsgTaskInfo {
	tmp := &SendMsgTaskInfo{}
	rpopStr := redisDeal.RedisDoRpop(redisKeys.GetAllSendMsgTaskList())
	if rpopStr != "" {
		jsoniter.UnmarshalFromString(rpopStr, &tmp)
	}
	time.Sleep(4 * time.Millisecond)
	return tmp
}

func LenSendMsgTaskInfo() int64 {
	return redisDeal.RedisDoLLen(redisKeys.GetAllSendMsgTaskList())
}

// ======================================================================================================================
// 自动发送任务列表
func SetAutoSendMsgTaskInfo(tmp *AutoSendMsgTaskInfo) {
	redisDeal.RedisSendLpush(redisKeys.GetAutoAllSendMsgTaskList(), tmp)
}

func GetAutoSendMsgTaskInfo() *AutoSendMsgTaskInfo {
	tmp := &AutoSendMsgTaskInfo{}
	rpopStr := redisDeal.RedisDoRpop(redisKeys.GetAutoAllSendMsgTaskList())
	if rpopStr != "" {
		jsoniter.UnmarshalFromString(rpopStr, &tmp)
	}
	time.Sleep(4 * time.Millisecond)
	return tmp
}

func LenAutoSendMsgTaskInfo() int64 {
	return redisDeal.RedisDoLLen(redisKeys.GetAutoAllSendMsgTaskList())
}

// =================================================
// 自动群发任务明细
func SetAutoSendMsgRecord(tmp *AutoSendMsgRecord) {
	redisDeal.RedisDoHSet(redisKeys.GetAutoSendMsgTaskInfoKey(tmp.Account), tmp.MessageId, tmp)
}

// 自动群发任务明细详情
func GetAutoSendMsgRecordInfo(account, messageId string) *AutoSendMsgRecord {
	tmp := &AutoSendMsgRecord{}
	v := redisDeal.RedisDoHGetSrt(redisKeys.GetAutoSendMsgTaskInfoKey(account), messageId)
	jsoniter.UnmarshalFromString(v, &tmp)
	return tmp
}

// 删除自动群发任务详情
func DelAutoSendMsgRecordInfo(account string) {
	redisDeal.RedisSendDel(redisKeys.GetAutoSendMsgTaskInfoKey(account))
}
