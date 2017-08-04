package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/2liang/mcache/modules/utils/setting"
	"strconv"
	"github.com/2liang/mcache/models/sqlmodel"
	"time"
)

var k_name = "keyController "

type GetKeyParams struct {
	CaseId	int		`form:"case_id" json:"case_id" binding:"required"`
	Name 	string 	`form:"name" json:"name"`
	Page	int		`form:"page" json:"page" binding:"required"`
	Limit 	int 	`form:"limit" json:"limit" binding:"required"`
}

type AddKeyParams struct {
	CaseId 		int 	`form:"case_id" json:"case_id" binding:"required"`
	Name 		string 	`form:"name" json:"name" binding:"required"`
	Desc 		string	`form:"desc" json:"desc" binding:"required"`
	Prefix 		string	`form:"prefix" json:"prefix" binding:"required"`
	KeyType 	string	`form:"key_type" json:"key_type" binding:"required"`
}

type UpdateKeyParams struct {
	Id 			int 	`form:"id" json:"id" binding:"required"`
	Name 		string 	`form:"name" json:"name" binding:"required"`
	Desc 		string	`form:"desc" json:"desc" binding:"required"`
	Prefix 		string	`form:"prefix" json:"prefix" binding:"required"`
	KeyType 	string	`form:"key_type" json:"key_type" binding:"required"`
}

type DeleteKeyParams struct {
	Id 			int 	`form:"id" json:"id" binding:"required"`
}

func GetKeyById(c *gin.Context) {
	id := c.Param("id")
	kid, err := strconv.Atoi(id)
	if err != nil {
		setting.SeeLog.Error(k_name + "id strconv.atoi error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	kd := new(sqlmodel.KeyData)
	kd.Id = kid
	res, err := kd.GetKeyById()
	if err != nil {
		setting.SeeLog.Error(k_name + "get key error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "获取key成功", "data": res[0]})
}

func GetKey(c *gin.Context) {
	var params GetKeyParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(k_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	if params.Limit > 20 {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "每页最多请求20条数据", "data": nil})
		return
	}

	kd := new(sqlmodel.KeyData)
	res, err := kd.GetKey(params.CaseId, params.Name, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "获取key成功", "data": res})
}

func AddKey(c *gin.Context) {
	var params AddKeyParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(k_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	kd := new(sqlmodel.KeyData)
	kd.CaseId = params.CaseId
	kd.Name = params.Name
	kd.Desc = params.Desc
	kd.Prefix = params.Prefix
	kd.KeyType = params.KeyType
	kd.CreateTime = time.Now().Unix()
	res, err := kd.AddKey()
	if err != nil {
		setting.SeeLog.Error(k_name + "insert into error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "添加key成功", "data": res})
}

func UpdateKey(c *gin.Context) {
	var params UpdateKeyParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(k_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	kd := new(sqlmodel.KeyData)
	kd.Id   = params.Id
	kd.Name = params.Name
	kd.Desc = params.Desc
	kd.Prefix = params.Prefix
	kd.KeyType = params.KeyType
	kd.ModifyTime = time.Now().Unix()
	res, err := kd.UpdateKey(params.Id)
	if err != nil {
		setting.SeeLog.Error(k_name + "update error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "更新key成功", "data": res})
}

func DeleteKey(c *gin.Context) {
	var params DeleteKeyParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error(k_name + "param error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	kd := new(sqlmodel.KeyData)
	kd.Id = params.Id
	res, err := kd.DeleteKey()
	if err != nil {
		setting.SeeLog.Error(k_name + "delete error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0, "msg": "删除key成功", "data": res})
}