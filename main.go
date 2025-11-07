package main

import (
	"fmt"

	"AccountingAssistant/database"
	"AccountingAssistant/handlers"
	"AccountingAssistant/middleware"
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
	transactionService := services.NewTransactionService(db)
	// 添加: 基于数据库的会话管理器
	sessionManager := services.NewDBSessionManager(db)

	authHandler := handlers.NewAuthHandler(userService, sessionManager)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	r := gin.Default()

	r.POST("/register", authHandler.RegisterUser)
	r.POST("/login", authHandler.LoginUser)
	// 需要认证的路由组（先应用会话中间件，再应用认证中间件）
	authGroup := r.Group("/")
	authGroup.Use(middleware.SessionMiddleware(sessionManager), middleware.AuthRequired())
	{
		authGroup.POST("/transaction", transactionHandler.RecordTransaction)
		authGroup.GET("/transactions", transactionHandler.GetTransactions)
		authGroup.POST("/logout", authHandler.LogoutUser) // 添加退出登录
	}
	r.Run(":8080")
}
