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
// NewAmountFromString 从金额字符串创建 Amount（返回非负的分值，符号由业务层决定）
func NewAmountFromString(str string) (Amount, error) {
	cents, err := ParseToCents(str)
	if err != nil {
		return Amount{}, err
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

	/*
		原方案问题: 按第一个点拆分，将第二部分点全部去除；可能会与用户意图相违背，多个不连续小数点时直接返回错误会更好
				if strings.Count(cleaned, ".") > 1 {
			    parts := strings.SplitN(cleaned, ".", 2) // 按第一个点分割
			    cleaned = parts[0] + "." + strings.ReplaceAll(parts[1], ".", "")
			}
	*/
	// 如果存在连续的多个小数点，将连续的点合并为一个（例如 👀 -> 👁）
	for strings.Contains(cleaned, "..") {
		cleaned = strings.ReplaceAll(cleaned, "..", ".")
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
// ParseToCents 将用户输入的金额字符串解析为非负的分（int64）。
// 该函数不处理业务符号（正负由上层决定）。
//
// 四舍五入策略说明：
// - 清理输入（去除千位分隔符、空格、正负号等），验证只包含数字和最多一个小数点。
// - 如果小数位少于两位，右侧补 0；如果小数位 >= 3，则依据第三位决定是否对两位小数进行四舍五入。
// - 例如："1.234" -> 1.23（第三位为 '4'，不进位）；"1.235" -> 1.24（第三位为 '5'，进位）。
//
// 设计理由：后端以字符串接收金额以避免浮点精度问题，utils 负责把字符串转换为
// 一个明确的整数分值（单位：分），而业务层（service/handler）负责将该绝对值与
// 交易类型（income/expense）结合以决定符号并执行业务校验。
func ParseToCents(str string) (int64, error) {
	cleanedStr := CleanAmountString(str)
	if err := ValidateAmountString(cleanedStr); err != nil {
		return 0, err
	}

	// 切分整数和小数
	parts := strings.SplitN(cleanedStr, ".", 2)
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// 补充: 处理类似 ".5" 的情况
	if integerPart == "" {
		integerPart = "0"
	}

	// 补全两位小数
	if len(decimalPart) < 2 {
		decimalPart += strings.Repeat("0", 2-len(decimalPart))
	}

	// 如果存在第三位小数并且 >= '5'，则对两位小数进行四舍五入
	roundUp := false
	if len(decimalPart) >= 3 {
		if decimalPart[2] >= '5' {
			roundUp = true
		}
		decimalPart = decimalPart[:2]
	}

	resultStr := integerPart + decimalPart
	result, err := strconv.ParseInt(resultStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidParameter
	}

	if roundUp {
		result += 1
	}

	// 返回非负分值，符号由业务层决定
	if result < 0 {
		result = -result
	}
	return result, nil
}

// 2. 分转字符串
func CentsToYuanString(cents int64) string {
	sign := "+"       // 正数加号
	absCents := cents // 求绝对值分
	if cents < 0 {
		sign = "-" // 负数的负号
		absCents = -cents
	}
	yuan := absCents / 100     // 分转元整数部分
	centPart := absCents % 100 // 分转元小数部分
	return sign + fmt.Sprintf("%d.%02d", yuan, centPart)
}
