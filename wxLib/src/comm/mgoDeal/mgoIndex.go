package mgoDeal

import (
	"comm/tableName"
	"fmt"
	"sync"
	"time"
)

var self_data_index map[string]bool = nil
var self_mutex sync.Mutex //锁 插入的时候必须顺序插入
func init() {
	self_data_index = map[string]bool{}
}

// 检查创建索引
func CheckMongoDBIndex(dbName string, collName string) {
	key := fmt.Sprintf("%s_%s", dbName, collName)
	self_mutex.Lock()
	ret := self_data_index[key]
	self_mutex.Unlock()
	if ret {
		return
	}
	self_mutex.Lock()
	//索引设置
	go SetMongoIndex(dbName, collName)
	time.Sleep(500 * time.Microsecond) //强行占20秒 免得索引没有创建号
	self_data_index[key] = true
	self_mutex.Unlock()
	return
}

// 检查和设置索引-wechat
func SetMongoIndex(dbName string, collName string) {
	switch collName {
	case tableName.GetTableAdminUserInfo(): //用户
		AdminUserInfo(dbName, collName)
	case tableName.GetTableAdminRoleInfo(): //角色
		AdminRoleInfo(dbName, collName)
	case tableName.GetTableAdminMenuInfo(): //菜单
		AdminMenuInfo(dbName, collName)
	case tableName.GetTableIpListInfo(): //ip
		IpListInfo(dbName, collName)
	case tableName.GetTableAccountInfoListInfo(): //账号表
		AccountInfo(dbName, collName)
	case tableName.GetTableSendMsgInfoListInfo(): //发送消息详情
		SendMsgInfo(dbName, collName)
	default:

	}
}
