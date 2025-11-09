package handlers

import (
	"AccountingAssistant/services"
	"AccountingAssistant/utils"
	"AccountingAssistant/web/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// "记录账单"要求结构体
type RecordRequest struct {
	Type     string `form:"type" binding:"required"`
	Amount   string `form:"amount" binding:"required"`
	Category string `form:"category" binding:"required"`
	Note     string `form:"note" binding:"required"`
}

// "更新账单"要求结构体
type UpdateTransactionRequest struct {
	Type     *string `form:"type"`     // 使用指针，nil表示不更新
	Amount   *string `form:"amount"`   // 使用指针，nil表示不更新
	Category *string `form:"category"` // 使用指针，nil表示不更新
	Note     *string `form:"note"`     // 使用指针，nil表示不更新
}

// 处理账单服务的对象
type TransactionHandler struct {
	transactionService *services.TransactionService
}

// TransactionHandler 负责将 HTTP 请求映射到 transactionService 的业务操作。
// 设计约定：handler 只做输入绑定与最小校验；关于金额的解析/四舍五入/符号应用（income/expense）由 service 层统一处理，handler 将金额字符串原样传递给 service。

// 创建"处理账单服务的对象"的方法
func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService} // 修复：明确指定字段名
}

// "记录账单"HTTP响应
func (h *TransactionHandler) RecordTransaction(c *gin.Context) {
	var req RecordRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleError(c, utils.ErrEmptyContent)
		return
	}
	// 把原始金额字符串交给 service 处理（包括解析、符号、校验）

	// 从会话中获取用户ID，而不是从请求参数
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}

	transactionId, err := h.transactionService.RecordTransaction(userID.(int64), req.Type, req.Amount, req.Category, req.Note)
	if err != nil {
		response.HandleError(c, err) // 使用统一的错误处理
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "记录成功",
		"transaction_id": transactionId,
	})
}

// "获取账单"HTTP响应
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// 从会话或JWT令牌中获取用户ID，而不是从URL参数
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	transactions, err := h.transactionService.GetTransactions(userID.(int64))
	if err != nil {
		response.HandleError(c, err) // 使用统一的错误处理
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "获取成功",
		"transactions": transactions,
	})
}

// "删除账单"HTTP响应
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		// 修复：使用自定义错误而不是直接传递底层错误
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}
	err = h.transactionService.DeleteTransaction(userID.(int64), int64(transactionID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// "更新账单"HTTP响应
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, utils.ErrNotLoggedIn)
		return
	}

	transactionIDStr := c.Param("id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}

	var req UpdateTransactionRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleError(c, utils.ErrInvalidParameter)
		return
	}

	// 如果提供了金额，必须同时提供 type，以便解析正负和规则
	if req.Amount != nil {
		if req.Type == nil {
			response.HandleError(c, utils.ErrInvalidParameter)
			return
		}
	}

	err = h.transactionService.UpdateTransaction(
		userID.(int64),
		int64(transactionID),
		req.Type,     // 可能是nil
		req.Amount,   // 可能是nil
		req.Category, // 可能是nil
		req.Note,     // 可能是nil
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}
