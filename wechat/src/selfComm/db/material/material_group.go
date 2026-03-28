package material

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 素材分组
type MaterialGroup struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`           //主键ID
	Name    string        `json:"name" bson:"name"`       //名称
	Type    int64         `json:"type" bson:"type"`       //类型 1-文字 2-图片 3-语音 4-视频 7-超链
	Creator string        `json:"creator" bson:"creator"` //创建者
	Itime   int64         `json:"itime" bson:"itime"`     //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`     //更新时间
}

// 新增
func AddMaterialGroup(tmp *MaterialGroup) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
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
func UpMaterialGroup(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelMaterialGroup(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListMaterialGroup(where bson.M, num int64) []*MaterialGroup {
	list := []*MaterialGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountMaterialGroup(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneMaterialGroup(where bson.M) *MaterialGroup {
	tmp := &MaterialGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdMaterialGroup(id string) *MaterialGroup {
	tmp := &MaterialGroup{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumMaterialGroup(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialGroupListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
