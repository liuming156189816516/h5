package fb

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Fb
type Fb struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`       //主键ID
	Phone   string        `json:"phone" bson:"phone"` // 手机号
	Fbclid  string        `json:"fbclid" bson:"fbclid"`
	Fbp     string        `json:"fbp" bson:"fbp"`
	PixelId string        `json:"pixel_id" bson:"pixel_id"`
	Creator string        `json:"creator" bson:"creator"` //创建者
	Itime   int64         `json:"itime" bson:"itime"`     //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`     //更新时间
}

// 新增
func AddFb(tmp *Fb) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
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
func UpFb(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelFb(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListFb(where bson.M, num int64) []*Fb {
	list := []*Fb{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountFb(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneFb(where bson.M) *Fb {
	tmp := &Fb{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	mgoDeal.QueryMongoOneData(db, tb, where, "-itime", &tmp)
	return tmp
}

// 根据id获取
func GetByIdFb(id string) *Fb {
	tmp := &Fb{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumFb(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFb()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
