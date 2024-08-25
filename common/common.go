package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"ginDemo/config"
	"ginDemo/models"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/big"
	"net/http"
	"strconv"
)

const (
	inviteCodeLength = 8
	charset          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func getLanguage(c *gin.Context) string {
	language, _ := c.Get("language")
	langStr, ok := language.(string)
	if !ok || langStr == "" {
		langStr = "EN"
	}
	return langStr
}

func ErrorJson(code uint, c *gin.Context) {
	langStr := getLanguage(c)
	msg, codeStatus := config.CustomErrorCodes(code, langStr)

	errorResponse := models.ErrorResponse{
		ErrorCode: code,
		Message:   msg,
	}
	c.JSON(codeStatus, errorResponse)
	// 终止后续处理
	c.Abort()
}

func WbeJson(data interface{}, c *gin.Context) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"jsonData_error": err.Error()})
		return
	}

	encryptedData, err := utils.Encrypt(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"encryptedData_error": err.Error()})
		return
	}
	response := models.Response{
		ErrorCode: 0,
		Message:   "OK",
		Data:      encryptedData,
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func ErrorSocketJson(conn *websocket.Conn, code uint, c *gin.Context) {
	langStr := getLanguage(c)
	msg, _ := config.CustomErrorCodes(code, langStr)

	errorResponse := models.ErrorResponse{
		ErrorCode: code,
		Message:   msg,
	}

	res, _ := json.Marshal(errorResponse)
	conn.WriteMessage(websocket.TextMessage, res)
}

func SocketJson(conn *websocket.Conn, data interface{}, errorCode uint) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encryptedData, err := utils.Encrypt(jsonData)
	if err != nil {
		return err
	}

	response := models.Response{
		ErrorCode: errorCode,
		Message:   "OK",
		Data:      encryptedData,
	}
	res, _ := json.Marshal(response)
	fmt.Println("json处理:", string(res))
	if err := conn.WriteMessage(websocket.TextMessage, res); err != nil {
		return err
	}
	return nil
}

func HandleError(conn *websocket.Conn, err error, code uint, c *gin.Context) bool {
	if err != nil {
		ErrorSocketJson(conn, code, c)
		return true
	}
	return false
}

func GenerateInviteCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[num.Int64()]
	}
	return string(code)
}

func StrToInt(param string) (int, error) {
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return paramInt, nil
}
