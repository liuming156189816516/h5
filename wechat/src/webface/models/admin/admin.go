package admin

import (
	"comm/comm"
	"comm/goError"
	"comm/mgoDeal"
	"comm/redisDeal"
	"comm/redisKeys"
	"comm/tableName"
	"comm/token"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/account"
	"selfComm/db/admin"
	"selfComm/db/material"
	"sort"
	"time"
	"utils"
	info "webface/webstru"
)

type AdminServer struct {
	Sess *comm.SessInfo // 当前的用户
}

// 获取uid
func (this *AdminServer) getUid() string {
	return this.Sess.Uid
}

// 用户-列表
func (this *AdminServer) GetAdminUserList(req *info.GetAdminUserListReq, rsp *info.GetAdminUserListRsp) *goError.ErrRsp {
	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminUserInfo()
	where := bson.M{}
	where["creator"] = this.Sess.Uid
	if req.Account != "" {
		where["account"] = req.Account
	}
	rsp.List = []*info.GetAdminUserListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "-itime", start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetAdminUserListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Uid = utils.GetString(p)
		}
		if p, ok := data["account"]; ok {
			tmp.Account = utils.GetString(p)
		}
		if p, ok := data["account_type"]; ok {
			tmp.AccountType = utils.GetInt64(p)
		}
		if p, ok := data["role_id"]; ok {
			tmp.RoleId = utils.GetInt64(p)
			roleInfo := admin.GetAdimnRoleInfo(tmp.RoleId)
			tmp.RoleName = roleInfo.Name
		}
		if p, ok := data["status"]; ok {
			tmp.Status = utils.GetInt64(p)
		}
		if p, ok := data["itime"]; ok {
			tmp.Itime = utils.GetInt64(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// 用户-操作
func (this *AdminServer) DoAdminUser(req *info.DoAdminUserReq, rsp *info.NullRsp) *goError.ErrRsp {
	if req.Ptype == 1 {
		//新增
		tmp := &admin.AdminUser{}
		tmp.Id = bson.NewObjectId()
		count := admin.GetCountAdminUser(bson.M{"account": req.Account})
		if count > 0 {
			return goError.UserExitErr
		}
		tmp.Account = req.Account
		tmp.AccountType = int64(2)
		tmp.Pwd = req.Pwd
		tmp.Status = req.Status
		tmp.RoleId = req.RoleId
		tmp.Creator = this.getUid()
		tmp.PwdStr = req.PwdStr
		tmp.TwoPwd = req.TwoPwd
		admin.AddAdminUser(tmp)
		if tmp.AccountType == 2 {
			go initAccount(tmp.Id.Hex())
		}
	}
	if req.Ptype == 2 {
		//编辑
		where := bson.M{}
		where["_id"] = bson.ObjectIdHex(req.Uid)
		update := bson.M{}
		if req.Status > 0 {
			update["status"] = req.Status
		}
		if req.RoleId > 0 {
			update["role_id"] = req.RoleId
		}
		if req.PwdStr != "" {
			update["pwd_str"] = req.PwdStr
		}
		if req.TwoPwd != "" {
			update["two_pwd"] = req.TwoPwd
		}
		if req.Pwd != "" {
			update["pwd"] = req.Pwd
		}
		token.DelToken(comm.AdminSrc, req.Uid)
		admin.UpAdminUser(where, update)
	}
	return nil
}

// 初始化数据
func initAccount(uid string) {
	//初始化素材分组
	for i := 1; i < 7; i++ {
		if i == 5 || i == 6 {
			continue
		}
		tmp := &material.MaterialGroup{}
		tmp.Id = bson.NewObjectId()
		tmp.Name = "未分组"
		tmp.Type = int64(i)
		material.AddMaterialGroup(tmp)
	}
	tmp2 := &account.AccountGroup{}
	tmp2.Id = bson.NewObjectId()
	tmp2.Name = "未分组"
	tmp2.Sort = 0 - time.Now().Unix()
	account.AddAccountGroup(tmp2)

}

// 登录
func (this *AdminServer) Login(req *info.AdminLoginReq, rsp *info.AdminLoginRsp) *goError.ErrRsp {
	userInfo := &admin.AdminUser{}
	if req.Account == "admin" {
		userInfo = admin.GetOneAdminUser(bson.M{"account": req.Account})
		if userInfo.Pwd != comm.MD5ToLower(req.Pwd) {
			return goError.UserPwdErr
		}
	} else {
		if req.AccountType == 2 {
			userInfo = admin.GetOneAdminUser(bson.M{"account": req.Account, "account_type": req.AccountType})
			if userInfo.Pwd != comm.MD5ToLower(req.Pwd) {
				return goError.UserPwdErr
			}
			if userInfo.Status == 2 {
				return goError.UserDisableErr
			}
		}
	}
	extime := time.Now().Unix() + 3*24*3600
	tokenInfo := &token.TokenInfo{
		Uid:         userInfo.Id.Hex(),
		Src:         comm.AdminSrc,
		Db:          comm.GetMgoDBName(),
		Extime:      extime,
		AccountType: userInfo.AccountType,
	}
	//保存token
	token := token.SaveToken(tokenInfo)
	if token == "" {
		return goError.GLOBAL_SYSTEMERROR
	}
	rsp.Token = token
	user := info.UserInfo{
		Uid:         userInfo.Id.Hex(),
		AccountType: userInfo.AccountType,
		Account:     userInfo.Account,
	}
	rsp.UserInfo = user
	return nil
}

// 退出登录
func (this *AdminServer) LoginOut(req *info.NullReq, rsp *info.NullRsp) *goError.ErrRsp {
	token.DelToken(comm.AdminSrc, this.getUid())
	admin.UpAdminUser(bson.M{"_id": bson.ObjectIdHex(this.getUid())}, bson.M{"online": int64(1)})
	return nil
}

// ======================================
// 递归菜单项
func GetTrees(allmenu []*info.MemuInfo, pid int64) []*info.MemuInfo {
	menus := []*info.MemuInfo{}
	mapmenus := map[int64]*info.MemuInfo{}
	for _, menu := range allmenu {
		if menu.Pid == 0 {
			menu.Children = []*info.MemuInfo{}
			menus = append(menus, menu)
			mapmenus[menu.MenuId] = menu
		}
	}
	for _, menu := range allmenu {
		if menu.Pid > 0 {
			me, ok := mapmenus[menu.Pid]
			if !ok {
				continue
			}
			menu.Children = []*info.MemuInfo{}
			me.Children = append(me.Children, menu)
		}
	}
	return menus
}

// 用户菜单
func (this *AdminServer) Menu(req *info.MenuReq, rsp *info.MenuRsp) *goError.ErrRsp {
	userInfo := admin.GetOneAdminUser(bson.M{"_id": bson.ObjectIdHex(this.getUid())})
	if userInfo.Id.Hex() == "" || userInfo.Id.Hex() != this.getUid() {
		return goError.GLOBAL_INVALIDPARAM
	}
	roleInfo := admin.GetAdimnRoleInfo(userInfo.RoleId)
	if roleInfo.RoleId != userInfo.RoleId /*|| roleInfo.Itime == 0*/ {
		return goError.GLOBAL_INVALIDPARAM
	}
	rolemap := map[int64]bool{}
	for _, m := range roleInfo.Menu {
		rolemap[m] = true
	}
	allNenu := admin.GetAdimnAllMenuList()
	list := info.MemuSort{}
	for _, data := range allNenu {
		if roleInfo.RoleId != 0 { //不超级管理员
			if _, ok := rolemap[data.MenuId]; !ok {
				continue
			}
		}
		tmp := &info.MemuInfo{}
		tmp.MenuId = data.MenuId
		tmp.Type = data.Type
		tmp.Pid = data.Pid
		tmp.Sort = data.Sort
		tmp.Url = data.Url
		mate := info.Mate{
			Title: data.Title,
			Icon:  data.Icon,
		}
		tmp.Title = mate
		//tmp.Title = data.Title
		tmp.ClassName = data.ClassName
		tmp.Icon = data.Icon
		list = append(list, tmp)
	}
	sort.Sort(list) //排序
	rsp.Memu = GetTrees(list, 0)
	return nil
}

// 菜单
func (this *AdminServer) AllMenu(req *info.MenuReq, rsp *info.MenuRsp) *goError.ErrRsp {
	allNenu := admin.GetAdimnAllMenuList()
	list := info.MemuSort{}
	for _, data := range allNenu {
		tmp := &info.MemuInfo{}
		tmp.MenuId = data.MenuId
		tmp.Type = data.Type
		tmp.Pid = data.Pid
		tmp.Sort = data.Sort
		tmp.Url = data.Url
		tmp.Status = data.Status
		tmp.ClassName = data.ClassName
		mate := info.Mate{
			Title: data.Title,
			Icon:  data.Icon,
		}
		tmp.Title = mate
		tmp.Icon = data.Icon
		list = append(list, tmp)
	}
	sort.Sort(list) //排序
	rsp.Memu = GetTrees(list, 0)

	return nil
}

// 编辑菜单
func (this *AdminServer) DoMenu(req *info.DoMenuReq, rsp *info.NullRsp) *goError.ErrRsp {
	if req.Ptype == 0 { //新增
		inData := &admin.AdminMenuInfo{}
		inData.Status = req.Status
		inData.Title = req.Title
		inData.Api = req.Api
		inData.Pid = req.Pid
		inData.Sort = req.Sort
		inData.Icon = req.Icon
		inData.Url = req.Url
		err := admin.AddAdminMenu(inData)
		if err != nil {
			return goError.GLOBAL_SYSTEMERROR
		}
	}
	if req.Ptype == 1 { //编辑
		data := map[string]interface{}{}
		data["status"] = req.Status
		data["icon"] = req.Icon
		data["title"] = req.Title
		data["class_name"] = req.ClassName
		data["pid"] = req.Pid
		data["sort"] = req.Sort
		data["url"] = req.Url
		data["api"] = req.Api
		err := admin.UpdateAdminMenuInfo(req.MenuId, data)
		if err != nil {
			return goError.GLOBAL_SYSTEMERROR
		}
	}
	if req.Ptype == 2 { //删除
		admin.DelAdminMenuInfo(req.MenuId)
	}
	return nil
}

// 角色列表
func (this *AdminServer) RoleList(req *info.RoleListReq, rsp *info.RoleListRsp) *goError.ErrRsp {

	db := comm.GetMgoDBName()
	tb := tableName.GetTableAdminRoleInfo()
	where := bson.M{}
	if req.Name != "" {
		where["name"] = req.Name
	}
	rsp.List = []*info.RoleListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.QueryMongoAll(db, tb, where, "role_id", start, req.Limit)
	if err != nil {
		return nil
	}
	for _, data := range all {
		tmp := &info.RoleListInfo{}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
		}
		if p, ok := data["role_id"]; ok {
			tmp.RoleId = utils.GetInt64(p)
		}
		if p, ok := data["desc"]; ok {
			tmp.Desc = utils.GetString(p)
		}
		if p, ok := data["menu"]; ok {
			menu := utils.GetString(p)
			jsoniter.UnmarshalFromString(menu, &tmp.Menu)
		}
		if tmp.RoleId == 0 {
			menu := admin.GetAdimnAllMenuList()
			tmp.Menu = []int64{}
			for _, m := range menu {
				if m.Status == 0 {
					tmp.Menu = append(tmp.Menu, m.MenuId)
				}
			}
		}
		if p, ok := data["itime"]; ok {
			tmp.Stime = utils.GetInt64(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

// 操作角色
func (this *AdminServer) DoRole(req *info.DoRoleReq, rsp *info.NullRsp) *goError.ErrRsp {

	if req.Type == 0 { //新增
		tmp := &admin.AdminRoleInfo{}
		key := redisKeys.GetAdminIncInfo()
		role_id, _ := redisDeal.RedisDoHincrby(key, comm.RoleId, 1)
		if role_id <= 0 {
			return goError.GLOBAL_SYSTEMERROR
		}
		tmp.RoleId = role_id
		tmp.Name = req.Name
		tmp.Desc = req.Desc
		tmp.Menu = req.Menu
		err := admin.AddAdminRole(tmp)
		if err != nil {
			return nil
		}
	}
	if req.Type == 1 { //修改
		up := map[string]interface{}{}
		up["name"] = req.Name
		up["desc"] = req.Desc
		up["menu"], _ = jsoniter.MarshalToString(req.Menu)
		err := admin.UpdateAdminRoleInfo(req.RoleId, up)
		if err != nil {
			return nil
		}
	}
	if req.Type == 2 { //删除
		admin.DelAdminRoleInfo(req.RoleId)
	}
	return nil
}
