package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/models"
	"AccountingAssistant/utils"
	"database/sql"
)

type TransactionService struct {
	masterDB *sql.DB
}

func NewTransactionService(masterDB *sql.DB) *TransactionService {
	return &TransactionService{masterDB: masterDB}
}

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

func (s *TransactionService) GetTransactions(userID int64) ([]models.DisplayTransaction, error) {

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return nil, err
	}
	defer userDB.Close()

	return database.GetTransaction(userDB)
}

func (s *TransactionService) DeleteTransaction(userID int64, transactionID int64) error {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	return database.DeleteTransaction(userDB, transactionID)
}

func (s *TransactionService) UpdateTransaction(userID int64, transactionID int64, updateType *string, updateAmount *string, updateCategory *string, updateNote *string) error {

	var centsPtr *int64
	if updateAmount != nil {
		// 解析字符串
		parsed, err := utils.ParseToCents(*updateAmount)
		if err != nil {
			return err
		}
		if updateType == nil {
			// 业务要求：当更新金额时必须提供 type，以便决定符号
			return utils.ErrInvalidParameter
		}
		if *updateType == "expense" {
			parsed = -parsed
		}
		centsPtr = &parsed
	}

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	return database.UpdateTransaction(userDB, transactionID, updateType, centsPtr, updateCategory, updateNote)
}
