package middleware

import (
	"AccountingAssistant/services"
	"AccountingAssistant/utils"
	"AccountingAssistant/web/response"

	"github.com/gin-gonic/gin"
)

// SessionMiddleware 从 Cookie 中检索 session_id 并验证，会将用户信息写入上下文。
// 在无法获取或验证会话时，使用统一的错误处理器返回错误信息。
func SessionMiddleware(sessionManager *services.DBSessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie获取sessionID（或Header）
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			response.HandleError(c, utils.ErrNotLoggedIn)
			c.Abort()
			return
		}

		userID, username, valid := sessionManager.ValidateSession(sessionID)
		if !valid {
			response.HandleError(c, utils.ErrInvalidSession)
			c.Abort()
			return
		}
		// 将会话信息存入上下文
		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			response.HandleError(c, utils.ErrNotLoggedIn)
			c.Abort()
			return
		}

		// 确保userID是int64类型
		if _, ok := userID.(int64); !ok {
			response.HandleError(c, utils.ErrInvalidSession)
			c.Abort()
			return
		}
		c.Next()
	}
}
