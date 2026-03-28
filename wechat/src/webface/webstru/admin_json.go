package info

type NullReq struct {
}

type NullRsp struct {
}

type GetAdminUserListReq struct {
	Account string `json:"account"`
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
}

type GetAdminUserListRsp struct {
	Total int64                   `json:"total"`
	List  []*GetAdminUserListInfo `json:"list"`
}
type GetAdminUserListInfo struct {
	Uid         string `json:"uid"`
	Account     string `json:"account"`      //账号
	AccountType int64  `json:"account_type"` //账号类型 1-管理员 2-主管 3-席位
	RoleId      int64  `json:"role_id"`      //角色id
	RoleName    string `json:"role_name"`    //角色名称
	Status      int64  `json:"status"`       //状态 1-启用 2-禁用
	InviteCode  string `json:"invite_code"`  //邀请码
	Itime       int64  `json:"itime"`        //创建时间
}

type DoAdminUserReq struct {
	Ptype   int64    `json:"ptype"`   // 1-新增 2-编辑 3-删除
	Uid     string   `json:"uid"`     //uid
	Account string   `json:"account"` //账号
	Pwd     string   `json:"pwd"`     //密码
	PwdStr  string   `json:"pwd_str"` //明文密码
	TwoPwd  string   `json:"two_pwd"` //二级密码
	Status  int64    `json:"status"`  //状态 1-启用 2-禁用
	RoleId  int64    `json:"role_id"` //角色id
	DelId   []string `json:"del_id"`  //删除id
}

// 登录
type AdminLoginReq struct {
	Account     string `json:"account"`
	Pwd         string `json:"pwd"`
	AccountType int64  `json:"account_type"` //2-主管 3-用户
}

type AdminLoginRsp struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"user_info"`
}
type UserInfo struct {
	Uid         string `json:"uid"`
	Account     string `json:"account"`      //账号
	AccountType int64  `json:"account_type"` //账号类型 1-管理员 2-主管 3-用户
}

type MemuInfo struct {
	MenuId    int64       `json:"menu_id"`
	Children  []*MemuInfo `json:"children"`
	Icon      string      `json:"icon"`
	Api       string      `json:"api"`
	Title     Mate        `json:"meta"`
	ClassName string      `json:"class_name"`
	Pid       int64       `json:"pid"`
	Status    int64       `json:"status"`
	Sort      int64       `json:"sort"`
	Type      int64       `json:"type"`
	Url       string      `json:"path"`
}

type Mate struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

// 获取code
type MenuReq struct {
	Pid int64 `json:"pid"`
}

type MenuRsp struct {
	Memu []*MemuInfo `json:"memu"`
}

type MemuSort []*MemuInfo

func (a MemuSort) Len() int           { return len(a) }
func (a MemuSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MemuSort) Less(i, j int) bool { return a[i].Sort < a[j].Sort }

type DoMenuReq struct {
	Ptype     int64  `json:"ptype"` //0 新增 1编辑 2删除
	MenuId    int64  `json:"menu_id"`
	Status    int64  `json:"status"` //0 启动 1禁用
	Icon      string `json:"icon"`
	Title     string `json:"title"`
	Pid       int64  `json:"pid"`
	Sort      int64  `json:"sort"`
	Url       string `json:"url"`
	Api       string `json:"api"`
	ClassName string `json:"class_name"`
}

type RoleListReq struct {
	Name  string `json:"name"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
}

type RoleListRsp struct {
	Total int64           `json:"total"`
	List  []*RoleListInfo `json:"list"`
}
type RoleListInfo struct {
	RoleId int64   `json:"role_id"`
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Stime  int64   `json:"stime"`
	Menu   []int64 `json:"menu"`
}

type DoRoleReq struct {
	Type   int64   `json:"type"` //0 新增 1编辑 2删除
	RoleId int64   `json:"role_id"`
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Menu   []int64 `json:"menu"`
}

//==============================================

type GetAppUserListReq struct {
	Account    string `json:"account"`
	Page       int64  `json:"page"`
	FuserName  string `json:"fuser_name"`
	Limit      int64  `json:"limit"`
	Ip         string `json:"ip"`
	InviteCode string `json:"invite_code"`
	Ch         string `json:"ch"` //投放渠道
}

type GetAppUserListRsp struct {
	Total int64                 `json:"total"`
	List  []*GetAppUserListInfo `json:"list"`
}
type GetAppUserListInfo struct {
	Uid        string `json:"uid"`
	Account    string `json:"account"`     //账号
	InviteCode string `json:"invite_code"` //邀请码
	FuserName  string `json:"fuser_name"`  //上级用户
	TuserName  string `json:"tuser_name"`  //顶级用户
	Level      int64  `json:"level"`       //等级
	Balance    int64  `json:"balance"`     //余额
	Itime      int64  `json:"itime"`       //注册时间
	Ip         string `json:"ip"`          //注册ip
	Status     int64  `json:"status"`      //状态 1-正常 2-禁用
	Ch         string `json:"ch"`          //投放渠道
}

type GetUserStatisReq struct {
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
	Account string `json:"account"`
}

type GetUserStatisRsp struct {
	Total int64                    `json:"total"`
	List  []*GetUserStatisListInfo `json:"list"`
}
type GetUserStatisListInfo struct {
	Uid        string `json:"uid"`
	Account    string `json:"account"`     //用户名
	InviteCode string `json:"invite_code"` //邀请码
	OnlineNum  int64  `json:"online_num"`  //在线账号
	DataNum    int64  `json:"data_num"`    //剩余数据
	Itime      int64  `json:"itime"`       //创建时间
}

type BlacklistReq struct {
	Ptype int64    `json:"ptype"` // 1-启用 2-禁用
	Ids   []string `json:"ids"`
}
