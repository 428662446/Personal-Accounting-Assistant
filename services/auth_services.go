package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/utils"
	"database/sql"
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
		return 0, utils.WrapError(utils.ErrRegisterFailed, err)
	}
	return userId, nil
}

func (s *UserService) Login(username, password string) (int64, error) {
	id, err := database.LoginUser(s.masterDB, username, password)
	if err != nil {
		return 0, utils.WrapError(utils.ErrLoginFailed, err)
	}
	return id, nil
}
