package handlers

import (
	"AccountingAssistant/services"
	"AccountingAssistant/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册请求结构体
type RegisterUserRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// 登录请求结构体
type LoginUserRequest struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// 报错记录：cannot define new methods on non-local type
// 修改
type AuthHandler struct {
	userService    *services.UserService
	sessionManager *services.DBSessionManager
}

func NewAuthHandler(userService *services.UserService, sessionManager *services.DBSessionManager) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionManager: sessionManager,
	}
}
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.HandleError(c, utils.ErrEmptyCredential)
		return
	}
	userId, err := h.userService.Register(req.UserName, req.Password)
	if err != nil {
		utils.HandleError(c, err) // 使用统一的错误处理
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功",
		"user_id": userId,
	})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.HandleError(c, utils.ErrEmptyCredential)
		return
	}
	userID, err := h.userService.Login(req.UserName, req.Password)
	if err != nil {
		utils.HandleError(c, err) // 使用统一的错误处理
		return
	}
	// 创建会话
	sessionID, err := h.sessionManager.CreateSession(userID, req.UserName)
	if err != nil {
		utils.HandleError(c, err) // 修改: 直接返回ErrCreateSessionFailed会丢失错误信息
		return
	}
	// 设置Cookie（浏览器自动保存）
	c.SetCookie("session_id", sessionID, 24*3600, "/", "", false, true) // 24小时过期

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"user_id": userID,
	})
}

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	// 从Cookie获取sessionID
	sessionID, err := c.Cookie("session_id")
	if err == nil {
		h.sessionManager.DeleteSession(sessionID)
	}

	// 清除Cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "退出成功",
	})
}
