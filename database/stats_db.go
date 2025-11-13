package database

import (
	"AccountingAssistant/models"
	"AccountingAssistant/utils"
	"database/sql"
)

// 统计相关业务
// 1. 获取总收入(coalesce意味合并)
func GetTotalIncome(userDB *sql.DB) (int64, error) {
	var result int64
	selectSQL := "SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE amount > 0 "
	err := userDB.QueryRow(selectSQL).Scan(&result)
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return result, nil
}

// 2. 获取总支出
func GetTotalExpenditure(userDB *sql.DB) (int64, error) {
	var result int64
	selectSQL := "SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE amount < 0 "
	err := userDB.QueryRow(selectSQL).Scan(&result)
	if err != nil {
		return 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return result, nil // 暂时返回负数
}

// 3. 获取净收入
func GetNetIncome(userDB *sql.DB) (int64, error) {
	totalexpenditure, err := GetTotalExpenditure(userDB)
	if err != nil {
		return 0, err
	}
	totalIncome, err := GetTotalIncome(userDB)
	if err != nil {
		return 0, err
	}
	return totalIncome + totalexpenditure, nil
}
// 4. 按月统计(返回当前月份的总收入/总支出/净收入，单位：分)
func GetMonthlyStats(userDB *sql.DB) (int64, int64, int64, error) {
	querySQL := `
SELECT
	COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END),0) AS total_income,
	COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END),0) AS total_expense,
	COALESCE(SUM(amount),0) AS net_income
FROM transactions
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m','now','localtime')
`
	var total_income int64
	var total_expense int64
	var net_income int64
	err := userDB.QueryRow(querySQL).Scan(&total_income, &total_expense, &net_income)
	if err != nil {
		return 0, 0, 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return total_income, total_expense, net_income, nil
}

// 5. 按周统计（当前周）
func GetWeeklyStats(userDB *sql.DB) (int64, int64, int64, error) {
	// 使用 ISO 周或 %W 需与前端约定；这里采用 SQLite 的 %Y-%W
	querySQL := `
SELECT
	COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END),0) AS total_income,
	COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END),0) AS total_expense,
	COALESCE(SUM(amount),0) AS net_income
FROM transactions
WHERE strftime('%Y-%W', created_at) = strftime('%Y-%W','now','localtime')
`
	var total_income int64
	var total_expense int64
	var net_income int64
	err := userDB.QueryRow(querySQL).Scan(&total_income, &total_expense, &net_income)
	if err != nil {
		return 0, 0, 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return total_income, total_expense, net_income, nil
}

// 按日统计（当天）
func GetDailyStats(userDB *sql.DB) (int64, int64, int64, error) {
	querySQL := `
SELECT
	COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END),0) AS total_income,
	COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END),0) AS total_expense,
	COALESCE(SUM(amount),0) AS net_income
FROM transactions
WHERE date(created_at) = date('now','localtime')
`
	var total_income int64
	var total_expense int64
	var net_income int64
	err := userDB.QueryRow(querySQL).Scan(&total_income, &total_expense, &net_income)
	if err != nil {
		return 0, 0, 0, utils.WrapError(utils.ErrQueryFailed, err)
	}
	return total_income, total_expense, net_income, nil
}

// 金额范围统计（按区间分组，返回每组的计数与总额）
func GetRangeAmountStats(userDB *sql.DB) ([]models.RangeAmountStat, error) {
	// 注意：amount 单位为“分”，阈值 10000 表示 100.00 元
	querySQL := `
SELECT amount_range, COUNT(*) AS transaction_count, COALESCE(SUM(amount),0) AS total_amount
FROM (
  SELECT
	CASE
	  WHEN amount >= 10000 THEN '大额收入(>=100)'
	  WHEN amount > 0 THEN '小额收入'
	  WHEN amount > -10000 THEN '小额支出'
	  ELSE '大额支出'
	END AS amount_range,
	amount
  FROM transactions
) t
GROUP BY amount_range
ORDER BY total_amount DESC
`
	rows, err := userDB.Query(querySQL)
	if err != nil {
		return nil, utils.WrapError(utils.ErrQueryFailed, err)
	}
	defer rows.Close()
	var rangeStats []models.RangeAmountStat
	for rows.Next() {
		var name string
		var transactionCount int
		var amount int64
		if err := rows.Scan(&name, &transactionCount, &amount); err != nil {
			return nil, utils.WrapError(utils.ErrReadFailed, err)
		}
		amountStr := utils.CentsToYuanString(amount)
		rs := models.RangeAmountStat{
			Name:             name,
			TransactionCount: transactionCount,
			Amount:           amount,
			AmountStr:        amountStr,
		}
		rangeStats = append(rangeStats, rs)
	}
	return rangeStats, nil
}
