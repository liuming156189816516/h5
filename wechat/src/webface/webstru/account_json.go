package info

type GetAccountGroupListReq struct {
	Name  string `json:"name"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type GetAccountGroupListRsp struct {
	Total int64                      `json:"total"`
	List  []*GetAccountGroupListInfo `json:"list"`
}
type GetAccountGroupListInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`       //名称
	Count     int64  `json:"count"`      //数量
	IsDefault int64  `json:"is_default"` //是否是默认分组 0-否 1-是
}

type DoAccountGroupReq struct {
	Ptype int64    `json:"ptype"`  // 1-新增 2-编辑 3-删除
	DelId []string `json:"del_id"` //删除id
	Id    string   `json:"id"`     //编辑id
	Name  string   `json:"name"`   //名称
}

type GetAccountInfoListReq struct {
	GroupId             string `json:"group_id"`
	Account             string `json:"account"`                //账号
	Status              int64  `json:"status"`                 //账号状态 1-离线 2-在线 3-登录中 4-登录失败 5-离线中
	AccountType         int64  `json:"account_type"`           //账号类型 1-个人号 2-商业号
	PlatformType        int64  `json:"platform_type"`          //平台类型 1-云控 2-APP
	Reason              string `json:"reason"`                 //离线原因
	PixelId             string `json:"pixel_id"`               //渠道
	ItimeStartTime      int64  `json:"itime_start_time"`       //入库开始时间
	ItimeEndTime        int64  `json:"itime_end_time"`         //入库结束时间
	FirstLoginStartTime int64  `json:"first_login_start_time"` //首次登录开始时间
	FirstLoginEndTime   int64  `json:"first_login_end_time"`   //首次登录结束时间
	OfflineStartTime    int64  `json:"offline_start_time"`     //离线开始时间
	OfflineEndTime      int64  `json:"offline_end_time"`       //离线结束时间
	Page                int64  `json:"page"`
	Limit               int64  `json:"limit"`

	IsProxyUser int64 `json:"is_proxy_user"` //是否反向代理账号 0-否 1-是
}

type GetAccountInfoListRsp struct {
	Total int64                     `json:"total"`
	List  []*GetAccountInfoListInfo `json:"list"`
}
type GetAccountInfoListInfo struct {
	Id             string `json:"id"`
	Head           string `json:"head"`             //头像
	Account        string `json:"account"`          //账号
	NickName       string `json:"nick_name"`        //昵称
	Status         int64  `json:"status"`           //账号状态 1-离线 2-在线 3-登录中 4-登录失败 5-离线中
	Reason         string `json:"reason"`           //离线原因
	AccountType    int64  `json:"account_type"`     //账号类型 1-个人号 2-商业号
	OfflineTime    int64  `json:"offline_time"`     //离线时间
	FirstLoginTime int64  `json:"first_login_time"` //首次登录时间
	Remark         string `json:"remark"`           //备注
	PixelId        string `json:"pixel_id"`         //渠道
	PlatformType   int64  `json:"platform_type"`    //平台类型 1-云控 2-APP
	Itime          int64  `json:"itime"`            //入库时间
	FuserName      string `json:"fuser_name"`       //上级用户
}

type DoUpRemarkReq struct {
	Accounts []string `json:"accounts"` //账号列表
	Remark   string   `json:"remark"`   //备注
}

type DoUpGroupReq struct {
	Accounts []string `json:"accounts"` //账号列表
	GroupId  string   `json:"group_id"`
}

type DoOutPutAccountReq struct {
	Accounts []string `json:"accounts"` //账号列表
}

type DoOutPutAccountRsp struct {
	Url string `json:"url"` //路径
}

type DoBatchDelAccountReq struct {
	Accounts []string `json:"accounts"` //账号列表
}

type DoBatchLogoutReq struct {
	Accounts []string `json:"accounts"` //账号列表
}

type SortGroupReq struct {
	List []string `json:"list"`
}

type DoBatchFastLoginReq struct {
	Accounts []string `json:"accounts"` //账号列表
}

type GetAccountFileListReq struct {
	Name      string `json:"name"`       //文件名
	StartTime int64  `json:"start_time"` //开始时间
	EndTime   int64  `json:"end_time"`   //结束时间
	Page      int64  `json:"page"`
	Limit     int64  `json:"limit"`
}

type GetAccountFileListRsp struct {
	Total int64                     `json:"total"`
	List  []*GetAccountFileListInfo `json:"list"`
}

type GetAccountFileListInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`         //文件名称
	AccType     string `json:"acc_type"`     //账号类型
	AccountType int64  `json:"account_type"` //账号分类 1-个人号 2-商业号
	Status      int64  `json:"status"`       //任务状态 1-导入中 2-已完成
	SuccessNum  int64  `json:"success_num"`  //导入成功
	FailNum     int64  `json:"fail_num"`     //导入失败
	Remark      string `json:"remark"`       //备注
	Itime       int64  `json:"itime"`        //导入时间
}

type DoBathDelAccountFileReq struct {
	Ids []string `json:"ids"`
}

type CheckAccountFileReq struct {
	ImportType int64 `form:"import_type"` //导入类型 1-channel格式 2-全参格式
}

type UpAccount struct {
	Phone string `json:"phone"`
}

type CheckAccountFileRsp struct {
	Name          string   `json:"name"`           //文件名称
	FailNumber    int      `json:"fail_number"`    //失败数量
	Url           string   `json:"url"`            //错误文件地址
	SuccessList   []string `json:"success_list"`   //成功集合
	SuccessNumber int      `json:"success_number"` //成功数量
}

type AddAccountReq struct {
	SuccessList []string `json:"success_list"`
	Name        string   `json:"name"`         //文件名称
	AccountType int64    `json:"account_type"` //账号分类 1-个人号 2-商业号
	ImportType  int64    `json:"import_type"`  //导入类型 1-channel格式
	GroupId     string   `json:"group_id"`     //分组id
	Remark      string   `json:"remark"`       //备注
}

type AddAccountRsp struct {
	Id string `json:"id"` //文件id
}

type GetAccountScheduleReq struct {
	Id string `json:"id"`
}

type GetAccountScheduleRsp struct {
	Fail     int64 `json:"fail"`
	Success  int64 `json:"success"`
	UpStatus int64 `json:"up_status"`
}

type GetAccountLogListReq struct {
	FileId string `json:"file_id"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
}

type GetAccountLogListRsp struct {
	Total int64                    `json:"total"`
	List  []*GetAccountLogListInfo `json:"list"`
}
type GetAccountLogListInfo struct {
	Id      string `json:"id"`
	Account string `json:"account"` //账号
	Status  int64  `json:"status"`  //入库状态 1-失败 2-成功
	Reason  string `json:"reason"`  //原因
}

type DoOutPutAccountLogReq struct {
	Ptype int64    `json:"ptype"` //导出类型 1-全部导出 2-导出入库失败的数据 3-导出入库成功的数据
	Ids   []string `json:"ids"`   //id列表
}

type DoOutPutAccountLogRsp struct {
	Url string `json:"url"` //路径
}

type DoFreedIpReq struct {
	Accounts []string `json:"accounts"` //账号列表
}

type DoBatchLoginReq struct {
	Accounts   []string `json:"accounts"`    //账号列表
	IpCategory int64    `json:"ip_category"` //ip类别 1-静态ip 2-动态ip
	IpType     int64    `json:"ip_type"`     //ip类型 1-ipV4 2-ipv6 3-动态住宅
	IpId       string   `json:"ip_id"`       //ipId
	Country    string   `json:"country"`     //国家
}
