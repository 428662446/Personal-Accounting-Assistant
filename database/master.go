package database

import (
	"AccountingAssistant/utils"
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	UploadDir = "database_files"
)

func InitMasterDB() (*sql.DB, error) {
	// 确保主数据文件路径存在
	if err := os.MkdirAll(UploadDir, 0755); err != nil { // 0755 是Linux文件权限：用户可读写执行，组和其他用户可读执行
		return nil, utils.WrapError(utils.ErrCreateDirFailed, err)
	}
	dbPath := filepath.Join(UploadDir, "master.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, utils.WrapError(utils.ErrDBConnFailed, err)
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
		return nil, utils.WrapError(utils.ErrCreateTableFailed, err)
	}

	// 新增：创建会话表
	createSessionTableSQL := `
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    expires DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);`
	_, err = db.Exec(createSessionTableSQL)
	if err != nil {
		db.Close()
		return nil, utils.WrapError(utils.ErrCreateTableFailed, err)
	}

	return db, nil
}
