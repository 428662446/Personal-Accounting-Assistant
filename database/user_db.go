package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func getUserDB(userId int64) (*sql.DB, error) {
	userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", userId))
	db, err := sql.Open("sqlite3", userDBPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func registerUser(masterDB *sql.DB, username, password string) (int64, error) {

	var existingID int64
	err := masterDB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&existingID)
	if err == nil {
		// 用户存在，检查个人数据库是否存在
		userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", existingID))
		if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
			// 个人数据库不存在
			return 0, fmt.Errorf("用户存在，个人数据库不存在: %w", err)
		}
		return 0, fmt.Errorf("用户名已存在")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	tx, err := masterDB.Begin() // 开始事务
	if err != nil {
		return 0, err
	}
	// 1.1 改进：确保在函数返回前处理事务
	defer func() {
		if err != nil {
			tx.Rollback() // 出错时回滚
		}
	}()

	insertSQL := "INSERT INTO users (username, password) VALUES (?, ?)"
	result, err := tx.Exec(insertSQL, username, password)
	if err != nil {
		tx.Rollback() // 出错时回滚
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// 创建个人数据库
	if err := createUserDatabase(userID); err != nil {
		return 0, err // defer 会处理回滚
	}

	return userID, tx.Commit() // 提交事务
}
func RegisterUser(masterDB *sql.DB, username, password string) (int64, error) {

	var existingID int64
	err := masterDB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&existingID)
	if err == nil {
		// 用户存在，检查个人数据库是否存在
		userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", existingID))
		if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
			// 个人数据库不存在
			return 0, fmt.Errorf("用户存在，个人数据库不存在: %w", err)
		}
		return 0, fmt.Errorf("用户名已存在")
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	tx, err := masterDB.Begin() // 开始事务
	if err != nil {
		return 0, err
	}
	// 1.1 改进：确保在函数返回前处理事务
	defer func() {
		if err != nil {
			tx.Rollback() // 出错时回滚
		}
	}()

	insertSQL := "INSERT INTO users (username, password) VALUES (?, ?)"
	result, err := tx.Exec(insertSQL, username, password)
	if err != nil {
		tx.Rollback() // 出错时回滚
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// 创建个人数据库
	if err := createUserDatabase(userID); err != nil {
		return 0, err // defer 会处理回滚
	}

	return userID, tx.Commit() // 提交事务
}

func LoginUser(masterDB *sql.DB, username string, password string) (int64, error) {
	var pd string
	var id int64
	selectSQL := "SELECT id, password FROM users WHERE username = ?"
	err := masterDB.QueryRow(selectSQL, username).Scan(&id, &pd)
	if err != nil {
		if err == sql.ErrNoRows { // 补
			return 0, fmt.Errorf("用户不存在: %v", err)
		}
		return 0, fmt.Errorf("用户登录失败: %v", err)
	}
	if pd != password {
		return 0, fmt.Errorf("密码错误")
	}
	fmt.Printf("%s 成功登录\n", username)

	return id, nil
}
