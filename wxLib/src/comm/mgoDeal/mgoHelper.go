package mgoDeal

import (
	"comm/goError"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2/bson"
	"mframe/stat"
	"time"
)

//更新 inc
func UpInc(dbName string, collName string, where map[string]interface{}, incData map[string]int64, setData map[string]interface{}) (int64, error) {
	//logs.Debug("Mongo UpInc db==%s tb==%s where==%s data==%+v ", dbName, collName, where, incData)
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return 0, errors.New("No Mgo DB")
	}
	defer session.Close()
	start := time.Now()
	Id := int64(0)
	if len(incData) <= 0 || len(where) <= 0 {
		return Id, errors.New("No Update item")
	}
	qDoc := bson.D{}
	for k, v := range where {
		qDoc = append(qDoc, bson.DocElem{Name: k, Value: v})
	}
	uDoc := bson.D{}
	if len(incData) > 0 {
		//inc
		incDoc := bson.D{}
		for k, v := range incData {
			incDoc = append(incDoc, bson.DocElem{Name: k, Value: v})
		}
		uDoc = append(uDoc, bson.DocElem{Name: "$inc", Value: incDoc})
	}

	//设置
	if len(setData) > 0 {
		setDoc := bson.D{}
		for k, v := range setData {
			setDoc = append(setDoc, bson.DocElem{Name: k, Value: v})
		}
		uDoc = append(uDoc, bson.DocElem{Name: "$set", Value: setDoc})
	}
	if len(uDoc) == 0 {
		return 0, nil
	}
	coll := session.DB(dbName).C(collName)
	cInfo, err := coll.Upsert(qDoc, uDoc)
	ret := 0
	if err != nil {
		stat.ReportStat(fmt.Sprintf("UpInc-%s", collName), 1, time.Now().Sub(start))
		return Id, err
	}
	stat.ReportStat(fmt.Sprintf("UpInc-%s", collName), ret, time.Now().Sub(start))
	CheckMongoDBIndex(dbName, collName) //检查索引
	Id = int64(cInfo.Updated)
	return Id, nil
}

//更新 inset
func UpInsert(dbName string, collName string, where map[string]interface{}, inData map[string]int64) (int64, error) {
	//logs.Debug("Mongo UpInsert db==%s tb==%s where==%s data==%+v ", dbName, collName, where, inData)
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return 0, errors.New("No Mgo DB")
	}
	defer session.Close()
	Id := int64(0)
	if len(inData) <= 0 || len(where) <= 0 {
		return Id, errors.New("No Update item")
	}
	qDoc := bson.D{}
	for k, v := range where {
		qDoc = append(qDoc, bson.DocElem{Name: k, Value: v})
	}
	inDoc := bson.D{}
	for k, v := range inData {
		inDoc = append(inDoc, bson.DocElem{Name: k, Value: v})
	}
	uDoc := bson.D{}
	uDoc = append(uDoc, bson.DocElem{Name: "$set", Value: inDoc})

	coll := session.DB(dbName).C(collName)
	cInfo, err := coll.Upsert(qDoc, uDoc)
	if err != nil {
		return Id, err
	}
	CheckMongoDBIndex(dbName, collName) //检查索引
	Id = int64(cInfo.Updated)
	return Id, nil
}

//更新
func Update(dbName string, collName string, updateData bson.M, whereData map[string]interface{}) error {
	//logs.Debug("Mongo Update db==%s tb==%s where==%s data==%+v ", dbName, collName, whereData, updateData)
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return errors.New("No MgoDb")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	up := bson.M{"$set": updateData}
	_, err := coll.UpdateAll(whereData, up)
	if err != nil {
		return err
	}
	return nil
}

//更新
func UpdateNum(dbName string, collName string, updateData bson.M, whereData bson.M) int {
	//logs.Debug("Mongo UpdateNum db==%s tb==%s where==%s data==%+v ", dbName, collName, whereData, updateData)
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return 0
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)

	start := time.Now()
	up := bson.M{"$set": updateData}
	get, err := coll.UpdateAll(whereData, up)
	if err != nil || get == nil {
		stat.ReportStat(fmt.Sprintf("UpdateNum-%s", collName), 1, time.Now().Sub(start))
		return 0
	}
	stat.ReportStat(fmt.Sprintf("UpdateNum-%s", collName), 0, time.Now().Sub(start))
	return get.Updated
}

