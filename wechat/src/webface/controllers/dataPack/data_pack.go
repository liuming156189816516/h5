package dataPack

import (
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"path"
	"webface/controllers"
	"webface/models/dataPack"
	info "webface/webstru"
)

//数据包
type DataPackController struct {
	controllers.AdminController
}

//获取数据包列表
func (this *DataPackController) GetDataPackList() {
	req := &info.GetDataPackListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &dataPack.DataPackServer{
		Sess: this.Sess,
	}
	rsp := &info.GetDataPackListRsp{}
	erro := member.GetDataPackList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ws数据上传
func (this *DataPackController) UpLoadFile() {
	f, h, err := this.GetFile("file") //获取上传的文件
	req := &info.UpLoadFileReq{}
	if err != nil {
		this.JsonResult(goError.NewGoError(400, "未提交文件"), nil)
		return
	}
	if err := this.ParseForm(req); err != nil {
		this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
		return
	}
	defer f.Close()

	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".txt": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		this.JsonResult(goError.NewGoError(400, "后缀名不符合上传要求"), nil)
		return
	}

	member := &dataPack.DataPackServer{
		Sess: this.Sess,
	}
	rsp := &info.UpLoadFileRsp{}
	//string 来源于数据包
	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		this.JsonResult(goError.NewGoError(400, "数据错误"), nil)
		return
	}
	erro := member.UpLoadFile(string(fileContent), req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//查看上传结果
func (this *DataPackController) GetSchedule() {
	req := &info.GetScheduleReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &dataPack.DataPackServer{
		Sess: this.Sess,
	}
	rsp := &info.GetScheduleRsp{}
	erro := member.GetSchedule(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//获取剩余数量
func (this *DataPackController) GetResidueNum() {
	req := &info.GetResidueNumReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &dataPack.DataPackServer{
		Sess: this.Sess,
	}
	rsp := &info.GetResidueNumRsp{}
	erro := member.GetResidueNum(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量删除
func (this *DataPackController) BathDel() {
	req := &info.BathDelReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &dataPack.DataPackServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.BathDel(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//导出全部数据和导出剩余数据
func (this *DataPackController) DoOutPutData() {
	req := &info.DoOutPutDataReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &dataPack.DataPackServer{
		Sess: this.Sess,
	}
	rsp := &info.DoOutPutDataRsp{}
	erro := member.DoOutPutData(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
