package response

import (
	"AccountingAssistant/utils"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError 统一的错误处理函数 - 现在在 utils 包
func HandleError(c *gin.Context, err error) {
	// 先检查是否是已知的标准库错误
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "数据不存在",
		})
		return
	}

	// 再检查是否是我们的自定义错误
	var appErr *utils.Error
	if errors.As(err, &appErr) {
		// 处理自定义错误类型（保持你原来的逻辑）
		switch appErr.Code {
		case utils.CodeUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "用户不存在",
			})
		case utils.CodeUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "用户名已存在",
			})
		case utils.CodeUserEmptyCredential:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "用户名、密码不能为空",
			})
		case utils.CodeUserDBNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "用户数据库不存在",
			})

		// 认证相关错误 11xx
		case utils.CodeAuthInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "密码错误",
			})
		case utils.CodeAuthNotLoggedIn:
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "未登录",
			})
		case utils.CodeAuthLoginFailed:
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "登录失败",
			})
		case utils.CodeCreateSessionFailed:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "创建会话失败",
			})
		case utils.CodeSessionNotFound, utils.CodeInvalidSession:
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "会话无效或已过期",
			})

		// 数据操作错误 12xx
		case utils.CodeDataEmptyContent:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "内容不能为空",
			})
		case utils.CodeDataInsertFailed, utils.CodeDataQueryFailed, utils.CodeDataReadFailed:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "数据操作失败，请稍后重试",
			})

		// 系统错误 13xx
		case utils.CodeSystemDBConnFailed, utils.CodeSystemCreateDirFailed,
			utils.CodeSystemCreateTableFailed, utils.CodeSystemEncryptFailed:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "系统错误，请联系管理员",
			})
		case utils.CodeFileNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "文件不存在",
			})

		// 业务操作错误 14xx
		case utils.CodeOperationRegisterFailed:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "注册失败",
			})
		case utils.CodeOperationRecordBillFailed:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "记录账单失败",
			})

		// 默认情况
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "系统错误",
			})

		}
		return
	}

	// 非自定义错误，尝试判断一些已知错误
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "记录不存在",
		})
		return
	}

}
