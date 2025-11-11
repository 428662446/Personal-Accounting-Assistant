package database

import (
	"AccountingAssistant/models"
	"AccountingAssistant/utils"
	"database/sql"
)

// CreateCategory 在用户数据库中新增一个类别，返回插入的 ID
func CreateCategory(userDB *sql.DB, name string) (int64, error) {
	insertSQL := "INSERT INTO categories (name) VALUES (?)"
	result, err := userDB.Exec(insertSQL, name)
	if err != nil {
		return 0, utils.WrapError(utils.ErrInsertFailed, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return id, nil
}

// GetCategories 返回用户数据库中所有类别
func GetCategories(userDB *sql.DB) ([]models.Category, error) {
	querySQL := "SELECT id, name, created_at FROM categories ORDER BY created_at DESC" // 修复SQL
	rows, err := userDB.Query(querySQL)
	if err != nil {
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, utils.WrapError(utils.ErrReadFailed, err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// DeleteCategory 删除类别（在事务中把相关交易的 category_id 置为 NULL）
func DeleteCategory(userDB *sql.DB, categoryID int64) error {
	// 开始事务
	tx, err := userDB.Begin()
	if err != nil {
		return utils.WrapError(utils.ErrDBConnFailed, err)
	}
	defer func() {
		if p := recover(); p != nil {
			// 回滚
			tx.Rollback()
			// 如果程序发生了panic（严重错误），就回滚事务，然后重新抛出这个panic让上层处理
			panic(p)
		}
	}()

	// 1. 更新相关交易的category_id
	updateSQL := "UPDATE transactions SET category_id = 0 WHERE category_id = ?"
	_, err = tx.Exec(updateSQL, categoryID)
	if err != nil {
		return utils.WrapError(utils.ErrUpdateFailed, err)
	}

	// 2. 删除类别
	deleteSQL := "DELETE FROM categories WHERE id = ?"
	_, err = tx.Exec(deleteSQL, categoryID)
	if err != nil {
		return utils.WrapError(utils.ErrDeleteFailed, err)
	}
	//
	return tx.Commit()
}

// UpdateCategory 更新类别名称
func UpdateCategory(userDB *sql.DB, categoryID int64, name string) error {
	updateSQL := "UPDATE categories SET name = ? WHERE id = ?"
	_, err := userDB.Exec(updateSQL, name, categoryID)
	if err != nil {
		return utils.WrapError(utils.ErrUpdateFailed, err)
	}
	return nil
}

// GetCategoryByID 返回 category 对象，若不存在返回 nil, nil
func GetCategoryByID(userDB *sql.DB, categoryID int64) (*models.Category, error) {
	var c models.Category
	err := userDB.QueryRow("SELECT id, name, created_at FROM categories WHERE id = ?", categoryID).Scan(&c.ID, &c.Name, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return &c, nil
}

func GetCategoryIdByName(userDB *sql.DB, name string) (int64, error) {
	var categoryID int64
	querySQL := "SELECT id FROM categories WHERE name = ?"
	err := userDB.QueryRow(querySQL, name).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // 不存在返回0
		}
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return categoryID, nil
}
