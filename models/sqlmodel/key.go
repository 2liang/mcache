package sqlmodel

import (
	"github.com/2liang/mcache/models/base"
	"errors"
	"strconv"
)

type KeyData struct {
	Id 			int		`xorm:"id int(11)" json:"id"`
	CaseId 		int 	`xorm:"case_id int(10)" json:"case_id"`
	Name 		string 	`xorm:"name varchar(250)" json:"name"`
	Desc 		string	`xorm:"desc varchar(250)" json:"desc"`
	Prefix 		string	`xorm:"prefix varchar(250)" json:"prefix"`
	KeyType 	string	`xorm:"key_type varchar(250)" json:"key_type"`
	CreateTime	int64	`xorm:"create_time int(11)" json:"create_time"`
	ModifyTime	int64	`xorm:"modify_time int(11)" json:"modify_time"`
}

func(kd *KeyData) GetKeyById() ([]KeyData, error) {
	db := base.DbCache.GetSlave()
	r := make([]KeyData, 0)
	if err := db.Table("keys").Where("id = ?", kd.Id).Find(&r); err != nil {
		return nil, err
	}
	if len(r) < 1 {
		return nil, errors.New("this key (" + strconv.Itoa(kd.Id) + ") does not exists!")
	}
	return r, nil
}

func(kd *KeyData) GetKey(cid int, name string, page int, limit int) ([]KeyData, error) {
	start := (page - 1) * limit
	db := base.DbCache.GetSlave()
	// 判断cid是否存在
	CaseData := new(CaseData)
	CaseInfo, err := CaseData.GetCaseById(cid)
	if err != nil {
		return nil, err
	}
	if len(CaseInfo) < 1 {
		return nil, errors.New("this is case(" + strconv.Itoa(cid) + ") does not exist!")
	}
	r := make([]KeyData, 0)
	if err := db.Table("keys").Where("case_id = ? AND name LIKE ?", cid, "%" + name + "%").Limit(limit, start).Find(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func(kd *KeyData) AddKey() (int64, error) {
	db := base.DbCache.GetMaster()
	// 判断实例是否存在
	cd := new(CaseData)
	caseInfo, err := cd.GetCaseById(kd.CaseId)
	if err != nil {
		return 1, err
	}
	if len(caseInfo) < 1 {
		return 1, errors.New("this is case(" + strconv.Itoa(kd.CaseId) + ") does not exist!")
	}
	r := make([]KeyData, 0)
	if err := db.Table("keys").Where("name = ?", kd.Name).Limit(1, 0).Find(&r); err != nil {
		return 1, err
	}
	if len(r) > 0 {
		return 1, errors.New("this is name(" + kd.Name + ") already exists!")
	}
	res, err := db.Table("keys").Insert(kd)
	return res, err
}

func(kd *KeyData) UpdateKey (id int) (int64, error) {
	db := base.DbCache.GetMaster()
	keyInfo, err := kd.GetKeyById()
	if err != nil {
		return 1, err
	}
	if len(keyInfo) < 1 {
		return 1, errors.New("this is key(" + strconv.Itoa(kd.Id) + ") does not exists!")
	}

	r := make([]KeyData, 0)
	if err := db.Table("keys").Where("name = ?", kd.Name).Limit(1, 0).Find(&r); err != nil {
		return 1, err
	}
	if len(r) > 0 {
		return 1, errors.New("this is name(" + kd.Name + ") already exists!")
	}

	res, err := db.Table("keys").Where("id = ?", id).Update(kd)
	if err != nil {
		return 1, err
	}
	return res, nil
}

func(kd *KeyData) DeleteKey () (int64, error) {
	db := base.DbCache.GetMaster()
	keyInfo, err := kd.GetKeyById()
	if err != nil {
		return 1, err
	}
	if len(keyInfo) < 1 {
		return 1, errors.New("this is key(" + strconv.Itoa(kd.Id) + ") does not exists!")
	}
	res, err := db.Table("keys").Where("id = ?", kd.Id).Delete(kd)
	return res, err
}
