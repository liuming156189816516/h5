package account

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 账号分组
type AccountGroup struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`           //主键ID
	Name    string        `json:"name" bson:"name"`       //分组名称
	Sort    int64         `json:"sort" bson:"sort"`       //排序字段
	Creator string        `json:"creator" bson:"creator"` //创建者
	Itime   int64         `json:"itime" bson:"itime"`     //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`     //更新时间
}

// 新增
func AddAccountGroup(tmp *AccountGroup) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
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
func UpAccountGroup(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelAccountGroup(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListAccountGroup(where bson.M, num int64) []*AccountGroup {
	list := []*AccountGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountAccountGroup(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneAccountGroup(where bson.M) *AccountGroup {
	tmp := &AccountGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdAccountGroup(id string) *AccountGroup {
	tmp := &AccountGroup{}
	if id == "" {
		return tmp
	}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumAccountGroup(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountGroupListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
