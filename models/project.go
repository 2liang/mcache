package models

import (
	"github.com/2liang/mcache/models/base"
	"github.com/2liang/mcache/modules/utils/setting"
	"errors"
)

type ProjectData struct {
	Id 		int			`xorm:"'id' int(11)" json:"id"`
	Name	string		`xorm:"'name' varchar(250)" json:"name"`
	Desc 	string		`xorm:"'desc' varchar(250)" json:"desc"`
	CreateTime	int64		`xorm:"'create_time' int(11)" json:"create_time"`
	ModifyTime	int64		`xorm:"'modify_time' int(11)" json:"modify_time"`
}

func (ap *ProjectData) GetProject (name string, page int, limit int) ([]ProjectData, error) {
	start := (page - 1) * limit
	db := base.DbCache.GetSlave()
	r := make([]ProjectData, 0)
	if err := db.Table("project").Where("name LIKE ?", "%"+name+"%").Limit(limit, start).Find(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func (ap *ProjectData) AddProject () (int64, error) {
	db := base.DbCache.GetMaster()
	setting.SeeLog.Error("name:" + ap.Name + "; desc:" + ap.Desc)
	r := make([]ProjectData, 0)
	if err := db.Table("project").Where("name = ?", ap.Name).Find(&r); err != nil {
		return 1, err
	}
	if len(r) > 0 {
		return 1, errors.New("this name (" + ap.Name + ") already exists")
	}
	sess := db.NewSession()
	// 判断是否已经存在
	res, err := sess.Table("project").Insert(ap)
	return res, err
}
