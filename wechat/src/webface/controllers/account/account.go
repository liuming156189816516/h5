package account

import (
	"comm/goError"
	"fmt"
	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"path/filepath"
	accountDB "selfComm/db/account"
	"selfComm/wxComm"
	"strings"
	"webface/controllers"
	"webface/models/account"
	info "webface/webstru"
)

type AccountController struct {
	controllers.AdminController
}

// 账号分组-列表
func (this *AccountController) GetAccountGroupList() {
	req := &info.GetAccountGroupListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAccountGroupListRsp{}
	erro := member.GetAccountGroupList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 账号分组-操作
func (this *AccountController) DoAccountGroup() {
	req := &info.DoAccountGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoAccountGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 账号列表
func (this *AccountController) GetAccountInfoList() {
	req := &info.GetAccountInfoListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAccountInfoListRsp{}
	erro := member.GetAccountInfoList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 移动至其他分组
func (this *AccountController) DoUpGroup() {
	req := &info.DoUpGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoUpGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量导出
func (this *AccountController) DoOutPutAccount() {
	req := &info.DoOutPutAccountReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.DoOutPutAccountRsp{}
	erro := member.DoOutPutAccount(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量删除
func (this *AccountController) DoBatchDelAccount() {
	req := &info.DoBatchDelAccountReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoBatchDelAccount(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 快速上线
func (this *AccountController) DoBatchFastLogin() {
	req := &info.DoBatchFastLoginReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoBatchFastLogin(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量下线
func (this *AccountController) DoBatchLogout() {
	req := &info.DoBatchLogoutReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoBatchLogout(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 分组排序
func (this *AccountController) SortGroup() {
	req := &info.SortGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.SortGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 批量上线
func (this *AccountController) DoBatchLogin() {
	req := &info.DoBatchLoginReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoBatchLogin(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 释放ip
func (this *AccountController) DoFreedIp() {
	req := &info.DoFreedIpReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoFreedIp(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//入库文件-列表
func (this *AccountController) GetAccountFileList() {
	req := &info.GetAccountFileListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAccountFileListRsp{}
	erro := member.GetAccountFileList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//入库文件-批量删除
func (this *AccountController) DoBathDelAccountFile() {
	req := &info.DoBathDelAccountFileReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoBathDelAccountFile(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

func (this *AccountController) CheckAccountFile() {
	f, h, err := this.GetFile("file")
	req := &info.CheckAccountFileReq{}

	if err != nil {
		this.JsonResult(goError.NewGoError(400, "未提交文件"), nil)
		return
	}
	defer f.Close()

	if err := this.ParseForm(req); err != nil {
		this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
		return
	}

	fmt.Println(jsoniter.MarshalToString(req))

	// 👉 校验后缀
	ext := strings.ToLower(filepath.Ext(h.Filename))
	if ext != ".zip" {
		this.JsonResult(goError.NewGoError(400, "只允许上传zip文件"), nil)
		return
	}

	// 👉 防止路径攻击
	filename := filepath.Base(h.Filename)
	fileId := bson.NewObjectId().Hex()
	tmpPath := beego.AppConfig.String("tmpPath") + fileId + "/"
	// 👉 保存目录（建议独立目录，避免冲突）
	saveDir := filepath.Join(tmpPath, strings.TrimSuffix(filename, ext))
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		this.JsonResult(goError.NewGoError(500, "创建目录失败"), nil)
		return
	}
	//defer os.RemoveAll(tmpPath)

	// 👉 保存 zip 文件
	savePath := filepath.Join(saveDir, filename)

	out, err := os.Create(savePath)
	if err != nil {
		this.JsonResult(goError.NewGoError(500, "文件创建失败"), nil)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, f)
	if err != nil {
		this.JsonResult(goError.NewGoError(500, "文件保存失败"), nil)
		return
	}
	rsp := &info.NullRsp{}
	req.FileId = fileId
	req.Name = h.Filename

	tmp1 := &accountDB.AccountFile{}
	tmp1.Id = bson.ObjectIdHex(req.FileId)
	tmp1.Name = req.Name
	tmp1.AccountType = req.AccountType
	tmp1.Remark = req.Remark
	tmp1.Status = 1
	accountDB.AddAccountFile(tmp1)

	go func(param *info.CheckAccountFileReq, saveDir, deleteDir string) {
		// ✅ 👉 直接调用你的公共方法
		accountJsons, err := wxComm.DoJsonUtils(saveDir)
		defer os.RemoveAll(deleteDir)
		if err != nil {
			param.Remark = "解压zip异常"
			return
		}

		if len(accountJsons) <= 0 {
			param.Remark = "zip为空"
			return
		}
		member := &account.AccountServer{
			Sess: this.Sess,
		}
		member.CheckAccountFile(accountJsons, param)
	}(req, saveDir, tmpPath)

	this.JsonResult(goError.SuccRsp, rsp)
}

//入库日志-列表
func (this *AccountController) GetAccountLogList() {
	req := &info.GetAccountLogListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.GetAccountLogListRsp{}
	erro := member.GetAccountLogList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

//批量导出
func (this *AccountController) DoOutPutAccountLog() {
	req := &info.DoOutPutAccountLogReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &account.AccountServer{
		Sess: this.Sess,
	}
	rsp := &info.DoOutPutAccountLogRsp{}
	erro := member.DoOutPutAccountLog(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
