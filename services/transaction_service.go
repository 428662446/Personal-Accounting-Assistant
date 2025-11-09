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
//     因此服务不在字段中持有某个用户的 DB 连接，避免长期占用资源。
//   - masterDB 字段保留以便在需要时访问主库（例如做用户存在性检查、跨用户聚合或事务协调）。
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

	var centsPtr *int64
	if updateAmount != nil {
		// 解析字符串
		cents, err := utils.ParseToCents(*updateAmount)
		if err != nil {
			return err
		}
		if updateType == nil {
			// 业务要求：当更新金额时必须提供 type，以便决定符号
			return utils.ErrInvalidParameter
		}
		if *updateType == "expense" {
			cents = -cents
		}
		centsPtr = &cents
	}

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	return database.UpdateTransaction(userDB, transactionID, updateType, centsPtr, updateCategory, updateNote)
}
