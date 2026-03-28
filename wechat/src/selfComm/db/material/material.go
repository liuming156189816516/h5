package material

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//素材
type Material struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`             //主键ID
	GroupId string        `json:"group_id" bson:"group_id"` //分组id
	Name    string        `json:"name" bson:"name"`         //标题
	Content string        `json:"content" bson:"content"`   //内容
	Type    int64         `json:"type" bson:"type"`         //类型  1-文字 2-图片 3-语音 4-视频 5-语音呼叫 6-名片
	Creator string        `json:"creator" bson:"creator"`   //创建者
	Itime   int64         `json:"itime" bson:"itime"`       //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`       //更新时间
}

//新增
func AddMaterial(tmp *Material) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
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
func UpMaterial(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

//删除
func DelMaterial(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

//获取列表
func GetListMaterial(where bson.M, num int64) []*Material {
	list := []*Material{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

//获取数量
func GetCountMaterial(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

//获取一个
func GetOneMaterial(where bson.M) *Material {
	tmp := &Material{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

//根据id获取
func GetByIdMaterial(id string) *Material {
	tmp := &Material{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

//更新
func UpNumMaterial(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableMaterialListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
