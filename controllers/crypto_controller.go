package controllers

import (
	"ginDemo/services"
	"github.com/gin-gonic/gin"
)

func RandomTransfer(c *gin.Context) {
	println("test:", 111)
	address := c.PostForm("address")

	services.Transfer(address)

}
