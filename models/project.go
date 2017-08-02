package models

import (
	"github.com/2liang/mcache/models/base"
	"time"
)

type AddProjectData struct {
	Id 		int			`xorm:"'id' int(11)" json:"id"`
	Name	string		`xorm:"'name' varchar(250)" json:"name"`
	Desc 	string		`xorm:"'name' varchar(250)" json:"desc"`
	CreateTime	int64		`xorm:"'create_time' int(11)" json:"create_time"`
	ModifyTime	int64		`xorm:"'modify_time' int(11)" json:"modify_time"`
}

func (ap *AddProjectData) AddProject () (int64, error) {
	var x = base.DbCache.GetMaster()
	sess := x.NewSession()
	ap.CreateTime = time.Now().Unix()
	res, err := sess.Table("project").Insert(ap)
	return res, err
}
