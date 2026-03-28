package mlog

import (
	"bytes"
	"runtime"
)

const (
	LL_NONE   = 0x00000000 //!<没有级别，不打印
	LL_EMERG  = 1 << iota
	LL_ALERT  = 1 << iota
	LL_CRIT   = 1 << iota
	LL_ERROR  = 1 << iota
	LL_WARN   = 1 << iota
	LL_NOTICE = 1 << iota
	LL_INFO   = 1 << iota
	LL_DEBUG  = 1 << iota
	LL_ANY    = 0xFFFFFFFF //!<强制打印
)

var logLevelStr = map[int64]string{
	LL_EMERG:  "EMERG",
	LL_ALERT:  "ALERT",
	LL_CRIT:   "CRIT",
	LL_ERROR:  "ERROR",
	LL_WARN:   "WARN",
	LL_NOTICE: "NOTICE",
	LL_INFO:   "INFO",
	LL_DEBUG:  "DEBUG"}

func GetGpid() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	//	n, _ := strconv.ParseUint(string(b), 10, 64)
	return string(b)
}

type Logger interface {
	Async()

	SetPrintGpid(b bool)
	// 是否打印毫秒级时间
	SetPrintMs(b bool)
	// 设置日志头函数层数
	SetLogFuncCallDepth(d int)
	WriteMsg(l int, format string, v ...interface{})
	Close()
}
