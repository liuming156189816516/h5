package mframe

import (
	"errors"
	"fmt"
	"mframe/stat"
	"mlog"
	"os"

	"runtime"
	"strconv"
	"strings"
	"time"
	"utils/baselib"

	"github.com/BurntSushi/toml"
)

const (
	LL_NONE = mlog.LL_NONE
	LL_RUN  = mlog.LL_EMERG | mlog.LL_ALERT | mlog.LL_CRIT | mlog.LL_ERROR | mlog.LL_NOTICE | mlog.LL_INFO
	LL_DBG  = mlog.LL_ANY
)

type FrameOption struct {
	Svrid        int
	LogPath string
	OnLoadConfig func()
}

var defFrameOption = &FrameOption{}

type LogConfig struct {
	logLevel    int64
	Level       string
	PrintGid    bool //日志是否打印协程ID
	PrintMS     bool //日志是否打印毫秒级时间
	TraceAllUid bool
}

type RunConfig struct {
	ServerName string
	ServerID   int32
}

type FrameConfig struct {
	DisableCheckService bool
	RpcCallTimeout      int32
	//UseRpcx             bool
	//LogPath             string
	Log LogConfig
}

var config = NewDefaultFrameConfig()
var server_config = RunConfig{ServerName: "goServer", ServerID: 0}

func GetServerConfig() *RunConfig {
	return &server_config
}
func GetServerName() string {
	return server_config.ServerName
}
func GetServerID() int32 {
	return server_config.ServerID
}

func SetServerID(sid int32) {
	 server_config.ServerID = sid
	 mlog.SetServerId(sid)
}

func GetCallTimeout() int32 {
	if config.RpcCallTimeout == 0 {
		config.RpcCallTimeout = DEFAULT_RPC_REQUEST_SECONDS
	}
	return config.RpcCallTimeout
}

func LogLevel() int64 {
	return config.Log.logLevel
}



func NewDefaultFrameConfig() *FrameConfig {
	return &FrameConfig{Log: LogConfig{logLevel: LL_RUN}}
}

func InitConfig(svrName string, opt ...*FrameOption) error {
	// 初始化
	if len(opt) > 0 {
		defFrameOption = opt[0]
	}
	cpun := runtime.NumCPU()
	if cpun >= 16 {
		cpun = 16
		runtime.GOMAXPROCS(cpun)
	}

	// 加载框架配置
	server_config.ServerName = svrName
	if defFrameOption.Svrid <= 0 {
		// path := filepath.Base(os.Args[0])
		// arr := strings.Split(path, ".")
		//
		// if len(arr) > 1 {
		// 	sid, err := strconv.Atoi(arr[len(arr)-1])
		// 	if err == nil {
		// 		server_config.ServerID = int32(sid)
		// 	}
		// }

		if len(os.Args) > 1 {
			sid, err := strconv.Atoi(os.Args[1])
			if err == nil {
				server_config.ServerID = int32(sid)
			}
		}
	} else {
		server_config.ServerID = int32(defFrameOption.Svrid)
	}

	//logSN := server_config.ServerName
	//
	//if strings.Index(path, "_dbg") > 0 {
	//	logSN = path
	//
	//}
	sleepTime := time.Millisecond * 20
	if defFrameOption.LogPath == ""{
		defFrameOption.LogPath = "./mlog"
	}
	mlog.InitServer(server_config.ServerName, server_config.ServerID, defFrameOption.LogPath, nil)

	// 加载启动配置
	if err := LoadBootConfig(); err != nil {
		mlog.Error("InitConfig error:%s", err)
		mlog.Error("Run exit")
		//等待一会, 让日志打印出去
		time.Sleep(sleepTime)
		return err
	}

	if err := LoadFrameConfig(); err != nil {
		//等待一会, 让日志打印出去
		time.Sleep(sleepTime)
		return err
	}
	baselib.RegisterReloadFunc(LoadFrameConfig)
	mlog.Trace("Init Server GOMAXPROCS:%d, NumCpu:%d", cpun, runtime.NumCPU())
	mlog.Trace("Init Server %+v, FrameConfig:%+v", server_config, config)

	mlog.SetServerId(GetServerID())

	stat.SetAdditionMsgReport(additionalMsgStat)
	return nil
}

