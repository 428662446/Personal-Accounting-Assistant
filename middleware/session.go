package middleware

import (
	"AccountingAssistant/services"

	"github.com/gin-gonic/gin"
)

// 完全不会，由deepseek提供
func SessionMiddleware(sessionManager *services.DBSessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie获取sessionID
		sessionID, err := c.Cookie("session_id")
		if err == nil {
			// 如果有sessionID，就验证它
			userID, username, valid := sessionManager.ValidateSession(sessionID)
			if valid {
				// 会话有效，保存用户信息
				c.Set("userID", userID)
				c.Set("username", username)
			}
		}
		// 不管有没有session，都继续处理（让AuthRequired决定是否需要登录）
		c.Next()
	}
}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "请先登录",
			})
			c.Abort() // 停止后继续处理
			return
		}

		// 确保userID是int64类型
		if _, ok := userID.(int64); !ok {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "无效的用户会话",
			})
			c.Abort()
			return
		}
		c.Next() // 验证通过，继续处理
	}
}
