package sendmsg

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 群发详情
type SendMsgInfo struct {
	Id            bson.ObjectId `json:"-" bson:"_id"`                         //主键ID
	Account       string        `json:"account" bson:"account"`               //账号
	AccountGroup  string        `json:"account_group" bson:"account_group"`   //账号分组
	AccountStatus int64         `json:"account_status" bson:"account_status"` //账号状态 1-离线 2-在线
	SucessNum     int64         `json:"sucess_num" bson:"sucess_num"`         //发送数量
	ArrivedNum    int64         `json:"arrived_num" bson:"arrived_num"`       //已送达
	Reason        string        `json:"reason" bson:"reason"`                 //原因
	Creator       string        `json:"creator" bson:"creator"`               //创建者
	Itime         int64         `json:"itime" bson:"itime"`                   //创建时间
	Ptime         int64         `json:"ptime" bson:"ptime"`                   //更新时间
}

// 新增
func AddSendMsgInfo(tmp *SendMsgInfo) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
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
func UpSendMsgInfo(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelSendMsgInfo(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListSendMsgInfo(where bson.M, num int64) []*SendMsgInfo {
	list := []*SendMsgInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountSendMsgInfo(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneSendMsgInfo(where bson.M) *SendMsgInfo {
	tmp := &SendMsgInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdSendMsgInfo(id string) *SendMsgInfo {
	tmp := &SendMsgInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumSendMsgInfo(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableSendMsgInfoListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}
