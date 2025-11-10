package utils

import (
	"strings"
	"testing"
)

func TestCleanAmountString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"普通数字", "123.45", "123.45"},
		{"带空格", " 123.45 ", "123.45"},
		{"千分位", "1,234.56", "1234.56"},
		{"多个小数点", "123..45", "123.45"},
		{"前导零", "000123.45", "000123.45"}, // 注意：不移除前导零，因为可能影响小数位 好像也不会影响但是感觉没必要 以后不去除有影响再说吧
		{"正号", "+123.45", "123.45"},
		{"负号", "-123.45", "123.45"},
		{"空字符串", "", ""},
		{"只有小数点", ".", "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanAmountString(tt.input)
			if result != tt.expected {
				t.Errorf("CleanAmountString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateAmountString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{"有效数字", "123.45", false},
		{"有效整数", "123", false},
		{"有效小数", "0.45", false},
		{"只有小数部分", ".45", false},
		{"多位小数", "123.456", false}, // 验证通过，解析时会处理
		{"负数", "-123.45", false},
		{"正数", "+123.45", false},
		{"空字符串", "", true},
		{"只有小数点", ".", true},
		{"字母字符", "abc", true},
		{"混合无效字符", "12a.34", true},
		{"多个小数点", "12.34.56", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmountString(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("ValidateAmountString(%q) error = %v, wantError = %v", tt.input, err, tt.hasError)
			}
		})
	}
}

func TestParseToCents(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int64
		expectError bool
	}{
		// 正常情况
		{"两位小数", "123.45", 12345, false},
		{"整数", "123", 12300, false},
		{"一位小数", "123.4", 12340, false},
		{"零", "0", 0, false},
		{"小数点后零", "123.00", 12300, false},
		{"只有小数部分", ".45", 45, false},
		{"只有小数点和零", "0.45", 45, false},

		// 负数
		{"负数小数", "-123.45", 12345, false},
		{"负数整数", "-123", 12300, false},
		{"负数一位小数", "-123.4", 12340, false},
		{"错误负数", "-abc", 0, true},
		{"无效负数字符", "-+-", 0, true},

		// 四舍五入测试
		{"四舍五入第三位4", "123.454", 12345, false},
		{"四舍五入第三位5", "123.455", 12346, false},
		{"四舍五入第三位6", "123.456", 12346, false},
		{"四舍五入进位", "999.999", 100000, false},

		// 边界情况
		{"大数", "999999.99", 99999999, false},
		{"很小的小数", "0.01", 1, false},

		// 错误情况
		{"空字符串", "", 0, true},
		{"无效字符", "abc", 0, true},
		{"多个小数点", "12.34.56", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseToCents(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseToCents(%q) expected error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseToCents(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseToCents(%q) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestCentsToYuanString(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"正数", 12345, "123.45"},
		{"零", 0, "0.00"},
		{"一位数", 5, "0.05"},
		{"两位数", 50, "0.50"},
		{"大数", 123456789, "1234567.89"},
		{"负数", -12345, "-123.45"},
		{"负小数", -5, "-0.05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CentsToYuanString(tt.input)
			if result != tt.expected {
				t.Errorf("CentsToYuanString(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAmountIntegration(t *testing.T) {
	// 集成测试：完整的字符串→分→字符串转换
	testCases := []string{
		"123.45",
		"0.01",
		"999.99",
		"123",
		"0.5",
		"123.456", // 测试四舍五入
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			// 字符串 → 分
			cents, err := ParseToCents(tc)
			if err != nil {
				t.Fatalf("ParseToCents(%q) failed: %v", tc, err)
			}

			// 分 → 字符串
			resultStr := CentsToYuanString(cents)

			// 验证往返一致性（注意：由于四舍五入，可能与原字符串不完全相同）
			t.Logf("Original: %s -> Cents: %d -> String: %s", tc, cents, resultStr)

			// 对于没有四舍五入的情况，应该能往返一致
			if !containsThirdDecimal(tc) {
				// 重新解析结果字符串应该得到相同的分值
				reparseCents, err := ParseToCents(resultStr)
				if err != nil {
					t.Fatalf("Reparse failed: %v", err)
				}
				if reparseCents != cents {
					t.Errorf("Round-trip inconsistency: %s -> %d -> %s -> %d",
						tc, cents, resultStr, reparseCents)
				}
			}
		})
	}
}

// 辅助函数：检查字符串是否有第三位小数

func containsThirdDecimal(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return false
	}
	return len(parts[1]) > 2
}
