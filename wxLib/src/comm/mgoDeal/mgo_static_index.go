package mgoDeal

import (
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
)

//后台账号信息
func AdminUserInfo(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	//user_key
	typeindex := mgo.Index{}
	typeindex.Key = append(typeindex.Key, "account")
	typeindex.Background = true
	typeindex.Unique = true
	indexs = append(indexs, typeindex)

	level_index := mgo.Index{}
	level_index.Key = append(level_index.Key, "level")
	level_index.Background = true
	indexs = append(indexs, level_index)

	ip_index := mgo.Index{}
	ip_index.Key = append(ip_index.Key, "ip")
	ip_index.Background = true
	indexs = append(indexs, ip_index)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}

func AdminRoleInfo(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	//user_key
	typeindex := mgo.Index{}
	typeindex.Key = append(typeindex.Key, "role_id")
	typeindex.Background = true
	typeindex.Unique = true
	indexs = append(indexs, typeindex)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}

func AdminMenuInfo(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	//user_key
	typeindex := mgo.Index{}
	typeindex.Key = append(typeindex.Key, "menu_id")
	typeindex.Background = true
	typeindex.Unique = true
	indexs = append(indexs, typeindex)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}

//ip列表
func IpListInfo(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	proxy_index := mgo.Index{}
	proxy_index.Key = append(proxy_index.Key, "proxy_ip")
	proxy_index.Key = append(proxy_index.Key, "proxy_user")
	proxy_index.Key = append(proxy_index.Key, "proxy_pwd")
	proxy_index.Background = true
	proxy_index.Unique = true
	indexs = append(indexs, proxy_index)

	ip_type_index := mgo.Index{}
	ip_type_index.Key = append(ip_type_index.Key, "ip_type")
	ip_type_index.Background = true
	indexs = append(indexs, ip_type_index)

	ip_category_index := mgo.Index{}
	ip_category_index.Key = append(ip_category_index.Key, "ip_category")
	ip_category_index.Background = true
	indexs = append(indexs, ip_category_index)

	disable_status_index := mgo.Index{}
	disable_status_index.Key = append(disable_status_index.Key, "disable_status")
	disable_status_index.Background = true
	indexs = append(indexs, disable_status_index)

	status_index := mgo.Index{}
	status_index.Key = append(status_index.Key, "status")
	status_index.Background = true
	indexs = append(indexs, status_index)

	country_index := mgo.Index{}
	country_index.Key = append(country_index.Key, "country")
	country_index.Background = true
	indexs = append(indexs, country_index)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}
//账号信息
func AccountInfo(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	groupId_index := mgo.Index{}
	groupId_index.Key = append(groupId_index.Key, "group_id")
	groupId_index.Background = true
	indexs = append(indexs, groupId_index)

	account_index := mgo.Index{}
	account_index.Key = append(account_index.Key, "account")
	account_index.Background = true
	account_index.Unique = true
	indexs = append(indexs, account_index)

	status_index := mgo.Index{}
	status_index.Key = append(status_index.Key, "status")
	status_index.Background = true
	indexs = append(indexs, status_index)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}

//发送消息详情
func SendMsgInfo(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	account_index := mgo.Index{}
	account_index.Key = append(account_index.Key, "account")
	account_index.Background = true
	account_index.Unique = true
	indexs = append(indexs, account_index)

	account_status_index := mgo.Index{}
	account_status_index.Key = append(account_status_index.Key, "account_status")
	account_status_index.Background = true
	indexs = append(indexs, account_status_index)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}


