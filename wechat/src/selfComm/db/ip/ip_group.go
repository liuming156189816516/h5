package ip

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//ip分组
type IpGroup struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`           //主键ID
	Name    string        `json:"name" bson:"name"`       //名称
	Creator string        `json:"creator" bson:"creator"` //创建者
	Itime   int64         `json:"itime" bson:"itime"`     //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`     //更新时间
}

//新增
func AddIpGroup(tmp *IpGroup) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	now := time.Now().Unix()
	tmp.Itime = now
	tmp.Ptime = now
	err := mgoDeal.InsertMgoData(db, tb, tmp)
	if err != nil {
		return err
	}
	return nil
}

//更新
func UpIpGroup(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

//删除
func DelIpGroup(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

//获取列表
func GetListIpGroup(where bson.M, num int64) []*IpGroup {
	list := []*IpGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

//获取数量
func GetCountIpGroup(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

//获取一个
func GetOneIpGroup(where bson.M) *IpGroup {
	tmp := &IpGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

//根据id获取
func GetByIdIpGroup(id string) *IpGroup {
	tmp := &IpGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

//更新
func UpNumIpGroup(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableIpGroupListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
