# 个人记账助手 (Personal Accounting Assistant)

> Token团队2025考核项目 - 一个功能差不多完整的多用户记账系统

## 快捷测试
```bash
./test.sh
```

## 简介

这是一个基于Go语言和Gin框架开发的个人记账助手，完成了题目要求的基础功能，在安全性、架构设计和扩展性方面做了深入思考和实践

## 核心特性

### 基础功能
- ✅ **用户管理** - 注册、登录、会话管理
- ✅ **收支记录** - 完整的CRUD操作（添加、查看、修改、删除）
- ✅ **类别管理** - 自定义收支类别
- ✅ **数据统计** - 日/周/月统计、金额范围分析
- ✅ **数据持久化** - SQLite本地存储，重启数据不丢失

### 进阶实现
- ✅ **多用户支持** - 完整的多用户隔离系统
- ✅ **高精度运算** - 基于分的金额存储，避免浮点误差
- ✅ **安全架构** - 会话管理、密码加密、错误处理
- ✅ **RESTful API** - 清晰的接口设计，易于扩展

## 技术栈

**后端框架**
- Go 1.21+ 
- Gin Web Framework
- SQLite3 数据库

**安全与架构**
- Bcrypt 密码加密
- 数据库会话管理
- 统一错误处理
- 中间件认证

## 架构设计思路
### 分层架构
```text
HTTP层 (handlers) → 业务层 (services) → 数据层 (database) → 数据库
```
### 扩展性考虑
**模块化设计** - 各功能模块独立，易于维护
**接口清晰** - 明确的层间边界，便于测试
**错误处理** - 统一错误码，便于前端处理
**数据隔离** - 为后续功能扩展预留空间

## 设计思考文档
在开发过程中，对关键问题的思考笔记
>未仔细核对实际实现完全一致

[金额处理设计思路](./log_txt/add_amount_processing.txt) - 精确计算方案
[类别管理架构](./log_txt/category_management_design%20.txt) - 数据关系设计
错误处理规范 - 用户体验优化

## 项目结构
AccountingAssistant/
├── database/ # 数据访问层
│ ├── master.go # 主数据库初始化
│ ├── user_db.go # 用户数据库操作
│ ├── transaction_db.go
│ ├── category_db.go
│ └── stats_db.go
├── handlers/ # HTTP处理器
│ ├── auth_handler.go
│ ├── transaction_handler.go
│ ├── category_handler.go
│ └── stats_handler.go
├── models/ # 数据模型
│ └── models.go
├── services/ # 业务逻辑层
│ ├── auth_services.go
│ ├── transaction_service.go
│ ├── category_services.go
│ ├── stats_service.go
│ └── session_service.go
├── utils/ # 工具包
│ ├── errors.go # 统一错误处理
│ ├── amount.go # 金额精确处理
│ ├── amount_test.go
│ └── password.go # 密码加密
└── web/
  ├── middleware/ # 中间件
  │     └── session.go
  └── response/ # HTTP响应处理
        └── response.go

## 安全设计亮点

### 1. 密码安全
```go
// 使用bcrypt进行密码哈希
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
```
### 2. 会话管理
- 基于数据库的会话存储
- 中间件统一验证
### 3. 数据隔离
```go
// 每个用户独立的数据库文件
func GetUserDB(userId int64) (*sql.DB, error) {
    userDBPath := filepath.Join("database_files", "usersdata", fmt.Sprintf("user_%d.db", userId))
    // ...
}
```
### 4. 统一错误处理
```go
// 分级错误处理机制
func HandleError(c *gin.Context, err error) {
    // 业务错误、系统错误、认证错误分类处理
}
```
## 金额精确处理设计
### 问题背景
浮点数在金融计算中存在精度问题：
```go
// 错误示例
0.1 + 0.2 = 0.30000000000000004
```
### 解决方案
- 字符串传输 - 前后端使用字符串传递金额
- 分单位存储 - 数据库中以分为单位存储整数
- 精确转换 - 四舍五入到分，避免累计误差

```go
// 核心转换逻辑
func ParseToCents(str string) (int64, error) {
    // 清理输入 → 验证格式 → 四舍五入 → 转换为分
}
```
### 处理流程
```text
用户输入 "123.456" 
→ 清理为 "123.456" 
→ 四舍五入 "123.46" 
→ 转换为分 12346
→ 存储为整数
```

## 开始
### 环境要求
- Go 1.21+
- SQLite3
### 安装运行(仅供参考)
```bash
# 克隆项目
git clone <repository-url>

# 进入目录
cd AccountingAssistant

# 运行测试
go test ./...

# 启动服务
go run main.go
```
### API示例
#### 用户注册

```http
POST /register
Content-Type: application/x-www-form-urlencoded

username=testuser&password=testpass
```
#### 记录账单

```http
POST /transaction
Content-Type: application/x-www-form-urlencoded
Cookie: session_id=xxx

type=expense&amount=123.45&category=餐饮&note=午餐
```