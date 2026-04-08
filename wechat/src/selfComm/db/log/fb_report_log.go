package log

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// FbReportLog 小莫上报
type FbReportLog struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`           //主键ID
	Ptype   int64         `json:"ptype" bson:"ptype"`     //1-内容查看-ViewContent 2-获取验证码-SubmitVerification 3-挂机成功-SessionSuccess
	Data    string        `json:"data" bson:"data"`       //上报
	Creator string        `json:"creator" bson:"creator"` //创建者
	Itime   int64         `json:"itime" bson:"itime"`     //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`     //更新时间
}

// 新增
func AddFbReportLog(tmp *FbReportLog) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
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
func UpFbReportLog(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelFbReportLog(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListFbReportLog(where bson.M, num int64) []*FbReportLog {
	list := []*FbReportLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountFbReportLog(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneFbReportLog(where bson.M) *FbReportLog {
	tmp := &FbReportLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdFbReportLog(id string) *FbReportLog {
	tmp := &FbReportLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumFbReportLog(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableFbReportLog()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
