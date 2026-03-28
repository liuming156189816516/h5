package comm

const (
	UserSrc  = "user"
	AdminSrc = "admin"

	UserTypeAdmin = 0 //员工
	UserTypeReg   = 1 //用户
	UserTypeUnReg = 2 //游客
)

// 用户信息
type SessInfo struct {
	Uid      string     `json:"uid"`
	Tuid     string     `json:"tuid"`
	DB       string     `json:"db"`
	ExtTime  int64      `json:"ext_time"`
	Pack     int64      `json:"pack"`
	Src      string     `json:"src"`
	Info     *LoginInfo `json:"info"`
	Model    string     `json:"model"`
	Func     string     `json:"func"`
	//Fbclid   string     `json:"fbclid"`
	//Pixellid string     `json:"pixellid"`
	AccountType int64  `json:"account_type"`   //1-管理员 2-主管 3-用户
}

// 登录信息
type LoginInfo struct {
	Ip  string `json:"ip"`
	Max string `json:"max"`
}
