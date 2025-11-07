package middleware

import (
	"AccountingAssistant/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 暂时看不懂
func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供令牌"})
			c.Abort()
			return
		}

		// 验证token，获取userID
		userID, err := userService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效令牌"})
			c.Abort()
			return
		}

		// 将userID存入上下文
		c.Set("userID", userID)
		c.Next()
	}
}
