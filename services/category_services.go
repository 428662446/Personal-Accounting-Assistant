package services

import (
	"AccountingAssistant/database"
	"AccountingAssistant/models"
	"database/sql"
)

// 类别服务
type CategoryService struct {
	masterDB *sql.DB
}

// 新建类别服务的方法
func NewCategoryService(masterDB *sql.DB) *CategoryService {
	return &CategoryService{masterDB: masterDB}
}

// 新建类别服务
func (s *CategoryService) CreateCategory(userID int64, name string) (int64, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return 0, err
	}
	defer userDB.Close()
	return database.CreateCategory(userDB, name)
}

// 删除类别服务
func (s *CategoryService) DeleteCategory(userID int64, categoryID int64) error {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()
	return database.DeleteCategory(userDB, categoryID)
}

// 获取类别服务
func (s *CategoryService) GetCategory(userID int64) ([]models.Category, error) {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return nil, err
	}
	defer userDB.Close()
	return database.GetCategories(userDB)
}

// 更新类别服务
func (s *CategoryService) UpdateCategory(userID int64, catecoryID int64, updateName string) error {
	userDB, err := database.GetUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()
	return database.UpdateCategory(userDB, catecoryID, updateName)
}
