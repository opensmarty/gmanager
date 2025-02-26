package system

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"gmanager/module/constants"
	"gmanager/utils/base"
)

type SysMenu struct {
	// columns START
	Id       int    `json:"id" gconv:"id,omitempty"`             // 主键
	Parentid int    `json:"parentid" gconv:"parentid,omitempty"` // 父id
	Name     string `json:"name" gconv:"name,omitempty"`         // 名称/11111
	Icon     string `json:"icon" gconv:"icon,omitempty"`         // 菜单图标
	Urlkey   string `json:"urlkey" gconv:"urlkey,omitempty"`     // 菜单key
	Url      string `json:"url" gconv:"url,omitempty"`           // 链接地址
	Perms    string `json:"perms" gconv:"perms,omitempty"`       // 授权(多个用逗号分隔，如：user:list,user:create)
	Status   int    `json:"status" gconv:"status,omitempty"`     // 状态//radio/2,隐藏,1,显示
	Type     int    `json:"type" gconv:"type,omitempty"`         // 类型//select/1,目录,2,菜单,3,按钮
	Sort     int    `json:"sort" gconv:"sort,omitempty"`         // 排序
	Level    int    `json:"level" gconv:"level,omitempty"`       // 级别
	// columns END

	base.BaseModel
}

func (model SysMenu) Get() SysMenu {
	if model.Id <= 0 {
		glog.Error(model.TableName() + " get id error")
		return SysMenu{}
	}

	var resData SysMenu
	err := model.dbModel("t").Where(" id = ?", model.Id).Fields(model.columns()).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SysMenu{}
	}

	return resData
}

func (model SysMenu) GetOne(form *base.BaseForm) SysMenu {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["id"] != "" {
		where += " and id = ? "
		params = append(params, gconv.Int(form.Params["id"]))
	}

	var resData SysMenu
	err := model.dbModel("t").Where(where, params).Fields(model.columns()).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SysMenu{}
	}

	return resData
}

func (model SysMenu) ListUser(userId int, userType int) []SysMenu {
	if userType == constants.UserTypeAdmin {
		return model.List(&base.BaseForm{})
	}

	var resData []SysMenu
	err := model.dbModel("t").Fields(model.columns()).LeftJoin(
		"sys_role_menu rm", "rm.menu_id = t.id ").LeftJoin(
		"sys_user_role ur", "ur.role_id = rm.role_id ").Where(
		"ur.user_id = ? ", userId).Structs(&resData)
	if err != nil {
		glog.Error(model.TableName()+" list error", err)
		return []SysMenu{}
	}

	return resData
}

func (model SysMenu) List(form *base.BaseForm) []SysMenu {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["name"] != "" {
		where += " and name like ? "
		params = append(params, "%"+form.Params["name"]+"%")
	}
	if form.Params != nil && form.Params["level"] != "" {
		where += " and level in (?) "
		params = append(params, gstr.Split(form.Params["level"], ","))
	}
	if gstr.Trim(form.OrderBy) == "" {
		form.OrderBy = " sort,id desc"
	}

	var resData []SysMenu
	err := model.dbModel("t").Fields(
		model.columns()).Where(where, params...).OrderBy(form.OrderBy).Structs(&resData)
	if err != nil {
		glog.Error(model.TableName()+" list error", err)
		return []SysMenu{}
	}

	return resData
}

func (model SysMenu) Page(form *base.BaseForm) []SysMenu {
	if form.Page <= 0 || form.Rows <= 0 {
		glog.Error(model.TableName()+" page param error", form.Page, form.Rows)
		return []SysMenu{}
	}

	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["name"] != "" {
		where += " and name like ? "
		params = append(params, "%"+form.Params["name"]+"%")
	}

	num, err := model.dbModel("t").Where(where, params...).Count()
	form.TotalSize = num
	form.TotalPage = num / form.Rows

	// 没有数据直接返回
	if num == 0 {
		form.TotalPage = 0
		form.TotalSize = 0
		return []SysMenu{}
	}

	var resData []SysMenu
	pageNum, pageSize := (form.Page-1)*form.Rows, form.Rows
	dbModel := model.dbModel("t").Fields(model.columns() + ",su1.real_name as updateName,su2.real_name as createName")
	dbModel = dbModel.LeftJoin("sys_user su1", " t.update_id = su1.id ")
	dbModel = dbModel.LeftJoin("sys_user su2", " t.update_id = su2.id ")
	err = dbModel.Where(where, params...).Limit(pageNum, pageSize).OrderBy(form.OrderBy).Structs(&resData)
	if err != nil {
		glog.Error(model.TableName()+" page list error", err)
		return []SysMenu{}
	}

	return resData
}

func (model SysMenu) Delete() int64 {
	if model.Id <= 0 {
		glog.Error(model.TableName() + " delete id error")
		return 0
	}

	r, err := model.dbModel().Where(" id = ?", model.Id).Delete()
	if err != nil {
		glog.Error(model.TableName()+" delete error", err)
		return 0
	}

	res, err2 := r.RowsAffected()
	if err2 != nil {
		glog.Error(model.TableName()+" delete res error", err2)
		return 0
	}

	LogSave(model, DELETE)
	return res
}

func (model SysMenu) Update() int64 {
	r, err := model.dbModel().Data(model).Where(" id = ?", model.Id).Update()
	if err != nil {
		glog.Error(model.TableName()+" update error", err)
		return 0
	}

	res, err2 := r.RowsAffected()
	if err2 != nil {
		glog.Error(model.TableName()+" update res error", err2)
		return 0
	}

	LogSave(model, UPDATE)
	return res
}

func (model *SysMenu) Insert() int64 {
	r, err := model.dbModel().Data(model).Insert()
	if err != nil {
		glog.Error(model.TableName()+" insert error", err)
		return 0
	}

	res, err2 := r.RowsAffected()
	if err2 != nil {
		glog.Error(model.TableName()+" insert res error", err2)
		return 0
	} else if res > 0 {
		lastId, err2 := r.LastInsertId()
		if err2 != nil {
			glog.Error(model.TableName()+" LastInsertId res error", err2)
			return 0
		} else {
			model.Id = gconv.Int(lastId)
		}
	}

	LogSave(model, INSERT)
	return res
}

func (model SysMenu) dbModel(alias ...string) *gdb.Model {
	var tmpAlias string
	if len(alias) > 0 {
		tmpAlias = " " + alias[0]
	}
	tableModel := g.DB().Table(model.TableName() + tmpAlias).Safe()
	return tableModel
}

func (model SysMenu) PkVal() int {
	return model.Id
}

func (model SysMenu) TableName() string {
	return "sys_menu"
}

func (model SysMenu) columns() string {
	sqlColumns := "t.id,t.parentid,t.name,t.icon,t.urlkey,t.url,t.perms,t.status,t.type,t.sort,t.level,t.enable,t.update_time as updateTime,t.update_id as updateId,t.create_time as createTime,t.create_id as createId"
	return sqlColumns
}
