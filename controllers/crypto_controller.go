package controllers

import (
	"ginDemo/services"
	"github.com/gin-gonic/gin"
)

func RandomTransfer(c *gin.Context) {

	address := c.PostForm("address")

	services.Transfer(address)

}
