package api

import (
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
	"proxy/ent/httpQury"
	"proxy/ent/info"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"serApi"
	"serApi/dllApi"
	"strconv"
	"strings"
	"utils"
)

// 挂断通话
func CallHangup(req *dllApi.CallHangupReq, rsp *dllApi.CallHangupRsp) error {
	para := &info.CallHangupParam{
		Account: req.Account,
		Target:  req.Target,
		Id:      req.Id,
	}
	ret := httpQury.CallHangup(para)
	rsp.Code = ret.Code
	rsp.Msg = ret.Message
	return nil
}

// 自主进群
func GroupJoin(req *dllApi.GroupJoinReq, rsp *dllApi.GroupJoinRsp) error {
	param := &info.GroupJoinParam{
		Account: req.Account,
		Code:    req.Code,
	}
	ret := httpQury.GroupJoin(param)
	rsp.Code = ret.Code
	rsp.Qid = utils.GetString(ret.Data)
	rsp.Message = ret.Message
	return nil
}

// 获取群的详细信息
func GroupInfo(req *dllApi.GroupInfoReq, rsp *dllApi.GroupInfoRsp) error {
	param := &info.GroupInfoParam{
		Account: req.Account,
		Group:   req.Qid,
	}
	ret := httpQury.GroupInfo(param)
	if ret.Code == 0 {
		groupData := info.GroupInfoData{}
		members := []string{}
		lids := []string{}
		toString, _ := jsoniter.MarshalToString(ret.Data)
		jsoniter.UnmarshalFromString(toString, &groupData)
		rsp.Creator = groupData.Creator
		for _, member := range groupData.Members {
			if member != "" {
				memberSplit := strings.Split(member, "|")
				for _, s := range memberSplit {
					if strings.Contains(s, "@s.whatsapp.net") {
						members = append(members, s)
					}
					if strings.Contains(s, "@lid") {
						lids = append(lids, s)
					}
				}
			}
		}
		rsp.Members = members
		rsp.Lids = lids
		rsp.Name = groupData.Name
		rsp.Id = groupData.Id
		rsp.Announcement = groupData.Announcement
	}
	return nil
}

// 查询注册
func PhoneQuery(req *dllApi.PhoneQueryReq, rsp *dllApi.PhoneQueryRsp) error {
	param := &info.PhoneQueryParam{
		Account: req.Account,
		Numbers: req.Numbers,
	}
	ret := httpQury.PhoneQuery(param)
	rsp.Code = ret.Code
	if ret.Code == 0 {
		nList := []info.PhoneQueryData{}
		toString, _ := jsoniter.MarshalToString(ret.Data)
		jsoniter.UnmarshalFromString(toString, &nList)
		upMap := make(map[string]string)
		upList := []string{}
		for _, nData := range nList {
			if nData.Exist {
				nData.Jid = strings.ReplaceAll(nData.Jid, "@s.whatsapp.net", "")
				upMap[nData.Jid] = nData.Number
				upList = append(upList, nData.Jid)
			}
		}
		rsp.UpList = upList
		rsp.QueryMap = upMap
	}
	return nil
}

// 消息发送
func MessageSend(req *dllApi.MessageSendReq, rsp *dllApi.MessageSendRsp) error {
	ret := &info.ResponseResult{}
	if req.Type == 5 {
		para := &info.CallPhoneParam{}
		para.Account = req.Account
		para.Target = req.Target
		ret = httpQury.CallPhone(para)
	} else {
		para := &info.MessageSendParam{}
		para.Account = req.Account
		para.Target = req.Target
		if req.Type == 1 {
			str2 := utils.RandStr2(8)
			textContent := str2 + req.Content
			para.Type = "text"
			para.Content = base64.StdEncoding.EncodeToString([]byte(textContent))
		}
		if req.Type == 2 {
			para.Type = "image"
			material := cache.Material{}
			jsoniter.UnmarshalFromString(req.Content, &material)
			image, w, h := wxComm.GetBase64ByUrl(req.Content)
			imageSmall := wxComm.GetBase64ByUrlSmall(req.Content, utils.StrToInt(w), utils.StrToInt(h))
			map1 := make(map[string]interface{})
			map1["thumb"] = imageSmall
			map1["image"] = image
			map1["width"] = utils.StrToInt64(w)
			map1["height"] = utils.StrToInt64(h)
			map1["text"] = base64.StdEncoding.EncodeToString([]byte(utils.RandStr2(8) + material.Remark))
			para.Content = map1
		}
		if req.Type == 7 {
			//视频链接
			para.Type = "adv"
			advertise := wxComm.Advertise{}
			jsoniter.UnmarshalFromString(req.Content, &advertise)
			image, _, _ := wxComm.GetBase64ByUrl(advertise.Img)
			map1 := make(map[string]interface{})
			map1["thumb"] = image
			map1["title"] = base64.StdEncoding.EncodeToString([]byte(advertise.Title))
			str2 := utils.RandStr2(8)
			textContent := str2 + advertise.Remark
			map1["body"] = base64.StdEncoding.EncodeToString([]byte(textContent))
			map1["text"] = base64.StdEncoding.EncodeToString([]byte(advertise.Remark))
			map1["url"] = base64.StdEncoding.EncodeToString([]byte(advertise.Url))
			para.Content = map1
		}
		ret = httpQury.MessageSend(para)
	}
	rsp.Code = ret.Code
	rsp.Msg = ret.Message
	utils.GetInt64(ret.Data)
	rsp.DataId = strconv.FormatFloat(utils.GetFloat64(ret.Data), 'f', 0, 64)
	return nil
}

