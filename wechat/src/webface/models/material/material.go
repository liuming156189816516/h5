package material

import (
	"comm/comm"
	"comm/goError"
	"comm/mgoDeal"
	"comm/tableName"
	"gopkg.in/mgo.v2/bson"
	"selfComm/db/admin"
	"selfComm/db/material"
	"utils"
	info "webface/webstru"
)

//素材
type MaterialServer struct {
	Sess *comm.SessInfo // 当前的用户
}

//获取uid
func (this *MaterialServer) getUid() string {
	return this.Sess.Uid
}

//素材分组-列表
func (this *MaterialServer) GetMaterialGroupList(req *info.GetMaterialGroupListReq, rsp *info.GetMaterialGroupListRsp) *goError.ErrRsp {
	uid := this.Sess.Uid
	user := admin.GetByIdAdminUser(uid)
	if user.AccountType == 3 {
		uid = user.Creator
	}
	db := comm.GetUserMgoDBName(uid)
	tb := tableName.GetTableMaterialGroupListInfo()
	where := bson.M{}
	where["type"] = req.Type
	if req.Name != "" {
		where["name"] = req.Name
	}
	rsp.List = []*info.GetMaterialGroupListInfo{}
	rsp.Total, _ = mgoDeal.QueryMongoCount(db, tb, where)
	if rsp.Total == 0 {
		return nil
	}
	start := (req.Page - 1) * req.Limit
	all, err := mgoDeal.RealQueryMongoAll(db, tb, where, "itime", start, req.Limit)
	if err != nil {
		return goError.GLOBAL_SYSTEMERROR
	}
	for _, data := range all {
		tmp := &info.GetMaterialGroupListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
			count := material.GetCountMaterial(bson.M{"group_id": tmp.Id})
			tmp.Count = count
		}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

//素材分组-操作
func (this *MaterialServer) DoMaterialGroup(req *info.DoMaterialGroupReq, rsp *info.NullRsp) *goError.ErrRsp {
	if req.Ptype == 1 {
		//新增
		if req.Name == "未分组" {
			return goError.MaterailGroupNameErr
		}
		tmp := &material.MaterialGroup{}
		tmp.Id = bson.NewObjectId()
		tmp.Name = req.Name
		tmp.Type = req.Type
		material.AddMaterialGroup(tmp)
	}
	if req.Ptype == 2 {
		//编辑
		if req.Name == "未分组" {
			return goError.MaterailGroupNameErr
		}
		where := bson.M{}
		where["_id"] = bson.ObjectIdHex(req.Id)
		update := bson.M{}
		update["name"] = req.Name
		material.UpMaterialGroup(where, update)
	}
	if req.Ptype == 3 {
		//删除
		for _, delId := range req.DelId {
			where := bson.M{}
			where["_id"] = bson.ObjectIdHex(delId)
			material.DelMaterialGroup(where)
			//删除分组下的素材
			material.DelMaterial(bson.M{"group_id": delId})
		}
	}
	return nil
}

//素材-列表
func (this *MaterialServer) GetMaterialList(req *info.GetMaterialListReq, rsp *info.GetMaterialListRsp) *goError.ErrRsp {
	uid := this.Sess.Uid
	user := admin.GetByIdAdminUser(uid)
	if user.AccountType == 3 {
		uid = user.Creator
	}
	db := comm.GetUserMgoDBName(uid)
	tb := tableName.GetTableMaterialListInfo()
	where := bson.M{}
	where["group_id"] = req.GroupId
	where["type"] = req.Type
	if req.Name != "" {
		where["name"] = req.Name
	}
	rsp.List = []*info.GetMaterialListInfo{}
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
		tmp := &info.GetMaterialListInfo{}
		if p, ok := data["_id"]; ok {
			tmp.Id = utils.GetString(p)
		}
		if p, ok := data["name"]; ok {
			tmp.Name = utils.GetString(p)
		}
		if p, ok := data["content"]; ok {
			tmp.Content = utils.GetString(p)
		}
		rsp.List = append(rsp.List, tmp)
	}
	return nil
}

//素材-操作
func (this *MaterialServer) DoMaterial(req *info.DoMaterialReq, rsp *info.NullRsp) *goError.ErrRsp {
	uid := this.Sess.Uid
	if req.Ptype == 1 {
		//新增
		tmp := &material.Material{}
		tmp.Id = bson.NewObjectId()
		tmp.Name = req.Name
		tmp.Content = req.Content
		tmp.Creator = uid
		tmp.GroupId = req.GroupId
		tmp.Type = req.Type
		material.AddMaterial(tmp)
	}
	if req.Ptype == 2 {
		//编辑
		where := bson.M{}
		where["_id"] = bson.ObjectIdHex(req.Id)
		update := bson.M{}
		update["name"] = req.Name
		update["content"] = req.Content
		material.UpMaterial(where, update)
	}
	if req.Ptype == 3 {
		//删除
		for _, delId := range req.DelId {
			where := bson.M{}
			where["_id"] = bson.ObjectIdHex(delId)
			material.DelMaterial(where)
		}
	}
	return nil
}

//素材-移动分组
func (this *MaterialServer) DoMoveGroup(req *info.DoMoveGroupReq, rsp *info.NullRsp) *goError.ErrRsp {
	//编辑
	idList := []bson.ObjectId{}
	for _, _id := range req.MaterialId {
		idList = append(idList, bson.ObjectIdHex(_id))
	}
	where := bson.M{}
	where["_id"] = bson.M{"$in": idList}
	update := bson.M{}
	update["group_id"] = req.GroupId
	material.UpMaterial(where, update)
	return nil
}

//素材-移出分组
func (this *MaterialServer) DoMoveOutGroup(req *info.DoMoveOutGroupReq, rsp *info.NullRsp) *goError.ErrRsp {
	//编辑
	idList := []bson.ObjectId{}
	for _, _id := range req.MaterialId {
		idList = append(idList, bson.ObjectIdHex(_id))
	}
	//查询未分组id
	groupInfo := material.GetOneMaterialGroup(bson.M{"name": "未分组", "type": req.Type})
	where := bson.M{}
	where["_id"] = bson.M{"$in": idList}
	update := bson.M{}
	update["group_id"] = groupInfo.Id.Hex()
	material.UpMaterial(where, update)
	return nil
}
