package wxComm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var API_VERSION = "v18.0"

var pixelTokens = map[string]string{
	"1733351437460530": "EAAMJGbgXFZBABRN0SsOVoRtBCEwya630fr2ZBRsMeZCcb6N4SPmMEwhjvaH7rI3Lv4x0nZAKGWFXkU69LTQuJjvZAVBroVdUxnJFEe8NBaqEtv3LdAogIb84z3lB7pu9Tbxdok8zsav8NDBIVNKsk8Y4IaHeu6oJEy8nQdI4uH38GdJEZCFr6IJxZCZCZCFOLSg1rNwZDZD",
	"794525866630709":  "EAAMJGbgXFZBABRIAIoFEXbGx4VvZAD7LfX4pfKZBXWzRCnO8HRByGLlUs3hfWlYXyZBLooONexvPLXDsj3iZCzm15ZBDZCc5ukHR2zN2uy7oNpucxQxzfQZBULyiAAGgitzMeRz2ESfQjBFRe1aE0k6IhMrlZAc3bUfspJuDOBvJFpaEXDr9MweMrsjZAvYnueRQZDZD",
}

type FbReportReq struct {
	Data []FbEvent `json:"data"`
}

type FbEvent struct {
	EventId      string     `json:"event_id"`
	EventName    string     `json:"event_name"`
	EventTime    int64      `json:"event_time"`
	ActionSource string     `json:"action_source"`
	UserData     UserData   `json:"user_data"`
	CustomData   CustomData `json:"custom_data"`
}

type UserData struct {
	Fbc string `json:"fbc"`
	Fbp string `json:"fbp"`
}

type CustomData struct {
	Type    string  `json:"type"`
	Amount  float64 `json:"amount"`
	OrderId string  `json:"orderId"`
}

func FbReport(eventID string, eventName, fbclid, Fbp, PixelId, ReqId string) {

	// ===== 构造请求数据 =====
	param := FbReportReq{}

	// 构造 event
	event := FbEvent{}
	event.EventId = eventName + eventID
	event.EventName = eventName
	event.EventTime = time.Now().Unix()
	event.ActionSource = "website"

	// 构造 user_data
	user := UserData{}
	//user.Fbc = "fb.1." + IntToStr(time.Now().UnixMilli()) + "." + fbclid
	user.Fbc = fmt.Sprintf("fb.1.%d.%s", time.Now().UnixNano()/1e6, fbclid)
	user.Fbp = Fbp

	// 赋值
	event.UserData = user

	if ReqId != "" {
		// 构造 custom_data
		custom := CustomData{}
		custom.Type = "USD"
		custom.Amount = 1.00
		custom.OrderId = ReqId
		event.CustomData = custom
	}

	// 放入 data 数组
	param.Data = []FbEvent{event}

	// ===== 序列化日志 =====
	reqStr, _ := jsoniter.MarshalToString(param)
	fmt.Println("FbReport req========>", reqStr)
	//global.GVA_LOG.Info("FbReport req =======>", zap.String("req", reqStr))

	// ===== API 地址 =====
	api := "https://graph.facebook.com/" + API_VERSION + "/" + PixelId + "/events?access_token=" + pixelTokens[PixelId]

	// ===== 请求头 =====
	//headMap := map[string]string{
	//	"Content-Type": "application/json",
	//}

	//// ===== 创建客户端 =====
	//client := NewClient(
	//	WithTimeout(60 * time.Second),
	//)
	//
	//// ===== 发起请求 =====
	//resp, err := client.Post(api, param, WithHeaders(headMap))
	//if err != nil {
	//	global.GVA_LOG.Error("FbReport post err =======>", zap.Error(err))
	//	return
	//}
	//
	//// ===== 打印返回=====
	//if resp != nil {
	//	global.GVA_LOG.Info("FbReport rsp =======>", zap.String("rsp", string(resp.Body)))
	//}

	client := &http.Client{}
	req, err := http.NewRequest("POST", api, strings.NewReader(reqStr))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
