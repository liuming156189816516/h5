package sendAlarm

import (
	"comm/comm"
	"golang.org/x/exp/errors/fmt"
	"net/url"
	"os"
)

var sysName = ""

//发送告警导群里面
func SendAlarm(msg string, proxy string) {
	if sysName == "" {
		sysName, _ = os.Hostname()
	}
	if sysName != "" {
		msg = fmt.Sprintf("机器:%s-%s", sysName, msg)
	}
	msg = url.QueryEscape(msg)
	url := "https://api.telegram.org/bot697282790:AAF0dpZeGo7OLatjSS78U3w7QjKe-DUSpdk/sendMessage?chat_id=-422425099&text=" + msg
	comm.HttpProxyRequest("GET", url, proxy, []byte{}, 10)
	return
}
