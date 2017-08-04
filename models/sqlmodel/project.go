package sqlmodel

import (
	"github.com/2liang/mcache/models/base"
	"github.com/2liang/mcache/modules/utils/setting"
	"errors"
	"strconv"
)

type ProjectData struct {
	Id 		int			`xorm:"'id' int(11)" json:"id"`
	Name	string		`xorm:"'name' varchar(250)" json:"name"`
	Desc 	string		`xorm:"'desc' varchar(250)" json:"desc"`
	CreateTime	int64		`xorm:"'create_time' int(11)" json:"create_time"`
	ModifyTime	int64		`xorm:"'modify_time' int(11)" json:"modify_time"`
}

func (ap *ProjectData) GetProjectByPid () ([]ProjectData, error) {
	db := base.DbCache.GetSlave()
	r := make([]ProjectData, 0)
	if err := db.Table("projects").Where("id = ?", ap.Id).Find(&r); err != nil {
		return nil, err
	}
	return r, nil
}

/**
 * 获取项目
 */
func (ap *ProjectData) GetProject (name string, page int, limit int) ([]ProjectData, error) {
	start := (page - 1) * limit
	db := base.DbCache.GetSlave()
	r := make([]ProjectData, 0)
	if err := db.Table("projects").Where("name LIKE ?", "%"+name+"%").Limit(limit, start).Find(&r); err != nil {
		return nil, err
	}
	return r, nil
}

/**
 * 添加项目
 */
func (ap *ProjectData) AddProject () (int64, error) {
	db := base.DbCache.GetMaster()
	r := make([]ProjectData, 0)
	// 判断是否已经存在
	if err := db.Table("projects").Where("name = ?", ap.Name).Limit(1, 0).Find(&r); err != nil {
		return 1, err
	}
	if len(r) > 0 {
		return 1, errors.New("this name (" + ap.Name + ") already exists")
	}
	res, err := db.Table("projects").Insert(ap)
	return res, err
}

/**
 * 更新项目
 */
func (ap *ProjectData) UpdateProject (id int) (int64, error) {
	db := base.DbCache.GetMaster()
	r := make([]ProjectData, 0)
	if err := db.Table("projects").Where("id = ?", id).Limit(1, 0).Find(&r); err != nil {
		return 1, err
	}
	if len(r) < 1 {
		return 1, errors.New("this project (" + strconv.Itoa(id) + ") does not exist")
	}
	res, err := db.Table("projects").Where("id = ?", id).Update(ap)
	if err != nil {
		setting.Logger.Error(err.Error())
		return res, err
	}

	return res, nil
}

/**
 * 删除项目
 */
func (ap *ProjectData) DeleteProject () (int64, error) {
	db := base.DbCache.GetSlave()

	res, err := db.Table("projects").Where("id = ?", ap.Id).Delete(ap)
	if err != nil {
		setting.Logger.Error(err.Error())
		return res, err
	}
	return res, nil
}