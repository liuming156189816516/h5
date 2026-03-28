package mlog

import (
	"fmt"
)

const playerLogDateSuffix = "ymd"

type playerLogMsg struct {
	uin string
	msg string
}
type PlayerLoggerMng struct {
	baseFilePrefix string
	traceAllUin    bool
	traceUin       map[string]bool

	mtLogger *TMultiLogger
}

var playerLoggerMng *PlayerLoggerMng

func GetPlayerLoggerMng() *PlayerLoggerMng {
	return playerLoggerMng
}
func newPlayerLoggerMng() *PlayerLoggerMng {
	plm := new(PlayerLoggerMng)
	//u, _ := user.Current()
	//user := ""
	//if u != nil {
	//	user = u.Username
	//}
	plm.baseFilePrefix = GetLogBasePath() + "/wxlogs/"

	plm.traceUin = make(map[string]bool)
	plm.mtLogger = NewMultiLogger(plm.baseFilePrefix, 200, 3, serverAll)
	plm.mtLogger.StartLogger()
	return plm
}

func (plm *PlayerLoggerMng) ClearTraceUin() {
	plm.traceUin = make(map[string]bool)
}

func (plm *PlayerLoggerMng) AddTraceUin(uin string) {
	plm.traceUin[uin] = true
}
func (plm *PlayerLoggerMng) CheckUin(uin string) bool {
	if plm.traceAllUin {
		return true
	}
	if uin == "" {
		return false
	}
	return plm.traceUin[uin]
}

func (plm *PlayerLoggerMng) WriteMsg(uin string, format string, v ...interface{}) {
	if plm == nil || plm.mtLogger == nil {
		return
	}
	plm.mtLogger.WriteLog(uin, format, v...)
}

func (plm *PlayerLoggerMng) Close() {
	if plm == nil {
		return
	}
	plm.mtLogger.Close()
}

func SetTraceAllUid(b bool) {
	playerLoggerMng.traceAllUin = b
	Trace("set traceAllUin %t", b)
}

func ClearTraceUid() {
	playerLoggerMng.ClearTraceUin()
}

func AddTraceUid(uin string) {
	playerLoggerMng.AddTraceUin(uin)
}

func IsNeedTraceUid(uin string) bool {

	return CheckPrintLogLevel(LL_DEBUG) || playerLoggerMng.CheckUin(uin) || IsDebug()
}

func LogWx(uin string, format string, v ...interface{}) {
	if uin == ""{
		uin = "unknown"
	}
	if IsNeedTraceUid(uin) {
		playerLoggerMng.WriteMsg(uin, format, v...)
	}
	if CheckPrintLogLevel(LL_DEBUG) && serverLogger != nil {
		log := fmt.Sprintf(format, v...)
		msg := fmt.Sprintf("[WX:%s]"+log, uin)
		serverLogger.WriteMsg(LL_DEBUG, msg)
	}
}
func LogWxError(uin string, format string, v ...interface{}) {
	if uin == ""{
		uin = "unknown"
	}
	msg := fmt.Sprintf("WX:%s|", uin)
	msg = msg + fmt.Sprintf(""+format, v...)
	if errorLogger != nil {
		errorLogger.WriteMsg(LL_ERROR, msg)
	}
	if serverLogger != nil {
		serverLogger.WriteMsg(LL_ERROR, msg)
	}
	if !IsNeedTraceUid(uin) {
		return
	}
	playerLoggerMng.WriteMsg(uin, format, v...)
}
