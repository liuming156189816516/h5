// Package tof IT服务-TOF助手，利用接口自动发送 邮件，rtx，微信，短信，使用前需要到 http://tof.oa.com/application/views/system_list.php# 申请权限
// 添加api: http://tof.oa.com/application/views/system.api.php?sysid=26320
// 添加发件人：http://tof.oa.com/application/views/system_sender.php?sysid=26320
// 权限审批可找rtx: IT服务-TOF助手(TOF)
// 更多详细信息可以参考:
// http://km.oa.com/knowledge/1937
// http://km.oa.com/knowledge/3417
// http://km.oa.com/group/599/articles/show/271103
package tof

/*
@version v1.0
@author nickzydeng
@copyright Copyright (c) 2018 Tencent Corporation, All Rights Reserved
@license http://opensource.org/licenses/gpl-license.php GNU Public License

You may not use this file except in compliance with the License.

Most recent version can be found at:
http://git.code.oa.com/going_proj/going_proj

Please see README.md for more information.
*/

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"going/cat"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"qqgame/baselib/logs"
)

var ConfPath = "../config/tof.toml"

// Conf tof配置格式
type Conf struct {
	Sysid      string          `toml:"Sysid"`
	Appkey     string          `toml:"Appkey"`
	Nettype    string          `toml:"Nettype"`             // 网络类型 devnet idc
	Host       string          `toml:"Host"` // tof域名在不同网络下不相同 idc: api.tof.oss.com 10.123.119.83; devnet: api.tof.oa.com 10.14.83.154
	Timeout    int64 			`toml:"Timeout"`              // 超时时间
	Ratelimit  int64           `toml:"Ratelimit"`              // 频率限制，默认每秒最多只能调20次，防止出事故
	MailAttr   int
	RtxAttr    int
	SmsAttr    int
	WeixinAttr int
	SuccAttr   int
	FailAttr   int
	LimitAttr  int
}

var conf = struct {
	Tof Conf
}{}
var once sync.Once
var limiter cat.Limiter

// GetConfig 返回tof配置信息
func GetConfig() *Conf {
	once.Do(func() {
		if _, e := toml.DecodeFile(ConfPath, &conf.Tof); e != nil {
			logs.TRACESVR("toml.DecodeFile %s Failed, err:", ConfPath, e)
		}
		if conf.Tof.Ratelimit > 0 {
			limiter = cat.NewLeakyBucketLimiter(conf.Tof.Ratelimit)
		}
	})

	return &conf.Tof
}

// SendMail 发送邮件 cc抄送, bcc密送, priority优先级 0:普通 1:高, emailType邮件类型 0:外部邮件 1:内部邮件, bodyFormat内容格式 0:文本 1:html格式
func SendMail(receiver, cc, bcc, sender, title, content string, priority, bodyFormat, emailType int) error {
	params := make(map[string]string)
	params["To"] = receiver
	params["From"] = sender
	params["CC"] = cc
	params["Bcc"] = bcc
	params["Title"] = title
	params["Content"] = content
	params["EmailType"] = strconv.Itoa(emailType)
	params["BodyFormat"] = strconv.Itoa(bodyFormat)
	params["Priority"] = strconv.Itoa(priority)
	return DoRequest(params, "mail")
}

// SendRTX 发送rtx
func SendRTX(receiver, sender, title, content string, priority int) error {
	params := make(map[string]string)
	params["Receiver"] = receiver
	params["Sender"] = sender
	params["Title"] = title
	params["MsgInfo"] = content
	params["Priority"] = strconv.Itoa(priority)
	return DoRequest(params, "rtx")
}

// SendWeiXin 发送微信 注意：这里的微信帐号只能是 alarm.weixin.oa.com 上申请的微信报警帐号，不是私人微信帐号
func SendWeiXin(receiver, sender, content string, priority int) error {
	params := make(map[string]string)
	params["Receiver"] = receiver
	params["Sender"] = sender
	params["MsgInfo"] = content
	params["Priority"] = strconv.Itoa(priority)
	return DoRequest(params, "weixin")
}

// DoRequest 发送http请求
func DoRequest(params map[string]string, op string) error {
	if GetConfig().Sysid == "" || GetConfig().Appkey == "" {
		return fmt.Errorf("sysid or appkey empty")
	}
	if GetConfig().Ratelimit > 0 && !limiter.Acquire() {
		return fmt.Errorf("over rate limit")
	}
	logs.DEBUGLOG("config:%+v", GetConfig())
	var query string
	switch op {
	case "mail":
		query = fmt.Sprintf("http://%s/api/v1/Message/SendMail", GetConfig().Host)
	case "rtx":
		query = fmt.Sprintf("http://%s/api/v1/Message/SendRTX", GetConfig().Host)
	case "weixin":
		query = fmt.Sprintf("http://%s/api/v1/Message/SendWeiXin", GetConfig().Host)
	default:
		return fmt.Errorf("not support op")
	}

	var err error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			return err
		}
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", query, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	magicNum := rand.Intn(8999) + 1000 // 生成一个四位数的随机数
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	req.Header.Set("appkey", GetConfig().Appkey)
	req.Header.Set("random", strconv.Itoa(magicNum))
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("host", GetConfig().Host)

	signature, err := buildSign(GetConfig().Sysid, magicNum, timestamp)
	if err != nil {
		return err
	}
	req.Header.Set("signature", signature)

	client := &http.Client{Timeout: time.Duration(2 * time.Second)}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Printf("errbyjain: %s.\n", err.Error())
		return err
	}

	body = &bytes.Buffer{}
	_, err = body.ReadFrom(rsp.Body)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	fmt.Println(string(body.Bytes()))

	var rspBody = &struct {
		Ret        int
		ErrCode    int
		ErrMsg     string
		StackTrace string
		Data       bool
	}{}
	err = json.Unmarshal(body.Bytes(), rspBody)
	if err != nil {
		return err
	}

	if rsp.StatusCode != 200 || rspBody.ErrCode != 0 {
		return fmt.Errorf("http status code:%d, rsp body:%+v", rsp.StatusCode, rspBody)
	}

	return nil
}

// buildSign 生成tof3头部签名,具体说明可以参见:http://km.oa.com/articles/show/338381
func buildSign(sysid string, magicNum int, timestamp string) (string, error) {
	key := fmt.Sprintf("%s%s", sysid, "--------")[0:8] // 生成appid，需要用-补全到8位
	plain := fmt.Sprintf("%s%d%s%s", "random", magicNum, "timestamp", timestamp)
	crypted, err := desEncrypt([]byte(plain), []byte(key))
	if err != nil {
		return "", err
	}
	signature := strings.ToUpper(hex.EncodeToString(crypted))
	return signature, nil
}

func desEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pkcs5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
