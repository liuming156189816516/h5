package natsRpc

import (
	"mframe/stat"
	"time"
	"fmt"
)


func ReportCallRpcStat(serviceName string, funcName string, result int32, processTime time.Duration) {
	key := fmt.Sprintf("call-nrpc-%s-%s", serviceName, funcName)
	stat.ReportStat(key, int(result), processTime)
}
func ReportDoRpcStat(serviceName string, funcName string, result int32, processTime time.Duration) {
	key := fmt.Sprintf("do-nrpc-%s-%s", funcName, serviceName)
	stat.ReportStat(key, int(result), processTime)
}
func ReportStat(object string, funcName string, result int, processTime time.Duration) {
	key := fmt.Sprintf("%s-%s", object, funcName)
	stat.ReportStat(key, int(result), processTime)
}

