package sendmsg

import (
	"sync"
)

var staticMapLock = sync.Mutex{}
var staticMap = make(map[string]string)

func GetStaticData(mapKey string) string {
	staticMapLock.Lock()
	defer staticMapLock.Unlock()
	return staticMap[mapKey]
}

func SetStaticData(key string, val string) {
	staticMapLock.Lock()
	defer staticMapLock.Unlock()
	staticMap[key] = val
}

func DelStaticData(key string) {
	staticMapLock.Lock()
	defer staticMapLock.Unlock()
	delete(staticMap, key)
}

type SendMsgDetail struct {
	Name         string `json:"name"`
	SendSort     int64  `json:"send_sort"`
	TemplateType int64  `json:"template_type"`                //1文字 2 图片 3 链接 4 语音
	TemplateId   string `json:"template_id"`                  //当is_channel = 1 时,这个留空就好 当is_channel = 0时,必传
	SendInterval int64  `json:"send_interval"`                // 当is_channel = 1 默认写5
	Content      string `json:"content"`                      //发送内容 当is_channel = 0 时,这个留空就好 , 当is_channel = 1时,必传
	ReplyType    int64  `json:"reply_type" bson:"reply_type"` //回复类型 1-正常 2-已读 3-已回复
}
