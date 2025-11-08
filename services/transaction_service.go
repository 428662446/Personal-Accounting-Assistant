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

func (s *TransactionService) RecordTransaction(userID int64, transactionType string, amount float64, category string, note string) (int64, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return 0, utils.WrapError(utils.ErrDBConnFailed, err)
	}
	defer userDB.Close()

	transactionID, err := database.RecordTransaction(userDB, transactionType, amount, category, note)
	if err != nil {
		return 0, utils.WrapError(utils.ErrRecordBillFailed, err)
	}
	return transactionID, nil
}

func (s *TransactionService) GetTransactions(userID int64) ([]models.Transaction, error) {

	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return nil, utils.WrapError(utils.ErrDBConnFailed, err)
	}
	defer userDB.Close()

	transactions, err := database.GetTransaction(userDB)
	if err != nil {
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return transactions, nil
}
