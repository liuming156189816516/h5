package mysqlapi

import (
	"fmt"
	"testing"
)


var api = &MySQLAPI{}
func init() {
	var err error
	api, err = NewMySQLAPI("root", "", "localhost", 3306, "dbmyproject")
	if err != nil {
		fmt.Println(err)
	}
}

func TestMySQLAPI_ExecQuery(t *testing.T) {

	rtn, _, err := api.ExecQuery("SELECT * FROM `tbsystemusermgr` ORDER BY `sUserPasswd` ASC")
	if err != nil {
		t.Error(err)
	}

	t.Logf("rtn: %+v", rtn)

}

func TestMySQLAPI_ExecUpdate(t *testing.T) {

	rtn, err := api.ExecUpdate("UPDATE `tbsystemusermgr` SET `iUserCharacter`=2 WHERE `sUserName`= ?", "yanwei")
	if err != nil {
		t.Error(err)
	}

	t.Logf("rtn: %+v", rtn)
}



func TestMySQLAPI_ExecTrans(t *testing.T) {

	stmt1 :=  StmtAndArgs{
		Stmt: "select * from tbsystemusermgr where `sUserName`=?",
		Args: "yanwei"}

	stmt := StmtAndArgs{
		Stmt: "select * from tbsystemusermgr",
		Args: nil}

	stmts := []StmtAndArgs{stmt1, stmt}

	rtn, err := api.ExecTrans(stmts)
	if err != nil {
		t.Error(err)
	}

	t.Logf("rtn: %+v", rtn)
}