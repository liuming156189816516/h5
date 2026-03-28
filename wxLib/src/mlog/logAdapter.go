package mlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"utils/baselib"
	"utils/baselib/crypto"
	"github.com/astaxie/beego"
)

const (
	maxFileNum  = 100
	maxSize     = 100 * 1024 * 1024
	maxDays     = 3
	channelSize = 10 * 1024
)

const logDateSuffix = "ymd"
const billDateSuffix = "ymd"

var serverLogger *ServerLogger
var errorLogger *ErrorLogger
var smsLogger *SmsLogger

// 线程安全的并发 logger 池（PlayerLogger 的并发池在 PlayerLoggerMng 内部）
var loggerMap = new(baselib.Map)     // 并发 系统logger池，key是 md5(filename)
var billLoggerMap = new(baselib.Map) // 并发 BillLogger池，key是 md5(tag)

var serverName = "Server"
var serverID = int32(1)
var serverAll = ""
var logBasePath = "./mlog"
var logPath = ""
var g_LogLevel = int64(LL_ANY)
var g_isPrintGpid bool = false
var g_isPrintMS bool = false

//==========1.所有的 LogAdapter 的初始化逻辑（包加载时自动执行）==========

var funcReport func(string, string, string)

var _debug_mode = false

func init() {
	for i := 0; i < len(os.Args); i++ {
		//	fmt.Printf("flag.Arg(%d) = %s\n", i, flag.Arg(i))
		if strings.ToLower(os.Args[i]) == "-d" || strings.ToLower(os.Args[i]) == "dbg" {
			_debug_mode = true
			break
		}
	}

}

func IsDebug() bool {
	return _debug_mode
}

func initLogBasePath(lp string) (string, string) {
	if lp == "" {
		return "..", serverName
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		return os.Args[0], serverName
	}
	dir = strings.Replace(dir, "\\", "/", -1) //将\替换成/
	dir = strings.Replace(dir, "/pack/", "/", -1) //将\替换成/
	if strings.HasSuffix(dir, "/bin") {
		dir = dir[:len(dir)-len("/bin")]
	}
	s := strings.LastIndex(dir, "/")
	if s < 0 {
		s = 0
	}
	svr := dir[s+1:]
	if strings.HasPrefix(dir, "/home/") {
		dir = dir[len("/home/"):]
	}

	dir = strings.Replace(dir, "/", "_", -1) //将\替换成/
	if strings.HasSuffix(lp, "/") {
		lp = lp[0 : len(lp)-1]
	}
	dir = lp + "/" + dir

	return dir, svr
}

func InitServer(sn string, sid int32, lp string, rfunc func(string, string, string)) {
	serverName = sn

	serverID = sid
	serverAll = fmt.Sprintf("_%s_%d", sn, sid)
	logBasePath = lp
	logPath, serverAll = initLogBasePath(lp)
	serverAll = "_" + serverAll
	//if !strings.HasSuffix(logPath, "/") {
	//	logPath = logPath + "/"
	//}

	funcReport = rfunc
	// 1.初始化各种等级的系统logger类
	serverLogger = GetServerLogger()
	errorLogger = GetErrorLogger()
	smsLogger = GetSmsLogger()

	// 2.不同tag的billLogger是在TRACEBILL的时候才生成的，所以不用在这里初始化

	// 3.初始化用户追踪logger类
	playerLoggerMng = newPlayerLoggerMng()
	Trace("================Start Server %s as id %d================", serverName, serverID)
}

func SetServerId(sid int32) {
	if sid == serverID {
		return
	}
	serverID = sid
	//从新加载log文件

	oldLogger := serverLogger
	serverLogger = CreateServerLogger()
	if oldLogger != nil {
		oldLogger.Close()
		oldLogger = nil
	}

}

//==========2.所有的 LogAdapter 的清理逻辑（业务逻辑手动调用）==========
func CloseAllLogAdapter() {
	playerLoggerMng.Close()
	loggerMap.Range(func(key, value interface{}) bool {
		if logger, ok := value.(*FileLogger); ok {
			logger.Close()
		}
		return true
	})
	loggerMap = new(baselib.Map)

	billLoggerMap.Range(func(key, value interface{}) bool {
		if billLogger, ok := value.(*BillLogger); ok {
			billLogger.Close()
		}
		return true
	})
	//billLoggerMap = new(baselib.Map)

}
func GetLogBasePath() string {
	mlogPath :=beego.AppConfig.String("mlogPath")
	if mlogPath == "" {
		return logBasePath
	}
	return mlogPath
}

//==========3.所有的 LogAdapter 的定义==========

