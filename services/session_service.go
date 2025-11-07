package services

import (
	"database/sql"
	"fmt"
	"time"
)

type Session struct { // 这是一个会话
	UserID   int64
	Username string
	Expires  time.Time // 这是过期时间
}

// 使用数据库存储会话
type DBSessionManager struct {
	masterDB *sql.DB
}

func NewDBSessionManager(masterDB *sql.DB) *DBSessionManager {
	return &DBSessionManager{masterDB: masterDB}
}

// 创建会话
func (sm *DBSessionManager) CreateSession(userID int64, username string) string {
	sessionID := generateSessionID()
	expires := time.Now().Add(24 * time.Hour)

	insertSQL := "INSERT INTO sessions (session_id, user_id, username, expires) VALUES(?, ?, ?, ?)"
	_, err := sm.masterDB.Exec(insertSQL, sessionID, userID, username, expires)
	if err != nil {
		return ""
	}
	return sessionID
}

// 验证会话
func (sm *DBSessionManager) ValidateSession(sessionID string) (int64, string, bool) {
	var userID int64
	var username string
	var expires time.Time
	selectSQL := "SELECT user_id, username, expires FROM sessions WHERE session_id = ? "
	err := sm.masterDB.QueryRow(selectSQL, sessionID).Scan(&userID, &username, &expires)
	if err != nil {
		return 0, "", false
	}
	if time.Now().After(expires) {
		return 0, "", false
	}
	return userID, username, true
}

// 删除会话
func (sm *DBSessionManager) DeleteSession(sessionID string) {
	sm.masterDB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
}

// 清理过期会话（deepseek建议）
func (sm *DBSessionManager) CleanupExpiredSessions() {
	sm.masterDB.Exec("DELETE FROM sessions WHERE expires < ?", time.Now())
}

// 简单的sessionID生成
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
