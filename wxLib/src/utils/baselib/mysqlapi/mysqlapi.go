package mysqlapi

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"qqgame/baselib/logs"
	"strings"
	"time"
)

const (
	DEFAULT_MAX_IDLE_CONNS = 4
	DEFAULT_MAX_OPEN_CONNS = 32

	MYSQL_MAX_IDLE_CONNS = 64
	MYSQL_MAX_OPEN_CONNS = 128
)

type MySQLAPI struct {
	OrmInstance orm.Ormer
}

/*
	@desc: 创建与mysql的长连接
	@params:
		uname: mysql username;
		ip:    mysql ip;
		port:  mysql port;
		dbname:mysql database name;

	@options params:
		params1: max idle connect number;
		params2: max open connect number; [ 注意idle/open connect number不能设置太大 ]

	@returns:
		*MySQLAPI: db handler;
		error: ...
*/
func NewMySQLAPI(uname, passwd, ip string, port int, dbname string, params ...int) (*MySQLAPI, error) {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	dataSrc := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		uname, passwd, ip, port, dbname)

	var maxIdleConns = int(0)
	var maxOpenConns = int(0)
	for i, v := range params {
		switch i {
		case 0:
			maxIdleConns = v
		case 1:
			maxOpenConns = v
		}
	}
	if maxIdleConns == 0 || maxIdleConns > MYSQL_MAX_IDLE_CONNS {
		maxIdleConns = DEFAULT_MAX_IDLE_CONNS
	}
	if maxOpenConns == 0 || maxOpenConns > MYSQL_MAX_OPEN_CONNS {
		maxOpenConns = DEFAULT_MAX_OPEN_CONNS
	}

	if err := orm.RegisterDataBase("default", "mysql", dataSrc, maxIdleConns, maxOpenConns); err != nil {
		logs.ERRORLOG("orm.RegisterDataBase failed! err: %s", err.Error())
		return nil, err
	}

	return &MySQLAPI{OrmInstance: orm.NewOrm()}, nil
}

