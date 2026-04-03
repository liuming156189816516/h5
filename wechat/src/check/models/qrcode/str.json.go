package qrcode

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

type CheckQrcodeTaskData struct {
	ProxyId     string `json:"proxy_id"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	DefaultGid  string `json:"default_gid"`
	AccountType int64  `json:"account_type"`
	PixelId     string `json:"pixel_id"` //kwai pixelId
	ClickId     string `json:"click_id"` //kwai clickid
	AreaCode    string `json:"area_code"`
	IsProxyUser int64  `json:"is_proxy_user"`
}

type QrcodeCheckData struct {
	Token    string `json:"token"`
	Device   string `json:"device"`
	Nickname string `json:"nickname"`
	Synckeys string `json:"synckeys"`
}
