package controllers

import (
	"comm/goError"
	"fmt"
	"github.com/astaxie/beego"
	"net"
	"strings"
	"sync"
	"time"
)

// 基本路由
type BaseController struct {
	beego.Controller
}
type jsonResult struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 请求频率限制
var localLocker map[string]int64
var lm sync.Mutex

func init() {
	localLocker = make(map[string]int64)
	go func() {
		time.Sleep(3600 * time.Second)
		cleanLock()
	}()
}
func tryLock(uid string, cmd string) bool {
	key := fmt.Sprintf("%s.%s", cmd, uid)
	lm.Lock()
	defer lm.Unlock()
	t := time.Now().UnixNano()
	//3秒检测上传群组的接口
	if _, ok := exceptionList5Second[cmd]; ok {
		if v, ok1 := localLocker[key]; ok1 {
			if v > t-int64(5*time.Second) {
				return false
			}
		}
	} else {
		if v, ok := localLocker[key]; ok {
			if v > t-int64(100*time.Millisecond) {
				return false
			}
		}
	}

	localLocker[key] = t
	return true
}

// 不检查频率限制
var exceptionList = map[string]bool{
	"AccountController-GetQrCode":        true,
	"AccountController-GetAccountResult": true,
}

// 限制5秒频率
var exceptionList5Second = map[string]bool{}

// 不检查登录限制
var noLoginMap = map[string]bool{
	"GetQrCode":        true,
	"GetAccountResult": true,
}

// 每天清理一下
func cleanLock() {
	lm.Lock()
	defer lm.Unlock()
	t := time.Now().UnixNano()
	n := int64(0)
	for key, mt := range localLocker {
		if mt+int64(60*time.Second) < t { //60秒以前的删除
			delete(localLocker, key)
		}
		n++
		if n > 100000 { //处理100000个就够了 以免堵住别的用户
			return
		}
	}
}

/*
*
返回json数据
*/
func (this *BaseController) JsonResult(rsp *goError.ErrRsp, data interface{}) {
	if data == nil {
		data = ""
	}
	this.Data["json"] = jsonResult{rsp.Ret, rsp.Msg, data}
	this.ServeJSON()
}

// 获取客户端真实IP地址
func (this *BaseController) GetRealIp() string {
	cdnIp := this.Ctx.Request.Header["Cdn-Src-Ip"]
	if len(cdnIp) > 0 {
		ipListStr := strings.Replace(cdnIp[0], " ", "", -1)
		ipList := strings.Split(ipListStr, ",")
		ip0 := net.ParseIP(ipList[0])
		if len(ipList) > 0 && ip0 != nil && !ip0.IsLoopback() {
			// r.RemoteAddr = ipList[0]
			/*xRealPort := r.Header["X-Cdn-Src-Port"]
			if len(xRealPort) > 0 {
				port := xRealPort[0]
				r.RemoteAddr = r.RemoteAddr + ":" + port
			}*/
			return ipList[0]
		}
	}
	xRealIp := this.Ctx.Request.Header["X-Real-Ip"]
	if len(xRealIp) > 0 {
		ipListStr := strings.Replace(this.Ctx.Request.Header["X-Real-Ip"][0], " ", "", -1)
		ipList := strings.Split(ipListStr, ",")
		ip0 := net.ParseIP(ipList[0])
		if len(ipList) > 0 && ip0 != nil && !ip0.IsLoopback() {
			// r.RemoteAddr = ipList[0]
			/*xRealPort := r.Header["X-Real-Port"]
			if len(xRealPort) > 0 {
				port := xRealPort[0]
				r.RemoteAddr = r.RemoteAddr + ":" + port
			}*/
			return ipList[0]
		}
	}
	return this.Ctx.Request.RemoteAddr
}
