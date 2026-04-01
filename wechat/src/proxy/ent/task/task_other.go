package task

import (
	"comm/comm"
	"comm/event"
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"natsRpc"
	"proxy/ent/httpQury"
	"proxy/ent/info"
	"selfComm/wxComm"
	"selfComm/wxComm/cache"
	"serApi"
	"strings"
	"time"
	"utils"
)

// 登录事件
func TaskUserLoginEventHandler(msg *natsRpc.NatsMsg) int32 {
	req := &event.TaskUserLoginEventReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	para := []*info.AccountLoginParam{}
	for _, accountStr := range req.Account {
		accInfo := cache.GetAccountInfo(accountStr)
		token := accInfo.Token
		proxyIp := cache.GetProxyIp(accountStr)
		proxy := info.AccountAddParamSocks5{
			Type: proxyIp.Type,
			Host: proxyIp.Host,
			Port: proxyIp.Port,
			User: proxyIp.User,
			Pwd:  proxyIp.Pwd,
		}
		business := true
		if accInfo.AccountType == 1 {
			business = false
		}
		param := info.AccountLoginParam{
			Account:               accInfo.Account,
			Token:                 token,
			Business:              business,
			Proxy:                 proxy,
			DisableDecryptMessage: int64(3),
			DisableNotifyReceipt:  int64(2),
			Callback:              serApi.ServerMessage,
		}
		para = append(para, &param)
	}
	ret := &info.ResponseResult{}
	if len(req.Account) == 1 {
		num := cache.GetReplacedAccount(req.Account[0])
		if num > 0 {
			ret.Code = -1200
			ret.Message = "账号抢登"
		} else {
			ret = httpQury.AccountLogin(para, "login")
		}
	} else {
		ret = httpQury.AccountLogin(para, "login")
	}
	ret.Req = req
	ret.TaskId = req.TaskId

	msg.ResponeSucc(ret)

	err = natsRpc.NrpcCallNoReply(comm.ServerData, -1, event.TaskTypeUserLoginEventBack, ret)
	if err != nil {
		logs.Error("请求%s:%d->%s 失败", comm.ServerData, -1, event.TaskTypeUserLoginEventBack)
	}

	return natsRpc.ESMR_SUCCEED
}

// 下线事件
func TaskTypeUserLogoutEventHandler(msg *natsRpc.NatsMsg) int32 {
	req := &event.TaskUserLogoutEventReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	ret := &info.ResponseResult{}
	ret = httpQury.AccountLogout(req.Account, "logout")
	ret.Req = req
	ret.TaskId = req.TaskId
	msg.ResponeSucc(ret)
	err = natsRpc.NrpcCallNoReply(comm.ServerData, -1, event.TaskTypeUserLogoutEventBack, ret)
	if err != nil {
		logs.Error("请求%s:%d->%s 失败", comm.ServerData, -1, event.TaskTypeUserLogoutEventBack)
	}
	return natsRpc.ESMR_SUCCEED
}

// 移除事件
func TaskTypeUserRemveEventHandler(msg *natsRpc.NatsMsg) int32 {
	req := &event.TaskUserRemoveEventReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	ret := &info.ResponseResult{}
	ret = httpQury.AccountRemove(req.Account, "logout")
	ret.Req = req
	ret.TaskId = req.TaskId
	msg.ResponeSucc(ret)
	err = natsRpc.NrpcCallNoReply(comm.ServerData, -1, event.TaskTypeUserRemoveEventBack, ret)
	if err != nil {
		logs.Error("请求%s:%d->%s 失败", comm.ServerData, -1, event.TaskTypeUserRemoveEventBack)
	}
	return natsRpc.ESMR_SUCCEED
}

