package account

import (
	"comm/comm"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 账号信息
type AccountInfo struct {
	Id                bson.ObjectId `json:"-" bson:"_id"`                                   //主键ID
	GroupId           string        `json:"group_id" bson:"group_id"`                       //分组id
	Head              string        `json:"head" bson:"head"`                               //头像
	Account           string        `json:"account" bson:"account"`                         //账号
	NickName          string        `json:"nick_name" bson:"nick_name"`                     //昵称
	Status            int64         `json:"status" bson:"status"`                           //账号状态 1-离线 2-在线 3-登录中 4-登录失败 5-离线中
	Reason            string        `json:"reason" bson:"reason"`                           //离线原因
	AccountType       int64         `json:"account_type" bson:"account_type"`               //账号类型 1-个人号 2-商业号
	PlatformType      int64         `json:"platform_type" bson:"platform_type"`             //平台类型 1-云控 2-APP
	OfflineTime       int64         `json:"offline_time" bson:"offline_time"`               //离线时间
	FirstLoginTime    int64         `json:"first_login_time" bson:"first_login_time"`       //首次登录时间
	Remark            string        `json:"remark" bson:"remark"`                           //备注
	Token             string        `json:"token" bson:"token"`                             //token
	Creator           string        `json:"creator" bson:"creator"`                         //创建者
	Itime             int64         `json:"itime" bson:"itime"`                             //创建时间
	Ptime             int64         `json:"ptime" bson:"ptime"`                             //更新时间
	PullDistributeNum int64         `json:"pull_distribute_num" bson:"pull_distribute_num"` //拉群分配数量
	PixelId           string        `json:"pixel_id" bson:"pixel_id"`                       //渠道id
	ProxyIp           string        `json:"proxy_ip" bson:"proxy_ip"`                       //代理ip
	IpTime            int64         `json:"ip_time" bson:"ip_time"`                         //分配ip时间
	AreaCode          string        `json:"area_code" bson:"area_code"`                     //区号
	Synckeys          string        `json:"synckeys" bson:"synckeys"`
	IsProxyUser       int64         `json:"is_proxy_user" bson:"is_proxy_user"` //是否反向代理账号 0=否，1-是
}

// 批量新增
func AddBatchAccountInfo(tmpList []interface{}) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	err := mgoDeal.InsertMgoData(db, tb, tmpList...)
	if err != nil {
		return err
	}
	return nil
}

// 新增
func AddAccountInfo(tmp *AccountInfo) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
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
func UpAccountInfo(where bson.M, up bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	up["ptime"] = time.Now().Unix()
	err := mgoDeal.Update(db, tb, up, where)
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DelAccountInfo(where bson.M) error {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	return nil
}

// 获取列表
func GetListAccountInfo(where bson.M, num int64) []*AccountInfo {
	list := []*AccountInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, "-itime", -1, num, &list)
	if err != nil {
		return list
	}
	return list
}

// 获取数量
func GetCountAccountInfo(where bson.M) int64 {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	count, _ := mgoDeal.QueryMongoCount(db, tb, where)
	return count
}

// 获取一个
func GetOneAccountInfo(where bson.M) *AccountInfo {
	tmp := &AccountInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	mgoDeal.QueryMongoOneData(db, tb, where, "", &tmp)
	return tmp
}

// 根据id获取
func GetByIdAccountInfo(id string) *AccountInfo {
	tmp := &AccountInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	mgoDeal.QueryMongoOneData(db, tb, bson.M{"_id": bson.ObjectIdHex(id)}, "", &tmp)
	return tmp
}

// 更新
func UpNumAccountInfo(where bson.M, up bson.M) int {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	up["ptime"] = time.Now().Unix()
	return mgoDeal.UpdateNum(db, tb, up, where)
}

// 获取列表
func GetListAccountInfoSort(where bson.M, num int64, sort string) []*AccountInfo {
	list := []*AccountInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAccountInfoListInfo()
	err := mgoDeal.QueryMongoAllData(db, tb, where, sort, -1, num, &list)
	if err != nil {
		return list
	}
	return list
}
