package taskConfig

import (
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"webface/controllers"
	taskConfig "webface/models/taskConfig"
	info "webface/webstru"
)

type TaskConfigController struct {
	controllers.AdminController
}

// 任务配置-获取
func (this *TaskConfigController) GetTaskConfigInfo() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &taskConfig.TaskConfigServer{
		Sess: this.Sess,
	}
	rsp := &info.GetTaskConfigInfoRsp{}
	erro := member.GetTaskConfigInfo(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 任务配置-修改
func (this *TaskConfigController) DoTaskConfigInfo() {
	req := &info.DoTaskConfigInfoReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &taskConfig.TaskConfigServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoTaskConfigInfo(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
