package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RecordTransaction(masterDB *sql.DB, Name string, Type string, Amount float64, Category string, Note string) error {
	var userId int
	selectSQL := "SELECT id FROM users WHERE username = ?"
	err := masterDB.QueryRow(selectSQL, Name).Scan(&userId)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}
	userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", userId))
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		return fmt.Errorf("用户个人数据库不存在")
	}

	userDB, err := getUserDB(int64(userId))
	if err != nil {
		return fmt.Errorf("连接用户数据库失败")
	}
	defer userDB.Close()
	insertSQL := "INSERT INTO transactions (userid, type, amount, category, note) VALUES (?, ?, ?, ?, ?)"
	_, err = userDB.Exec(insertSQL, userId, Type, Amount, Category, Note)
	if err != nil {
		return fmt.Errorf("插入交易失败: %v", err)
	}
	return nil
}

func ReadTransaction(masterDB *sql.DB, Name string) error {
	var userId int
	selectSQL := "SELECT id FROM users WHERE username = ?"
	err := masterDB.QueryRow(selectSQL, Name).Scan(&userId)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}
	userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", userId))
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		return fmt.Errorf("用户个人数据库不存在")
	}

	userDB, err := getUserDB(int64(userId))
	if err != nil {
		return fmt.Errorf("连接用户数据库失败")
	}
	defer userDB.Close()

	rows, err := userDB.Query("SELECT id, type, amount, category, note, created_at FROM transactions ORDER BY created_at DESC")
	if err != nil {
		return fmt.Errorf("查询交易失败: %v", err)
	}
	defer rows.Close()

	fmt.Printf("用户 %s 的账单:\n", Name)
	for rows.Next() {
		var id int
		var t string
		var amount float64
		var category string
		var note string
		var createdAt string
		if err := rows.Scan(&id, &t, &amount, &category, &note, &createdAt); err != nil {
			return fmt.Errorf("读取行失败: %v", err)
		}
		fmt.Printf("id=%d type=%s amount=%v category=%s note=%s at=%s\n", id, t, amount, category, note, createdAt)
	}
	return nil
}
