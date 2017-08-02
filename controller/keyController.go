package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetKey(c *gin.Context) {
	c.JSON(http.StatusOK, "get_key")
}

func AddKey(c *gin.Context) {
	c.JSON(http.StatusOK, "add_key")
}

func UpdateKey(c *gin.Context) {
	c.JSON(http.StatusOK, "update_key")
}

func DeleteKey(c *gin.Context) {
	c.JSON(http.StatusOK, "delete_key")
}