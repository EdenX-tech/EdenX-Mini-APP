package controllers

import (
	"ginDemo/common"
	"ginDemo/models"
	"ginDemo/services"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func RandomTransfer(c *gin.Context) {
	var request models.TransferRequest
	if err := c.BindJSON(&request); err != nil {
		common.ErrorJson(1000, c)
		return
	}
	address := request.Address
	println("address:::", address)
	isEarnStatus := randomTransferStatus()
	transInfo := models.ReturnEarnResult{
		IsEarn: isEarnStatus,
		Amount: 0,
	}
	println("isEarnStatus:", isEarnStatus)
	if isEarnStatus {
		transferAmount := randomEarn()
		transStatus := services.Transfer(address, transferAmount)
		if transStatus {
			transInfo.Amount = transferAmount
		} else {
			common.ErrorJson(1004, c)
			return
		}
	} else {
		transInfo.Amount = 0
	}

	common.WbeJson(transInfo, c)
}

func randomTransferStatus() bool {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	if random.Intn(2) == 0 {
		// 50% 的概率执行此代码块
		return true
	} else {
		// 50% 的概率执行此代码块
		return false
	}
}

func randomEarn() uint64 {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// 生成一个在 0.01 到 0.09 之间的随机浮点数
	randomAmount := 0.01 + random.Float64()*(0.09-0.01)

	amount := uint64(randomAmount * 1000000)
	println("randomEarn:", amount)
	return amount
}
