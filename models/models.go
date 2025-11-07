package models

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Transaction struct {
	ID        int     `json:"id"`
	UserID    int     `json:"userid"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	Note      string  `json:"note"`
	CreatedAt string  `json:"created_at"`
}
