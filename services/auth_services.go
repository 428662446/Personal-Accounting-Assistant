package services

import (
	"AccountingAssistant/database"
	"database/sql"
	"fmt"
)

type UserService struct {
	masterDB *sql.DB
}

func NewUserService(masterDB *sql.DB) *UserService {
	return &UserService{masterDB: masterDB}
}
func (s *UserService) Register(username, password string) (int64, error) {
	userId, err := database.RegisterUser(s.masterDB, username, password)
	if err != nil {
		return 0, fmt.Errorf("注册失败: %v", err)
	}
	return userId, nil
}

func (s *UserService) Login(username, password string) (int64, error) {
	id, err := database.LoginUser(s.masterDB, username, password)
	if err != nil {
		return 0, fmt.Errorf("登录失败: %v", err)
	}
	return id, nil
}
func (s *UserService) ValidateToken(token string) (int64, error) {
	return 0, nil
}
