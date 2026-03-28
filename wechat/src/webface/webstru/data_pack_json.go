package info

type GetDataPackListReq struct {
	Name  string `json:"name"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type GetDataPackListRsp struct {
	Total int64                  `json:"total"`
	List  []*GetDataPackListInfo `json:"list"`
}
type GetDataPackListInfo struct {
	Id              string `json:"id"`
	Name            string `json:"name"`              //数据名称
	UpNum           int64  `json:"up_num"`            //上传数据
	InvalidNum      int64  `json:"invalid_num"`       //无效数据
	SourceRepeatNum int64  `json:"source_repeat_num"` //源重复数据
	RepeatNum       int64  `json:"repeat_num"`        //账号内重复
	IntoNum         int64  `json:"into_num"`          //入库数量
	ResidueNum      int64  `json:"residue_num"`       //剩余数量
	ErrNum          int64  `json:"err_num"`           //异常数量
	UpStatus        int64  `json:"up_status"`         //上传状态 1-上传中 2-上传完成
	Itime           int64  `json:"itime"`             //创建时间
}

type UpLoadFileReq struct {
	Name     string `form:"name"`
	IntoType int64  `form:"into_type"` //入库方式 1-已使用过的数据不入库 2-已使用过的数据要入库
	Ptype    int64  `form:"ptype"`     //1-新增 2-补充
	Id       string `form:"id"`        //补充id
}

type UpLoadFileRsp struct {
	Id string `json:"id"`
}

type GetScheduleReq struct {
	Id string `json:"id"`
}

type GetScheduleRsp struct {
	Fail     int64 `json:"fail"`
	Success  int64 `json:"success"`
	UpStatus int64 `json:"up_status"`
}

type GetResidueNumReq struct {
	Id    string `json:"id"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type GetResidueNumRsp struct {
	Total int64    `json:"total"`
	List  []string `json:"list"`
}

type DoOutPutDataReq struct {
	Id     string `json:"id"`      //id
	Type   int64  `json:"type"`    //1-导出全部数据 2-导出剩余数据 3-异常数据
	TwoPwd string `json:"two_pwd"` //二级密码
}

type DoOutPutDataRsp struct {
	Url string `json:"url"` //路径
}

type BathDelReq struct {
	Ids []string `json:"ids"`
}
