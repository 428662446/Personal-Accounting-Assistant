package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/utils"
	"database/sql"
)

// 保持与主数据库连接、方便会话管理
type StatService struct {
	masterDB *sql.DB
}

// 新建统计服务的方法
func NewStatService(masterDB *sql.DB) *StatService {
	return &StatService{masterDB: masterDB}
}

// 统计服务
func (s *StatService) GetTotalIncome(userID int64) (string, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return "", err
	}
	defer userDB.Close()

	cents, err := database.GetTotalIncome(userDB)
	if err != nil {
		return "", err
	}
	amountStr := utils.CentsToYuanString(cents)
	return amountStr, nil
}

func (s *StatService) GetTotalExpenditure(userID int64) (string, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return "", err
	}
	defer userDB.Close()

	cents, err := database.GetTotalExpenditure(userDB)
	if err != nil {
		return "", err
	}
	amountStr := utils.CentsToYuanString(cents)
	return amountStr, nil
}

func (s *StatService) GetNetIncome(userID int64) (string, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return "", err
	}
	defer userDB.Close()

	cents, err := database.GetNetIncome(userDB)
	if err != nil {
		return "", err
	}
	amountStr := utils.CentsToYuanString(cents)
	return amountStr, nil
}
