package token

import (
	"comm/encrypt"
	"comm/redisDeal"
	"comm/redisKeys"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"math/rand"
	"time"
)

type TokenInfo struct {
	Src         string `json:"src"`
	Uid         string `json:"uid"`
	Tuid        string `json:"tuid"`
	Db          string `json:"db"`
	Extime      int64  `json:"extime"`
	Ip          string `json:"ip"`
	Pack        int64  `json:"pack"`
	RandNum     int64  `json:"rand_num"`
	Token       string `json:"token"`
	AccountType int64  `json:"account_type"`
}

func getUserKey(src, uid string) string {
	return fmt.Sprintf("%s-%s", src, uid)
}

//保存token
func SaveToken(info *TokenInfo) string {

	//if info.AccountType == 1 {
	//	DelToken(info.Src, info.Uid)
	//}
	userKey := getUserKey(info.Src, info.Uid)
	info.RandNum = int64(rand.Int31())
	sesskeyStr := fmt.Sprintf("%s|%s|%s|%d|%d", info.Src, info.Uid, info.Db, info.Extime, info.RandNum)
	logs.Debug("SaveLogin sesskeyStr = %s", sesskeyStr)
	sessionKey := encrypt.EncryptString(sesskeyStr, string(encrypt.AES256cbcKey), string(encrypt.AES256cbcIv), false)
	logs.Debug("userKey = %s sessionKey = %s", userKey, sessionKey)
	info.Token = sessionKey

	redisKey := redisKeys.GetUserTokenInfo()
	redisDeal.RedisSendHSet(redisKey, info.Token, info,60*60*24)
	redisDeal.RedisSendHSet(redisKey, userKey, info.Token,60*60*24)

	//做互相提 就在这里 删除老的 token 就行

	return info.Token
}

//检查token
func CheckToken(token string) *TokenInfo {
	info := &TokenInfo{}
	if token == "" {
		return info
	}
	redisKey := redisKeys.GetUserTokenInfo()
	str := redisDeal.RedisDoHGetSrt(redisKey, token)
	if str == "" {
		return info
	}
	err := jsoniter.UnmarshalFromString(str, &info)
	if err != nil {
		return &TokenInfo{}
	}
	if info.Token != token || info.Uid == "" {
		return &TokenInfo{}
	}
	if info.Extime < time.Now().Unix() {
		return &TokenInfo{}
	}
	return info
}

//删除token
func DelToken(src string, uid string) {
	userKey := getUserKey(src, uid)
	redisKey := redisKeys.GetUserTokenInfo()
	token := redisDeal.RedisDoHGetSrt(redisKey, userKey)
	redisDeal.RedisSendHDel(redisKey, userKey)
	if token != "" {
		redisDeal.RedisSendHDel(redisKey, token)
	}
}

//获取token
func GetToken(src string, uid string) string {
	userKey := getUserKey(src, uid)
	redisKey := redisKeys.GetUserTokenInfo()
	token := redisDeal.RedisDoHGetSrt(redisKey, userKey)
	return token
}