/*
	@desc:
		检查传入语句是否存在风险SQL操作；
		SQL语句只允许insert/select/update/replace/delete

	@params:
		string: SQL string
	@return:
		bool
*/
func (api *MySQLAPI) isValidStmt(s string) bool {
	if len(s) == 0 {
		return false
	}
	stPermitActions := []string{"INSERT", "SELECT", "UPDATE", "REPLACE", "DELETE"}
	for _, v := range stPermitActions {
		if strings.HasPrefix(s, v) || strings.HasPrefix(s, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

/*
	@desc: 执行mysql select查询语句
	@params:
		stmt: SQL查询语句，注意预处理防止SQL注入
	@returns:
		[]orm.Params: orm数组对象。注意，所有字段值均为string类型；
		int: 查询结果条数；
		error:
*/
func (api *MySQLAPI) ExecQuery(stmt string, args ...interface{}) ([]orm.Params, int, error) {
	var result []orm.Params

	if !api.isValidStmt(stmt) {
		return result, -1, fmt.Errorf("SQL stmt is not allowed, sql: %s", stmt)
	}

	rownums, err := api.OrmInstance.Raw(stmt, args).Values(&result)
	if err != nil {
		logs.ERRORLOG("ExecQuery failed! err: %s", err.Error())
		return result, -1, err
	}
	return result, int(rownums), nil
}

/*
	@desc: 执行mysql update/insert/delete更新语句
	@params:
		stmt: SQL查询语句，注意预处理防止SQL注入。
	@returns:
		int: 影响mysql表行数
		int: 最后插入的id
		error:
*/
func (api *MySQLAPI) ExecUpdate(stmt string, args ...interface{}) (int, int, error) {
	if !api.isValidStmt(stmt) {
		return -1, -1, fmt.Errorf("SQL stmt is not allowed, sql: %s", stmt)
	}

	res, err := api.OrmInstance.Raw(stmt, args).Exec()
	if err != nil {
		logs.ERRORLOG("ExecUpdate failed! err: %s", err.Error())
		return -1, -1, err
	}

	numrow, _ := res.RowsAffected()
	lastId, _ := res.LastInsertId()
	return int(numrow), int(lastId), nil
}

/*
	@desc: 获取取初始化后得到的ormer对象，可用些对象完成原生操作
	@returns:
		(&orm.Ormer, error)
	@reference:
		https://beego.me/docs/mvc/model/overview.md
*/
func (api *MySQLAPI) GetOrmInstance() (orm.Ormer, error) {
	if api != nil {
		return api.OrmInstance, nil
	}
	return nil, fmt.Errorf("mysql api is not initialized")
}

/*
	@desc: 一组事务操作的封装
*/
func (api *MySQLAPI) ExecTransBegin() error {
	return api.OrmInstance.Begin()
}

func (api *MySQLAPI) ExecTransUpdate(stmt string, args ...interface{}) (int, int, error) {
	res, err := api.OrmInstance.Raw(stmt, args).Exec()
	if err != nil {
		return -1, -1, err
	}

	rowafct, _ := res.RowsAffected()
	lastid, _ := res.LastInsertId()
	return int(rowafct), int(lastid), nil
}

func (api *MySQLAPI) ExecTransCommit() error {
	return api.OrmInstance.Commit()
}

func (api *MySQLAPI) ExecTransRollback() error {
	return api.OrmInstance.Rollback()
}

/*---end---*/

type StmtAndArgs struct {
	Stmt string      // SQL语句
	Args interface{} // SQL参数
}

/*
	@desc: 执行mysql事务
	@params:
		stmt: []StmtAndArgs, 一组打包的SQL事务查询语句，注意预处理防止SQL注入。
	@returns:
		int: 影响mysql表行数
		error:
*/
func (api *MySQLAPI) ExecTrans(stmtArray []StmtAndArgs) (int, error) {
	var errorinfo error

	err := api.OrmInstance.Begin()
	if err != nil {
		logs.ERRORLOG("o.Begin() failed, err: %s", err.Error())
	}

	isExecFailed := false
	for i := 0; i < len(stmtArray); i++ {

		if !api.isValidStmt(stmtArray[i].Stmt) {
			return -1, fmt.Errorf("SQL stmt is not allowed, sql: %s", stmtArray[i].Stmt)
		}

		var tmpres []orm.Params

		if stmtArray[i].Args == nil || stmtArray[i].Args == "" {
			if _, err := api.OrmInstance.Raw(stmtArray[i].Stmt).Values(&tmpres); err != nil {
				logs.ERRORLOG("ExecTrans failed! current exec stmt: %s, args: %+v err: %s",
					stmtArray[i].Stmt, stmtArray[i].Args, err.Error())

				isExecFailed = true
				errorinfo = err
				break
			}
		} else {
			if _, err := api.OrmInstance.Raw(stmtArray[i].Stmt, stmtArray[i].Args).Values(&tmpres); err != nil {
				logs.ERRORLOG("ExecTrans failed! current exec stmt: %s, args: %+v err: %s",
					stmtArray[i].Stmt, stmtArray[i].Args, err.Error())

				isExecFailed = true
				errorinfo = err
				break
			}
		}
	}

	if isExecFailed {
		api.OrmInstance.Rollback()
		return -1, errorinfo
	} else {
		api.OrmInstance.Commit()
		return 0, nil
	}
}

/*
	@desc: 采用redis来做一个全局的锁服务；
			确保在ttl(5s)时间内全局只允许一个服务执行mysql更新操作
			todo: 先调用此函数初始化redis
	@params:
		redisIP: redis IP address;
		redisport: redis port;
		passwd:  redis passwd
	@returns:
		int: 影响mysql表行数；
		error:
*/
var gRedisPool *redis.Pool

func (api *MySQLAPI) EnableGlobalLock(redisIP string, redisport int, passwd string) {
	gRedisPool = api.initRedis(redisIP, redisport, passwd)
}

func (api *MySQLAPI) GlobalLockExist() error {
	if gRedisPool == nil {
		logs.ERRORLOG("you should call EnableMySQLGlobalLock() first")
		return fmt.Errorf("please call EnableMySQLGlobalLock() first")
	}
	return nil
}

/*
	@desc: 执行mysql update/insert/delete更新语句，保证在5秒内全局只有一个服务可以执行更新操作
	@params:
		stmt: SQL查询语句，注意预处理防止SQL注入；
		lockname: 锁名称。例如，传入用户的QQ号；
	@returns:
		int: 影响mysql表行数；
		error:
*/
func (api *MySQLAPI) ExecUpdateWithLock(lockname, stmt string, args ...interface{}) (int, int, error) {
	if gRedisPool == nil {
		logs.ERRORLOG("you should call EnableMySQLGlobalLock() first")
		return -10, 0, fmt.Errorf("please call EnableMySQLGlobalLock() first")
	}

	if err := api.LockGlobal(lockname); err != nil {
		logs.ERRORLOG("ExecUpdateWithLock() add global lock failed! err: %s", err.Error())
		return -1, -1, err
	}

	ret, lastid, err := api.ExecUpdate(stmt, args)
	if err != nil {
		logs.ERRORLOG("ExecUpdateWithLock failed! err: %s", err.Error())
		//return -2, err
	}

	if err := api.UnLockGlobal(lockname); err != nil {
		logs.ERRORLOG("ExecUpdateWithLock() unlock failed! lock would be auto expire for 10 seconds. err: %s", err.Error())
		return -2, -1, err
	}
	return ret, lastid, nil
}

/*
	@desc: 执行mysql事务更新，保证在5秒内全局只有一个服务可以执行更新操作
	@params:
		stmt: []StmtAndArgs, 一组打包的SQL事务查询语句，注意预处理防止SQL注入；
		lockname: 锁名称。例如，传入用户的QQ号；
	@returns:
		int: 影响mysql表行数；
		error:
*/
func (api *MySQLAPI) ExecTransWithLock(lockname string, stmtArr []StmtAndArgs) (int, error) {
	if gRedisPool == nil {
		logs.ERRORLOG("you should call EnableMySQLGlobalLock() first")
		return -10, fmt.Errorf("please call EnableMySQLGlobalLock() first")
	}

	if err := api.LockGlobal(lockname); err != nil {
		logs.ERRORLOG("ExecTransWithLock() failed! add global lock failed! err: %s", err.Error())
		return -1, err
	}

	ret, err := api.ExecTrans(stmtArr)
	if err != nil {
		logs.ERRORLOG("ExecTransWithLock failed! err: %s", err.Error())
		//return -2, err
	}

	if err := api.UnLockGlobal(lockname); err != nil {
		logs.ERRORLOG("ExecTransWithLock() failed! lock would be auto expire for 10 seconds. err: %s", err.Error())
		return -2, err
	}
	return ret, nil
}

const (
	REDIS_MAX_IDLE   = 2
	REDIS_MAX_ACTIVE = 16

	REDIS_LOCK_KEY_PREFIX   = "mysqlLock_%s" // 全局锁key name
	REDIS_LOCK_KEYVALUE_TTL = 5              // 全局锁，锁定默认的最大时长为5秒
)

/*
	初始化redis连接池，采用redis作为全局的更新锁
*/
func (api *MySQLAPI) initRedis(ip string, port int, passwd string, params ...int) *redis.Pool {
	var maxIdleConns = int(0)
	var maxOpenConns = int(0)
	for i, v := range params {
		switch i {
		case 0:
			maxIdleConns = v
		case 1:
			maxOpenConns = v
		}
	}
	if maxIdleConns == 0 {
		maxIdleConns = REDIS_MAX_IDLE
	}
	if maxOpenConns == 0 {
		maxOpenConns = REDIS_MAX_ACTIVE
	}

	server := fmt.Sprintf("%s:%d", ip, port)
	return &redis.Pool{
		MaxIdle:     maxIdleConns,
		MaxActive:   maxOpenConns, //when zero,there's no limit. https://godoc.org/github.com/garyburd/redigo/redis#Pool
		IdleTimeout: time.Duration(240 * time.Second),
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if _, err = c.Do("AUTH", passwd); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

/*
	加全局锁。用于保证一次完整的事务操作，防止重复提交。默认锁定5秒，5秒后锁自动解除。
*/
func (api *MySQLAPI) LockGlobal(lkname string) error {
	lk := fmt.Sprintf(REDIS_LOCK_KEY_PREFIX, lkname)
	conn := gRedisPool.Get()
	defer conn.Close()

	value := "LOCK"
	_, err := redis.String(conn.Send("SET", lk, value, "EX", REDIS_LOCK_KEYVALUE_TTL, "NX"))
	if err != nil {
		if err == redis.ErrNil {
			logs.DEBUGLOG("Lock is owned by others")
			return fmt.Errorf("Lock is owned by others")
		} else {
			logs.DEBUGLOG("GetUserLock():Do Redis Cmd=SET %s has err=%v", lk, err)
			return fmt.Errorf("Do Redis Cmd=SET %s has err=%v", lk, err)
		}
	}
	return nil
	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//原来版本保存
	/*cmdTouch := fmt.Sprintf("TOUCH %s", lk)
	ret, err := redis.Int(conn.Do(cmdTouch))
	if err != nil {
		logs.ERRORLOG("GlobalLock failed! cmd: %s, err: %s", cmdTouch, err.Error())
		return err
	}

	if ret == 0 { // 未被全局锁定
		value := "LOCK"
		cmd := fmt.Sprintf("SET %s %s", lk, value)
		cmd2 := fmt.Sprintf("EXPIRE %s %d", lk, REDIS_LOCK_KEYVALUE_TTL)
		if _, err := conn.Do(cmd); err != nil {
			logs.ERRORLOG("GlobalLock failed! cmd: %s, err: %s", cmd, err.Error())
			return err
		}
		if _, err := conn.Do(cmd2); err != nil {
			logs.ERRORLOG("GlobalLock failed! cmd: %s, err: %s", cmd2, err.Error())
			return err
		}
		return nil
	} else {
		return fmt.Errorf("already locked by another task, please wait")
	}*/
	//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
}

/*
	解除全局锁
*/
func (api *MySQLAPI) UnLockGlobal(lkname string) error {
	lk := fmt.Sprintf(REDIS_LOCK_KEY_PREFIX, lkname)
	conn := gRedisPool.Get()
	defer conn.Close()

	cmd := fmt.Sprintf("DEL %s", lk)
	if _, err := conn.Do("DEL", lk); err != nil {
		logs.ERRORLOG("GlobalUnLock failed! cmd: %s, err: %s", cmd, err.Error())
		return err
	}
	return nil
}