// 创建关联的验证码
func VfcodeCreate(req *dllApi.VfcodeCreateReq, rsp *dllApi.VfcodeCreateRsp) error {
	/*proxy := info.AccountAddParamSocks5{
		Type: req.Proxy.Type,
		Host: req.Proxy.Host,
		Port: req.Proxy.Port,
		Pwd:  req.Proxy.Pwd,
		User: req.Proxy.User,
	}*/
	para := &info.VfcodeCreateParam{
		Account: req.Id,
		Code:    req.Code,
		//Proxy:    proxy,
		Business: false,
		//Platform: "android",
		Phone:    strings.ReplaceAll(req.Id, "_", ""),
		Callback: serApi.ServerMessage,
	}
	if req.AccountType == 2 {
		para.Business = true
	}
	ret := httpQury.VfcodeCreate(para)
	if ret.Code == 0 {
		rsp.QrCode = utils.GetString(ret.Data)
	}
	return nil
}

// 检测关联的验证码
func VfcodeCheck(req *dllApi.VfcodeCheckReq, rsp *dllApi.VfcodeCheckRsp) error {
	para := &info.VfcodeCheckParam{
		Account: req.Id,
	}
	ret := httpQury.VfcodeCheck(para)
	rsp.Data = ret.Data
	rsp.Code = ret.Code
	return nil
}

// 检测账号
func AccountCheck(req *dllApi.AccountCheckReq, rsp *dllApi.AccountCheckRsp) error {
	para := &info.AccountCheckParam{
		Account: req.Account,
	}
	ret := httpQury.AccountCheck(para)
	if ret.Code == 0 {
		accountCheckData := &info.AccountCheckData{}
		gStr, _ := jsoniter.MarshalToString(ret.Data)
		jsoniter.UnmarshalFromString(gStr, &accountCheckData)
		rsp.ErrMsg = accountCheckData.MessageError
	}
	return nil
}

// 消息发送
func MessageSendAsyn(req *dllApi.MessageSendAsynReq, rsp *dllApi.MessageSendAsynRsp) error {
	if req.Type == 5 {
		para := &info.CallPhoneParam{}
		para.Account = req.Account
		para.Target = req.Target
		httpQury.CallPhone(para)
	} else {
		para := &info.MessageSendParam{}
		para.Account = req.Account
		para.Target = req.Target
		if req.Type == 1 {
			str2 := utils.RandStr2(8)
			textContent := str2 + req.Content
			para.Type = "text"
			para.Content = base64.StdEncoding.EncodeToString([]byte(textContent))
		}
		if req.Type == 2 {
			para.Type = "image"
			material := cache.Material{}
			jsoniter.UnmarshalFromString(req.Content, &material)
			image, w, h := wxComm.GetBase64ByUrl(req.Content)
			imageSmall := wxComm.GetBase64ByUrlSmall(req.Content, utils.StrToInt(w), utils.StrToInt(h))
			map1 := make(map[string]interface{})
			map1["thumb"] = imageSmall
			map1["image"] = image
			map1["width"] = utils.StrToInt64(w)
			map1["height"] = utils.StrToInt64(h)
			map1["text"] = base64.StdEncoding.EncodeToString([]byte(utils.RandStr2(8) + material.Remark))
			para.Content = map1
		}
		if req.Type == 7 {
			//视频链接
			para.Type = "adv"
			advertise := wxComm.Advertise{}
			jsoniter.UnmarshalFromString(req.Content, &advertise)
			image, _, _ := wxComm.GetBase64ByUrl(advertise.Img)
			map1 := make(map[string]interface{})
			map1["thumb"] = image
			map1["title"] = base64.StdEncoding.EncodeToString([]byte(advertise.Title))
			str2 := utils.RandStr2(8)
			textContent := str2 + advertise.Remark
			map1["body"] = base64.StdEncoding.EncodeToString([]byte(textContent))
			map1["text"] = base64.StdEncoding.EncodeToString([]byte(advertise.Remark))
			map1["url"] = base64.StdEncoding.EncodeToString([]byte(advertise.Url))
			para.Content = map1
		}
		httpQury.MessageSendAsyn(para)
	}
	return nil
}

// 创建二维码
func QrcodeCreate(req *dllApi.QrcodeCreateReq, rsp *dllApi.QrcodeCreateRsp) error {
	proxy := info.AccountAddParamSocks5{
		Type: req.Proxy.Type,
		Host: req.Proxy.Host,
		Port: req.Proxy.Port,
		Pwd:  req.Proxy.Pwd,
		User: req.Proxy.User,
	}
	para := &info.QrcodeCreateParam{
		Account:  req.Id,
		Proxy:    proxy,
		Business: false,
		Platform: "pc",
		Callback: serApi.ServerMessage,
	}
	if req.AccountType == 2 {
		para.Business = true
	}
	ret := httpQury.QrcodeCreate(para)
	rsp.Data = ret.Data
	rsp.Code = ret.Code
	return nil
}

// 检测二维码
func QrcodeCheck(req *dllApi.QrcodeCheckReq, rsp *dllApi.QrcodeCheckRsp) error {
	para := &info.QrcodeCheckParam{
		Account: req.Id,
	}
	ret := httpQury.QrcodeCheck(para)
	rsp.Data = ret.Data
	rsp.Code = ret.Code
	return nil
}
