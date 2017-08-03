package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/2liang/mcache/modules/utils/setting"
	"github.com/2liang/mcache/models/sqlmodel"
	"time"
)

var p_name = "projectController "

type GetProjectParams struct {
	Page 	int	`form:"page" json:"page" binding:"required"`
	Limit 	int	`form:"limit" json:"limit" binding:"required"`
	Name 	string `form:"name" json:"name"`
}

type AddProjectParams struct {
	Name 	string 	`form:"name" json:"name" binding:"required"`
	Desc 	string 	`form:"desc" json:"desc" binding:"required"`
}

type UpdateProjectParams struct {
	Id 		int `form:"id" json:"id" binding:"required"`
	Name 	string 	`form:"name" json:"name" binding:"required"`
	Desc 	string 	`form:"desc" json:"desc" binding:"required"`
}

type DeleteProjectParams struct {
	Id 		int `form:"id" json:"id" binding:"required"`
}

func GetProject(c *gin.Context) {
	var params GetProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(p_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": params})
		return
	}
	data := new(sqlmodel.ProjectData)
	res, err := data.GetProject(params.Name, params.Page, params.Limit)
	if err != nil {
		setting.SeeLog.Error(p_name + "get data error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "获取项目成功", "data": res})
	}
}

func AddProject(c *gin.Context) {
	var params AddProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(p_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	// 开始操作
	data := new(sqlmodel.ProjectData)
	data.Name = params.Name
	data.Desc = params.Desc
	data.CreateTime = time.Now().Unix()
	res, err := data.AddProject()
	if err != nil {
		setting.SeeLog.Error(p_name + "insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "新增项目成功", "data": res})
	}
}

func UpdateProject(c *gin.Context) {
	var params UpdateProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(p_name + "params error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}

	data := new(sqlmodel.ProjectData)
	data.Name = params.Name
	data.Desc = params.Desc
	data.ModifyTime = time.Now().Unix()

	res, err := data.UpdateProject(params.Id)
	if err != nil {
		setting.SeeLog.Error(p_name + "update error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "更新项目信息成功", "data": res})
}

func DeleteProject(c *gin.Context) {
	var params DeleteProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(p_name + "params error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	data := new(sqlmodel.ProjectData)
	data.Id = params.Id
	res, err := data.DeleteProject()
	if err != nil {
		setting.SeeLog.Error(p_name + "delete error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "删除项目成功", "data": res})
}
