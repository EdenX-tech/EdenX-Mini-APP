package middlewares

import (
	"fmt"
	"ginDemo/common"
	"ginDemo/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		log.Println(authHeader)
		if authHeader == "" {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}

		tokenString := parts[1]

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

		// 将用户信息存储在上下文中
		c.Set("user", user)
		c.Next()
	}
}
