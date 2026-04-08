package log

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// FbLog 给fb上报
type FbLog struct {
	Id        bson.ObjectId `json:"-" bson:"_id"`                 //主键ID
	EventID   string        `json:"event_id" bson:"event_id"`     //事件ID
	EventName string        `json:"event_name" bson:"event_name"` //事件名称
	Phone     string        `json:"phone" bson:"phone"`           //手机号
	Fbclid    string        `json:"fbclid" bson:"fbclid"`         //fbc
	Fbp       string        `json:"fbp" bson:"fbp"`               //fbc
	PixelId   string        `json:"pixel_id" bson:"pixel_id"`     //像素ID
	Creator   string        `json:"creator" bson:"creator"`       //创建者
	Itime     int64         `json:"itime" bson:"itime"`           //创建时间
	Ptime     int64         `json:"ptime" bson:"ptime"`           //更新时间
}

// 新增
func AddFbLog(tmp *FbLog) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
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
func UpFbLog(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelFbLog(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListFbLog(where bson.M, num int64) []*FbLog {
	list := []*FbLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountFbLog(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneFbLog(where bson.M) *FbLog {
	tmp := &FbLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdFbLog(id string) *FbLog {
	tmp := &FbLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumFbLog(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbLog()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
