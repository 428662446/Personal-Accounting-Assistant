package database

import (
	"AccountingAssistant/utils"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	UsersDataDir = "database_files/usersdata"
)

func RegisterUser(masterDB *sql.DB, username, password string) (int64, error) {

	_, err := GetUserIDByUsername(masterDB, username)
	if err == nil {
		return 0, utils.ErrUserAlreadyExists
	}
	// 哈希密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, utils.WrapError(utils.ErrEncryptFailed, err)
	}

	tx, err := masterDB.Begin() // 开始事务
	if err != nil {
		return 0, utils.WrapError(utils.ErrDBConnFailed, err)
	}
	// 改进：确保在函数返回前处理事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() // 出错时回滚
			panic(p)      // 重新抛出 panic
		}
	}()
	// 插入用户表单
	insertSQL := "INSERT INTO users (username, password) VALUES (?, ?)"
	result, err := tx.Exec(insertSQL, username, hashedPassword)
	if err != nil {
		tx.Rollback() // 出错时回滚
		return 0, utils.WrapError(utils.ErrInsertFailed, err)
	}
	// 获取新用户ID
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	// 创建个人数据库
	if err := createUserDatabase(userID); err != nil {
		tx.Rollback()
		return 0, err // 错误已经在createUserDatabase中包装过了
	}

	return userID, tx.Commit() // 提交事务
}

func LoginUser(masterDB *sql.DB, username string, password string) (int64, error) {
	userID, err := GetUserIDByUsername(masterDB, username)
	if err != nil {
		return 0, err // 错误已经在GetUserIDByUsername中处理过了
	}
	// 用户存在，检查个人数据库是否存在
	_, err = EnsureUserDatabase(userID)
	if err != nil {
		// 个人数据库不存在
		return 0, err // // 错误已经在EnsureUserDatabase中处理过了

	}

	// 验证密码
	var hashedPassword string
	selectSQL := "SELECT password FROM users WHERE id = ?"
	err = masterDB.QueryRow(selectSQL, userID).Scan(&hashedPassword)
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}

	if err := utils.VerifyPassword(hashedPassword, password); err != nil {
		return 0, utils.ErrInvalidPassword
	}

	return userID, nil
}
func createUserDatabase(userID int64) error {
	// 确保用户数据文件路径存在
	path, err := EnsureUserDatabase(userID)
	if err != nil {
		return utils.WrapError(utils.ErrCreateDirFailed, err)
	}
	// 打开用户数据库
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return utils.WrapError(utils.ErrDBConnFailed, err)
	}
	defer db.Close()

	// 创建记账表
	createTableSQL := `
CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL,
	amount INTEGER NOT NULL,
	category TEXT NOT NULL,
	note TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);` // 已修改表单金额类型

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return utils.WrapError(utils.ErrCreateTableFailed, err)
	}
	return nil
}

func GetUserDB(userId int64) (*sql.DB, error) {
	userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", userId))
	db, err := sql.Open("sqlite3", userDBPath)
	if err != nil {
		return nil, utils.WrapError(utils.ErrDBConnFailed, err)
	}
	return db, nil
}
func GetUserIDByUsername(masterDB *sql.DB, username string) (int64, error) {
	var userID int64
	err := masterDB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.ErrUserNotFound // 用户不存在是业务逻辑的一部分，不需要底层错误信息
		}
		return 0, utils.WrapError(utils.ErrQueryFailed, err) // 其他错误
	}
	return userID, nil
}
func EnsureUserDatabase(userID int64) (string, error) {
	if err := os.MkdirAll(UsersDataDir, 0755); err != nil {
		return "", utils.WrapError(utils.ErrCreateDirFailed, err)
	}
	userDBPath := filepath.Join(UsersDataDir, fmt.Sprintf("user_%d.db", userID))
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		return "", utils.ErrUserDBNotFound // 文件不存在是业务逻辑
	}
	return userDBPath, nil
}
