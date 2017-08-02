package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCase(c *gin.Context) {
	c.JSON(http.StatusOK, "get_case")
}

func AddCase(c *gin.Context) {
	c.JSON(http.StatusOK, "add_case")
}

func UpdateCase(c *gin.Context) {
	c.JSON(http.StatusOK, "update_case")
}

func DeleteCase(c *gin.Context) {
	c.JSON(http.StatusOK, "delete_case")
}