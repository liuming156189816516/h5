package admin

import (
	"comm/comm"
	"comm/goError"
	"comm/mgoDeal"
	"comm/redisDeal"
	"comm/redisKeys"
	"comm/tableName"
	"github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"time"
	"utils"
)

type AdminRoleInfo struct {
	RoleId int64   `json:"role_id"`
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Itime  int64   `json:"itime"`
	Menu   []int64 `json:"menu"`
}

// 获取用户角色
func GetAdimnRoleInfo(role_id int64) *AdminRoleInfo { //永久
	roleInfo := &AdminRoleInfo{}
	if role_id < 0 {
		return roleInfo
	}
	//如果是角色 0的 页面全部都有
	defer func() {
		if roleInfo.RoleId == 0 && roleInfo.Itime > 0 { //角色 0点的
			roleInfo.Menu = []int64{}
			allmenu := GetAdimnAllMenuList()
			for _, me := range allmenu {
				roleInfo.Menu = append(roleInfo.Menu, me.MenuId)
			}
		}
	}()
	key := redisKeys.GetAdminRoleInfo()
	str := redisDeal.RedisDoHGetSrt(key, role_id)
	if str != "" { //直接查数据
		err := jsoniter.UnmarshalFromString(str, &roleInfo)
		if err == nil && roleInfo.RoleId == role_id && roleInfo.Itime > 0 {
			return roleInfo
		}
	}
	//查数据库
	find := bson.M{}
	find["role_id"] = role_id
	roleInfo = adminRoleInfoFid(find)
	if roleInfo == nil || roleInfo.RoleId != role_id || roleInfo.Itime == 0 {
		return roleInfo
	}
	//写redis
	redisDeal.RedisDoHSet(key, role_id, roleInfo)
	return roleInfo
}

// 添加角色
func AddAdminRole(info *AdminRoleInfo) error {
	itime := time.Now().Unix()
	//写入db
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminRoleInfo()
	menu, _ := jsoniter.MarshalToString(info.Menu)
	insert := bson.D{
		{Name: "role_id", Value: info.RoleId},
		{Name: "name", Value: info.Name},
		{Name: "desc", Value: info.Desc},
		{Name: "menu", Value: menu},
		{Name: "itime", Value: itime},
		{Name: "ctime", Value: utils.TimeIntToString(itime)},
	}
	err := mgoDeal.InsertMgoData(db, tb, insert)
	if err != nil {
		return goError.NewError("创建失败")
	}
	return err
}

// 数据库查询
func adminRoleInfoFid(find bson.M) *AdminRoleInfo {
	roleInfo := &AdminRoleInfo{}

	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminRoleInfo()
	data, err := mgoDeal.RealQueryMongoOne(db, tb, find, "")
	if err != nil {
		//logs.Error("查询错误%+v", err)
		return roleInfo
	}
	if p, ok := data["role_id"]; ok {
		roleInfo.RoleId = utils.GetInt64(p)
	}
	if p, ok := data["name"]; ok {
		roleInfo.Name = utils.GetString(p)
	}
	if p, ok := data["desc"]; ok {
		roleInfo.Desc = utils.GetString(p)
	}
	roleInfo.Menu = []int64{}
	if p, ok := data["menu"]; ok {
		menu := utils.GetString(p)
		if menu != "" {
			jsoniter.UnmarshalFromString(menu, &roleInfo.Menu)
		}
	}

	if p, ok := data["itime"]; ok {
		roleInfo.Itime = utils.GetInt64(p)
	}
	return roleInfo
}

// 更新用户信息
func UpdateAdminRoleInfo(role_id int64, data map[string]interface{}) error {
	//if role_id < 0 {
	//	return goError.NewError("no user")
	//}
	where := bson.M{"role_id": role_id}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminRoleInfo()
	err := mgoDeal.Update(db, tb, data, where)
	if err != nil {
		return err
	}
	key := redisKeys.GetAdminRoleInfo()
	redisDeal.RedisSendHDel(key, role_id)
	return nil
}

// 更新用户信息
func DelAdminRoleInfo(role_id int64) error {
	//if role_id <= 0 {
	//	return goError.NewError("no user")
	//}
	where := bson.M{"role_id": role_id}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminRoleInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	key := redisKeys.GetAdminRoleInfo()
	redisDeal.RedisSendHDel(key, role_id)
	return nil
}
