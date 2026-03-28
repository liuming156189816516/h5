// goServer.go 框架主逻辑，包括注册Rpc/Seed、Run等接口
package mframe

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"mlog"

	"runtime"
	"strings"
	"sync"
	"time"
	"utils"
	"utils/baselib"
)

const DEFAULT_RPC_REQUEST_SECONDS = 10

const MYSQL_ENCRY = 1

var G_CloseChan chan struct{}
var g_deferFuncs []func()

var _debug_mode = false
var once = &sync.Once{}

func init() {

	fmt.Printf("\n+++++===============+++++\n")

	for i := 0; i < len(os.Args); i++ {
		//	fmt.Printf("flag.Arg(%d) = %s\n", i, flag.Arg(i))
		if strings.ToLower(os.Args[i]) == "-d" || strings.ToLower(os.Args[i]) == "dbg" {
			_debug_mode = true
			break
		}
	}

	fmt.Printf("_debug_mode = %t\n", _debug_mode)
	fmt.Println("+++++===============+++++\n")
}

func IsDebug() bool {
	return _debug_mode
}
func IsTool(t string) bool {
	for i := 0; i < len(os.Args); i++ {
		//	fmt.Printf("flag.Arg(%d) = %s\n", i, flag.Arg(i))
		if strings.ToLower(os.Args[i]) == t {
			return true
		}
	}
	return false

}

func printMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	mlog.Trace("PrintMemStats ===================Begin================")
	mlog.Trace("mem.Alloc: %d, %d Mb", mem.Alloc, mem.Alloc/1024/1024)
	mlog.Trace("mem.TotalAlloc: %d, %d Mb", mem.TotalAlloc, mem.TotalAlloc/1024/1024)
	mlog.Trace("mem.HeapAlloc: %d, %d Mb", mem.HeapAlloc, mem.HeapAlloc/1024/1024)
	mlog.Trace("mem.HeapSys: %d, %d Mb", mem.HeapSys, mem.HeapSys/1024/1024)

	mlog.Trace("MemStats: %+v", &mem)
	mlog.Trace("PrintMemStats ===================End================")
}
func ProfileToFile(file string, profileTime int) {
	mlog.Trace("Start CPUProfile")
	if err := baselib.ProfileToFile(file+"_cpu", profileTime); err != nil {
		mlog.Trace("Create file %s failed, error: %+v", file+"_cpu", err)
	}
	mlog.Trace("Finish CPUProfile")
	mlog.Trace("Start Memory Profile")
	if err := baselib.ProfileHeapMemoryToFile(file + "_memory"); err != nil {
		mlog.Trace("Create file %s failed, error: %+v", file+"_memory", err)
	}
	mlog.Trace("Finish Memory Profile")
}
func Defer(f func()) {

	g_deferFuncs = append(g_deferFuncs, f)
	mlog.Trace("append Defer %+v", g_deferFuncs)
}

func onStopServer([]byte) {
	baselib.SendStopSignal()
}
func Start(svrName string, opt ...*FrameOption) {

	if err := InitConfig(svrName, opt...); err != nil {
		log.Fatalf("InitConfig error:%s", err.Error())
		return
	}

	mlog.Trace(">>>>>>>>>>>>>>> start  ================")
	runinfo := fmt.Sprintf("%d %s %d %s ... start\n", os.Getpid(),
		GetServerName(), GetServerID(), time.Now().Format(utils.TimeFormat_Full))
	//fmt.Printf("================Start Server %s as id %d================\n", GetServerName(), GetServerID())
	ioutil.WriteFile("run.id", []byte(runinfo), 0666)

	var startTime time.Time
	startTime = time.Now()
	mlog.Trace(">>>>>>>>>>>>>>> Start Run %s-%d at %s ================", GetServerName(), GetServerID(), startTime.Format(baselib.TimeFormatMilli))

	go func() {
		sig := baselib.NewSignalHandler()
		timewait := time.NewTicker((200 * time.Millisecond))
		defer timewait.Stop()

		for {
			select {
			case <-sig.ReloadSignal():
				fmt.Println("Receive Reload Signal, Reloading")
				mlog.Trace("Receive Reload Signal, Reloading")

			case <-sig.ProfSignal():
				fmt.Println("Receive Prof Signal, Start Profile")
				mlog.Trace("Receive Prof Signal, Start Profile")
				printMemStats()
				go ProfileToFile("profile", 20)
			case <-timewait.C:
			}
		}
	}()

	return
}

func stop() {
	startTime := time.Now()
	mlog.Trace(">>>>>>>>>>>>>>> Start Stop %s-%d at %s ================", GetServerName(), GetServerID(), startTime.Format(baselib.TimeFormatMilli))
	for _, f := range g_deferFuncs {
		if f != nil {
			f()
		}
	}
	//// 必须按顺序来，先将nats的订阅全部取消, 使得新请求不会再过来
	//UnloadAllSubject()
	//stopRpcxServer()

	//log.Printf("UnloadAllSubject\n")
	//延迟一秒, 处理完所有的消息
	time.Sleep(1000 * time.Millisecond)
	// 然后处理剩下的请求之后，停止Service的协程
	//StopService()
	log.Printf("Server %s-%d stop!!!", GetServerName(), GetServerID())
	mlog.Trace("Server %s-%d stop!!!", GetServerName(), GetServerID())
	// 最后Flush日志缓存
	mlog.CloseAllLogAdapter() // 注册清理操作，将并发管道中还没有打印的日志全部打印完才退出
	log.Printf("CloseAllLogAdapter\n")
	//mlog.CloseAllSyslogAdapter()
	buf, _ := ioutil.ReadFile("run.id")
	runinfo := fmt.Sprintf("%s\n%d %s %d %s ... stop\n", buf, os.Getpid(),
		GetServerName(), GetServerID(), time.Now().Format(utils.TimeFormat_Full))
	ioutil.WriteFile("run.id", []byte(runinfo), os.ModeAppend)
	now := time.Now()
	mlog.Trace(">>>>>>>>>>>>>>> Succ Stop %s-%d at %s,cost: %d ms ================",
		GetServerName(), GetServerID(), now.Format(baselib.TimeFormatMilli), now.Sub(startTime)/time.Millisecond)

}

func Stop() {
	//once.Do(stop)
}
