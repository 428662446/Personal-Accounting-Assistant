package main

import (
	"fmt"

	"AccountingAssistant/database"
	"AccountingAssistant/handlers"
	"AccountingAssistant/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 初始化主数据库
	db, err := database.InitMasterDB()
	if err != nil {
		fmt.Printf("主数据库初始化失败： %v", err)
		return
	}
	defer db.Close()
	// 创建服务实例
	userService := services.NewUserService(db)
	authHandler := handlers.NewAuthHandler(userService)
	// 请求
	r := gin.Default()
	r.POST("/user/Register", authHandler.RegisterUser)
	r.GET("/user/Login", authHandler.LoginUser)
}
