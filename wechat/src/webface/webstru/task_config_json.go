package info

import "selfComm/wxComm/cache"

type DoTaskConfigInfoReq struct {
	DataPackId   string           `json:"data_pack_id"`  //数据包id
	MaterialList []cache.Material `json:"material_list"` //素材列表
}

type GetTaskConfigInfoRsp struct {
	DataPackId   string           `json:"data_pack_id"`   //数据包id
	DataPackName string           `json:"data_pack_name"` //数据包名称
	MaterialList []cache.Material `json:"material_list"`  //素材列表
}
