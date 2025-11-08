package handlers

import (
	"AccountingAssistant/services"
	"AccountingAssistant/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecordRequest struct {
	Type     string  `form:"type" binding:"required"`
	Amount   float64 `form:"amount" binding:"required"`
	Category string  `form:"category" binding:"required"`
	Note     string  `form:"note" binding:"required"`
}
type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService} // 修复：明确指定字段名
}
func (h *TransactionHandler) RecordTransaction(c *gin.Context) {
	var req RecordRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.HandleError(c, utils.ErrEmptyContent)
		return
	}

	// 从会话中获取用户ID，而不是从请求参数
	userID, exists := c.Get("userID")
	if !exists {
		utils.HandleError(c, utils.ErrNotLoggedIn)
		return
	}

	transactionId, err := h.transactionService.RecordTransaction(userID.(int64), req.Type, req.Amount, req.Category, req.Note)
	if err != nil {
		utils.HandleError(c, err) // 使用统一的错误处理
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
		utils.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	transactions, err := h.transactionService.GetTransactions(userID.(int64))
	if err != nil {
		utils.HandleError(c, err) // 使用统一的错误处理
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "获取成功",
		"transactions": transactions,
	})
}
