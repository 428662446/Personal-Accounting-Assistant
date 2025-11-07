package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitMasterDB() (*sql.DB, error) {
	// 确保主数据文件路径存在
	if err := os.MkdirAll("database_files", 0755); err != nil { // 0755 是Linux文件权限：用户可读写执行，组和其他用户可读执行
		return nil, fmt.Errorf("创建 database 目录失败: %w", err) // %w 是错误包装，保留原始错误信息
	}
	dbPath := filepath.Join("database_files", "master.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建用户表
	createTableSQL := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}
	fmt.Println("主数据库初始化完成")

	return db, nil
}

func createUserDatabase(userID int64) error {
	// 确保用户数据文件路径存在
	usersDir := filepath.Join("database_files", "usersdata")
	if err := os.MkdirAll(usersDir, 0755); err != nil {
		return fmt.Errorf("创建 usersdata 目录失败: %w", err)
	}
	path := filepath.Join(usersDir, fmt.Sprintf("user_%d.db", userID))
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	// 创建记账表
	createTableSQL := `
CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL,
	amount REAL NOT NULL,
	category TEXT NOT NULL,
	note TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	fmt.Printf("用户 %d 的数据库创建成功\n", userID)
	return nil
}
