package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/2liang/mcache/modules/utils/setting"
	"time"
	"strings"
	"github.com/2liang/mcache/models/sqlmodel"
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

func GetCase(c *gin.Context) {
	c.JSON(http.StatusOK, "get_case")
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
	c.JSON(http.StatusOK, "update_case")
}

func DeleteCase(c *gin.Context) {
	c.JSON(http.StatusOK, "delete_case")
}