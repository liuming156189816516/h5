package autoSendmsg

import (
	"github.com/astaxie/beego/logs"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"time"
)

func TaskRun() {
	go func() {
		//自动群发任务
		time.Sleep(2 * time.Second)
		var ticker = time.NewTicker(2 * time.Second)
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

	// 每天00:05重置监控号
	go func() {

		for {

			now := time.Now()

			// 今天 00:05
			next := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				0,
				5,
				0,
				0,
				now.Location(),
			)

			// 如果已经过了00:05
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			// 等待到执行时间
			time.Sleep(next.Sub(now))

			func() {

				defer func() {
					if err := recover(); err != nil {
						logs.Info("panic:", err)
					}
				}()

				// 重置监控号
				resetMonitoringAccount()

			}()
		}

	}()
}

// 每天00点05分重置监控号
func resetMonitoringAccount() {
	for _, ws := range cache.MonitoringAccount {
		cache.SetMonitoringAccount(ws)
	}
}

// 自动群发任务
func doAutoSendMsg(t time.Time, lock string) {
	defer DelStaticData(lock)
	SetStaticData(lock, lock)

	lkey := cache.LenAutoSendMsgTaskInfo()
	if lkey == 0 { //没有任务
		return
	}

	num := 200
	for i := 0; i < num; i++ {
		taskInfo := cache.GetAutoSendMsgTaskInfo()
		if taskInfo.Account == "" {
			return
		}
		go wxComm.AutoSendMsg(taskInfo.Account, taskInfo.SessionId, taskInfo.Node)
	}
}