func LoadBootConfig() error {

	return nil
}
func _loadFrameConfig() (*FrameConfig, error) {
	newConf := NewDefaultFrameConfig()

	filename := "./conf/frame.toml"
	_, err := toml.DecodeFile(filename, newConf)
	if err == nil {
		mlog.Trace("load LocalFile %s config: %+v", filename, newConf)
		return newConf, err
	}
	if !os.IsNotExist(err) {
		mlog.Error("DecodeFile:%s failed:%s", filename, err.Error())
	}

	return newConf, err
}

func LoadFrameConfig() error {

	newConf, _ := _loadFrameConfig()
	//mlog.Trace("Load Config:%+v from: %s", newConf, filename)
	if newConf.RpcCallTimeout <= 0 {
		newConf.RpcCallTimeout = DEFAULT_RPC_REQUEST_SECONDS
	}

	config = newConf
	LoadLogConfig(&config.Log)
	if defFrameOption.OnLoadConfig != nil {
		defFrameOption.OnLoadConfig()
	}
	return nil
}
func LoadSystemConfig() error {

	return nil
}
func LoadLogConfig(Log *LogConfig) {
	switch strings.ToLower(Log.Level) {
	case "dbg", "DBG":
		Log.logLevel = LL_DBG
	default:
		Log.logLevel = LL_RUN
	}

	mlog.SetLogConfig(Log.logLevel, Log.PrintGid, Log.PrintMS)
	// 加载UinDebug配置
	LoadUinConfig(Log)

}

type UidTraceCfg struct {
	Uids []uint64
}

func LoadUinConfig(Log *LogConfig) error {
	mlog.ClearTraceUid()
	mlog.SetTraceAllUid(Log.TraceAllUid)
	if Log.TraceAllUid {
		mlog.Trace("Log.TraceAllUid")
		return nil
}
//	newConf := &UidTraceCfg{}
//	fkey := "trace_uids.toml"
//
//	filename := fmt.Sprintf("../LocalConfig/%s", fkey)
//	_, err := toml.DecodeFile(filename, newConf)
//	if err == nil {
//		for _, u := range newConf.Uids {
//			if u > 10000 {
//				mlog.AddTraceUid(uint64(u))
//				mlog.Trace("AddTraceUid %d Local", u)
//			}
//		}
//		return nil
//	}
//
	return nil
}

func LoadServerLocalConfig(newConf interface{}, args ...string) error {
	if newConf == nil {
		return errors.New("No Parameter")
	}

	sname := GetServerName()

	filename := fmt.Sprintf("../LocalConfig/%s.toml", sname)
	_, err := toml.DecodeFile(filename, newConf)
	if err == nil {
		mlog.Trace("load LocalFile %s config: %+v", filename, newConf)
		return err
	}
	if !os.IsNotExist(err) {
		mlog.Error("DecodeFile:%s failed:%s", filename, err.Error())
	}

	/*
		//===============从mysql读取==================
		filename = fmt.Sprintf("%s_%d.toml", sname, GetServerID())
		err = LoadConfigFromMysql(filename, newConf)
		if err == nil {
			mlog.Trace("load Mysql %s config: %+v", filename, newConf)
			return err
		}

		filename = fmt.Sprintf("%s.toml", sname)
		err = LoadConfigFromMysql(filename, newConf)
		if err == nil {
			mlog.Trace("load Mysql %s config: %+v", filename, newConf)
			return err
		}
		if len(args) > 0 {
			sname = args[0]
			filename = fmt.Sprintf("%s.toml", sname)
			err = LoadConfigFromMysql(filename, newConf)
			if err == nil {
				mlog.Trace("load Mysql %s config: %+v", filename, newConf)
				return err
			}
		}
	*/
	//err = errors.Errorf("no config for  %s_%d.toml", GetServerName(), GetServerID()))
	mlog.Trace("load Mysql %s failed: %+v", filename, err.Error())
	//mlog.Trace("Load Config:%+v from: %s", newConf, filename)
	//tempConfig := &BaseLocalConfig{}
	//_, err = toml.DecodeFile(filename, tempConfig)
	//if tempConfig.Log != nil {
	//	mlog.Trace("Load LogConfig OnApp:%+v", tempConfig.Log)
	//	LoadLogConfig(tempConfig.Log)
	//}

	return err
}
