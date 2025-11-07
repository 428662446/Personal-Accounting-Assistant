package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/models"
	"database/sql"
	"fmt"
)

type TransactionService struct {
	masterDB *sql.DB
}

func NewTransactionService(masterDB *sql.DB) *TransactionService {
	return &TransactionService{masterDB: masterDB}
}

func (s *TransactionService) RecordTransaction(username string, transactionType string, amount float64, category string, note string) (int64, error) {
	userID, err := database.GetUserIDByUsername(s.masterDB, username)
	if err != nil {
		return 0, fmt.Errorf("用户不存在")
	}
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return 0, fmt.Errorf("连接用户数据库失败")
	}
	defer userDB.Close()

	transactionID, err := database.RecordTransaction(userDB, transactionType, amount, category, note)
	if err != nil {
		return 0, fmt.Errorf("记录交易失败: %v", err)
	}
	return transactionID, nil
}

func (s *TransactionService) GetTransactions(userID int64) ([]models.Transaction, error) {

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return nil, fmt.Errorf("连接用户数据库失败")
	}
	defer userDB.Close()

	transactions, err := database.GetTransaction(userDB)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
