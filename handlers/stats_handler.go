package handlers

import (
	"AccountingAssistant/services"
	"AccountingAssistant/utils"
	"AccountingAssistant/web/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理统计服务 对象
type StatHandler struct {
	statService *services.StatService
}

// 新建 处理统计服务 对象的方法
func NewStatHandler(statService *services.StatService) *StatHandler {
	return &StatHandler{
		statService: statService,
	}
}

// 对应的HTTP响应
func (h *StatHandler) GetSummary(c *gin.Context) {
	// 从会话或JWT令牌中获取用户ID，而不是从URL参数
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	totalIncome, err := h.statService.GetTotalIncome(userID.(int64))
	if err != nil {
		response.HandleError(c, err) // 使用统一的错误处理
		return
	}

	totalExpenditure, err := h.statService.GetTotalExpenditure(userID.(int64))
	if err != nil {
		response.HandleError(c, err) // 使用统一的错误处理
		return
	}
	totalNetIncome, err := h.statService.GetNetIncome(userID.(int64))
	if err != nil {
		response.HandleError(c, err) // 使用统一的错误处理
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{ // 统一使用 "data" 字段包装统计结果
			"total_income":      totalIncome,
			"total_expenditure": totalExpenditure,
			"total_net_income":  totalNetIncome,
		},
	})
}

func (h *StatHandler) GetMonthlyStats(c *gin.Context) {
	// 从会话中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	totalIncome, totalExpenditure, totalNetIncome, err := h.statService.GetMonthlyStats(userID.(int64))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"total_income":      totalIncome,
			"total_expenditure": totalExpenditure,
			"total_net_income":  totalNetIncome,
		},
	})
}

func (h *StatHandler) GetWeeklyStats(c *gin.Context) {
	// 从会话中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	totalIncome, totalExpenditure, totalNetIncome, err := h.statService.GetWeeklyStats(userID.(int64))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"total_income":      totalIncome,
			"total_expenditure": totalExpenditure,
			"total_net_income":  totalNetIncome,
		},
	})
}

func (h *StatHandler) GetDailyStats(c *gin.Context) {
	// 从会话中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	totalIncome, totalExpenditure, totalNetIncome, err := h.statService.GetDailyStats(userID.(int64))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"total_income":      totalIncome,
			"total_expenditure": totalExpenditure,
			"total_net_income":  totalNetIncome,
		},
	})
}

func (h *StatHandler) GetRangeAmountStats(c *gin.Context) {
	// 从会话中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	rangeAmountStats, err := h.statService.GetRangeAmountStats(userID.(int64))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"amount_range_stats": rangeAmountStats,
		},
	})
}
