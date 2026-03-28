package natsRpc

import (
	"github.com/BurntSushi/toml"
	"mlog"
	"os"
	"errors"
)

type NatsConfig struct {
	User     string
	Password string
	//Token    string
	Timeout int
	Secure  bool
	Servers []string
}

var nats_config = &NatsConfig{}

func GetNatsConfig() *NatsConfig {
	return nats_config
}

//func checkNatsConfigEqual(nConfig *NatsConfig, mConfig *NatsConfig) bool {
//	if nConfig == nil || mConfig == nil {
//		return false
//	}
//
//	for _, c := range nConfig.Servers {
//		if !utils.InArrayString(mConfig.Servers, c) {
//			return false
//		}
//	}
//	return true
//}
func LoadConfig(filepath string) error {
	newConf := &NatsConfig{}
	_, err := toml.DecodeFile(filepath, newConf)
	if err != nil {
		if !os.IsNotExist(err) {
			mlog.Error("DecodeFile:%s failed:%s", filepath, err.Error())
		}
		return err
	}
	mlog.Debug("newConf:%+v", newConf)
	nats_config = newConf
	return nil
}
//beego 加载配置
func LoadBeegoConfig(newConf *NatsConfig) error {
	mlog.Debug("newConf:%+v", newConf)
	nats_config = newConf
	if len(nats_config.Servers) == 0 {
		return errors.New("配置没有")
	}
	return nil
}

