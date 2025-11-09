package utils

import "fmt"

/*
使用说明
技术错误（数据库、文件、网络等）→ 用 WrapError 保留错误详情
业务错误（用户不存在、密码错误等）→ 直接返回预定义错误
*/
// 错误模型
type Error struct {
	Code    string
	Message string
	Err     error
}

// 实现Error() string方法，这样自定义错误类型就实现了Go的error接口，可以当作普通错误使用
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// 允许使用 errors.Unwrap() 函数获取嵌套的错误
// 实现Unwrap() error方法，这样可以利用Go 1.13引入的错误包装机制，使用 errors.Is 和 error.As 等函数
func (e *Error) Unwrap() error {
	return e.Err
}

// 包装错误函数
func WrapError(errType *Error, err error) error {
	// 由于 *Error 实现了 Error() string 方法，它就满足了 error 接口
	return &Error{
		Code:    errType.Code,
		Message: errType.Message,
		Err:     err,
	}
}

// 错误码分类(感谢AI分类)
const (
	// 用户相关错误 10xx
	CodeUserNotFound        = "1001"
	CodeUserAlreadyExists   = "1002"
	CodeUserEmptyCredential = "1003"
	CodeUserDBNotFound      = "1004"

	// 认证相关错误 11xx
	CodeAuthInvalidPassword = "1101"
	CodeAuthNotLoggedIn     = "1102"
	CodeAuthLoginFailed     = "1103"
	CodeSessionNotFound     = "1104"
	CodeInvalidSession      = "1105"
	CodeCreateSessionFailed = "1106"

	// 数据操作错误 12xx
	CodeDataInsertFailed = "1201"
	CodeDataQueryFailed  = "1202"
	CodeDataReadFailed   = "1203"
	CodeDataEmptyContent = "1204"
	CodeDataDeleteFailed = "1205" // 增
	CodeDataUpdateFailed = "1206" // 增

	// 系统错误 13xx
	CodeSystemDBConnFailed      = "1301"
	CodeSystemCreateDirFailed   = "1302"
	CodeSystemCreateTableFailed = "1303"
	CodeSystemEncryptFailed     = "1304"
	CodeSystemGetPasswordFailed = "1305"
	CodeFileNotFound            = "1306"

	// 业务操作错误 14xx
	CodeOperationRegisterFailed   = "1401"
	CodeOperationRecordBillFailed = "1402"
	CodeOperationGetBillFailed    = "1403" // 增
	CodeOperationDeleteBillFailed = "1404" // 增

	// 增：参数处理错误 15xx
	CodeInvalidParameter = "1501"

	// 账单相关错误 16xx
	CodeAmountInvalidFormat    = "1601"
	CodeAmountTooLarge         = "1602"
	CodeAmountZero             = "1603"
	CodeInvalidTransactionType = "1604"
	CodeTransactionNotFound    = "1605"
)

// 预定义错误(错误码 错误消息)
// 用户相关
var (
	ErrUserNotFound      = &Error{Code: CodeUserNotFound, Message: "用户不存在"}
	ErrUserAlreadyExists = &Error{Code: CodeUserAlreadyExists, Message: "用户名已存在"}
	ErrEmptyCredential   = &Error{Code: CodeUserEmptyCredential, Message: "用户名、密码不能为空"}
	ErrUserDBNotFound    = &Error{Code: CodeUserDBNotFound, Message: "用户数据库不存在"}
)

// 认证相关
var (
	ErrInvalidPassword     = &Error{Code: CodeAuthInvalidPassword, Message: "密码错误"}
	ErrNotLoggedIn         = &Error{Code: CodeAuthNotLoggedIn, Message: "未登录"}
	ErrLoginFailed         = &Error{Code: CodeAuthLoginFailed, Message: "登录失败"}
	ErrCreateSessionFailed = &Error{Code: CodeCreateSessionFailed, Message: "创建会话失败"}
	ErrSessionNotFound     = &Error{Code: CodeSessionNotFound, Message: "会话不存在"}
	ErrInvalidSession      = &Error{Code: CodeInvalidSession, Message: "无效会话"}
)

// 数据操作相关
var (
	ErrInsertFailed = &Error{Code: CodeDataInsertFailed, Message: "插入数据失败"}
	ErrQueryFailed  = &Error{Code: CodeDataQueryFailed, Message: "查询数据失败"}
	ErrReadFailed   = &Error{Code: CodeDataReadFailed, Message: "读取数据失败"}
	ErrEmptyContent = &Error{Code: CodeDataEmptyContent, Message: "内容不能为空"}
	ErrDeleteFailed = &Error{Code: CodeDataDeleteFailed, Message: "删除数据失败"} // 增
	ErrUpdateFailed = &Error{Code: CodeDataUpdateFailed, Message: "更新数据失败"} // 增
)

// 系统相关
var (
	ErrDBConnFailed      = &Error{Code: CodeSystemDBConnFailed, Message: "数据库连接失败"}
	ErrCreateDirFailed   = &Error{Code: CodeSystemCreateDirFailed, Message: "创建目录失败"}
	ErrCreateTableFailed = &Error{Code: CodeSystemCreateTableFailed, Message: "创建表失败"}
	ErrEncryptFailed     = &Error{Code: CodeSystemEncryptFailed, Message: "加密失败"}
	ErrGetPasswordFailed = &Error{Code: CodeSystemGetPasswordFailed, Message: "获取密码失败"}
	ErrFileNotFound      = &Error{Code: CodeFileNotFound, Message: "文件不存在"}
)

// 业务操作相关
var (
	ErrRegisterFailed   = &Error{Code: CodeOperationRegisterFailed, Message: "注册失败"}
	ErrRecordBillFailed = &Error{Code: CodeOperationRecordBillFailed, Message: "记录账单失败"}
	ErrGetBillFailed    = &Error{Code: CodeOperationGetBillFailed, Message: "获取账单失败"}    // 增
	ErrDeleteBillFailed = &Error{Code: CodeOperationDeleteBillFailed, Message: "删除账单失败"} // 增
)

// 参数处理相关
var (
	ErrInvalidParameter = &Error{Code: CodeInvalidParameter, Message: "参数错误"}
)

// 账单相关
var (
	ErrAmountInvalidFormat    = &Error{Code: CodeAmountInvalidFormat, Message: "金额格式错误"}
	ErrAmountTooLarge         = &Error{Code: CodeAmountTooLarge, Message: "金额过大"}
	ErrAmountZero             = &Error{Code: CodeAmountZero, Message: "金额不能为零"}
	ErrInvalidTransactionType = &Error{Code: CodeInvalidTransactionType, Message: "无效的账单类型"}
	ErrTransactionNotFound    = &Error{Code: CodeTransactionNotFound, Message: "账单不存在"}
)
