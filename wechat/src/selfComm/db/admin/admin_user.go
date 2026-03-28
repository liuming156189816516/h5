package admin

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/redisDeal"
	"comm/redisKeys"
	"comm/tableName"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 初始化 一个root账号
func init() {
	go func() {
		time.Sleep(5 * time.Second)
		info := GetOneAdminUser(bson.M{"account": "admin"})
		if info.Account != "admin" {
			tmp := &AdminUser{}
			tmp.Id = bson.NewObjectId()
			tmp.Account = "admin"
			tmp.AccountType = 1
			tmp.Pwd = comm.Md5("admin")
			AddAdminUser(tmp)
		}
	}()
}

// 用户
type AdminUser struct {
	Id          bson.ObjectId `json:"-" bson:"_id"`                     //主键ID
	Account     string        `json:"account" bson:"account"`           //账号
	PwdStr      string        `json:"pwd_str" bson:"pwd_str"`           //明文密码
	TwoPwd      string        `json:"two_pwd" bson:"two_pwd"`           //二级密码
	AccountType int64         `json:"account_type" bson:"account_type"` //账号类型 1-管理员 2-主管 3-用户
	Pwd         string        `json:"pwd" bson:"pwd"`                   //密码
	Status      int64         `json:"status" bson:"status"`             //状态 1-启用 2-禁用
	RoleId      int64         `json:"role_id" bson:"role_id"`           //角色id
	Creator     string        `json:"creator" bson:"creator"`           //创建者
	Itime       int64         `json:"itime" bson:"itime"`               //创建时间
	Ptime       int64         `json:"ptime" bson:"ptime"`               //更新时间

	Head string `json:"head" bson:"head"` //头像
}

// 新增
func AddAdminUser(tmp *AdminUser) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	now := time.Now().Unix()
	tmp.Itime = now
	tmp.Ptime = now
	err := mgoDeal.InsertMgoData(db, tb, tmp)
	if err != nil {
		return err
	}
	redisDeal.RedisSendDel(redisKeys.GetAdminUserInfo())
	return nil
}

// 更新
func UpAdminUser(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	redisDeal.RedisSendDel(redisKeys.GetAdminUserInfo())
	return nil
}

// 删除
func DelAdminUser(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	redisDeal.RedisSendDel(redisKeys.GetAdminUserInfo())
	return nil
}

// 获取列表
func GetListAdminUser(where bson.M, num int64) []*AdminUser {
	list := []*AdminUser{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountAdminUser(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneAdminUser(where bson.M) *AdminUser {
	tmp := &AdminUser{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	err := mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	if err != nil {
		return tmp
	}
	return tmp
}

// 更新
func UpNumAdminUser(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	up["ptime"] = time.Now().Unix()
	redisDeal.RedisSendDel(redisKeys.GetAdminUserInfo())
	return mgoDeal.UpdateNum(db, tb, up, where)
}

// 根据id获取
func GetByIdAdminUser(id string) *AdminUser {
	tmp := &AdminUser{}
	if id == "" {
		return tmp
	}
	key := redisKeys.GetAdminUserInfo()
	str := redisDeal.RedisDoHGetSrt(key, id)
	if str != "" { //直接查数据
		jsoniter.UnmarshalFromString(str, &tmp)
		return tmp
	}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func IncAdminUser(where bson.M, inc map[string]int64) {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	mgoDeal.UpInc(db, tb, where, inc, nil)
	redisDeal.RedisSendDel(redisKeys.GetAdminUserInfo())
}
