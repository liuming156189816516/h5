package wxComm

import (
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"selfComm/wxComm/wxHttp"
	"time"
)

type Ant struct {
	SessionId string      `json:"sessionId"`
	Channel   string      `json:"channel"`
	Data      interface{} `json:"data"`
}

type ImportJsonRsp struct {
	Ok bool `json:"ok"`
}

// 导入json号
func ImportJson(account string, param interface{}) (*ImportJsonRsp, error) {
	api := "https://tls.v168.vip/api/login/import-json"
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"
	channel := beego.AppConfig.String("channel")
	ant := &Ant{
		SessionId: GenerateUniqueString(),
		Channel:   channel,
		Data:      param,
	}

	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Headers: headerMap,
		Method:  "post",
		Content: ant,
		Timeout: 30,
	})
	ret := &ImportJsonRsp{}
	if rsp.Err == nil {
		//logs.Info("ImportJson result: " + string(rsp.Body))
		jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	} else {
		logs.Info("ImportJson err  account: " + account + rsp.Err.Error())
		return ret, rsp.Err
	}

	return ret, nil
}

// GenerateUniqueString 使用随机字节+纳秒时间戳生成不重复字符串
func GenerateUniqueString() string {
	// 生成 16 字节随机数
	randomBytes := make([]byte, 5)
	_, _ = rand.Read(randomBytes)
	randomPart := hex.EncodeToString(randomBytes)
	// 附加纳秒时间戳
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	// 组合并返回 (可进一步缩短长度)
	return randomPart + timestamp
}
