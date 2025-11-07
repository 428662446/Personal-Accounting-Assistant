package handlers

import (
	"AccountingAssistant/services"
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
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return (&AuthHandler{userService})
}
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "用户名、密码不能为空",
		})
		return
	}
	userId, err := h.userService.Register(req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	userId, err := h.userService.Login(req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"user_id": userId,
	})
}
