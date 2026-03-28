package event

import (
	"comm/redisDeal"
	"comm/redisKeys"
	"fmt"
	"github.com/json-iterator/go"
	"time"
)

// 任务类型
const (
	TaskTypeUserLogin  = 1 //用户登录
	TaskTypeUserLogout = 2 //下线
	TaskTypeShopName   = 3 //设置商店名称
	TaskTypeUserMessageResult      = 5  //处理消息发送结果
	TaskTypeUserSendMsg            = 6  //发送消息
	TaskTypeUserMessageCallBack    = 7  //处理消息回调
	TaskTypeUserAcceptCallCallBack = 8  //处理接听电话回调
	TaskTypeUserRemove         = 9  //账号移除
	TaskTypeUserAccountHeadimg = 10 //设置头像
	TaskTypeUserContactHeadimg = 11 //获取头像
	TaskTypeVerifyTwostep      = 13 //设置商店名称
	TaskTypeZhaGroup           = 17 //炸群
)

// 任务数据
type TaskInfoData struct {
	Phone    string
	TaskType int64
	Data     interface{}
}

// 任务发送结构体
type TaskInfo struct {
	Phone    string
	TaskType int64
	Data     map[string]interface{}
}

// 添加登录任务
func AddLoginTask(taskType int64, phone string, data interface{}) {

	str, ok := data.(string)
	if ok { //字符串 需要转换成功结构体
		tmp := map[string]interface{}{}
		jsoniter.UnmarshalFromString(str, &tmp)
		data = tmp
	}

	//写redis 里面  等任务 进程 完成
	task := TaskInfoData{
		TaskType: taskType,
		Data:     data,
		Phone:    phone,
	}
	key := redisKeys.GetAllTaskLoginList()
	redisDeal.RedisSendLpush(key, task)
	time.Sleep(2 * time.Millisecond)
}

// 添加任务 datas 必须是 task_event 文件里面的 结构体
func AddTask(taskType int64, phone string, data interface{}) {
	str, ok := data.(string)
	if ok { //字符串 需要转换成功结构体
		tmp := map[string]interface{}{}
		jsoniter.UnmarshalFromString(str, &tmp)
		data = tmp
	}

	//写redis 里面  等任务 进程 完成
	task := TaskInfoData{
		TaskType: taskType,
		Data:     data,
		Phone:    phone,
	}
	key := redisKeys.GetAllTaskList()
	redisDeal.RedisSendLpush(key, task)
	time.Sleep(2 * time.Millisecond)

}

// 获取任务id
func GetTaskId(task_type int64) string {
	if task_type == 0 {
		return ""
	}
	key := redisKeys.GetAdminIncInfo()
	i, err := redisDeal.RedisDoHincrby(key, fmt.Sprintf("task_%d", task_type), 1)
	if err != nil || i <= 0 {
		return ""
	}
	return fmt.Sprintf("%d_%d", task_type, i)
}

// 处理消息结果回调
func AddMessageTaskResult(taskType int64, phone string, data interface{}) {

	str, ok := data.(string)
	if ok { //字符串 需要转换成功结构体
		tmp := map[string]interface{}{}
		jsoniter.UnmarshalFromString(str, &tmp)
		data = tmp
	}

	//写redis 里面  等任务 进程 完成
	task := TaskInfoData{
		TaskType: taskType,
		Data:     data,
		Phone:    phone,
	}
	key := redisKeys.GetAllTaskMessageResultList()
	redisDeal.RedisSendLpush(key, task)
	time.Sleep(2 * time.Millisecond)
}

// 处理对方接听电话
func AddAcceptCallTask(taskType int64, phone string, data interface{}) {

	str, ok := data.(string)
	if ok { //字符串 需要转换成功结构体
		tmp := map[string]interface{}{}
		jsoniter.UnmarshalFromString(str, &tmp)
		data = tmp
	}

	//写redis 里面  等任务 进程 完成
	task := TaskInfoData{
		TaskType: taskType,
		Data:     data,
		Phone:    phone,
	}
	key := redisKeys.GetAllTaskAcceptCallList()
	redisDeal.RedisSendLpush(key, task)
	time.Sleep(2 * time.Millisecond)
}

// 处理炸群任务
func AddZhaGroupTask(taskType int64, phone string, data interface{}) {
	str, ok := data.(string)
	if ok { //字符串 需要转换成功结构体
		tmp := map[string]interface{}{}
		jsoniter.UnmarshalFromString(str, &tmp)
		data = tmp
	}

	//写redis 里面  等任务 进程 完成
	task := TaskInfoData{
		TaskType: taskType,
		Data:     data,
		Phone:    phone,
	}
	key := redisKeys.GetAllTaskZhaGroupList()
	redisDeal.RedisSendLpush(key, task)
	time.Sleep(2 * time.Millisecond)
}

// 添加任务 datas 必须是 task_event 文件里面的 结构体
func AddSendmsgTask(taskType int64, phone string, data interface{}) {
	str, ok := data.(string)
	if ok { //字符串 需要转换成功结构体
		tmp := map[string]interface{}{}
		jsoniter.UnmarshalFromString(str, &tmp)
		data = tmp
	}

	//写redis 里面  等任务 进程 完成
	task := TaskInfoData{
		TaskType: taskType,
		Data:     data,
		Phone:    phone,
	}
	key := redisKeys.GetAllTaskSendmsgList()
	redisDeal.RedisSendLpush(key, task)
	time.Sleep(2 * time.Millisecond)

}

