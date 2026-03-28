package server

import (
	"comm/event"
	"comm/redisDeal"
	"comm/redisKeys"
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"task/models/CallTask"
	"time"
)

// 任务
func RunTask() { //发送任务的 池
	go func() {
		time.Sleep(2 * time.Second)
		ticker := time.NewTicker(2 * time.Second) //2秒取一次任务
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				go doTask(t)              //任务事件
				go doMessageTaskResult(t) //处理消息发送结果
				go doAcceptCallTask(t)    //处理接听电话回调
				go doLoginTask(t)         //登录事件
				go doTaskZhaGroup(t)      //处理炸群
				go doSendmsg(t)           //处理私发
			}
		}
	}()
}

func doLoginTask(t time.Time) {
	key := redisKeys.GetAllTaskLoginList()
	lkey := redisDeal.RedisDoLLen(key)
	logs.Info("doLoginTask lkey=======>", lkey)
	if lkey == 0 { //没有任务
		return
	}
	num := 1
	if lkey > 10 {
		num = 3
	}
	if lkey > 100 {
		num = 20
	}
	if lkey > 1000 {
		num = 200
	}
	if lkey > 10000 {
		num = 500
	}
	for i := 0; i < num; i++ {
		time.Sleep(4 * time.Millisecond)
		str := redisDeal.RedisDoRpop(key)
		if str == "" {
			continue
		}
		data := &event.TaskInfo{}
		err := jsoniter.UnmarshalFromString(str, &data)
		if err != nil {
			continue
		}
		go CallTask.CallUserOneTask(data.Phone, data.TaskType, data.Data)
	}
}

func doTask(t time.Time) {
	key := redisKeys.GetAllTaskList()
	lkey := redisDeal.RedisDoLLen(key)
	logs.Info("doTask lkey=======>", lkey)
	if lkey == 0 { //没有任务
		return
	}
	num := 1
	if lkey > 10 {
		num = 3
	}
	if lkey > 100 {
		num = 20
	}
	if lkey > 1000 {
		num = 200
	}
	if lkey > 10000 {
		num = 500
	}
	for i := 0; i < num; i++ {
		time.Sleep(4 * time.Millisecond)
		str := redisDeal.RedisDoRpop(key)
		if str == "" {
			continue
		}
		data := &event.TaskInfo{}
		err := jsoniter.UnmarshalFromString(str, &data)
		if err != nil {
			continue
		}
		go CallTask.CallUserOneTask(data.Phone, data.TaskType, data.Data)
	}
}

// 炸群
func doTaskZhaGroup(t time.Time) {
	key := redisKeys.GetAllTaskZhaGroupList()
	lkey := redisDeal.RedisDoLLen(key)
	logs.Info("lkey=======>", lkey)
	if lkey == 0 { //没有任务
		return
	}
	//上正式要把下面改成5
	num := 3
	//if lkey > 10 {
	//	num = 3
	//}
	//if lkey > 100 {
	//	num = 20
	//}
	//if lkey > 1000 {
	//	num = 200
	//}
	//if lkey > 10000 {
	//	num = 500
	//}
	for i := 0; i < num; i++ {
		time.Sleep(4 * time.Millisecond)
		str := redisDeal.RedisDoRpop(key)
		if str == "" {
			continue
		}
		data := &event.TaskInfo{}
		err := jsoniter.UnmarshalFromString(str, &data)
		if err != nil {
			continue
		}
		go CallTask.CallUserOneTask(data.Phone, data.TaskType, data.Data)
	}
}

func doMessageTaskResult(t time.Time) {
	key := redisKeys.GetAllTaskMessageResultList()
	lkey := redisDeal.RedisDoLLen(key)
	if lkey == 0 { //没有任务
		return
	}
	num := 1
	if lkey > 10 {
		num = 3
	}
	if lkey > 100 {
		num = 20
	}
	if lkey > 1000 {
		num = 200
	}
	if lkey > 10000 {
		num = 500
	}
	for i := 0; i < num; i++ {
		time.Sleep(4 * time.Millisecond)
		str := redisDeal.RedisDoRpop(key)
		if str == "" {
			continue
		}
		data := &event.TaskInfo{}
		err := jsoniter.UnmarshalFromString(str, &data)
		if err != nil {
			continue
		}
		go CallTask.CallDataTask(data.Phone, data.TaskType, data.Data)
	}
}

// 处理接听电话回调
func doAcceptCallTask(t time.Time) {
	key := redisKeys.GetAllTaskAcceptCallList()
	lkey := redisDeal.RedisDoLLen(key)
	if lkey == 0 { //没有任务
		return
	}
	num := 1
	if lkey > 10 {
		num = 3
	}
	if lkey > 100 {
		num = 20
	}
	if lkey > 1000 {
		num = 200
	}
	if lkey > 10000 {
		num = 500
	}
	for i := 0; i < num; i++ {
		time.Sleep(4 * time.Millisecond)
		str := redisDeal.RedisDoRpop(key)
		if str == "" {
			continue
		}
		data := &event.TaskInfo{}
		err := jsoniter.UnmarshalFromString(str, &data)
		if err != nil {
			continue
		}
		go CallTask.CallDataTask(data.Phone, data.TaskType, data.Data)
	}
}

// 处理私发任务
func doSendmsg(t time.Time) {
	key := redisKeys.GetAllTaskSendmsgList()
	lkey := redisDeal.RedisDoLLen(key)
	if lkey == 0 { //没有任务
		return
	}
	logs.Info("doSendmsg lkey==========>", lkey)
	num := 30
	/*if lkey > 10 {
		num = 3
	}
	if lkey > 100 {
		num = 20
	}
	if lkey > 1000 {
		num = 200
	}
	if lkey > 10000 {
		num = 500
	}*/
	for i := 0; i < num; i++ {
		time.Sleep(4 * time.Millisecond)
		str := redisDeal.RedisDoRpop(key)
		if str == "" {
			continue
		}
		data := &event.TaskInfo{}
		err := jsoniter.UnmarshalFromString(str, &data)
		if err != nil {
			continue
		}
		go CallTask.CallUserOneTask(data.Phone, data.TaskType, data.Data)
	}
}
