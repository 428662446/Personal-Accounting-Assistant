package models

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Transaction struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`   // "income" 或 "expense"
	Amount     int64  `json:"amount"` // 已修改金额存储类型
	CategoryID int64  `json:"category"`
	Note       string `json:"note"`
	CreatedAt  string `json:"created_at"`
}

type DisplayTransaction struct {
	ID           int64  `json:"id"`
	Type         string `json:"type"` // "income" 或 "expense"
	Amount       string `json:"amount"`
	CategoryName string `json:"category_name"`
	Note         string `json:"note"`
	CreatedAt    string `json:"created_at"`
}

type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
