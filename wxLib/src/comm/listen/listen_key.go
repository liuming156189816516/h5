package listen

import "natsRpc"

//广播 事件
const (
	TaskUpPlanStatusListen = "TaskUpPlanStatusListen" //修改方案执行状态

	TaskUpGroupPlanStatusListen = "TaskUpGroupPlanStatusListen" //修改群方案执行状态

	TaskUpRobotPlanStatusListen = "TaskUpRobotPlanStatusListen" //修改机器人方案执行状态
)

//修改方案某个微信号执行状态
type TaskUpPlanStatusListenReq struct {
	PlanTaskId string //方案任务id
	Status     int64  //1 启动 2关闭
}

func (s *TaskUpPlanStatusListenReq) Send() {
	natsRpc.NotifyReloadConfig(TaskUpPlanStatusListen, s)
}

//删除
type TaskUpGroupPlanStatusListenReq struct {
	GroupTaskId string //方案任务id
}

func (s *TaskUpGroupPlanStatusListenReq) Send() {
	natsRpc.NotifyReloadConfig(TaskUpGroupPlanStatusListen, s)
}

type TaskUpRobotPlanStatusListenReq struct {
	RobotTaskId string
}

//关闭
func (s *TaskUpRobotPlanStatusListenReq) Send() {
	natsRpc.NotifyReloadConfig(TaskUpRobotPlanStatusListen, s)
}



//广播 事件
const (
	QQTaskUpPlanStatusListen = "QQTaskUpPlanStatusListen" //修改方案执行状态

	QQTaskUpGroupPlanStatusListen = "QQTaskUpGroupPlanStatusListen" //修改群方案执行状态

	QQTaskUpRobotPlanStatusListen = "QQTaskUpRobotPlanStatusListen" //修改机器人方案执行状态
)

//修改方案某个微信号执行状态
type QQTaskUpPlanStatusListenReq struct {
	PlanTaskId string //方案任务id
	Status     int64  //1 启动 2关闭
}

func (s *QQTaskUpPlanStatusListenReq) Send() {
	natsRpc.NotifyReloadConfig(QQTaskUpPlanStatusListen, s)
}

//删除
type QQTaskUpGroupPlanStatusListenReq struct {
	GroupTaskId string //方案任务id
}

func (s *QQTaskUpGroupPlanStatusListenReq) Send() {
	natsRpc.NotifyReloadConfig(QQTaskUpGroupPlanStatusListen, s)
}

type QQTaskUpRobotPlanStatusListenReq struct {
	RobotTaskId string
}

//关闭
func (s *QQTaskUpRobotPlanStatusListenReq) Send() {
	natsRpc.NotifyReloadConfig(QQTaskUpRobotPlanStatusListen, s)
}
