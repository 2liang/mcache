package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/2liang/mcache/modules/utils/setting"
	"time"
	"strings"
	"github.com/2liang/mcache/models/sqlmodel"
	"strconv"
)

var c_name  = "caseController "

type GetCaseParams struct {
	Pid 	int 	`form:"pid" json:"pid" binding:"required"`
	Page 	int 	`form:"page" json:"page" binding:"required"`
	Limit 	int 	`form:"limit" json:"limit" binding:"required"`
	Name 	string 	`form:"name" json:"name"`
}

type AddCaseParams struct {
	Pid 	int		`form:"pid" json:"pid" binding:"required"`
	Name	string	`form:"name" json:"name" binding:"required"`
	Desc 	string  `form:"desc" json:"desc"`
	Type 	string	`form:"type" json:"type" binding:"required"`
	MasterHost	string	`form:"master_host" json:"master_host" binding:"required"`
	SlaveHost	string	`form:"slave_host" json:"slave_host"`
	Port 		int		`form:"port" json:"port"`
}

type UpdateCaseParams struct {
	Id 			int		`form:"id" json:"id" binding:"required"`
	Name		string	`form:"name" json:"name" binding:"required"`
	Desc 		string  `form:"desc" json:"desc"`
	Type 		string	`form:"type" json:"type" binding:"required"`
	MasterHost	string	`form:"master_host" json:"master_host" binding:"required"`
	SlaveHost	string	`form:"slave_host" json:"slave_host"`
	Port 		int		`form:"port" json:"port"`
}

type DeleteCaseParams struct {
	Id 		int	`form:"id" json:"id" binding:"required"`
}

func GetCaseById(c *gin.Context) {
	id := c.Param("id")
	cid, err := strconv.Atoi(id)
	if err != nil {
		setting.SeeLog.Error(c_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": id})
		return
	}

	data := new(sqlmodel.CaseData)
	res, err := data.GetCaseById(cid)
	if err != nil {
		setting.SeeLog.Error(c_name + "get case error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": res})
		return
	}
	if len(res) < 1 {
		setting.SeeLog.Info(c_name + "this case (" + id + ") does not exist!")
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "this case (" + id + ") does not exist!", "data": res})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "获取实例成功", "data": res[0]})
	}
}

func GetCase(c *gin.Context) {
	var params GetCaseParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(c_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": params})
		return
	}
	data := new(sqlmodel.CaseData)
	res, err := data.GetCase(params.Pid, params.Name, params.Page, params.Limit)
	if err != nil {
		setting.SeeLog.Error(c_name + "get data error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "获取实例成功", "data": res})
	}
}

func AddCase(c *gin.Context) {
	var params AddCaseParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(c_name + " param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": params})
		return
	}
	CaseData := new(sqlmodel.CaseData)
	if strings.ToUpper(params.Type) != "REDIS" && strings.ToUpper(params.Type) != "MC" {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "type参数不合法，仅支持：redis,mc", "data": nil})
		return
	}
	if params.SlaveHost == "" {
		params.SlaveHost = params.MasterHost
	}
	CaseData.ProjectId 	= params.Pid
	CaseData.Name		= params.Name
	CaseData.Desc		= params.Desc
	CaseData.Type		= strings.ToUpper(params.Type)
	CaseData.MasterHost = params.MasterHost
	CaseData.SlaveHost	= params.MasterHost
	CaseData.SlaveHost	= params.SlaveHost
	CaseData.Port		= params.Port
	CaseData.CreateTime	= time.Now().Unix()
	res, err := CaseData.AddCase()
	if err != nil {
		setting.SeeLog.Error(c_name + " insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": params})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "添加实例成功", "data": res})
}

func UpdateCase(c *gin.Context) {
	var params UpdateCaseParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(c_name + " param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": params})
		return
	}
	CaseData := new(sqlmodel.CaseData)
	if strings.ToUpper(params.Type) != "REDIS" && strings.ToUpper(params.Type) != "MC" {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "type参数不合法，仅支持：redis,mc", "data": nil})
		return
	}
	if params.SlaveHost == "" {
		params.SlaveHost = params.MasterHost
	}

	CaseData.Name		= params.Name
	CaseData.Desc		= params.Desc
	CaseData.Type		= strings.ToUpper(params.Type)
	CaseData.MasterHost = params.MasterHost
	CaseData.SlaveHost	= params.MasterHost
	CaseData.SlaveHost	= params.SlaveHost
	CaseData.Port		= params.Port
	CaseData.ModifyTime = time.Now().Unix()
	res, err := CaseData.UpdateCase(params.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": res})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "更新实例成功", "data": res})
	}
}

func DeleteCase(c *gin.Context) {
	var params DeleteCaseParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(c_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}

	CaseData := new(sqlmodel.CaseData)
	CaseData.Id = params.Id
	res, err := CaseData.DeleteCase()
	if err != nil {
		setting.SeeLog.Error(c_name + "delete error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": res})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "删除实例成功", "data": res})
}