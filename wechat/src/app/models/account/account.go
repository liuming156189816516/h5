package account

import (
	info "app/webstru"
	"comm/comm"
	"comm/goError"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/ip"
	"selfComm/db/log"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"serApi/dllApi"
	"strings"
)

// 数据包
type AccountServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *AccountServer) getUid() string {
	return this.Sess.Uid
}

// 获取验证码
func (this *AccountServer) GetQrCode(req *info.GetQrCodeReq, rsp *info.GetQrCodeRsp) *goError.ErrRsp {

	if len(req.AreaCode+req.Account) != 12 || len(req.AreaCode+req.Account) != 13 || req.AreaCode != "55" {
		return goError.GLOBAL_INVALIDPARAM
	}

	//kwai发送访问回调
	tmp := &log.FbReportLog{}
	tmp.Id = bson.NewObjectId()
	tmp.Ptype = 2
	reqStr, _ := jsoniter.MarshalToString(req)
	tmp.Data = reqStr
	log.AddFbReportLog(tmp)

	if req.PixelId == wxComm.PixId && req.ClickId != "" {
		go func() {
			wxComm.KwaiPlace(req.ClickId, "EVENT_BUTTON_CLICK")
		}()
	}
	phone := req.AreaCode + req.Account
	result := wxComm.QrcodeUtils(phone)
	if result.PairingCode != "" {
		rsp.Code = result.PairingCode
		accInfo := &cache.AccountInfo{
			Account:      phone,
			AccountType:  1,
			PlatformType: 2,
			PixelId:      req.PixelId,
			ClickId:      req.ClickId,
		}
		cache.SetAccountInfo(phone, accInfo)
	} else {
		return goError.AccountCodeLoginErr
	}
	return nil
}

// 获取ws挂机状态
func (this *AccountServer) GetAccountResult(req *info.GetAccountResultReq, rsp *info.GetAccountResultRsp) *goError.ErrRsp {
	status := cache.GetAccountStatus(req.AreaCode + req.Account)
	if status == 3 {
		rsp.Status = 1
	} else if status == 2 {
		rsp.Status = 2
	} else {
		rsp.Status = 3
	}
	return nil
}

// 获取验证码
func (this *AccountServer) GetLgQrCode(req *info.GetQrCodeReq, rsp *info.GetQrCodeRsp) *goError.ErrRsp {

	//kwai发送访问回调
	tmp := &log.FbReportLog{}
	tmp.Id = bson.NewObjectId()
	tmp.Ptype = 2
	reqStr, _ := jsoniter.MarshalToString(req)
	tmp.Data = reqStr
	log.AddFbReportLog(tmp)

	if req.PixelId == wxComm.PixId && req.ClickId != "" {
		go func() {
			wxComm.KwaiPlace(req.ClickId, "EVENT_BUTTON_CLICK")
		}()
	}

	tmpProxy := &cache.AccountSocks5Info{}
	lockIpId := ""
	//获取ip
	lockIp := ip.GetOneLockIp()
	if lockIp.ProxyIp == "" {
		return goError.IpOperationErr
	}
	split := strings.Split(lockIp.ProxyIp, ":")
	tmpProxy.User = lockIp.ProxyUser
	tmpProxy.Pwd = lockIp.ProxyPwd
	tmpProxy.Type = lockIp.ProxyType
	tmpProxy.Host = split[0]
	tmpProxy.Port = split[1]
	lockIpId = lockIp.Id.Hex()
	proxy := dllApi.AccountAddParamSocks5{}
	proxy.User = tmpProxy.User
	proxy.Pwd = tmpProxy.Pwd
	proxy.Host = tmpProxy.Host
	proxy.Port = tmpProxy.Port
	proxy.Type = tmpProxy.Type
	uuid := req.AreaCode + req.Account + "_"
	dreq := &dllApi.VfcodeCreateReq{
		Id:          uuid,
		Code:        "77777777",
		Proxy:       proxy,
		AccountType: 1,
	}
	dRsp, _ := dllApi.VfcodeCreate(dreq, -1, true, 30)
	if dRsp != nil && dRsp.QrCode != "" {
		rsp.Code = dRsp.QrCode
		taskData := info.CheckQrcodeTaskData{}
		taskData.User = tmpProxy.User
		taskData.Pwd = tmpProxy.Pwd
		taskData.Host = tmpProxy.Host
		taskData.Port = tmpProxy.Port
		taskData.Type = tmpProxy.Type
		taskData.ProxyId = lockIpId
		taskData.AccountType = 1
		taskData.AreaCode = req.AreaCode
		taskData.PixelId = req.PixelId
		taskData.ClickId = req.ClickId
		cache.SetCheckQrcodeTask(uuid, taskData)
	} else {
		return goError.AccountCodeLoginErr
	}

	return nil
}
