package controllers

import (
	"encoding/json"
	"errors"
	"ginDemo/common"
	"ginDemo/models"
	"ginDemo/services"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func checkInviteeId(userId uint) bool {
	_, err := services.SelectInviteeId(userId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	return false
}

func UserInviteList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		common.ErrorJson(1000, c)
		return
	}
	currentUser := user.(*models.User)

	pageInfo := c.Query("page_and_per_page")

	if pageInfo == "" {
		common.ErrorJson(1000, c)
		c.Abort()
		return
	}

	// 解码
	decryptedJson, err := utils.Decrypt(pageInfo)
	if err != nil {
		common.ErrorJson(1000, c)
		return
	}

	var pageAndPerPage models.PageAndPerPage

	if err := json.Unmarshal(decryptedJson, &pageAndPerPage); err != nil {
		common.ErrorJson(1000, c)
		return
	}

	userId := currentUser.ID
	list, err := services.InviteList(userId, pageAndPerPage.Page, pageAndPerPage.PerPage)
	if err != nil {
		common.ErrorJson(1002, c)
		return
	}

	if len(list) == 0 {
		list = make([]*models.ResponseInviteDetail, 0)
	}
	common.WbeJson(list, c)
}
