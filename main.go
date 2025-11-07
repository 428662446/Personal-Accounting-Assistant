package main

import (
	"fmt"
	"log"

	"AccountingAssistant/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.InitMasterDB()
	if err != nil {
		log.Fatalf("主数据库创建失败: %v", err)
	}
	defer db.Close() // 重要

	// 注册用户
	userID, err := database.RegisterUser(db, "叶叶", "1234")
	if err != nil {
		fmt.Printf("注册失败: %v\n", err)
	} else {
		fmt.Printf("注册成功, 用户ID: %d\n", userID)
	}

	// 测试重复注册
	_, err = database.RegisterUser(db, "叶叶", "4321")
	if err != nil {
		fmt.Printf("预期中的注册失败: %v\n", err)
	}

	// 测试登录
	fmt.Println("\n--- 测试登录 ---")

	// 正确登录
	loginID, err := database.LoginUser(db, "叶叶", "1234")
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
	} else {
		fmt.Printf("登录成功, 用户ID: %d\n", loginID)
	}

	// 错误密码
	_, err = database.LoginUser(db, "叶叶", "wrongpassword")
	if err != nil {
		fmt.Printf("预期中的登录失败: %v\n", err)
	}

	// 不存在的用户
	_, err = database.LoginUser(db, "不存在的用户", "1234")
	if err != nil {
		fmt.Printf("预期中的登录失败: %v\n", err)
	}

	// 记录账单
	err = database.RecordTransaction(db, "叶叶", "支出", 22.2, "标签", "备注")
	if err != nil {
		fmt.Printf("账单创建失败: %v", err)
	}
	// 查询账单
	err = database.ReadTransaction(db, "叶叶")
	if err != nil {
		fmt.Printf("查询失败: %v", err)
	}
}
