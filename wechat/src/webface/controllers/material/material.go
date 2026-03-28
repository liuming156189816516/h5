package material

import (
	"comm/comm"
	"comm/cos"
	"comm/goError"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/json-iterator/go"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"selfComm/wxComm"
	"strings"
	"webface/controllers"
	"webface/models/material"
	"webface/webstru"
)

type MaterialController struct {
	controllers.AdminController
}

// 素材分组-列表
func (this *MaterialController) GetMaterialGroupList() {
	req := &info.GetMaterialGroupListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &material.MaterialServer{
		Sess: this.Sess,
	}
	rsp := &info.GetMaterialGroupListRsp{}
	erro := member.GetMaterialGroupList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 素材分组-操作
func (this *MaterialController) DoMaterialGroup() {
	req := &info.DoMaterialGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &material.MaterialServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoMaterialGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 素材-列表
func (this *MaterialController) GetMaterialList() {
	req := &info.GetMaterialListReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &material.MaterialServer{
		Sess: this.Sess,
	}
	rsp := &info.GetMaterialListRsp{}
	erro := member.GetMaterialList(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 素材-操作
func (this *MaterialController) DoMaterial() {
	req := &info.DoMaterialReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &material.MaterialServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoMaterial(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 素材-移动分组
func (this *MaterialController) DoMoveGroup() {
	req := &info.DoMoveGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &material.MaterialServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoMoveGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 素材-移出分组
func (this *MaterialController) DoMoveOutGroup() {
	req := &info.DoMoveOutGroupReq{}
	if len(this.Ctx.Input.RequestBody) != 0 {
		err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, &req)
		if err != nil {
			this.JsonResult(goError.GLOBAL_INVALIDPARAM, nil)
			return
		}
	}
	member := &material.MaterialServer{
		Sess: this.Sess,
	}
	rsp := &info.NullRsp{}
	erro := member.DoMoveOutGroup(req, rsp)
	if erro != nil {
		this.JsonResult(erro, nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}

// 素材-上传
func (this *MaterialController) Upload() {
	rsp := info.ApiResponseUpload{}
	f, h, err := this.GetFile("file") //获取上传的文件
	if err != nil {
		this.JsonResult(goError.NewGoError(400, "未提交文件"), nil)
		return
	}
	defer f.Close()
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".mp4":  true,
		".mp3":  true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		this.JsonResult(goError.NewGoError(400, "后缀名不符合上传要求"), nil)
		return
	}
	filePath := ""
	fileName := ""
	fileContent, _ := ioutil.ReadAll(f)
	fileName = comm.Md5FromBin(fileContent) + ext
	tmpPath := beego.AppConfig.String("tmpPath")
	filePath = tmpPath + fileName
	err = this.SaveToFile("file", filePath)
	if err != nil {
		this.JsonResult(goError.NewGoError(400, "服务器错误"), nil)
		return
	}
	defer os.Remove(filePath)
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		//打开图片文件
		// 打开原始图片文件
		inFile, err := os.Open(filePath)
		if err != nil {
			this.JsonResult(goError.NewGoError(400, fmt.Sprintf("Error opening image file: %v", err)), nil)
			return
		}
		defer inFile.Close()
		// 解码图片
		srcImg, _, err := image.Decode(inFile)
		if err != nil {
			this.JsonResult(goError.NewGoError(400, fmt.Sprintf("Error decoding image: %v", err)), nil)
			return
		}

		// 创建一个新的图片文件
		jpgFilePath := strings.ReplaceAll(filePath, ext, ".jpg")
		fileName = strings.ReplaceAll(filePath, ext, ".jpg")

		// 创建输出文件
		outFile, err := os.Create(jpgFilePath)
		if err != nil {
			this.JsonResult(goError.NewGoError(400, fmt.Sprintf("Error creating output file: %v", err)), nil)
			return
		}
		defer outFile.Close()
		defer os.Remove(jpgFilePath)

		// 将图片编码为JPG格式
		//err = jpeg.Encode(outFile, srcImg, &jpeg.Options{Quality: 75})
		//if err != nil {
		//	this.JsonResult(goError.NewGoError(400, fmt.Sprintf("Error encoding image: %v", err)), nil)
		//	return
		//}
		wxComm.JfifEncode(outFile, srcImg, &jpeg.Options{Quality: 75})
		filePath = jpgFilePath
	}
	fileUrl := cos.UploadAwsFile(filePath, fileName)
	rsp.Url = fileUrl
	if rsp.Url == "" {
		this.JsonResult(goError.NewGoError(400, "上传异常"), nil)
		return
	}
	this.JsonResult(goError.SuccRsp, rsp)
}
