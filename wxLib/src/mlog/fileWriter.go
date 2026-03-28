/*
   基础的写文件类
*/

package mlog

import (
	"bytes"

	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	DEFAULT_LOG_FILE_MAX_SIZE = 1024 * 1024 * 10 // 滚动日志的默认文件最大尺寸
	DEFAULT_LOG_FILE_MAX_LINE = 1000000          // 滚动日志的默认文件最大行数
	DEFALUT_LOG_FILE_MAX_NUM  = 5                // 滚动日志默认最大文件数
)

// 自定义日志级别，不同级别的值在不同的bit位上，方便进行级别控
const (
	LF_NONE   = 0x00000000 //!<没有级别，不打印
	LF_FATAL  = 0x00000100
	LF_ERROR  = 0x00000200
	LF_WARN   = 0x00000400
	LF_DEBUG  = 0x00000800
	LF_INFO   = 0x00001000
	LF_NOTICE = 0x00002000
	LF_ANY    = 0xFFFFFFFF //!<强制打印
)

const (
	DST_NONE = "none"
	DST_YMD  = "ymd"
	DST_YMDH = "ymdh"
)

// 本地日志接口类
type FileWriter struct {
	*log.Logger
	mw *MuxWriter

	FileName string `json:"filename"` // 日志文件名

	DateSuffixType string `json:"datesuffix"`

	Rotate bool `json:"rotate"` // 是否滚动日志

	MaxFileNum int `json:"maxfilenum"` //最大滚动文件个数

	MaxLines          int `json:"maxlines"` // 文件最大行数
	maxlines_curlines int

	MaxSize         int `json:"maxsize"` // 文件最大尺寸
	maxsize_cursize int

	Daily          bool `json:"daily"` // 是否按天滚动（该功能开启后，即使DateSuffixType配置成文件名固定不变，按天也会滚动文件）
	daily_opendate int

	MaxDays int64 `json:"maxdays"` // 日志最大保留天数

	startLock sync.Mutex // 文件检查锁。只有一个 log 可以写同一个文件

	Level int64 `json:"level"` // 日志等级
}

// an *os.File writer with locker.
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

// write to os.File.
func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

// set os.File in writer.
func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

// FileWriter构造函数
func NewFileWriter() *FileWriter {
	w := &FileWriter{
		FileName:       "",
		DateSuffixType: DST_NONE,
		MaxFileNum:     DEFALUT_LOG_FILE_MAX_NUM,
		MaxLines:       DEFAULT_LOG_FILE_MAX_LINE,
		MaxSize:        DEFAULT_LOG_FILE_MAX_SIZE,
		Daily:          true,
		MaxDays:        7,
		Rotate:         true,
		Level:          LF_ANY,
	}

	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)

	// set MuxWriter as FileLogger's io.Writer
	w.Logger = log.New(w.mw, "", 0)

	return w
}

// 初始化日志打印类
// json配置形如：
//	{
//	"filename":"logs/server.log",
//	"datesuffix":"none",
//	"maxfilenum":5,
//	"maxlines":10000,
//	"maxsize":1<<30,
//	"daily":true,
//	"maxdays":15,
//	"rotate":true
//	}
func (w *FileWriter) Init(jsonCfg string) error {
	if jsonCfg != "" {
		err := jsoniter.Unmarshal([]byte(jsonCfg), w)
		if err != nil {
			return err
		}
	}

	if len(w.FileName) == 0 {
		return errors.New("jsonCfg must have filename")
	}
	err := w.startLogger() // 打开文件，更新文件信息
	return err
}

// 打开文件，更新文件信息
func (w *FileWriter) startLogger() error {
	fd, err := w.openLogFile() // 1.打开文件
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)          // 2.关闭旧文件，设置新文件句柄
	return w.initFileInfo() // 3.更新文件信息
}