// 设置商店名称
func TaskTypeShopNameEventHandler(msg *natsRpc.NatsMsg) int32 {
	req := &event.TaskTypeShopNameEventReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	ret := &info.ResponseResult{}
	if req.IsCheck {
		para1 := &info.AccountCheckParam{}
		para1.Account = req.Account
		ret1 := httpQury.AccountCheck(para1)
		toString, _ := jsoniter.MarshalToString(ret1.Data)
		accountCheckData := info.AccountCheckData{}
		jsoniter.UnmarshalFromString(toString, &accountCheckData)
		if ret1.Code == 0 && accountCheckData.ShopNameExisted == false {
			para := &info.ShopNameParam{}
			para.Account = req.Account
			para.Name = req.NickName
			ret = httpQury.ShopName(para)
		}
	} else {
		para := &info.ShopNameParam{}
		para.Account = req.Account
		para.Name = req.NickName
		ret = httpQury.ShopName(para)
	}
	ret.Req = req
	ret.TaskId = req.TaskId
	msg.ResponeSucc(ret)
	err = natsRpc.NrpcCallNoReply(comm.ServerData, -1, event.TaskTypeShopNameEventBack, ret)
	if err != nil {
		logs.Error("请求%s:%d->%s 失败", comm.ServerData, -1, event.TaskTypeShopNameEventBack)
	}
	return natsRpc.ESMR_SUCCEED
}

// 发送消息
func TaskTypeSendMsgEventHandler(msg *natsRpc.NatsMsg) int32 {
	req := &event.TaskTypeSendMsgEventReq{}
	err := json.Unmarshal(msg.MsgData, &req)
	if err != nil {
		logs.Error("解析错误", err)
		msg.Response(-1000, "参数错误")
		return natsRpc.ESMR_SUCCEED
	}
	ret := &info.ResponseResult{}
	if req.IsUp == true {
		//获取数据
		target := ""
		lid := ""
		targetStr := cache.SpopDataPackList(req.DataPackId)
		split := strings.Split(targetStr, "-")
		if len(split) > 0 {
			target = split[0]
		}
		if len(split) > 1 {
			lid = split[1]
		}
		if target == "" || lid == "" {
			ret.Code = -2
			ret.Message = "粉丝数据不足"
		} else {
			//添加群发任务明细
			req.Target = target
			req.Lid = lid
			//获取素材
			mList := []cache.Material{}
			jsoniter.UnmarshalFromString(req.MaterialListStr, &mList)
			material := mList[0]

			req.Type = material.Type
			req.Content = material.Content
			req.Remark = material.Remark

			tmp := cache.SendMsgRecord{}
			tmp.RecordId = req.Account + "_" + lid
			tmp.Account = req.Account
			tmp.Target = target
			tmp.Lid = lid
			tmp.Content = req.Content
			tmp.SendTime = time.Now().Unix()
			tmp.DataPackId = req.DataPackId
			tmp.MaterialList = mList
			tmp.MaterialType = material.Type
			cache.SetSendMsgRecord(&tmp)
			cache.SetSendMsgPhoneLid(tmp.Account, tmp.Target, tmp.Lid)
		}
	}
	if ret.Code == 0 {
		if req.Type == 5 {
			para := &info.CallPhoneParam{}
			para.Account = req.Account
			para.Target = req.Target
			ret = httpQury.CallPhone(para)
		} else {

			para := &info.MessageSendParam{}
			para.Account = req.Account
			para.Target = req.Lid + ".1:0@lid"

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
				advertise := info.Advertise{}
				jsoniter.UnmarshalFromString(req.Content, &advertise)
				//视频链接
				para.Type = "adv"
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
	}
	ret.Req = req
	ret.TaskId = req.TaskId
	msg.ResponeSucc(ret)

	err = natsRpc.NrpcCallNoReply(comm.ServerData, -1, event.TaskTypeSendMsgEventBack, ret)
	if err != nil {
		logs.Error("请求%s:%d->%s 失败", comm.ServerData, -1, event.TaskTypeSendMsgEventBack)
	}
	return natsRpc.ESMR_SUCCEED
}
