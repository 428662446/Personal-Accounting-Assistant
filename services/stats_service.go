package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/models"
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

// 按月统计
func (s *StatService) GetMonthlyStats(userID int64) (string, string, string, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return "", "", "", err
	}
	defer userDB.Close()

	total_income_cents, total_expense_cents, net_income_cents, err := database.GetMonthlyStats(userDB)
	if err != nil {
		return "", "", "", err
	}
	total_income_string := utils.CentsToYuanString(total_income_cents)
	total_expense_string := utils.CentsToYuanString(total_expense_cents)
	net_income_string := utils.CentsToYuanString(net_income_cents)

	return total_income_string, total_expense_string, net_income_string, nil
}

// 按周统计
func (s *StatService) GetWeeklyStats(userID int64) (string, string, string, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return "", "", "", err
	}
	defer userDB.Close()

	total_income_cents, total_expense_cents, net_income_cents, err := database.GetWeeklyStats(userDB)
	if err != nil {
		return "", "", "", err
	}
	total_income_string := utils.CentsToYuanString(total_income_cents)
	total_expense_string := utils.CentsToYuanString(total_expense_cents)
	net_income_string := utils.CentsToYuanString(net_income_cents)

	return total_income_string, total_expense_string, net_income_string, nil
}

// 按日统计
func (s *StatService) GetDailyStats(userID int64) (string, string, string, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return "", "", "", err
	}
	defer userDB.Close()

	total_income_cents, total_expense_cents, net_income_cents, err := database.GetDailyStats(userDB)
	if err != nil {
		return "", "", "", err
	}
	total_income_string := utils.CentsToYuanString(total_income_cents)
	total_expense_string := utils.CentsToYuanString(total_expense_cents)
	net_income_string := utils.CentsToYuanString(net_income_cents)

	return total_income_string, total_expense_string, net_income_string, nil
}

// 金额范围统计
func (s *StatService) GetRangeAmountStats(userID int64) ([]models.RangeAmountStat, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return nil, err
	}
	defer userDB.Close()

	rangeAmountStats, err := database.GetRangeAmountStats(userDB)
	if err != nil {
		return nil, err
	}
	return rangeAmountStats, nil
}
