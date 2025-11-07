package handlers

import (
	"AccountingAssistant/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecordRequest struct {
	Username string  `form:"username" binding:"required"`
	Type     string  `form:"type" binding:"required"`
	Amount   float64 `form:"amount" binding:"required"`
	Category string  `form:"category" binding:"required"`
	Note     string  `form:"note" binding:"required"`
}
type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return (&TransactionHandler{transactionService})
}
func (h *TransactionHandler) RecordTransaction(c *gin.Context) {
	var req RecordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "账单内容不能为空",
		})
		return
	}
	transactionId, err := h.transactionService.RecordTransaction(req.Username, req.Type, req.Amount, req.Category, req.Note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "记录成功",
		"transaction_id": transactionId,
	})
}
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// 从会话或JWT令牌中获取用户ID，而不是从URL参数
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "未登录",
		})
		return
	}
	transactions, err := h.transactionService.GetTransactions(userID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "获取成功",
		"transactions": transactions,
	})
}
