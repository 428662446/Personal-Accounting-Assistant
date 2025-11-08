package services

import (
	"AccountingAssistant/utils"
	"database/sql"
	"fmt"
	"time"
)

const (
	SessionTimeout = 24 * time.Hour
)

// 使用数据库存储会话
type DBSessionManager struct {
	masterDB *sql.DB
}

func NewDBSessionManager(masterDB *sql.DB) *DBSessionManager {
	return &DBSessionManager{masterDB: masterDB}
}

// 创建会话
func (sm *DBSessionManager) CreateSession(userID int64, username string) (string, error) {
	sessionID := generateSessionID()
	expires := time.Now().Add(SessionTimeout)

	insertSQL := "INSERT INTO sessions (session_id, user_id, username, expires) VALUES(?, ?, ?, ?)"
	_, err := sm.masterDB.Exec(insertSQL, sessionID, userID, username, expires)
	if err != nil {
		return "", utils.WrapError(utils.ErrCreateSessionFailed, err)
	}
	return sessionID, nil
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
	_, err := sm.masterDB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		// 可以记录日志，但不需要返回错误
		fmt.Printf("删除会话失败: %v\n", err)
	}
}

// 清理过期会话（deepseek建议）
func (sm *DBSessionManager) CleanupExpiredSessions() {
	_, err := sm.masterDB.Exec("DELETE FROM sessions WHERE expires < ?", time.Now())
	if err != nil {
		fmt.Printf("清理过期会话失败: %v\n", err)
	}
}

// 简单的sessionID生成
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