//删除
func MgoDelete(dbName string, collName string, whereData bson.M) error {
	//logs.Debug("Mongo MgoDelete db==%s tb==%s where==%s data==%+v ", dbName, collName, whereData)
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return errors.New("No MgoDb")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	_, err := coll.RemoveAll(whereData)
	if err != nil {
		logs.Error("Delete MongoDB %s:%s Error %+v", dbName, collName, err)
		return err
	}
	return nil
}

//删除整张表
func MgoDrop(dbName string, collName string) error {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return errors.New("No MgoDb")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	err := coll.DropCollection()
	if err != nil {
		logs.Error("Drop MongoDB %s:%s Error %+v", dbName, collName, err)
		return err
	}
	return nil
}

//删除数据库
func MgoDropDB(dbName string) error {
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s, no session", dbName)
		return errors.New("No MgoDb")
	}
	defer session.Close()
	db1 := session.DB(dbName)
	err := db1.DropDatabase()
	if err != nil {
		logs.Error("Drop MongoDB %s Error %+v", dbName, err)
		return err
	}
	return nil
}

//插入
func InsertMgoData(dbName string, collName string, insertData ...interface{}) error {
	//logs.Debug("Mongo InsertMgoData db==%s tb==%s  data==%+v ", dbName, collName, insertData[0])
	if len(insertData) == 0 {
		return nil
	}
	session := GetRealSession()
	if session == nil {
		logs.Error("GetRealSession  %s:%s, no session", dbName, collName)
		return errors.New("No MgoDb")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	err := coll.Insert(insertData[0])
	if err != nil {
		return err
	}
	CheckMongoDBIndex(dbName, collName) //检查索引
	tmpInsert := []interface{}{}
	for k, in := range insertData {
		if k == 0 {
			continue
		}
		tmpInsert = append(tmpInsert, in)
	}
	start := time.Now()
	if len(tmpInsert) > 0 {
		err := coll.Insert(tmpInsert...)
		if err != nil {
			stat.ReportStat(fmt.Sprintf("InsertMgoData-%s", collName), 1, time.Now().Sub(start))
			return err
		}
	}
	stat.ReportStat(fmt.Sprintf("InsertMgoData-%s", collName), 0, time.Now().Sub(start))
	return nil
}

//更新 inset
func UpOrInsert(dbName string, collName string, where map[string]interface{}, inData map[string]interface{}) (int64, error) {
	//logs.Debug("Mongo UpOrInsert db==%s tb==%s where==%s data==%+v ", dbName, collName, where,inData)
	session := GetRealSession()
	if session == nil {
		return 0, errors.New("No Mgo DB")
	}
	defer session.Close()
	Id := int64(0)
	if len(inData) <= 0 || len(where) <= 0 {
		return Id, errors.New("No Update item")
	}
	qDoc := bson.D{}
	for k, v := range where {
		qDoc = append(qDoc, bson.DocElem{Name: k, Value: v})
	}
	inDoc := bson.D{}
	for k, v := range inData {
		inDoc = append(inDoc, bson.DocElem{Name: k, Value: v})
	}
	start := time.Now()
	uDoc := bson.D{}
	uDoc = append(uDoc, bson.DocElem{Name: "$set", Value: inDoc})
	coll := session.DB(dbName).C(collName)
	cInfo, err := coll.Upsert(qDoc, uDoc)
	if err != nil {
		stat.ReportStat(fmt.Sprintf("UpOrInsert-%s", collName), 1, time.Now().Sub(start))
		return Id, err
	}
	stat.ReportStat(fmt.Sprintf("UpOrInsert-%s", collName), 0, time.Now().Sub(start))
	CheckMongoDBIndex(dbName, collName) //检查索引
	Id = int64(cInfo.Updated)
	return Id, nil
}

//走写库
func QueryMongoOne(dbName string, collName string, findData bson.M, sortStr string) (bson.M, error) {
	//logs.Debug("Mongo QueryMongoOne db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession   %s:%s, no session", dbName, collName)
		return nil, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	start := time.Now()
	ret := bson.M{}
	coll := session.DB(dbName).C(collName)
	if sortStr == "" {
		err := coll.Find(findData).Limit(1).One(&ret)
		return ret, err
	}
	//0
	err := coll.Find(findData).Sort(sortStr).Limit(1).One(&ret)
	if err != nil {
		stat.ReportStat(fmt.Sprintf("QueryMongoOne-%s", collName), 1, time.Now().Sub(start))
		return ret, err
	}
	stat.ReportStat(fmt.Sprintf("QueryMongoOne-%s", collName), 0, time.Now().Sub(start))
	return ret, err
}

func QueryMongoOneData(dbName string, collName string, findData bson.M, sortStr string, in interface{}) error {
	//logs.Debug("Mongo QueryMongoOneData db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		return goError.NewError("没有 mongo session")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	if sortStr == "" {
		err := coll.Find(findData).Limit(1).One(in)
		return err
	}
	start := time.Now()
	//0
	err := coll.Find(findData).Sort(sortStr).Limit(1).One(in)
	if err != nil {
		stat.ReportStat(fmt.Sprintf("QueryMongoOneData-%s", collName), 0, time.Now().Sub(start))
		return err
	}
	stat.ReportStat(fmt.Sprintf("QueryMongoOneData-%s", collName), 0, time.Now().Sub(start))
	return err
}

