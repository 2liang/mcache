//Package routers provide the all routers
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/2liang/mcache/controller"
	"github.com/2liang/mcache/controller"
)

func Init(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/public", http.Dir("public"))
	//心跳
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "msg": "ok", "data": "pong"})
		ctx.Abort()
	})

	// 项目
	project_r := r.Group("/project")
	{
		// 获取项目
		project_r.GET("/get", controller.GetProject)
		// 添加项目
		project_r.POST("/add", controller.AddProject)
		// 更新项目
		project_r.POST("/update", controller.UpdateProject)
		// 删除项目
		project_r.POST("/delete", controller.DeleteProject)
	}

	// 实例
	case_r := r.Group("/case")
	{
		// 获取实例
		case_r.GET("/get", controller.GetCase)
		// 添加实例
		case_r.POST("/add", controller.AddCase)
		// 更新实例
		case_r.POST("/update", controller.UpdateCase)
		// 删除实例
		case_r.POST("/delete", controller.DeleteCase)
	}

	// key
	key_r := r.Group("/key")
	{
		// 获取key
		key_r.GET("/get", controller.GetKey)
		// 添加key
		key_r.POST("/add", controller.AddKey)
		// 更新key
		key_r.POST("/update", controller.UpdateKey)
		// 删除key
		key_r.POST("/delete", controller.DeleteKey)
	}
}
