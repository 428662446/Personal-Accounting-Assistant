package main

import (
	"AccountingAssistant/database"
	"AccountingAssistant/utils"
	"os"
	"testing"
)

// TestMain 会在所有测试之前运行，用于设置测试环境
func TestMain(m *testing.M) {
	// 设置测试模式
	os.Setenv("TEST_MODE", "true")

	// 运行测试
	code := m.Run()

	// 退出
	os.Exit(code)
}

// TestAmountUtils 测试金额工具函数
func TestAmountUtils(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"123.45", 12345},
		{"0.01", 1},
		{"100", 10000},
		{"99.99", 9999},
	}

	for _, tt := range tests {
		cents, err := utils.ParseToCents(tt.input)
		if err != nil {
			t.Errorf("ParseToCents(%q) failed: %v", tt.input, err)
			continue
		}
		if cents != tt.expected {
			t.Errorf("ParseToCents(%q) = %d, want %d", tt.input, cents, tt.expected)
		}
	}
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	db, err := database.InitMasterDB()
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 简单的ping测试
	err = db.Ping()
	if err != nil {
		t.Errorf("数据库连接测试失败: %v", err)
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 测试无效金额格式
	_, err := utils.ParseToCents("abc")
	if err == nil {
		t.Error("预期无效金额格式应该返回错误")
	}

	// 测试空金额
	_, err = utils.ParseToCents("")
	if err == nil {
		t.Error("预期空金额应该返回错误")
	}
}

// TestAmountRoundTrip 测试金额往返转换
func TestAmountRoundTrip(t *testing.T) {
	testCases := []string{
		"123.45",
		"0.50",
		"999.99",
		"1.00",
	}

	for _, tc := range testCases {
		cents, err := utils.ParseToCents(tc)
		if err != nil {
			t.Errorf("ParseToCents(%q) failed: %v", tc, err)
			continue
		}

		resultStr := utils.CentsToYuanString(cents)

		// 重新解析应该得到相同的分值
		reparseCents, err := utils.ParseToCents(resultStr)
		if err != nil {
			t.Errorf("Reparse failed for %s: %v", resultStr, err)
			continue
		}

		if reparseCents != cents {
			t.Errorf("Round-trip failed: %s -> %d -> %s -> %d",
				tc, cents, resultStr, reparseCents)
		}
	}
}
