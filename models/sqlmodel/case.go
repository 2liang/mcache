package sqlmodel

import (
	"github.com/2liang/mcache/models/base"
	"errors"
	"strconv"
)

type CaseData struct {
	Id 			int		`xorm:"'id' int(11)" json:"id"`
	ProjectId 	int 	`xorm:"'project_id' int(11)" json:"project_id"`
	Name 		string 	`xorm:"'name' varchar(250)" json:"name"`
	Desc 		string 	`xorm:"'desc' varchar(250)" json:"desc"`
	Type 		string 	`xorm:"'type' varchar(250)" json:"type"`
	MasterHost 	string 	`xorm:"'master_host' varchar(250)" json:"master_host"`
	SlaveHost 	string 	`xorm:"'slave_host' varchar(250)" json:"slave_host"`
	Port 		int 	`xorm:"'port' int(11)" json:"port"`
	CreateTime 	int64 	`xorm:"'create_time' int(11)" json:"create_time"`
	ModifyTime 	int64 	`xorm:"'modify_time' int(11)" json:"modify_time"`
}

func (cd *CaseData) GetCase (pid int, name string, page int, limit int) ([]CaseData, error) {
	// 判断项目是否存在
	ProjectData := new(ProjectData)
	ProjectData.Id = pid
	res, err := ProjectData.GetProjectByPid()
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("this project (" + strconv.Itoa(pid) + ") does not exist!")
	}
	start := (page - 1) * limit
	db := base.DbCache.GetSlave()
	r := make([]CaseData, 0)
	if err := db.Table("cases").Where("project_id = ? AND name LIKE ?", pid, "%" + name + "%").Limit(limit, start).Find(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func (cd *CaseData) GetCaseById(id int) ([]CaseData, error) {
	db := base.DbCache.GetSlave()
	r := make([]CaseData, 0)
	if err := db.Table("cases").Where("id = ?", id).Limit(1, 0).Find(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func (cd *CaseData) AddCase () (int64, error) {
	db := base.DbCache.GetMaster()
	// 判断项目是否存在
	ProjectData := new(ProjectData)
	ProjectData.Id = cd.ProjectId
	res, err := ProjectData.GetProjectByPid()
	if err != nil {
		return 1, err
	}

	if len(res) < 1 {
		return 1, errors.New("this project (" + strconv.Itoa(cd.ProjectId) + ") does not exist!")
	}

	r := make([]CaseData, 0)
	if err := db.Table("cases").Where("name = ?", cd.Name).Limit(1, 0).Find(&r); err != nil {
		return 1, err
	}

	if len(r) > 0 {
		return 1, errors.New("this name (" + cd.Name + ") already exists")
	}

	insert, err := db.Table("cases").Insert(cd)
	return insert, err
}

func (cd *CaseData) UpdateCase(id int) (int64, error) {
	db := base.DbCache.GetMaster()

	// 判断实例是否存在
	res, err := cd.GetCaseById(id)
	if err != nil {
		return 1, err
	}

	if len(res) < 1 {
		return 1, errors.New("this case (" + strconv.Itoa(id) + ") does not exist!")
	}

	updateRes, err := db.Table("cases").Where("id = ?", id).Update(cd)
	if err != nil {
		return 1, err
	}

	return updateRes, nil
}

func (cd *CaseData) DeleteCase() (int64, error) {
	db := base.DbCache.GetMaster()

	// 判断实例是否存在
	res, err := cd.GetCaseById(cd.Id)
	if err != nil {
		return 1, err
	}
	if len(res) < 1 {
		return 1, errors.New("this case (" + strconv.Itoa(cd.Id) + ") does not exist!")
	}

	delRes, err := db.Table("cases").Where("id = ?", cd.Id).Delete(cd)
	if err != nil {
		return 1, err
	}
	return delRes, nil
}