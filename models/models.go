package models

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Transaction struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`   // "income" 或 "expense"
	Amount    int64  `json:"amount"` // 已修改金额存储类型
	Category  string `json:"category"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

type DisplayTransaction struct {
	ID        int    `json:"id"`
	Type      string `json:"type"` // "income" 或 "expense"
	Amount    string `json:"amount"`
	Category  string `json:"category"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}
