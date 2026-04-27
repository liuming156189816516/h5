package account

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//账号入库日志
type AccountLog struct {
	Id      bson.ObjectId `json:"-" bson:"_id"`           //主键ID
	FileId  string        `json:"file_id" bson:"file_id"` //文件id
	Account string        `json:"account" bson:"account"` //账号
	Status  int64         `json:"status" bson:"status"`   //入库状态 1-失败 2-成功
	Reason  string        `json:"reason" bson:"reason"`   //原因
	Token   string        `json:"token" bson:"token"`     //token
	Creator string        `json:"creator" bson:"creator"` //创建者
	Itime   int64         `json:"itime" bson:"itime"`     //创建时间
	Ptime   int64         `json:"ptime" bson:"ptime"`     //更新时间
}

//新增
func AddAccountLog(tmp *AccountLog) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	now := time.Now().Unix()
	tmp.Itime = now
	tmp.Ptime = now
	err := mgoDeal.InsertMgoData(db, tb, tmp)
	if err != nil {
		return err
	}
	return nil
}

//批量新增
func AddBatchAccountLog(tmpList []interface{}) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	err := mgoDeal.InsertMgoData(db, tb, tmpList...)
	if err != nil {
		return err
	}
	return nil
}

//更新
func UpAccountLog(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

//删除
func DelAccountLog(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

//获取列表
func GetListAccountLog(where bson.M, num int64) []*AccountLog {
	list := []*AccountLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

//获取数量
func GetCountAccountLog(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

//获取一个
func GetOneAccountLog(where bson.M) *AccountLog {
	tmp := &AccountLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

//根据id获取
func GetByIdAccountLog(id string) *AccountLog {
	tmp := &AccountLog{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

//更新
func UpNumAccountLog(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountLogListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
