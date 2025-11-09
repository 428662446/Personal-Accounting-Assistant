package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/models"
	"database/sql"
)

type TransactionService struct {
	masterDB *sql.DB
}

func NewTransactionService(masterDB *sql.DB) *TransactionService {
	return &TransactionService{masterDB: masterDB}
}

func (s *TransactionService) RecordTransaction(userID int64, transactionType string, amount int64, category string, note string) (int64, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return 0, err
	}
	defer userDB.Close()

	return database.RecordTransaction(userDB, transactionType, amount, category, note)
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

func (s *TransactionService) UpdateTransaction(userID int64, transactionID int64, updateType *string, updateAmount *int64, updateCategory *string, updateNote *string) error {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	return database.UpdateTransaction(userDB, transactionID, updateType, updateAmount, updateCategory, updateNote)
}
