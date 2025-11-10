package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/models"
	"AccountingAssistant/utils"
	"database/sql"
)

// TransactionService 提供与账单（transactions）相关的业务操作。
// 设计说明：
//   - masterDB 保存主数据库的连接（用于用户表、会话或全局元数据）。
//   - 每个用户的数据保存在独立的 per-user sqlite 文件中（通过 database.GetUserDB 按需打开），
type TransactionService struct {
	masterDB *sql.DB
}

// NewTransactionService 构造 TransactionService 并注入主数据库连接。
// 注意：传入的 masterDB 用于访问用户表/会话等主库资源；具体的 per-user DB 仍由
// database.GetUserDB 在每次请求时按需打开与关闭。
func NewTransactionService(masterDB *sql.DB) *TransactionService {
	return &TransactionService{masterDB: masterDB}
}

// "记录账单"服务
func (s *TransactionService) RecordTransaction(userID int64, transactionType string, amountStr string, category string, note string) (int64, error) {
	// 解析金额字符串（utils负责清洗和四舍五入到分），业务层负责根据 type 应用符号
	cents, err := utils.ParseToCents(amountStr)
	if err != nil {
		return 0, err
	}
	if transactionType == "expense" {
		cents = -cents
	}

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return 0, err
	}
	defer userDB.Close()

	return database.RecordTransaction(userDB, transactionType, cents, category, note)
}

// "获取账单"服务
func (s *TransactionService) GetTransactions(userID int64) ([]models.DisplayTransaction, error) {

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return nil, err
	}
	defer userDB.Close()

	return database.GetTransaction(userDB)
}

// "删除账单"服务
func (s *TransactionService) DeleteTransaction(userID int64, transactionID int64) error {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	return database.DeleteTransaction(userDB, transactionID)
}

// "更新账单"服务
func (s *TransactionService) UpdateTransaction(userID int64, transactionID int64, updateType *string, updateAmount *string, updateCategory *string, updateNote *string) error {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	// 获取原交易信息(如果没有更新Type，就会获取原来的)
	existingTransaction, err := database.GetTransactionByID(userDB, transactionID)
	if err != nil {
		return err
	}

	var centsPtr *int64
	var finalType *string

	// 确定最终的类型
	// (其实就一个if update type 则 字符串转分 时根据updateType;else 根据 原type;但是这样AI说这样写更清晰)
	finalTransactionType := existingTransaction.Type
	if updateType != nil {
		finalTransactionType = *updateType
		finalType = updateType
	}

	// 验证类型有效性
	if finalTransactionType != "income" && finalTransactionType != "expense" {
		return utils.ErrInvalidTransactionType
	}

	// 处理金额更新
	if updateAmount != nil {
		// 解析新金额
		cents, err := utils.ParseToCents(*updateAmount)
		if err != nil {
			return err
		}

		// 根据最终类型应用符号
		if finalTransactionType == "expense" {
			cents = -cents
		}

		// 业务校验
		if cents == 0 {
			return utils.ErrInvalidTransactionType
		}

		centsPtr = &cents
	} else if updateType != nil {
		// 只更新类型：调整原金额的符号以匹配新类型
		var adjustedAmount int64
		if finalTransactionType == "income" {
			// 确保金额为正
			adjustedAmount = abs(existingTransaction.Amount)
		} else {
			// 确保金额为负
			adjustedAmount = -abs(existingTransaction.Amount)
		}

		// 如果符号确实改变了，才更新金额
		if adjustedAmount != existingTransaction.Amount {
			centsPtr = &adjustedAmount
		}
	}

	return database.UpdateTransaction(userDB, transactionID, finalType, centsPtr, updateCategory, updateNote)
}

// 辅助函数：获取绝对值
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
