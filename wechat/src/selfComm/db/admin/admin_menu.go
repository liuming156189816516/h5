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
	"sync"
	"time"
	"utils"
)

type AdminMenuInfo struct {
	MenuId    int64  `json:"menu_id"`
	Type      int64  `json:"type"` //0 菜单 1按钮
	Title     string `json:"title"`
	Pid       int64  `json:"pid"`
	Url       string `json:"url"` //路由
	Sort      int64  `json:"sort"`
	Icon      string `json:"icon"` //图标
	Api       string `json:"api"`  // 按钮才需要
	ClassName string `json:"class_name"`
	Status    int64  `json:"status"` // 0启用 1禁用
}

var AllMenuInfo = sync.Map{}

// 获取菜单
func GetAdimnAllMenuList() []*AdminMenuInfo { //永久
	allMenu := []*AdminMenuInfo{}
	key := redisKeys.GetAdminMenuInfo()
	str := redisDeal.RedisDoGetStr(key)
	if str != "" { //直接查数据
		err := jsoniter.UnmarshalFromString(str, &allMenu)
		if err == nil && len(allMenu) != 0 {
			return allMenu
		}
	}
	allMenu = adminMenuInfoFid()
	if allMenu == nil || len(allMenu) == 0 {
		return allMenu
	}
	//写redis
	redisDeal.RedisSendSet(key, allMenu)
	return allMenu
}

// 获取菜单
func GetAdimnMenuInfo(menu_id int64) *AdminMenuInfo { //永久
	roleInfo := &AdminMenuInfo{}

	data, ok := AllMenuInfo.Load(menu_id)
	if !ok {
		return roleInfo
	}
	tmp, ok := data.(*AdminMenuInfo)
	if ok && tmp.MenuId == menu_id {
		return tmp
	}
	all := GetAdimnAllMenuList()
	for _, info := range all {
		if info.MenuId == menu_id {
			AllMenuInfo.Store(menu_id, info)
			return info
		}
	}
	return roleInfo
}

// 添加菜单
func AddAdminMenu(info *AdminMenuInfo) error {
	itime := time.Now().Unix()
	key := redisKeys.GetAdminIncInfo()
	if info.MenuId == 0 {
		menu_id, _ := redisDeal.RedisDoHincrby(key, comm.MenuId, 1)
		if menu_id < 10000 {
			menu_id, _ = redisDeal.RedisDoHincrby(key, comm.MenuId, 10000)
			if menu_id < 10000 { //
				return goError.NewError("错误")
			}
		}
		info.MenuId = menu_id
	}
	//写入db
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminMenuInfo()
	insert := bson.D{
		{Name: "menu_id", Value: info.MenuId},
		{Name: "type", Value: info.Type},
		{Name: "title", Value: info.Title},
		{Name: "url", Value: info.Url},
		{Name: "pid", Value: info.Pid},
		{Name: "sort", Value: info.Sort},
		{Name: "status", Value: info.Status},
		{Name: "icon", Value: info.Icon},
		{Name: "class_name", Value: info.ClassName},
		{Name: "api", Value: info.Api},
		{Name: "itime", Value: itime},
		{Name: "ctime", Value: utils.TimeIntToString(itime)},
	}
	err := mgoDeal.InsertMgoData(db, tb, insert)
	if err != nil {
		return goError.NewError("创建失败")
	}

	key = redisKeys.GetAdminMenuInfo()
	redisDeal.RedisSendDel(key)
	return err
}

// 数据库查询
func adminMenuInfoFid() []*AdminMenuInfo {
	allMenu := []*AdminMenuInfo{}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminMenuInfo()
	alldata, err := mgoDeal.RealQueryMongoAll(db, tb, bson.M{"status": 0}, "", -1, -1)
	if err != nil {
		return allMenu
	}
	for _, data := range alldata {
		tmp := &AdminMenuInfo{}
		if p, ok := data["menu_id"]; ok {
			tmp.MenuId = utils.GetInt64(p)
		}
		if p, ok := data["type"]; ok {
			tmp.Type = utils.GetInt64(p)
		}
		if p, ok := data["title"]; ok {
			tmp.Title = utils.GetString(p)
		}
		if p, ok := data["class_name"]; ok {
			tmp.ClassName = utils.GetString(p)
		}
		if p, ok := data["pid"]; ok {
			tmp.Pid = utils.GetInt64(p)
		}
		if p, ok := data["status"]; ok {
			tmp.Status = utils.GetInt64(p)
		}
		if p, ok := data["url"]; ok {
			tmp.Url = utils.GetString(p)
		}
		if p, ok := data["sort"]; ok {
			tmp.Sort = utils.GetInt64(p)
		}
		if p, ok := data["icon"]; ok {
			tmp.Icon = utils.GetString(p)
		}
		if p, ok := data["api"]; ok {
			tmp.Api = utils.GetString(p)
		}
		allMenu = append(allMenu, tmp)
	}
	return allMenu
}

// 更新菜单
func UpdateAdminMenuInfo(menu_id int64, data map[string]interface{}) error {
	if menu_id < 0 {
		return goError.NewError("no user")
	}
	where := bson.M{"menu_id": menu_id}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminMenuInfo()
	err := mgoDeal.Update(db, tb, data, where)
	if err != nil {
		return err
	}
	key := redisKeys.GetAdminMenuInfo()
	redisDeal.RedisSendDel(key)
	return nil
}

// 删除菜单
func DelAdminMenuInfo(menu_id int64) error {
	if menu_id < 0 {
		return goError.NewError("no user")
	}
	where := bson.M{"menu_id": menu_id}
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminMenuInfo()
	err := mgoDeal.MgoDelete(db, tb, where)
	if err != nil {
		return err
	}
	key := redisKeys.GetAdminMenuInfo()
	redisDeal.RedisSendDel(key)
	return nil
}
