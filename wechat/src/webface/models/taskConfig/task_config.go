package taskConfig

import (
	"comm/comm"
	"comm/goError"
	"selfComm/db/dataPack"
	"selfComm/wxComm/cache"
	info "webface/webstru"
)

// 配置
type TaskConfigServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *TaskConfigServer) getUid() string {
	return this.Sess.Uid
}

// 任务配置-获取
func (this *TaskConfigServer) GetTaskConfigInfo(req *info.GetTaskConfigInfoReq, rsp *info.GetTaskConfigInfoRsp) *goError.ErrRsp {
	config := cache.GetTaskConfig("global")
	rsp.MaterialList = config.MaterialList
	rsp.DataPackId = config.DataPackId
	pack := dataPack.GetByIdDataPack(config.DataPackId)
	rsp.DataPackName = pack.Name
	return nil
}

// 任务配置-修改
func (this *TaskConfigServer) DoTaskConfigInfo(req *info.DoTaskConfigInfoReq, rsp *info.NullRsp) *goError.ErrRsp {
	taskConfigInfo := &cache.TaskConfigInfo{}
	taskConfigInfo.DataPackId = req.DataPackId
	taskConfigInfo.MaterialList = req.MaterialList
	cache.SetTaskConfig("global", taskConfigInfo)
	return nil
}
