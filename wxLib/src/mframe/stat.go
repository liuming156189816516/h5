package mframe

import (
	"mframe/stat"
	"os"
	"fmt"
	"time"
)

var hostName string

func init() {

	// 消息注册
	hostName, _ = os.Hostname()
}
func ReportCallRpcStat(serviceName string, funcName string, result int32, processTime time.Duration) {
	key := fmt.Sprintf("call-nats-%s-%s", serviceName, funcName)
	stat.ReportStat(key, int(result), processTime)
}
func ReportDoRpcStat(serviceName string, funcName string, result int32, processTime time.Duration) {
	key := fmt.Sprintf("do-nats-%s-%s", funcName, serviceName)
	stat.ReportStat(key, int(result), processTime)
}
func ReportStat(object string, funcName string, result int, processTime time.Duration) {
	key := fmt.Sprintf("%s-%s", object, funcName)
	stat.ReportStat(key, int(result), processTime)
}

/*
func statisicDataReport(req *monitor.StatisticReqMsgPara) (*monitor.StatisticRspMsgPara, error) {
b := &msg.PbMsg{req}
	// SendAndRecv
	session := NewEmptySession()
	rpc := NewRpc(session)
	rsp, err := rpc.NatsCall(session.GetUin(), "monitor", "StatisicDataReport", rb, 2)
	if err != nil {
		logs.Trace("monitor StatisicDataReport failed:%s", err)
		return nil, errors.New("NatsCall Error")
	}
	// 转成PbMsg
	if rsp, ok := rsp.(*msg.PbMsg); ok {
		// 得到PbMsg里的Pb结构
		if pb, ok := rsp.Pb.(*monitor.StatisticRspMsgPara); ok {
			return pb, nil
		} else {
			return nil, errors.New("Rsp Error") //
		}
	} else {
		logs.Trace("Recv Error Msg")
		return nil, errors.New("Rsp Error") //
	}
}
*/

func additionalMsgStat(key string, value *stat.MsgStatData, avgProcessTime time.Duration, avgSuccProcessTime time.Duration) {
	// 组包上报统计数据

	//value.Host = hostName
	//value.ServerName = GetServerName()
	//value.ServerId = GetServerID()
	//value.Key = key
	//
	//value.AvgProcessTime = avgProcessTime
	//value.AvgSuccProcessTime = avgSuccProcessTime
	//
	//reqMsg := &NatsMsg{}
	//reqMsg.Sess = Session{}
	//reqMsg.MsgBody.Func = "ReportState"
	//reqMsg.MsgBody.Check = time.Now().Format("2006-01-02 15:04:05.999")
	//
	//reqMsg.MsgBody.Param, _ = jsoniter.Marshal(value)
	//
	//reqMsg.Sess.SvrFE = GetServerName()
	//reqMsg.Sess.SvrID = GetServerID()
	//reqMsg.Sess.Time = time.Now().Unix()
	////	req.Sess.Route = fmt.Sprint("%s->%s.%d", req.Sess.Route, GetServerName(), GetServerID())
	//subj := GenReqSubject("monitor", "ReportState", -1)
	//data, err := jsoniter.Marshal(reqMsg)
	//if err != nil {
	//	logs.LogError("NatsCall Marshal req (%+v) failed:%s", reqMsg, err.Error())
	//	return
	//}
	//g_natsConn.Publish(subj, data)
}
