package account

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//账号入库文件
type AccountFile struct {
	Id          bson.ObjectId `json:"-" bson:"_id"`                     //主键ID
	Name        string        `json:"name" bson:"name"`                 //文件名称
	AccountType int64         `json:"account_type" bson:"account_type"` //账号分类 1-个人号 2-商业号
	Remark      string        `json:"remark" bson:"remark"`             //备注
	Status      int64         `json:"status" bson:"status"`             //任务状态 1-导入中 2-已完成
	SuccessNum  int64         `json:"success_num" bson:"success_num"`   //成功数量
	FailNum     int64         `json:"fail_num" bson:"fail_num"`         //失败数量
	Creator     string        `json:"creator" bson:"creator"`           //创建者
	Itime       int64         `json:"itime" bson:"itime"`               //创建时间
	Ptime       int64         `json:"ptime" bson:"ptime"`               //更新时间
}

//新增
func AddAccountFile(tmp *AccountFile) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
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
func UpAccountFile(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

//删除
func DelAccountFile(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

//获取列表
func GetListAccountFile(where bson.M, num int64) []*AccountFile {
	list := []*AccountFile{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

//获取数量
func GetCountAccountFile(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

//获取一个
func GetOneAccountFile(where bson.M) *AccountFile {
	tmp := &AccountFile{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

//根据id获取
func GetByIdAccountFile(id string) *AccountFile {
	tmp := &AccountFile{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

//更新
func UpNumAccountFile(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountFileListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