// 从Logger集合中获取负责file的Logger
// 如果集合中没有，则创建对应的Logger对象并返回
// 【注意】：该函数不适用于bill的Logger的快速初始化，因为同一个bill的Logger的文件名会随着时间改变，导致map的key改变
func GetLogger(fileName string, dateSuffix string, maxFileNum int, maxFileSize int, maxDays int, channelSize int, logLevel int64) (*FileLogger, error) {
	var log interface{}
	var ok bool

	fileName = GetLogBasePath() + fileName

	key := crypto.Md5Str([]byte(fileName))

	if log, ok = loggerMap.Load(key); !ok {
		if log, ok = loggerMap.Load(key); !ok {
			logger := NewFileLogger(int64(channelSize))
			_ = logger.Init(fileName, dateSuffix, maxFileNum, maxFileSize, maxDays, logLevel)
			logger.Async()
			log = logger
			loggerMap.Store(key, log)
		}
	}

	if l, ok := log.(*FileLogger); ok {
		l.SetLogFuncCallDepth(2) // LogWarn=>FileLogger.WriteMsg=>FileLogger.writeMsg，第3层才是SMSLOG被调用的文件位置和行数
		l.SetPrintGpid(g_isPrintGpid)
		l.SetPrintMs(g_isPrintMS)
		return l, nil
	} else {
		return nil, fmt.Errorf("key:%s, value not a FileLogger", fileName)
	}
}
func RemoveLogger(fileName string) {
	loggerMap.Delete(fileName)
}
func IsPrintGID() bool {
	return g_isPrintGpid || IsDebug()
}
func IsPrintMS() bool {
	return g_isPrintMS || IsDebug()
}
func GetLogLevel() int64 {
	return g_LogLevel
}

func SetLogConfig(logPrintLevel int64, isPrintGpid bool, isPrintMS bool) {
	g_LogLevel = logPrintLevel
	g_isPrintGpid = isPrintGpid
	g_isPrintMS = isPrintMS
	if serverLogger != nil {
		serverLogger.SetPrintGpid(g_isPrintGpid)
		serverLogger.SetPrintMs(g_isPrintMS)
	}
	if errorLogger != nil {
		errorLogger.SetPrintGpid(g_isPrintGpid)
		errorLogger.SetPrintMs(g_isPrintMS)
	}
	if smsLogger != nil {
		smsLogger.SetPrintGpid(g_isPrintGpid)
		smsLogger.SetPrintMs(g_isPrintMS)
	}
}
func CheckPrintLogLevel(level int64) bool {
	if level&g_LogLevel != 0 {
		return true
	}
	if IsDebug() {
		return true
	}
	return false
}
func CheckLogDebug() bool {
	return CheckPrintLogLevel(LL_DEBUG)
}

type ServerLogger struct {
	*FileLogger
}

func GetServerLogger() *ServerLogger {
	if serverLogger != nil {
		return serverLogger
	}
	serverLogger = CreateServerLogger()
	return serverLogger
}

func CreateServerLogger() *ServerLogger {
	slogger := new(ServerLogger)
	var err error

	//path := fmt.Sprintf("/log/%s_%d", serverName, serverID)
	path := fmt.Sprintf("/%s", serverName)
	if slogger.FileLogger, err = GetLogger(path, logDateSuffix, maxFileNum, maxSize, maxDays, channelSize, LF_ANY); err == nil {
		return slogger
	} else {
		return nil
	}
}
func Trace(format string, v ...interface{}) {
	if serverLogger != nil {
		msg := fmt.Sprintf("[%s.%d]|", serverName, serverID)
		msg = msg + fmt.Sprintf(format, v...)

		serverLogger.WriteMsg(LL_ANY, msg)
	}
}
func Debug(format string, v ...interface{}) {
	if serverLogger != nil && CheckPrintLogLevel(LL_DEBUG) {
		msg := fmt.Sprintf("[%s.%d]|", serverName, serverID)
		msg = msg + fmt.Sprintf(format, v...)
		serverLogger.WriteMsg(LL_DEBUG, msg)
	}
}
func LogInfo(format string, v ...interface{}) {
	if serverLogger != nil && CheckPrintLogLevel(LL_INFO) {
		msg := fmt.Sprintf("[%s.%d]|", serverName, serverID)
		msg = msg + fmt.Sprintf(format, v...)
		serverLogger.WriteMsg(LL_INFO, msg)
	}
}

type ErrorLogger struct {
	*FileLogger
}

