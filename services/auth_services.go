package services

import (
	"AccountingAssistant/database"
	"database/sql"
)

type UserService struct {
	masterDB *sql.DB
}

func NewUserService(masterDB *sql.DB) *UserService {
	return &UserService{masterDB: masterDB}
}
func (s *UserService) Register(username, password string) (int64, error) {
	return database.RegisterUser(s.masterDB, username, password)
}

func (s *UserService) Login(username, password string) (int64, error) {
	return database.LoginUser(s.masterDB, username, password)
}
