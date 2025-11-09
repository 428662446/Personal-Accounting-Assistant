package utils

/*
在utils包中添加金额处理文件amount.go：
0.1  Amount 金额类型：
内部存储storedValue为分（整数）、safeValue用于运算
0.2 清理输入cleanAmountString()
0.3 验证金额合法validateAmountFormat()
0.4 元转分 parseToCents()
0.5 分转元 CentsToparse()
0.6 从字符串创建金额 NewAmountFromString()
0.7 从分创建金额 NewAmountFromCents()
*/
import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// 定义金额类型(私有数据成员)
type Amount struct {
	storedValue int64    // 存储值，单位为“分”
	safeValue   *big.Int // 运算值，防溢出
}

// 金额类型创建方法(构造函数)
// 1. 分转金额
func NewAmountFromCents(cents int64) Amount {
	return Amount{
		storedValue: cents,
		safeValue:   big.NewInt(cents),
	}
}

// 2. 字符串转金额
func NewAmountFromString(str string, transactionType string) (Amount, error) {
	cents, err := ParseToCents(str, transactionType)
	if err != nil {
		return Amount{}, err // ❌!!!!暂时未完善错误处理
	}
	return Amount{
		storedValue: cents,
		safeValue:   big.NewInt(cents),
	}, nil
}

// 金额类型的方法(成员函数)
func (a Amount) ToYuanString() string {
	return CentsToYuanString(a.storedValue)
}
func (a Amount) ToCents() int64 {
	return a.storedValue
}
func (a Amount) SafeValue() *big.Int {
	return a.safeValue
}

// 字符清理
func CleanAmountString(input string) string {

	cleaned := strings.TrimSpace(input) // 用于去除字符串首尾的空白字符（如空格、制表符、换行符等）
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	// 处理多个小数点
	if strings.Count(cleaned, ".") > 1 {
		parts := strings.SplitN(cleaned, ".", 2) // The count determines the number of substrings to return
		cleaned = parts[0] + strings.ReplaceAll(parts[1], ".", "")
	}
	return cleaned
}

// 字符合法检验
func ValidateAmountString(str string) error {
	cleaned := CleanAmountString(str)

	// 基础检查
	if cleaned == "" || cleaned == "." {
		return ErrEmptyContent
	}

	// 具体的字符检查：只能包含数字和最多一个小数点
	hasDot := false
	for _, ch := range cleaned {
		if ch == '.' {
			if hasDot {
				return ErrInvalidParameter
			}
			hasDot = true
			continue
		}
		if ch < '0' || ch > '9' {
			return ErrInvalidParameter
		}
	}
	return nil
}

// 字符串转分工具
// 1. 字符串转分
func ParseToCents(str string, transactionType string) (int64, error) {
	cleanedStr := CleanAmountString(str)
	err := ValidateAmountString(cleanedStr)
	if err != nil {
		return 0, err // 在ValidateAmountString已经处理
	}
	// 区分正负数
	var isNegative bool
	switch transactionType {
	case "income":
		isNegative = false // 收入强制为正
	case "expense":
		isNegative = true // 支出强制为负
	}
	// 切分整数和小数
	parts := strings.SplitN(cleanedStr, ".", 2)
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 { // 如果有小数
		decimalPart = parts[1]
	}
	// 处理类似 ".5" 的情况
	if integerPart == "" {
		integerPart = "0"
	}
	// 补全两位小数
	if len(decimalPart) < 2 {
		decimalPart += strings.Repeat("0", 2-len(decimalPart)) // Repeat returns a new string consisting of count copies of the string s
	} else if len(decimalPart) > 2 {
		decimalPart = decimalPart[:2] // 直接截断
	}
	resultStr := integerPart + decimalPart
	result, err := strconv.ParseInt(resultStr, 10, 64) // ParseInt解释给定基数（0,2到36）和位大小（0到64）的字符串s，并返回相应的值
	if err != nil {
		return 0, ErrInvalidParameter
	}
	if isNegative {
		result = -result
	}

	return result, nil
}

// 2. 分转字符串
func CentsToYuanString(cents int64) string {
	yuan := cents / 100
	centPart := cents % 100
	if centPart < 0 {
		centPart = -centPart
	}
	return fmt.Sprintf("%d.%02d", yuan, centPart) // 不用再合并字符串再返回了
}
