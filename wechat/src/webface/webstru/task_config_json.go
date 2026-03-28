package info

import "selfComm/wxComm/cache"

type DoTaskConfigInfoReq struct {
	Ptype         int64            `json:"ptype"`           //1-WS-拉群 2-WS-拉粉 3-WS-私发 5-WS-AI
	MarketGroupId string           `json:"market_group_id"` //营销分组id
	Link          string           `json:"link"`            //链接
	DataPackId    string           `json:"data_pack_id"`    //数据包id
	MaterialList  []cache.Material `json:"material_list"`   //素材列表
}

type GetTaskConfigInfoReq struct {
	Ptype int64 `json:"ptype"` //1-WS-拉群 2-WS-拉粉 3-WS-私发 5-WS-AI
}

type GetTaskConfigInfoRsp struct {
	MarketGroupId   string           `json:"market_group_id"`   //营销分组id
	MarketGroupName string           `json:"market_group_name"` //营销分组名称
	DataPackId      string           `json:"data_pack_id"`      //数据包id
	DataPackName    string           `json:"data_pack_name"`    //数据包名称
	Link            string           `json:"link"`              //链接
	MaterialList    []cache.Material `json:"material_list"`     //素材列表
}

type GetMarketGroupListReq struct {
	Ptype int64 `json:"ptype"` //0-默认 1-Ai配置
}

type GetMarketGroupListRsp struct {
	List []*GetMarketGroupListInfo `json:"list1"`
}
type GetMarketGroupListInfo struct {
	GroupId   string `json:"group_id"`
	Name      string `json:"name"`       //名称
	Count     int64  `json:"count"`      //数量
	OnlineNum int64  `json:"online_num"` //在线数量
}
