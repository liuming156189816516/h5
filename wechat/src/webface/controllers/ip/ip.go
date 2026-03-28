package ip

import (
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"path"
	"webface/controllers"
	"webface/models/ip"
	info "webface/webstru"
)

type IpController struct {
	controllers.AdminController
}

//ip分组-列表
func (this *IpController) GetGroupList() {
	req := &info.GetIpGroupListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetIpGroupListRsp{}
	erro := member.GetIpGroupList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ip分组-操作
func (this *IpController) DoIpGroup() {
	req := &info.DoIpGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoIpGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ip-列表
func (this *IpController) GetIpList() {
	req := &info.GetIpListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetIpListRsp{}
	erro := member.GetIpList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//检查文件
func (this *IpController) CheckFile() {
	f, h, err := this.GetFile("file") //获取上传的文件
	req := &info.CheckFileReq{}
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

	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.CheckFileRsp{}
	//string 来源于数据包
	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		this.JsonResult(goError.NewGoError(400, "数据错误"), nil)
		return
	}
	erro := member.CheckFile(string(fileContent), req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ip-批量入库
func (this *IpController) AddIp() {
	req := &info.AddIpReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.AddIpRsp{}
	erro := member.AddIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-设置到期时间
func (this *IpController) DoExpireTime() {
	req := &info.DoExpireTimeReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoExpireTime(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-分配规则
func (this *IpController) DoAllotNum() {
	req := &info.DoAllotNumReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoAllotNum(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-移动分组
func (this *IpController) DoMoveIpGroup() {
	req := &info.DoMoveIpGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoMoveIpGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-网络检测
func (this *IpController) DoCheckStatus() {
	req := &info.DoCheckStatusReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoCheckStatus(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-ip启动分配
func (this *IpController) DoStartDistribution() {
	req := &info.DoStartDistributionReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoStartDistribution(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-ip禁用分配
func (this *IpController) DoDisableAllocation() {
	req := &info.DoDisableAllocationReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoDisableAllocation(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-批量删除
func (this *IpController) DoBatchDel() {
	req := &info.DoBatchDelReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoBatchDel(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-修改国家
func (this *IpController) DoUpCountry() {
	req := &info.DoUpCountryReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoUpCountry(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量操作-导出
func (this *IpController) DoOutPutIp() {
	req := &info.DoOutPutIpReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.DoOutPutIpRsp{}
	erro := member.DoOutPutIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//获取ipv4分配
func (this *IpController) GetIpV4Allot() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetIpV4AllotRsp{}
	erro := member.GetIpV4Allot(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//获取ipv6分配
func (this *IpController) GetIpV6Allot() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetIpV6AllotRsp{}
	erro := member.GetIpV6Allot(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//获取动态ip分配
func (this *IpController) GetIpDynamicAllot() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetIpDynamicAllotRsp{}
	erro := member.GetIpDynamicAllot(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//获取国家
func (this *IpController) GetCountryList() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetCountryListRsp{}
	erro := member.GetCountryList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//编辑备注
func (this *IpController) DoIpRemark() {
	req := &info.DoIpRemarkReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.DoIpRemarkRsp{}
	erro := member.DoIpRemark(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ip-动态
func (this *IpController) GetDynamicIp() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetDynamicIpRsp{}
	erro := member.GetDynamicIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ip-静态
func (this *IpController) GetStaticIp() {
	req := &info.GetStaticIpReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetStaticIpRsp{}
	erro := member.GetStaticIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//ip-校正工具
func (this *IpController) DoResetIp() {
	req := &info.NullReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoResetIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//获取分配的ip
func (this *IpController) GetUseList() {
	req := &info.GetUseListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &ip.IpServer{
		Sess: this.Sess,
	}
	rsp := &info.GetUseListRsp{}
	erro := member.GetUseList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