//mongo 查询
func QueryMongoAll(dbName string, collName string, findData bson.M, sortStr string, start int64, findNum int64, selectItme ...bson.M) ([]bson.M, error) {
	//logs.Debug("Mongo QueryMongoAll db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession  %s:%s, no session", dbName, collName)
		return nil, goError.NewError("没有 mongo session")
	}
	st := time.Now()
	defer session.Close()
	ret := []bson.M{}
	coll := session.DB(dbName).C(collName)
	query := coll.Find(findData)
	if len(selectItme) > 0 {
		query = query.Select(selectItme[0])
	}
	if sortStr != "" {
		query = query.Sort(sortStr)
	}
	if start > -1 {
		query = query.Skip(int(start))
	}
	if findNum > -1 {
		query = query.Limit(int(findNum))
	}
	err := query.All(&ret)
	if err != nil {
		stat.ReportStat(fmt.Sprintf("QueryMongoAll-%s", collName), 0, time.Now().Sub(st))
		return ret, err
	}
	stat.ReportStat(fmt.Sprintf("QueryMongoAll-%s", collName), 0, time.Now().Sub(st))
	return ret, err
}

//mongo
func QueryMongoAllData(dbName string, collName string, findData bson.M, sortStr string, start int64, findNum int64, in interface{}, selectItme ...bson.M) error {
	//logs.Debug("Mongo QueryMongoAllData db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession  %s:%s, no session", dbName, collName)
		return goError.NewError("没有 mongo session")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	query := coll.Find(findData)
	if len(selectItme) > 0 {
		query = query.Select(selectItme[0])
	}
	if sortStr != "" {
		query = query.Sort(sortStr)
	}
	if start > -1 {
		query = query.Skip(int(start))
	}
	if findNum > -1 {
		query = query.Limit(int(findNum))
	}
	st := time.Now()
	err := query.All(in)
	if err != nil {
		stat.ReportStat(fmt.Sprintf("QueryMongoAllData-%s", collName), 1, time.Now().Sub(st))
		return err
	}
	stat.ReportStat(fmt.Sprintf("QueryMongoAllData-%s", collName), 0, time.Now().Sub(st))
	return err
}

//走写库
func RealQueryMongoOne(dbName string, collName string, findData bson.M, sortStr string) (bson.M, error) {
	//logs.Debug("Mongo RealQueryMongoOne db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession   %s:%s, no session", dbName, collName)
		return nil, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	ret := bson.M{}
	coll := session.DB(dbName).C(collName)
	if sortStr == "" {
		err := coll.Find(findData).Limit(1).One(&ret)
		return ret, err
	}
	//0
	err := coll.Find(findData).Sort(sortStr).Limit(1).One(&ret)
	return ret, err
}

