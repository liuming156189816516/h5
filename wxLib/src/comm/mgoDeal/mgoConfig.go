package mgoDeal

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
	"net"
	"time"
)

//"user:pass@localhost:27017"
type MgoCfg struct {
	Url      []string //localhost:27017
	Query    string
	User     string
	Password string
	PoolNum  int
	Ssl      bool
	SyncTask bool
}

var mgo_config = &MgoCfg{}

var g_MongoDB *mgo.Session = nil     // 读
var g_RealMongoDB *mgo.Session = nil // 主库 只写

func StartMgoDb(cfg *MgoCfg) error {
	mgo_config = cfg
	return InitMgoSession()
}

func (mcfg *MgoCfg) GetUrl() string {
	add := ""
	for id, url := range mgo_config.Url {
		if id == 0 {
			add = url
		} else {
			add += fmt.Sprintf(",%s", url)
		}
	}
	url := fmt.Sprintf("%s:%s@%s", mcfg.User, mcfg.Password, add)
	if mcfg.Query != "" {
		url = url + "?" + mgo_config.Query
	}
	return url
}
func GetMgoCoinfig() *MgoCfg {
	return mgo_config
}

func GetMongoRealConnUrl() string {
	add := ""
	for id, url := range mgo_config.Url {
		if id == 0 {
			add = url
		} else {
			add += fmt.Sprintf(",%s", url)
		}
	}
	url := fmt.Sprintf("%s:%s@%s", mgo_config.User, mgo_config.Password, add)
	if mgo_config.Query != "" {
		url = url + "?" + mgo_config.Query
	}
	return url
}

func Dial(cfg *MgoCfg) (*mgo.Session, error) {
	url := cfg.GetUrl()
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	dailInfo, err := mgo.ParseURL(url)
	if err != nil {
		return nil, err
	}

	if cfg.Ssl {
		dailInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsconfig)
			if err != nil {
				logs.Error("tls.Dial(tcp, %s, %+v, failed:%s", addr.String(), tlsconfig, err.Error())
			}
			return conn, err
		}
	}
	dailInfo.Timeout = 10 * time.Second
	sess, err := mgo.DialWithInfo(dailInfo)
	if err == nil {
		sess.SetSyncTimeout(1 * time.Minute)
		sess.SetSocketTimeout(1 * time.Minute)
	}
	return sess, err
}
func InitMgoSession() error {
	gsess, err := Dial(GetMgoCoinfig())
	if err != nil {
		logs.Error(" mgo.Dial(%s) failed: %s", GetMongoRealConnUrl(), err.Error())
		return err
	}
	g_MongoDB = gsess
	g_MongoDB.SetMode(mgo.SecondaryPreferred, true)
	//g_MongoDB.SetMode(mgo.PrimaryPreferred, true)
	g_MongoDB.SetPoolLimit(GetMgoCoinfig().PoolNum)
	err = InitOnlyRealMgoSession()
	if err != nil {
		return err
	}
	return nil
}

func InitOnlyRealMgoSession() error {
	gsess, err := Dial(GetMgoCoinfig())
	if err != nil {
		logs.Error(" mgo.Dial(%s) failed: %s", GetMongoRealConnUrl(), err.Error())
		return err
	}
	g_RealMongoDB = gsess
	g_RealMongoDB.SetMode(mgo.Primary, true)
	g_RealMongoDB.SetPoolLimit(4) //默认最大4个直连写库
	logs.Debug(" mgo.Dial RealMgo (%s) succ", GetMongoRealConnUrl())
	return nil
}

//这个是写
func GetRealSession() *mgo.Session {
	if g_RealMongoDB == nil {
		logs.Error("mongo 获取失败 致命错误")
		return nil
	}
	return g_RealMongoDB.Copy()
}

//这个是读
func GetOnlyReadSession() *mgo.Session {
	if g_MongoDB == nil {
		logs.Error("mongo 获取失败  致命错误")
		return nil
	}
	return g_MongoDB.Copy()
}
