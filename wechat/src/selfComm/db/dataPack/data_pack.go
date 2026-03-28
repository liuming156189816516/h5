package dataPack

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 数据包
type DataPack struct {
	Id              bson.ObjectId `json:"-" bson:"_id"`                               //主键ID
	Name            string        `json:"name" bson:"name"`                           //数据名称
	UpNum           int64         `json:"up_num" bson:"up_num"`                       //上传数据
	SourceRepeatNum int64         `json:"source_repeat_num" bson:"source_repeat_num"` //源重复数据
	RepeatNum       int64         `json:"repeat_num" bson:"repeat_num"`               //账号内重复
	IntoNum         int64         `json:"into_num" bson:"into_num"`                   //入库数量
	ResidueNum      int64         `json:"residue_num" bson:"residue_num"`             //剩余数量
	UpStatus        int64         `json:"up_status" bson:"up_status"`                 //上传状态 1-上传中 2-上传完成
	InvalidNum      int64         `json:"invalid_num" bson:"invalid_num"`             //无效数量
	Creator         string        `json:"creator" bson:"creator"`                     //创建者
	Itime           int64         `json:"itime" bson:"itime"`                         //创建时间
	Ptime           int64         `json:"ptime" bson:"ptime"`                         //更新时间
}

// 新增
func AddDataPack(tmp *DataPack) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
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
func UpDataPack(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelDataPack(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListDataPack(where bson.M, num int64) []*DataPack {
	list := []*DataPack{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountDataPack(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneDataPack(where bson.M) *DataPack {
	tmp := &DataPack{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdDataPack(id string) *DataPack {
	tmp := &DataPack{}
	if id == "" {
		return tmp
	}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumDataPack(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableDataPackListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
