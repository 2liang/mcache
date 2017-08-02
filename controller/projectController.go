package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/2liang/mcache/modules/utils/setting"
	"github.com/2liang/mcache/models"
)

type GetProjectParams struct {

}

type AddProjectParams struct {
	Name 	string 	`form:"name" json:"name" binding:"required"`
	Desc 	string 	`form:"desc" json:"desc" binding:"required"`
}

func GetProject(c *gin.Context) {
	c.JSON(http.StatusOK, "get_project")
}

func AddProject(c *gin.Context) {
	var params AddProjectParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error("param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	}
	// 开始操作
	data := new(models.AddProjectData)
	data.Name = params.Name
	data.Desc = params.Desc
	_, err := data.AddProject()
	if err != nil {
		setting.SeeLog.Error("insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "新增项目成功", "data": nil})
}

func UpdateProject(c *gin.Context) {
	c.JSON(http.StatusOK, "update_project")
}

func DeleteProject(c *gin.Context) {
	c.JSON(http.StatusOK, "delete_project")
}
