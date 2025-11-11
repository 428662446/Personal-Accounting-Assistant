package handlers

import (
	"AccountingAssistant/services"
	"AccountingAssistant/utils"
	"AccountingAssistant/web/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// 类别要求结构体
type CategoryRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	var r CategoryRequest
	if err := c.ShouldBind(&r); err != nil {
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}
	categoryID, err := h.categoryService.CreateCategory(userID.(int64), r.Name)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "添加成功",
		"category_id": categoryID,
	})

}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	// 从会话或JWT令牌中获取用户ID，而不是从URL参数
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	categories, err := h.categoryService.GetCategory(userID.(int64))
	if err != nil {
		response.HandleError(c, err) // 使用统一的错误处理
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "获取成功",
		"categories": categories,
	})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		// 修复：使用自定义错误而不是直接传递底层错误
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}
	err = h.categoryService.DeleteCategory(userID.(int64), int64(categoryID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}

	err = h.categoryService.UpdateCategory(userID.(int64), int64(categoryID), req.Name)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}
