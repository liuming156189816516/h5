package info

type ResponseResult struct {
	Code    int64
	Success bool
	Message string
	Data    interface{}

	Account string
	Req     interface{}
	TaskId  string
}

type ResponseResultNew struct {
	Ret  int64       `json:"errno"`
	Msg  string      `json:"errmsg"`
	Data interface{} `json:"data"`
}

type AccountLoginParam struct {
	Account               string                `json:"account"`
	NickName              string                `json:"nickname"`
	Token                 interface{}           `json:"token"`
	Callback              string                `json:"callback"`
	Business              bool                  `json:"business"`
	DisableDecryptMessage int64                 `json:"disable_decrypt_message"`
	DisableNotifyReceipt  int64                 `json:"disable_notify_receipt"`
	Proxy                 AccountAddParamSocks5 `json:"proxy"`
}

type AccountAddParamSocks5 struct {
	Type string `json:"type"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type ShopNameParam struct {
	Account string `json:"account"`
	Name    string `json:"name"`
}

type AccountCheckParam struct {
	Account string `json:"account"`
}

type AccountCheckData struct {
	MessageError    string `json:"message_error"`
	ShopNameExisted bool   `json:"shop_name_existed"`
}

type MessageSendParam struct {
	Account string `json:"account"`
	Target  string `json:"target"`
	Type    string `json:"type"`
	Forward bool   `json:"forward"` //false ==代表不是转发发 true 代表是转发
	//Args    []interface{} `json:"args"`
	Content interface{} `json:"content"`
}

type CallPhoneParam struct {
	Account string `json:"account"`
	Target  string `json:"target"`
}

type CallHangupParam struct {
	Account string `json:"account"`
	Target  string `json:"target"`
	Id      string `json:"id"`
}

type GroupData struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type GroupJoinParam struct {
	Account string      `json:"account"`
	Code    interface{} `json:"code"`
}

type PhoneQueryParam struct {
	Account string   `json:"account"`
	Numbers []string `json:"phones"`
}

type PhoneQueryData struct {
	Number string `json:"phone"`
	Jid    string `json:"jid"`
	Exist  bool   `json:"exist"`
}

type UpAccount struct {
	IdPriKey      string `json:"id_pri_key"`
	IdPubKey      string `json:"id_pub_key"`
	ChatPrikey    string `json:"chat_pri_key"`
	ChatPubKey    string `json:"chat_pub_key"`
	Mcc           string `json:"mcc"`
	Mnc           string `json:"mnc"`
	Locale        string `json:"locale"`
	Language      string `json:"language"`
	Phone         string `json:"phone"`
	WaUuid        string `json:"wa_uuid"`
	FbUuid        string `json:"fb_uuid"`
	FnUuid        string `json:"fn_uuid"`
	Version       string `json:"version"`
	PhoneId       string `json:"phone_id"`
	Model         string `json:"model"`
	OsBuildNumber string `json:"os_build_number"`
	OsVersion     string `json:"os_version"`
	DeviceBoard   string `json:"device_board"`
	Manufacturer  string `json:"manufacturer"`
}

type GroupInfoParam struct {
	Account string `json:"account"`
	Group   string `json:"group"`
}

type GroupInfoData struct {
	Creator      string   `json:"creator"` //创建者
	Members      []string `json:"members"`
	Name         string   `json:"name"`         //群名称
	Id           string   `json:"id"`           //群id
	Announcement bool     `json:"announcement"` //为true代表只能由管理员发言
}

type VfcodeCreateParam struct {
	Account  string                `json:"account"`
	Platform string                `json:"platform"`
	Business bool                  `json:"business"`
	Phone    string                `json:"phone"`
	Callback string                `json:"callback"`
	Proxy    AccountAddParamSocks5 `json:"proxy"`
}

type VfcodeCheckParam struct {
	Account string `json:"account"`
}

type Advertise struct {
	Title   string `json:"title"`   //标题
	Img     string `json:"img"`     //图片
	Url     string `json:"url"`     //链接
	Remark  string `json:"remark"`  //描述
	Content string `json:"content"` //内容
	//IsShow  bool   `json:"isShow"`  //展示广告
	//IsBtn   bool   `json:"isBtn"`   //展示按钮
	UrlType int64 `json:"url_type"` //url类型 0-视频链接 1-广告链接 3-按钮链接
}

type Message struct {
	Decrypt struct {
		Root  string `json:"root"`
		Chain string `json:"chain"`
		Index int64  `json:"index"`
	} `json:"decrypt"`
	Encrypt struct {
		Root  string `json:"root"`
		Chain string `json:"chain"`
		Index int64  `json:"index"`
	} `json:"encrypt"`
	T                       int64         `json:"t"` // Timestamp
	SignedPrekeyID          int64         `json:"signed_prekey_id"`
	PrekeyID                int64         `json:"prekey_id"`
	RegistrationID          int64         `json:"registration_id"`
	Target                  string        `json:"target"`
	TheirIdentityKey        string        `json:"their_identitykey"`
	OurBaseKeyPublicKey     string        `json:"ourbasekey_publickey"`
	OurRatchetKeyPrivateKey string        `json:"ourratchetkey_privatekey"`
	OurRatchetKeyPublicKey  string        `json:"ourratchetkey_publickey"`
	Groups                  []interface{} `json:"groups"` // Could be an array of group objects
}