// 打开文件
func (w *FileWriter) openLogFile() (*os.File, error) {
	fName := w.getRealLogName(0)

	dir := filepath.Dir(fName)
	_ = os.MkdirAll(dir, 0777)

	fd, err := os.OpenFile(fName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err == nil {
		err = os.Chmod(fName, 0777)
	}
	return fd, err
}

func (w *FileWriter) initFileInfo() error {
	fd := w.mw.fd
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.maxsize_cursize = int(fInfo.Size())
	w.daily_opendate = time.Now().Day()
	w.maxlines_curlines = 0
	if fInfo.Size() > 0 {
		count, err := w.getCurlines()
		if err != nil {
			return err
		}
		w.maxlines_curlines = count
	}
	return nil
}

// 对外提供的开启指定级别log的接口
func (w *FileWriter) SetLogLevelOn(logPrintLevel int64) {
	w.Level |= logPrintLevel
}

// 对外提供的关闭指定级别log开关的接口
func (w *FileWriter) SetLogLevelOff(logPrintLevel int64) {
	w.Level &= ^logPrintLevel
}

// 对外提供的直接设置log级别的接口
func (w *FileWriter) SetLogLevel(logPrintLevel int64) {
	w.Level = logPrintLevel
}

// 检查日志level是否满足打印条件
func (w *FileWriter) checkPrintLogLevel(level int64) bool {
	if level&w.Level != 0 {
		return true
	}

	return false
}

// 打印日志
func (w *FileWriter) WriteMsg(level int64, msg string) error {
	if ok := w.checkPrintLogLevel(level); !ok {
		return nil
	}
	n := len(msg)
	if w.doCheck(n) == nil {
		w.Logger.Println(msg)
	}

	return nil
}

// 检查日志文件
// 如果满足日志滚动条件，则滚动日志、并打开新日志文件；否则清理句柄、重新打开日志文件
func (w *FileWriter) doCheck(size int) error {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if w.Rotate && ((w.MaxLines > 0 && w.maxlines_curlines >= w.MaxLines) ||
		(w.MaxSize > 0 && w.maxsize_cursize >= w.MaxSize) ||
		(w.Daily && time.Now().Day() != w.daily_opendate)) {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileWriter(%q): %s\n", w.FileName, err)
			return err
		}
	} else {
		fd, err := w.openLogFile() // 1.打开文件
		if err != nil {
			fmt.Fprintf(os.Stderr, "FileWriter(%q) openLogFile: %s\n", w.FileName, err)
			return err
		}
		w.mw.SetFd(fd) // 2.关闭旧文件，设置新文件句柄
	}
	w.maxlines_curlines++
	w.maxsize_cursize += size
	return nil
}

// 获取当前文件的行数
func (w *FileWriter) getCurlines() (int, error) {
	fName := w.getRealLogName(0)
	fd, err := os.Open(fName)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func (w *FileWriter) getRealLogName(logIdx int) string {
	fName := ""
	switch w.DateSuffixType {
	case DST_YMD:
		if logIdx > 0 {
			fName = w.FileName + "." + time.Now().Format("20060102") + "_" + strconv.Itoa(logIdx) + ".txt"
		} else {
			fName = w.FileName + "." + time.Now().Format("20060102") + ".txt"
		}
	case DST_YMDH:
		if logIdx > 0 {
			fName = w.FileName + "." + time.Now().Format("20060102_15") + "_" + strconv.Itoa(logIdx) + ".txt"
		} else {
			fName = w.FileName + "." + time.Now().Format("20060102_15") + ".txt"
		}
	case DST_NONE:
		fallthrough //必须明确添加fallthrough关键字，才会执行紧跟的下一句
	default:
		if logIdx > 0 {
			fName = w.FileName + "_" + strconv.Itoa(logIdx) + ".txt"
		} else {
			fName = w.FileName + ".txt"
		}
	}

	return fName
}

// 滚动日志，备份的文件形如 xx.log.2006-01-02_15#2
func (w *FileWriter) DoRotate() error {
	fName := w.getRealLogName(0)
	_, err := os.Lstat(fName)
	if err == nil { // 文件存在
		needRemoveFile := ""
		needRemoveFile = w.getRealLogName(w.MaxFileNum - 1)
		if _, err := os.Lstat(needRemoveFile); err == nil {
			err := os.Remove(needRemoveFile)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}

		// block FileLogger's io.Writer
		w.mw.Lock()
		defer w.mw.Unlock()

		w.mw.fd.Close()

		for i := w.MaxFileNum - 2; i >= 0; i-- {
			oldFileWriter := w.getRealLogName(i)
			if _, err := os.Lstat(oldFileWriter); err == nil {
				newFileWriter := w.getRealLogName(i + 1)
				err = os.Rename(oldFileWriter, newFileWriter)
				if err != nil {
					return fmt.Errorf("Rotate: %s\n", err)
				}
			}
		}

		// 更新文件后需要重新打开句柄
		err = w.startLogger()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s\n", err)
		}

	} else {
		// 文件也可能不存在，因为w.getRealLogName(0) 可能得到一个新的文件名，这个时候直接创建新文件即可
		err = w.startLogger()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s\n", err)
		}
	}

	go w.deleteOldLog()

	return nil
}

// 删除旧日志文件（超过MaxDays没有改动的文件）
func (w *FileWriter) deleteOldLog() {
	// 如果配置成0的话，则不删除旧文件
	if w.MaxDays <= 0 {
		return
	}

	fName := w.getRealLogName(0)
	dir := filepath.Dir(fName)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
				fmt.Println(returnErr)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.MaxDays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(fName)) {
				os.Remove(path)
			}
		}
		return
	})
}

// 关闭文件
func (w *FileWriter) Destroy() {
	w.mw.fd.Close()
}

// flush缓存到硬盘
func (w *FileWriter) Flush() {
	w.mw.fd.Sync()
}
