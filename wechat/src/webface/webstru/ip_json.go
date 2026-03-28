package info

type GetIpGroupListReq struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type GetIpGroupListRsp struct {
	Total int64                 `json:"total"`
	List  []*GetIpGroupListInfo `json:"list"`
}
type GetIpGroupListInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`  //名称
	Count int64  `json:"count"` //数量
}

type DoIpGroupReq struct {
	Ptype int64    `json:"ptype"`  // 1-新增 2-编辑 3-删除
	DelId []string `json:"del_id"` //删除id
	Id    string   `json:"id"`     //编辑id
	Name  string   `json:"name"`   //名称
}

type GetIpListReq struct {
	ProxyIp       string `json:"proxy_ip"`       //代理ip
	StartTime     int64  `json:"start_time"`     //开始时间
	EndTime       int64  `json:"end_time"`       //结束时间
	Status        int64  `json:"status"`         //网络状态 1-正常 2-异常 3-检测中 4-未检测 5-冻结
	Sort          string `json:"sort"`           //排序  正序 user_num 倒序 -user_num
	IpCategory    int64  `json:"ip_category"`    //ip类别 1-静态ip 2-动态ip
	ExpireStatus  int64  `json:"expire_status"`  //到期状态 1-正常 2-到期 3-即将到期
	DisableStatus int64  `json:"disable_status"` //禁用状态 1-启用 2-禁用
	GroupId       string `json:"group_id"`       //分组id
	Page          int64  `json:"page"`
	Limit         int64  `json:"limit"`
}

type GetIpListRsp struct {
	Total int64            `json:"total"`
	List  []*GetIpListInfo `json:"list"`
}
type GetIpListInfo struct {
	Id            string `json:"id"`
	ProxyIp       string `json:"proxy_ip"`       //代理ip
	ProxyUser     string `json:"proxy_user"`     //用户名
	Status        int64  `json:"status"`         //网络状态 1-正常 2-异常 3-检测中 4-未检测 5-冻结
	AllotNum      int64  `json:"allot_num"`      //分配总数
	UserNum       int64  `json:"user_num"`       //已分配
	IpType        int64  `json:"ip_type"`        //ip类型 1-ipV4 2-ipv6
	IpCategory    int64  `json:"ip_category"`    //ip类别 1-静态ip 2-动态ip
	ExpireStatus  int64  `json:"expire_status"`  //到期状态 1-正常 2-到期 3-即将到期
	Country       string `json:"country"`        //国家
	DisableStatus int64  `json:"disable_status"` //禁用状态 1-启用 2-禁用
	ExpireTime    int64  `json:"expire_time"`    //到期时间
	Reason        string `json:"reason"`         //原因
	Remark        string `json:"remark"`         //备注
	Creator       string `json:"creator"`        //创建者
	Itime         int64  `json:"itime"`          //创建时间
	Ptime         int64  `json:"ptime"`          //更新时间
}

type CheckFileReq struct {
	Ptype int64 `form:"ptype"` //1-静态ip 2-动态ip
}

type CheckFileRsp struct {
	FailNumber    int      `json:"fail_number"`
	Url           string   `json:"url"`
	SuccessList   []string `json:"success_list"`
	SuccessNumber int      `json:"success_number"`
}

type AddIpReq struct {
	SuccessList []string `json:"success_list"`
	IpType      int64    `json:"ip_type"`     //ip类型 1-ipV4 2-ipv6
	IpCategory  int64    `json:"ip_category"` //ip类别 1-静态ip 2-动态ip
	GroupId     string   `json:"group_id"`    //分组id
	Country     string   `json:"country"`     //国家
	ExpireTime  int64    `json:"expire_time"` //到期时间
	AllotNum    int64    `json:"allot_num"`   //分配总数
}

type AddIpRsp struct {
	FailNumber    int `json:"fail_number"`
	SuccessNumber int `json:"success_number"`
}

type DoExpireTimeReq struct {
	Ids        []string `json:"ids"`         //id列表
	ExpireTime int64    `json:"expire_time"` //到期时间
}

type DoAllotNumReq struct {
	Ids      []string `json:"ids"`       //id列表
	AllotNum int64    `json:"allot_num"` //分配总数
}

type DoMoveIpGroupReq struct {
	Ids     []string `json:"ids"`      //id列表
	GroupId string   `json:"group_id"` //分组id
}

type DoCheckStatusReq struct {
	Ids []string `json:"ids"` //id列表
}

type DoDisableAllocationReq struct {
	Ids []string `json:"ids"` //id列表
}

type DoStartDistributionReq struct {
	Ids []string `json:"ids"` //id列表
}

type DoBatchDelReq struct {
	Ids []string `json:"ids"` //id列表
}

type DoUpCountryReq struct {
	Ids     []string `json:"ids"` //id列表
	Country string   `json:"country"`
}

type DoOutPutIpReq struct {
	Ids []string `json:"ids"` //id列表
}

type DoOutPutIpRsp struct {
	Url string `json:"url"` //路径
}

type GetIpV4AllotRsp struct {
	TotalCount int64 `json:"total_count"` //分配总数
	UseNum     int64 `json:"use_num"`     //已分配
	NoUserNum  int64 `json:"no_user_num"` //未分配
}

type GetIpV6AllotRsp struct {
	TotalCount int64 `json:"total_count"` //分配总数
	UseNum     int64 `json:"use_num"`     //已分配
	NoUserNum  int64 `json:"no_user_num"` //未分配
}

type GetIpDynamicAllotRsp struct {
	TotalCount int64 `json:"total_count"` //分配总数
}

type GetCountryListRsp struct {
	CountryList []string `json:"country_list"`
}

type DoIpRemarkReq struct {
	Id     string `json:"id"`     //id
	Remark string `json:"remark"` //备注
}

type DoIpRemarkRsp struct {
	Remark string `json:"remark"` //备注
}

type GetDynamicIpRsp struct {
	List []*GetDynamicIpInfo `json:"list"`
}

type GetDynamicIpInfo struct {
	IpId    string `json:"ip_id"`
	Country string `json:"country"`
}

type GetStaticIpReq struct {
	IpType int64 `json:"ip_type"` //ip类型 1-ipV4 2-ipv6
}

type GetStaticIpRsp struct {
	List []*GetStaticIpInfo `json:"list"`
}

type GetStaticIpInfo struct {
	Country string `json:"country"`
}

type GetUseListReq struct {
	ProxyIp   string `json:"proxy_ip"`   //代理ip
	StartTime int64  `json:"start_time"` //开始时间
	EndTime   int64  `json:"end_time"`   //结束时间
	Id        string `json:"id"`
	Page      int64  `json:"page"`
	Limit     int64  `json:"limit"`
}

type GetUseListRsp struct {
	Total int64             `json:"total"`
	List  []*GetUseListInfo `json:"list"`
}
type GetUseListInfo struct {
	ProxyIp string `json:"proxy_ip"` //代理ip
	Account string `json:"account"`  //WS账号
	IpTime  int64  `json:"ip_time"`  //分配时间
}
