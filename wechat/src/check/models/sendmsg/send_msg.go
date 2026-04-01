package sendmsg

import (
	"comm/event"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/sendmsg"
	"selfComm/wxComm/cache"
	"time"
)

/*func TaskRun() {

	go func() {
		//生产群发队列
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(10 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				val := GetStaticData("createSendTask")
				if val != "" {
					//正在执行
					continue
				}
				go createSendTask(t, "createSendTask")
			}
		}
	}()

	go func() {
		//群发任务
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(10 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				val := GetStaticData("doSendMsg")
				if val != "" {
					//正在执行
					continue
				}
				go doSendMsg(t, "doSendMsg")
			}
		}
	}()
}*/

func TaskRun() {
	time.Sleep(2 * time.Second)
	createSendTask(time.Now(), "createSendTask")

	time.Sleep(10 * time.Second)
	doSendMsg(time.Now(), "doSendMsg")
}

func createSendTask(t time.Time, lock string) {
	defer DelStaticData(lock)
	SetStaticData(lock, lock)

	config := cache.GetTaskConfig("global")
	if config.DataPackId == "" || len(config.MaterialList) <= 0 {
		logs.Info("createSendTask 任务配置为空")
		return
	}

	if cache.ScardDataPackListCount(config.DataPackId) <= 0 {
		logs.Info("createSendTask 粉丝数据不足")
		return
	}

	sendTaskList := sendmsg.GetListSendMsgInfo(bson.M{"account_status": int64(2)}, -1)
	for _, sendTask := range sendTaskList {
		tmp := &cache.SendMsgTaskInfo{}
		tmp.Account = sendTask.Account
		tmp.DataPackId = config.DataPackId
		tmp.MaterialList = config.MaterialList
		cache.SetSendMsgTaskInfo(tmp)
	}
	lkey := cache.LenSendMsgTaskInfo()
	logs.Info("createSendTask==========>num: ", lkey)
}

// 群发任务
func doSendMsg(t time.Time, lock string) {
	defer DelStaticData(lock)
	SetStaticData(lock, lock)

	lkey := cache.LenSendMsgTaskInfo()
	if lkey == 0 { //没有任务
		return
	}

	num := 100
	for i := 0; i < num; i++ {
		taskInfo := cache.GetSendMsgTaskInfo()
		dataNum := cache.ScardDataPackListCount(taskInfo.DataPackId)
		if dataNum <= 0 {
			return
		}
		materialListStr, _ := jsoniter.MarshalToString(taskInfo.MaterialList)
		e := &event.TaskTypeSendMsgEventReq{
			Account:         taskInfo.Account,
			DataPackId:      taskInfo.DataPackId,
			IsUp:            true,
			MaterialListStr: materialListStr,
		}
		event.AddSendmsgTask(event.TaskTypeUserSendMsg, taskInfo.Account, e)
	}
}
