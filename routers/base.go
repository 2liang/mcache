//Package routers provide the all routers
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/public", http.Dir("public"))
	//心跳
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": 0, "msg": "ok", "data": "pong"})
		ctx.Abort()
	})

	//
	control_r := r.Group("/control")
	{
		control_r.GET("/get")
	}
}
