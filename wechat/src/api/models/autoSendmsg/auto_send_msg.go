package autoSendmsg

import (
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"time"
)

func TaskRun() {
	go func() {
		//自动群发任务
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(10 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				val := GetStaticData("doAutoSendMsg")
				if val != "" {
					//正在执行
					continue
				}
				if cache.GetTaskStatus() == "1" {
					continue
				}
				go doAutoSendMsg(t, "doAutoSendMsg")
			}
		}
	}()
}

// 自动群发任务
func doAutoSendMsg(t time.Time, lock string) {
	defer DelStaticData(lock)
	SetStaticData(lock, lock)

	lkey := cache.LenAutoSendMsgTaskInfo()
	if lkey == 0 { //没有任务
		return
	}

	num := 50
	for i := 0; i < num; i++ {
		taskInfo := cache.GetAutoSendMsgTaskInfo()
		if taskInfo.Account == "" {
			return
		}
		go wxComm.AutoSendMsg(taskInfo.Account, taskInfo.SessionId, taskInfo.Node)
	}
}
