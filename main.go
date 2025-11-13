package main

import (
	"fmt"

	"AccountingAssistant/database"
	"AccountingAssistant/handlers"
	"AccountingAssistant/web/middleware"

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
	// 创建服务实例并注入主数据库连接。注意：transactionService 会按需打开每个用户的 per-user DB，
	// master DB 在 service 中仅用于访问用户表/会话等全局元数据，不用于持久化某个用户的事务数据。
	userService := services.NewUserService(db)
	transactionService := services.NewTransactionService(db)
	statService := services.NewStatService(db)
	categoryServic := services.NewCategoryService(db)
	// 添加: 基于数据库的会话管理器
	sessionManager := services.NewDBSessionManager(db)

	authHandler := handlers.NewAuthHandler(userService, sessionManager)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	statHandler := handlers.NewStatHandler(statService)
	categoryHandler := handlers.NewCategoryHandler(categoryServic)
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

		authGroup.POST("/category", categoryHandler.CreateCategory)
		authGroup.GET("/categories", categoryHandler.GetCategory)
		authGroup.PUT("/category/:id", categoryHandler.UpdateCategory)    // 更新特定类别
		authGroup.DELETE("/category/:id", categoryHandler.DeleteCategory) // 删除特定类别

		authGroup.GET("/stats/summary", statHandler.GetSummary)
		authGroup.GET("/stats/monthly", statHandler.GetMonthlyStats)
		authGroup.GET("/stats/weekly", statHandler.GetWeeklyStats)
		authGroup.GET("/stats/daily", statHandler.GetDailyStats)
		authGroup.GET("/stats/range_amount", statHandler.GetRangeAmountStats)
	}
	r.Run(":8080")
}
