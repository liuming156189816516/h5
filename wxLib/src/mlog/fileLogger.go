/*
   本地日志类，支持并发写
*/

package mlog

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"
)

type logMsg struct {
	level int64
	msg   string
}

var logMsgPool *sync.Pool

func init() {
	logMsgPool = &sync.Pool{
		New: func() interface{} {
			return &logMsg{}
		},
	}
}

// 日志类
type FileLogger struct {
	loggerFuncCallDepth int          // 函数调用层数
	asynchronous        bool         // 并发写开关
	msg                 chan *logMsg // 并发接收消息的channel
	output              *FileWriter  // 写日志类
	isClosing           bool         // 是否在关闭中
	//closed              chan struct{}
	IsPrintGpid    bool //是否打印协程ID
	IsPrintMSecond bool //是否打印毫秒级别
	IsNeedReport   bool
}

// Logger类构造函数
func NewFileLogger(channelLen int64) *FileLogger {
	lp := new(FileLogger)
	lp.IsPrintGpid = false
	lp.IsPrintMSecond = false

	lp.loggerFuncCallDepth = 1 // 默认函数调用层数为2层：FileLogger.WriteMsg=>FileLogger.writeMsg，第2层才是 WriteMsg 被调用的文件位置和行数
	lp.asynchronous = false
	lp.msg = make(chan *logMsg, channelLen)
	lp.output = NewFileWriter()
	lp.isClosing = false
	//	lp.closed = make(chan struct{})

	return lp
}

func (lp *FileLogger) Init(fileName string, dateSuffix string, maxFileNum int, maxSize int, maxDays int, logLevel int64) error {
	format := "{\"filename\":\"%s\",\"datesuffix\":\"%s\",\"maxfilenum\":%d,\"maxsize\":%d,\"maxdays\":%d,\"level\":%d}"
	jsonConfig := fmt.Sprintf(format, fileName, dateSuffix, maxFileNum, maxSize, maxDays, logLevel)
	return lp.output.Init(jsonConfig)
}

// 设置为异步并发打印
func (lp *FileLogger) Async() {
	if lp.asynchronous == false {
		lp.asynchronous = true
		go lp.startLogger()
	}

}

// 开启并发logger协程类，通过channel从主线程获取msg并打印
// 当收到关闭请求后，将channel未打印的msg全部打印完，然后清理并退出
func (lp *FileLogger) startLogger() {
	if !lp.asynchronous {
		return
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for !lp.isClosing {
		select {
		case bm := <-lp.msg:
			_ = lp.output.WriteMsg(bm.level, bm.msg)
			logMsgPool.Put(bm)
		case <-ticker.C:
		}
	}
	for len(lp.msg) > 0 {
		bm := <-lp.msg
		_ = lp.output.WriteMsg(bm.level, bm.msg)
		logMsgPool.Put(bm)
	}
	lp.output.Flush()
	lp.output.Destroy()
}

func (lp *FileLogger) checkLogLevel(l int64) bool {
	if lp == nil {
		return false
	}
	return lp.output.checkPrintLogLevel(l)
}

// 设置日志头函数层数
func (lp *FileLogger) SetLogFuncCallDepth(d int) {
	lp.loggerFuncCallDepth = d
}
func (lp *FileLogger) SetPrintGpid(b bool) {
	lp.IsPrintGpid = b
}

// 是否打印毫秒级时间
func (lp *FileLogger) SetPrintMs(b bool) {
	lp.IsPrintMSecond = b
}

// 打印日志
func (lp *FileLogger) writerMsg(logLevel int64, msg string) error {

	if lp.asynchronous {
		lm := logMsgPool.Get().(*logMsg)
		lm.level = logLevel
		lm.msg = msg
		select {
		case lp.msg <- lm:
		default:
			//失败返还内存
			logMsgPool.Put(lm)
		}
	} else {
		_ = lp.output.WriteMsg(logLevel, msg)
	}
	return nil
}

// LF_ANY日志
func (lp *FileLogger) WriteMsg(level int64, format string, v ...interface{}) {
	head, ok := logLevelStr[level]
	var logStr string
	if ok {
		logStr = fmt.Sprintf("["+head+"]"+format, v...)
	} else {
		logStr = fmt.Sprintf(format, v...)
	}
	if lp.IsPrintGpid {
		logStr = "[Svrid:"+fmt.Sprintf("%d",serverID)+",Gpid:" + GetGpid() + "]" + logStr
	}
	var msg string
	var nowTime string
	if lp.IsPrintMSecond {
		nowTime = time.Now().Format("2006-01-02 15:04:05.000")
	} else {
		nowTime = time.Now().Format("2006-01-02 15:04:05")
	}
	if lp.loggerFuncCallDepth > 0 {
		_, file, line, ok := runtime.Caller(lp.loggerFuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, fileName := path.Split(file)
		msg = fmt.Sprintf("[%s]|[%s:%d]|%s", nowTime, fileName, line, logStr)
	} else {
		msg = fmt.Sprintf("[%s]|%s", nowTime, logStr)
	}

	lp.writerMsg(level, msg)
	if lp.IsNeedReport && funcReport != nil {
		funcReport(head, "", msg)
	}
}

// LF_ANY日志
//func (lp *FileLogger) WriteJson(v interface{}) {
//	msg, _ := jsoniter.MarshalToString(v)
//	lp.writerMsg(LL_ANY, msg)
//
//}

func (lp *FileLogger) Flush() {
	lp.output.Flush()
}

func (lp *FileLogger) Close() {
	if lp.asynchronous {
		lp.isClosing = true
		//<-lp.closed
		lp.writerMsg(0, "Server Tell me To Close")
	}
}
