package middlewares

import (
	"ginDemo/common"
	"github.com/gin-gonic/gin"
)

func SystemLanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		language := c.GetHeader("Accept-Language")
		if language == "" {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Accept-Language header required"})
			common.ErrorJson(1000, c)
			c.Abort()
			return
		}

		var systemLanguage string
		switch language {
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
		c.Set("language", systemLanguage)
		c.Next()
	}
}
