package ip

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// ip
type Ip struct {
	Id            bson.ObjectId `json:"-" bson:"_id"`                         //主键ID
	GroupId       string        `json:"group_id" bson:"group_id"`             //分组id
	ProxyIp       string        `json:"proxy_ip" bson:"proxy_ip"`             //代理ip
	ProxyUser     string        `json:"proxy_user" bson:"proxy_user"`         //代理用户
	ProxyPwd      string        `json:"proxy_pwd" bson:"proxy_pwd"`           //代理密码
	ProxyType     string        `json:"proxy_type" bson:"proxy_type"`         //代理ip类型 socks5 http
	Status        int64         `json:"status" bson:"status"`                 //网络状态 1-正常 2-异常 3-检测中 4-未检测 5-冻结
	AllotNum      int64         `json:"allot_num" bson:"allot_num"`           //分配总数
	UserNum       int64         `json:"user_num" bson:"user_num"`             //已分配
	IpType        int64         `json:"ip_type" bson:"ip_type"`               //ip类型 1-ipV4 2-ipv6
	IpCategory    int64         `json:"ip_category" bson:"ip_category"`       //ip类别 1-静态ip 2-动态ip
	ExpireStatus  int64         `json:"expire_status" bson:"expire_status"`   //到期状态 1-正常 2-到期 3-即将到期
	Country       string        `json:"country" bson:"country"`               //国家
	DisableStatus int64         `json:"disable_status" bson:"disable_status"` //禁用状态 1-启用 2-禁用
	ExpireTime    int64         `json:"expire_time" bson:"expire_time"`       //到期时间
	Remark        string        `json:"remark" bson:"remark"`                 //备注
	Reason        string        `json:"reason" bson:"reason"`                 //原因
	Creator       string        `json:"creator" bson:"creator"`               //创建者
	Itime         int64         `json:"itime" bson:"itime"`                   //创建时间
	Ptime         int64         `json:"ptime" bson:"ptime"`                   //更新时间
	AllotTime     int64         `json:"allot_time" bson:"allot_time"`         //ip分配时间
}

// 新增
func AddIp(tmp *Ip) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	now := time.Now().Unix()
	tmp.Itime = now
	tmp.Ptime = now
	err := mgoDeal.InsertMgoData(db, tb, tmp)
	if err != nil {
		return err
	}
	return nil
}

// 更新
func UpIp(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelIp(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListIp(where bson.M, num int64) []*Ip {
	list := []*Ip{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "_id", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountIp(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneIp(where bson.M) *Ip {
	tmp := &Ip{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdIp(id string) *Ip {
	tmp := &Ip{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumIp(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}

// 获取一个
func GetOneLockIp() *Ip {
	tmp := &Ip{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpListInfo()
	err := mgoDeal.QueryMongoOneData(db, tb, bson.M{"status": int64(1)}, "allot_time", &tmp)
	if err == nil {
		UpIp(bson.M{"_id": tmp.Id}, bson.M{"allot_time": time.Now().Unix()})
	}
	return tmp
}