//mongo 查询 走读库
func RealQueryMongoAll(dbName string, collName string, findData bson.M, sortStr string, start int64, findNum int64, selectItme ...bson.M) ([]bson.M, error) {
	//logs.Debug("Mongo RealQueryMongoAll db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession  %s:%s, no session", dbName, collName)
		return nil, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	ret := []bson.M{}
	coll := session.DB(dbName).C(collName)
	query := coll.Find(findData)
	if len(selectItme) > 0 {
		query = query.Select(selectItme[0])
	}
	if sortStr != "" {
		query = query.Sort(sortStr)
	}
	if start > -1 {
		query = query.Skip(int(start))
	}
	if findNum > -1 {
		query = query.Limit(int(findNum))
	}
	err := query.All(&ret)
	return ret, err
}

//sum 操作
func QueryMongoSum(dbName string, tableName string, whereData []bson.M) (map[string]int64, error) {
	//logs.Debug("Mongo QueryMongoSum db==%s tb==%s where==%s ", dbName, tableName, whereData)
	ret := map[string]int64{}
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession  %s:%s, no session", dbName, tableName)
		return ret, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	coll := session.DB(dbName).C(tableName)
	err := coll.Pipe(whereData).One(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

//sum 全部 操作
func QueryMongoSumAll(dbName string, tableName string, whereData []bson.M) ([]bson.M, error) {
	//logs.Debug("Mongo QueryMongoSumAll db==%s tb==%s where==%s ", dbName, tableName, whereData)
	ret := []bson.M{}
	session := GetOnlyReadSession()
	if session == nil {
		logs.Error("GetOnlyReadSession  %s:%s, no session", dbName, tableName)
		return ret, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	coll := session.DB(dbName).C(tableName)
	err := coll.Pipe(whereData).All(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

//查询去重
func QueryMongoDist(dbName string, collName string, findData bson.M, dist string) ([]string, error) {
	//logs.Debug("Mongo QueryMongoDist db==%s tb==%s where==%s ", dbName, collName, findData)
	ret := []string{}
	session := GetOnlyReadSession()
	if session == nil {
		return ret, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	err := coll.Find(findData).Distinct(dist, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

//查询count
func QueryMongoCount(dbName string, collName string, findData bson.M) (int64, error) {
	//logs.Debug("Mongo QueryMongoCount db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetOnlyReadSession()
	if session == nil {
		return 0, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	coll := session.DB(dbName).C(collName)
	nCount, err := coll.Find(findData).Count()
	if err != nil {
		return 0, err
	}
	return int64(nCount), nil
}

//删除表
func DelTable(dbName string, collName string) error {
	session := GetOnlyReadSession()
	if session == nil {
		return errors.New("No MgoDb")
	}
	defer session.Close()
	return session.DB(dbName).C(collName).DropCollection()
}

//走写库查询，唯一使用
func NewQueryMongoAll(dbName string, collName string, findData bson.M, sortStr string, start int64, findNum int64, selectItme ...bson.M) ([]bson.M, error) {
	//logs.Debug("Mongo RealQueryMongoAll db==%s tb==%s where==%s ", dbName, collName, findData)
	session := GetRealSession()
	if session == nil {
		logs.Error("GetOnlyReadSession  %s:%s, no session", dbName, collName)
		return nil, goError.NewError("没有 mongo session")
	}
	defer session.Close()
	ret := []bson.M{}
	coll := session.DB(dbName).C(collName)
	query := coll.Find(findData)
	if len(selectItme) > 0 {
		query = query.Select(selectItme[0])
	}
	if sortStr != "" {
		query = query.Sort(sortStr)
	}
	if start > -1 {
		query = query.Skip(int(start))
	}
	if findNum > -1 {
		query = query.Limit(int(findNum))
	}
	err := query.All(&ret)
	return ret, err
}