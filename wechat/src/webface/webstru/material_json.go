package info

type GetMaterialGroupListReq struct {
	Name  string `json:"name"`
	Type  int64  `json:"type"` //1-文本 2-图片 3-语音 4-视频
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type GetMaterialGroupListRsp struct {
	Total int64                       `json:"total"`
	List  []*GetMaterialGroupListInfo `json:"list"`
}
type GetMaterialGroupListInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`  //名称
	Count int64  `json:"count"` //数量
}

type DoMaterialGroupReq struct {
	Ptype int64    `json:"ptype"`  // 1-新增 2-编辑 3-删除
	DelId []string `json:"del_id"` //删除id
	Id    string   `json:"id"`     //编辑id
	Name  string   `json:"name"`   //名称
	Type  int64    `json:"type"`   //类型 1-文字 2-图片 3-语音 4-视频
}

type GetMaterialListReq struct {
	GroupId string `json:"group_id"` //分组id
	Type    int64  `json:"type"`     //类型
	Name    string `json:"name"`     //名称
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
}

type GetMaterialListRsp struct {
	Total int64                  `json:"total"`
	List  []*GetMaterialListInfo `json:"list"`
}
type GetMaterialListInfo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`    //标题
	Content string `json:"content"` //内容
}

type DoMaterialReq struct {
	Ptype   int64    `json:"ptype"`    // 1-新增 2-编辑 3-删除
	DelId   []string `json:"del_id"`   //删除id
	Id      string   `json:"id"`       //编辑id
	Name    string   `json:"name"`     //名称
	Content string   `json:"content"`  //内容
	GroupId string   `json:"group_id"` //分组id
	Type    int64    `json:"type"`     //类型
}

type DoMoveGroupReq struct {
	MaterialId []string `json:"material_id"`
	GroupId    string   `json:"group_id"`
}

type DoMoveOutGroupReq struct {
	MaterialId []string `json:"material_id"`
	Type       int64    `json:"type"`
}

type ApiResponseUpload struct {
	Url string `json:"url"`
}
