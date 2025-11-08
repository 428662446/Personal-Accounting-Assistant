package database

import (
	"AccountingAssistant/models"
	"AccountingAssistant/utils"
	"database/sql"
)

func RecordTransaction(userDB *sql.DB, Type string, Amount float64, Category string, Note string) (int64, error) {

	insertSQL := "INSERT INTO transactions (type, amount, category, note) VALUES (?, ?, ?, ?)"
	result, err := userDB.Exec(insertSQL, Type, Amount, Category, Note)
	if err != nil {
		return 0, utils.WrapError(utils.ErrInsertFailed, err)
	}
	transactionId, err := result.LastInsertId()
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return transactionId, nil
}

func GetTransaction(userDB *sql.DB) ([]models.Transaction, error) {
	rows, err := userDB.Query("SELECT id, type, amount, category, note, created_at FROM transactions ORDER BY created_at DESC")
	if err != nil {
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	defer rows.Close()

	var Transaction []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Type, &t.Amount, &t.Category, &t.Note, &t.CreatedAt); err != nil {
			return nil, utils.WrapError(utils.ErrReadFailed, err)
		}
		Transaction = append(Transaction, t)
	}

	return Transaction, nil
}
