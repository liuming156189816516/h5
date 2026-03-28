package mgoDeal

import (
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
)

//
func TableConfig(dbName string, collName string) {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession %s:%s, no session", dbName, collName)
		return
	}
	defer session.Close()
	session.SetSafe(nil)
	indexs := []mgo.Index{}

	//phone
	vip_index := mgo.Index{}
	vip_index.Key = append(vip_index.Key, "key")
	vip_index.Background = true
	vip_index.Unique = true
	indexs = append(indexs, vip_index)

	coll := session.DB(dbName).C(collName)

	for _, index := range indexs {
		coll.EnsureIndex(index)
	}

}
