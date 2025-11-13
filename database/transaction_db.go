package database

import (
	"AccountingAssistant/models"
	"AccountingAssistant/utils"
	"database/sql"
	"strings"
)

// CRUD数据库操作
// 1. 记录账单（CategoryID 可以为 nil）
func RecordTransaction(userDB *sql.DB, Type string, Amount int64, CategoryID *int64, Note string) (int64, error) {

	insertSQL := "INSERT INTO transactions (type, amount, category_id, note) VALUES (?, ?, ?, ?)"
	var cid interface{}
	if CategoryID == nil {
		cid = nil
	} else {
		cid = *CategoryID
	}
	result, err := userDB.Exec(insertSQL, Type, Amount, cid, Note)
	if err != nil {
		return 0, utils.WrapError(utils.ErrInsertFailed, err)
	}
	transactionId, err := result.LastInsertId()
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return transactionId, nil
}

// 2. 获取账单（含类别名，未分类显示为 "其他"）
func GetTransaction(userDB *sql.DB) ([]models.DisplayTransaction, error) {
	querySQL := `
SELECT
	t.id, t.type, t.amount,
	COALESCE(c.name, '其他') as category_name,
	t.note, t.created_at
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
ORDER BY t.created_at DESC
`
	rows, err := userDB.Query(querySQL)
	if err != nil {
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	defer rows.Close()

	var Transactions []models.DisplayTransaction
	var cents int64
	for rows.Next() {
		var t models.DisplayTransaction
		if err := rows.Scan(&t.ID, &t.Type, &cents, &t.CategoryName, &t.Note, &t.CreatedAt); err != nil {
			return nil, utils.WrapError(utils.ErrReadFailed, err)
		}
		t.Amount = utils.CentsToYuanString(cents)
		Transactions = append(Transactions, t)
	}

	return Transactions, nil
}

// 3. 删除账单
func DeleteTransaction(userDB *sql.DB, transactionID int64) error {
	deleteSQL := "DELETE FROM transactions WHERE id = ?"
	_, err := userDB.Exec(deleteSQL, transactionID)
	if err != nil {
		return utils.WrapError(utils.ErrDeleteFailed, err)
	}
	return nil
}

// 4. 更新账单
func UpdateTransaction(userDB *sql.DB, transactionID int64, updateType *string, updateAmount *int64, updateCategoryID *int64, updateNote *string) error {

	// 构建动态SQL
	var queryParts []string
	var args []interface{} // 可以放入不同类型

	if updateType != nil {
		queryParts = append(queryParts, "type = ?")
		args = append(args, *updateType)
	}
	if updateAmount != nil {
		queryParts = append(queryParts, "amount = ?")
		args = append(args, *updateAmount)
	}
	if updateCategoryID != nil {
		queryParts = append(queryParts, "category_id = ?")
		args = append(args, *updateCategoryID)
	}
	if updateNote != nil {
		queryParts = append(queryParts, "note = ?")
		args = append(args, *updateNote)
	}

	// 如果没有要更新的字段
	if len(queryParts) == 0 {
		return nil // 或者返回一个错误，表示没有字段需要更新，但是我懒
	}

	// 添加WHERE条件
	queryParts = append(queryParts, "id = ?")
	args = append(args, transactionID)

	// 构建完整SQL ∑( 口 ||
	query := "UPDATE transactions SET " + strings.Join(queryParts[:len(queryParts)-1], ", ") + " WHERE " + queryParts[len(queryParts)-1]

	_, err := userDB.Exec(query, args...)
	if err != nil {
		return utils.WrapError(utils.ErrUpdateFailed, err)
	}
	return nil
}

// 添加: 获取单个交易的函数
func GetTransactionByID(userDB *sql.DB, transactionID int64) (*models.Transaction, error) {
	var transaction models.Transaction
	err := userDB.QueryRow(
		"SELECT id, type, amount, category_id, note, created_at FROM transactions WHERE id = ?",
		transactionID,
	).Scan(&transaction.ID, &transaction.Type, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrTransactionNotFound
		}
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return &transaction, nil
}
