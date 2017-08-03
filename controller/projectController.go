package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/2liang/mcache/modules/utils/setting"
	"github.com/2liang/mcache/models"
	"time"
)

type GetProjectParams struct {
	Page 	int	`form:"page" json:"page" binding:"required"`
	Limit 	int	`form:"limit" json:"limit" binding:"required"`
	Name 	string `form:"name" json:"name"`
}

type AddProjectParams struct {
	Name 	string 	`form:"name" json:"name" binding:"required"`
	Desc 	string 	`form:"desc" json:"desc" binding:"required"`
}

func GetProject(c *gin.Context) {
	var params GetProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error("insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": params})
		return
	}
	data := new(models.ProjectData)
	res, err := data.GetProject(params.Name, params.Page, params.Limit)
	if err != nil {
		setting.SeeLog.Error("insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "获取项目成功", "data": res})
	}
}

func AddProject(c *gin.Context) {
	var params AddProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error("param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	}
	// 开始操作
	data := new(models.ProjectData)
	data.Name = params.Name
	data.Desc = params.Desc
	data.CreateTime = time.Now().Unix()
	res, err := data.AddProject()
	if err != nil {
		setting.SeeLog.Error("insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "新增项目成功", "data": res})
	}
}

func UpdateProject(c *gin.Context) {
	c.JSON(http.StatusOK, "update_project")
}

func DeleteProject(c *gin.Context) {
	c.JSON(http.StatusOK, "delete_project")
}