func GetErrorLogger() *ErrorLogger {
	if errorLogger != nil {
		return errorLogger
	}
	errorLogger = new(ErrorLogger)
	var err error

	if errorLogger.FileLogger, err = GetLogger("/error", logDateSuffix, maxFileNum, maxSize, maxDays, channelSize, LF_ANY); err == nil {
		errorLogger.FileLogger.IsNeedReport = true
		return errorLogger
	} else {
		return nil
	}
}
func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		msg := fmt.Sprintf("[%s.%d]|", serverName, serverID)
		msg = msg + fmt.Sprintf(format, v...)
		errorLogger.WriteMsg(LL_ERROR, msg)
	}
	if serverLogger != nil {
		serverLogger.WriteMsg(LL_ERROR, format, v...)
	}
}

type SmsLogger struct {
	*FileLogger
}

func GetSmsLogger() *SmsLogger {
	if smsLogger != nil {
		return smsLogger
	}
	smsLogger = new(SmsLogger)
	var err error
	if smsLogger.FileLogger, err = GetLogger("/warn", logDateSuffix, maxFileNum, maxSize, maxDays, channelSize, LF_ANY); err == nil {
		errorLogger.FileLogger.IsNeedReport = true
		return smsLogger
	} else {
		return nil
	}
}
func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[%s.%d]|", serverName, serverID)
	msg = msg + fmt.Sprintf(format, v...)
	if smsLogger != nil {
		smsLogger.WriteMsg(LL_WARN, msg)
	}
	if errorLogger != nil {
		errorLogger.WriteMsg(LL_WARN, msg)
	}
	if serverLogger != nil {
		serverLogger.WriteMsg(LL_WARN, msg)
	}
}

type BillLogger struct {
	*FileLogger
	curMonth time.Month
	tag      string
}

func GetBillLogger(tag string) (*BillLogger, error) {
	var log interface{}
	var ok bool

	key := crypto.Md5Str([]byte(tag))

	if log, ok = billLoggerMap.Load(key); !ok {
		if log, ok = billLoggerMap.Load(key); !ok {
			log = newBillLogger(tag)
			billLoggerMap.Store(key, log)
		}
	}

	if l, ok := log.(*BillLogger); ok {
		l.SetPrintGpid(g_isPrintGpid)
		l.SetPrintMs(g_isPrintMS)
		return l, nil
	} else {
		return nil, fmt.Errorf("key:%s, value not a BillLogger", tag)
	}
}

func newBillLogger(tag string) *BillLogger {
	bl := new(BillLogger)

	bl.tag = tag
	nowTime := time.Now()
	bl.curMonth = nowTime.Month()
	billFileName := fmt.Sprintf("%s/bills/%s/%s", GetLogBasePath(), nowTime.Format("2006-01"), bl.tag)

	bl.FileLogger = NewFileLogger(channelSize)
	bl.SetLogFuncCallDepth(0)
	bl.SetPrintGpid(false)
	bl.SetPrintMs(false)
	_ = bl.FileLogger.Init(billFileName, billDateSuffix, maxFileNum, maxSize, maxDays, LF_ANY)
	bl.FileLogger.Async()

	return bl
}

func (bl *BillLogger) CheckBillPath() {
	nowTime := time.Now()
	month := nowTime.Month()
	if month != bl.curMonth {
		billFileName := fmt.Sprintf("%s/bills/%s/%s", GetLogBasePath(), nowTime.Format("2006-01"), bl.tag)
		_ = bl.FileLogger.Init(billFileName, billDateSuffix, maxFileNum, maxSize, maxDays, LF_ANY)
		bl.curMonth = month
	}
}

func (bl *BillLogger) CheckAlarmPath() {
	if bl.curMonth > 0 {
		billFileName := fmt.Sprintf("%s/alarm/%s", GetLogBasePath(), bl.tag)
		_ = bl.FileLogger.Init(billFileName, billDateSuffix, maxFileNum, maxSize, maxDays, LF_ANY)
		bl.curMonth = 0
	}
}

func WriteBill(tag string, format string, v ...interface{}) {
	if billLogger, err := GetBillLogger(tag); err == nil {
		billLogger.CheckBillPath()
		billLogger.WriteMsg(LL_ANY, format, v...)
	}
}

//func WriteJsonBill(tag string, v interface{}) {
//	if v == nil {
//		return
//	}
//	if billLogger, err := GetBillLogger(tag); err == nil {
//		billLogger.CheckBillPath()
//		billLogger.WriteJson(v)
//	}
//}

func WriteAlarm(tag string, format string, v ...interface{}) {
	if billLogger, err := GetBillLogger(tag); err == nil {
		billLogger.IsPrintGpid = false
		billLogger.CheckAlarmPath()
		billLogger.WriteMsg(LL_ANY, format, v...)
	}
}

type PlayerLogger struct {
	*FileLogger
	uin uint64
}
