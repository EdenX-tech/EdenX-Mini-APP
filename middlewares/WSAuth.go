package middlewares

import (
	"encoding/json"
	"fmt"
	"ginDemo/common"
	"ginDemo/models"
	"ginDemo/services"
	"ginDemo/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

func WSAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Query("auth")

		if authHeader == "" {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}
		
		// 解码
		decryptedJson, err := utils.Decrypt(authHeader)
		if err != nil {
			common.ErrorJson(1000, c)
		}
		// 打印解密后的数据以进行调试
		var wsAuthInput models.WsAuthInput

		if err := json.Unmarshal(decryptedJson, &wsAuthInput); err != nil {
			common.ErrorJson(1000, c)
			return
		}

		if wsAuthInput.Authorization == "" || wsAuthInput.AcceptLanguage == "" {
			common.ErrorJson(1000, c)
			return
		}

		parts := strings.SplitN(wsAuthInput.Authorization, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}
		//
		tokenString := parts[1]
		//
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil || !token.Valid {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}

		userID := uint(claims["user_id"].(float64))
		user, err := services.GetUserByID(userID)

		if err != nil {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}

		var systemLanguage string
		switch wsAuthInput.AcceptLanguage {
		case "zh-CN":
			systemLanguage = "CN"
			break
		case "en-US":
			systemLanguage = "EN"
			break
		default:
			systemLanguage = "EN"
			break
		}
		// 将用户信息存储在上下文中
		c.Set("user", user)
		c.Set("language", systemLanguage)

		c.Next()
	}
}
