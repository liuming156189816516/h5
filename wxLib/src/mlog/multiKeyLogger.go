package mlog

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type keyLogMsg struct {
	key string
	msg string
}
type writerItem struct {
	actTime time.Time
	*FileWriter
}

func (wi *writerItem) isExpired() bool {
	return time.Now().Sub(wi.actTime) > 30*time.Second
}
func (wi *writerItem) active() {
	wi.actTime = time.Now()
}

type TMultiLogger struct {
	msgList       chan *keyLogMsg // 并发接收消息的channel
	isClosing     bool
	chanClosed    chan struct{}
	mapOutput     map[string]*writerItem // 写日志类
	CacheKeyCount uint32
	FileNameHead  string
	FileNameTail  string

	FuncCallDepth int
}

func NewMultiLogger(logPath string, keyCount uint32, callDepth int, fileTail string) *TMultiLogger {
	l := &TMultiLogger{}
	l.FileNameHead = logPath
	l.FileNameTail = fileTail
	l.CacheKeyCount = keyCount
	l.FuncCallDepth = callDepth

	l.msgList = make(chan *keyLogMsg, channelSize)
	l.chanClosed = make(chan struct{})
	l.mapOutput = make(map[string]*writerItem)
	l.isClosing = false

	return l
}
func (ml *TMultiLogger) startLogger() {
	//defer close(ml.chanClosed)
	//defer close(ml.msgList)
	t := time.Now()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for !ml.isClosing {
		select {
		case bm := <-ml.msgList:
			ml.outputMsg(bm)

		case <-ticker.C:
			if time.Now().Sub(t) > 3*time.Minute {
				//	WriteBill("log_cache", "%s-%d|%s|%d", serverName, serverID, ml.FileNameHead, len(ml.mapOutput))
				ml.clearWriter(true)
				t = time.Now()
			} else if len(ml.mapOutput) > int(ml.CacheKeyCount) {
				ml.clearWriter(false)
				//t = time.Now()
			}
		}
	}
	for len(ml.msgList) > 0 {
		bm := <-ml.msgList
		ml.outputMsg(bm)
	}
}
func (ml *TMultiLogger) StartLogger() {
	if ml == nil {
		return
	}
	go ml.startLogger()
}
func (ml *TMultiLogger) Close() {
	if ml == nil {
		return
	}
	if ml.isClosing == true {
		return
	}
	ml.isClosing = true
	//<-ml.chanClosed
	logMsg := &keyLogMsg{msg: "Server Tell me To Close"}
	ml.msgList <- logMsg
	Trace("close logger %s", ml.FileNameHead)
}

func (ml *TMultiLogger) clearWriter(force bool) {
	if ml == nil {
		return
	}
	if !force && len(ml.mapOutput) < int(ml.CacheKeyCount) {
		return
	}
	n := 0
	for k, v := range ml.mapOutput {
		if force || v.isExpired() {
			v.Flush()
			v.Destroy()
			delete(ml.mapOutput, k)
			n++
		}
	}
	if n > 0 {
		Debug("clear %d  log hander, CacheKeyCount: %d, used: %d", n, ml.CacheKeyCount, len(ml.mapOutput))
	}

}
func (ml *TMultiLogger) getWriter(key string) *writerItem {
	if w, ok := ml.mapOutput[key]; ok {
		return w
	}
	ml.clearWriter(false)
	it := &writerItem{}
	it.active()

	it.FileWriter = NewFileWriter()
	it.FileName = fmt.Sprintf(ml.FileNameHead+"%s"+ml.FileNameTail, key)
	it.DateSuffixType = "ymd"

	it.MaxFileNum = maxFileNum
	it.MaxSize = maxSize
	it.MaxDays = maxDays
	it.Level = LF_ANY

	if err := it.Init(""); err != nil {
		return nil
	}
	ml.mapOutput[key] = it
	return it
}

func (ml *TMultiLogger) outputMsg(logMsg *keyLogMsg) {
	if ml == nil || logMsg == nil {
		return
	}

	it := ml.getWriter(logMsg.key)
	if it == nil {
		return
	}
	it.active()
	it.WriteMsg(LF_ANY, logMsg.msg)
}

func (ml *TMultiLogger) WriteLog(key string, format string, v ...interface{}) {
	if ml == nil || ml.isClosing {
		return
	}
	if len(ml.msgList) == cap(ml.msgList) {
		return
	}

	logMsg := &keyLogMsg{}
	logMsg.key = key
	logStr := fmt.Sprintf(format, v...)
	if g_isPrintGpid {
		logStr = "[Svrid:"+fmt.Sprintf("%d",serverID)+",Gpid:" + GetGpid() + "]" + logStr
	}
	var nowTime string
	if g_isPrintMS {
		nowTime = time.Now().Format("2006-01-02 15:04:05.000")
	} else {
		nowTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if ml.FuncCallDepth > 0 {
		_, file, line, ok := runtime.Caller(ml.FuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, fileName := path.Split(file)
		logMsg.msg = fmt.Sprintf("[%s]|[%s:%d]|%s", nowTime, fileName, line, logStr)
	} else {
		logMsg.msg = fmt.Sprintf("[%s]|%s", nowTime, logStr)
	}
	select {
	case ml.msgList <- logMsg:
	default:

	}

}
