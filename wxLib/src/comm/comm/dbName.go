package comm

var DBName = string("ws")


//默认库
func GetMgoDBName() string {
	return DBName
}

//用户库
func GetUserMgoDBName(uid string) string {
	return DBName + "_" + uid
}
