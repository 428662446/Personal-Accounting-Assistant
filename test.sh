#!/bin/bash

# 个人记账助手 API 测试脚本
# 作者：428662446
# 描述：完整的 API 端到端测试(这是AI写的，只对接口期望的字段名做了检查更改)

set -e  # 遇到错误立即退出

echo "=========================================="
echo "   个人记账助手 API 测试脚本"
echo "=========================================="

# 配置
BASE_URL="http://localhost:8080"
COOKIE_FILE="test_cookies.txt"

# 清理之前的 cookie 文件
rm -f $COOKIE_FILE

# 颜色输出函数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# 检查服务是否运行
check_server() {
    print_info "检查服务是否运行在 $BASE_URL..."
    if curl -s --head $BASE_URL > /dev/null; then
        print_success "服务运行正常"
    else
        print_error "服务未运行，请先启动服务(运行go run main.go)：go run main.go"
        sleep 30
        exit 1
    fi
}

# 测试用户注册
test_register() {
    print_info "1. 测试用户注册..."
    
    # 测试正常注册
    response=$(echo "username=testuser_$(date +%s)&password=testpass123" | \
        curl -s -w "%{http_code}" -X POST $BASE_URL/register \
        -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
        --data-binary @-)
    
    status_code=${response: -3}
    body=${response%???}
    
    if [ "$status_code" -eq 200 ]; then
        print_success "用户注册成功"
    else
        print_error "用户注册失败: $body"
    fi
}

# 测试用户登录
test_login() {
    print_info "2. 测试用户登录..."
    
    # 测试错误密码
    print_info "  尝试错误密码登录..."
    response=$(echo "username=testuser&password=wrongpass" | \
        curl -s -w "%{http_code}" -X POST $BASE_URL/login \
        -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
        --data-binary @- \
        -c $COOKIE_FILE -b $COOKIE_FILE)
    
    status_code=${response: -3}
    if [ "$status_code" -ne 200 ]; then
        print_success "错误密码被正确拒绝"
    fi
    
    # 测试正确密码
    print_info "  尝试正确密码登录..."
    response=$(echo "username=testuser&password=testpass123" | \
        curl -s -w "%{http_code}" -X POST $BASE_URL/login \
        -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
        --data-binary @- \
        -c $COOKIE_FILE -b $COOKIE_FILE)
    
    status_code=${response: -3}
    body=${response%???}
    
    if [ "$status_code" -eq 200 ]; then
        print_success "用户登录成功"
        echo "   响应: $body"
    else
        print_error "用户登录失败"
    fi
}

# 测试分类管理
test_categories() {
    print_info "3. 测试分类管理..."
    
    # 创建分类
    categories=("餐饮" "交通出行" "购物" "娱乐" "医疗")
    for category in "${categories[@]}"; do
        print_info "  创建分类: $category"
        response=$(echo "name=$category" | \
            curl -s -w "%{http_code}" -X POST $BASE_URL/category \
            -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
            --data-binary @- \
            -b $COOKIE_FILE)
        
        status_code=${response: -3}
        if [ "$status_code" -eq 200 ]; then
            print_success "  分类 '$category' 创建成功"
        else
            print_error "  分类 '$category' 创建失败"
        fi
    done
    
    # 获取分类列表
    print_info "  获取分类列表..."
    response=$(curl -s -w "%{http_code}" -X GET $BASE_URL/categories \
        -b $COOKIE_FILE)
    
    status_code=${response: -3}
    body=${response%???}
    
    if [ "$status_code" -eq 200 ]; then
        print_success "获取分类列表成功"
        echo "   分类列表: $body"
    else
        print_error "获取分类列表失败"
    fi
}

test_transactions() {
    print_info "4. 测试交易记录..."
    
    # 记录几笔交易
    transactions=(
        "type=income&amount=5000.00&category=工资&note=月工资收入"
        "type=expense&amount=25.50&category=餐饮&note=午餐"
        "type=expense&amount=8.00&category=交通出行&note=地铁费"
        "type=expense&amount=120.00&category=购物&note=买衣服"
        "type=expense&amount=60.00&category=娱乐&note=看电影"
        "type=expense&amount=200.00&category=医疗&note=买药"
    )
    
    for transaction in "${transactions[@]}"; do
        print_info "  记录交易..."
        response=$(echo "$transaction" | \
            curl -s -w "%{http_code}" -X POST $BASE_URL/transaction \
            -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
            --data-binary @- \
            -b $COOKIE_FILE)
        
        status_code=${response: -3}
        body=${response%???}
        
        if [ "$status_code" -eq 200 ]; then
            print_success "  交易记录成功"
        else
            print_error "  交易记录失败: $body"
        fi
    done
    
    # 获取交易列表
    print_info "  获取交易列表..."
    response=$(curl -s -w "%{http_code}" -X GET $BASE_URL/transactions \
        -b $COOKIE_FILE)
    
    status_code=${response: -3}
    body=${response%???}
    
    if [ "$status_code" -eq 200 ]; then
        print_success "获取交易列表成功"
        # 尝试解析交易数量
        if command -v jq >/dev/null 2>&1; then
            count=$(echo "$body" | jq '.transactions | length' 2>/dev/null || echo "未知")
            echo "   交易数量: $count"
        else
            echo "   响应数据: $body"
        fi
    else
        print_error "获取交易列表失败"
    fi
}

# 测试统计功能
test_statistics() {
    print_info "5. 测试统计功能..."
    
    endpoints=(
        "/stats/summary"
        "/stats/monthly"
        "/stats/weekly"
        "/stats/daily"
        "/stats/range_amount"
    )
    
    for endpoint in "${endpoints[@]}"; do
        print_info "  测试 $endpoint ..."
        response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL$endpoint" \
            -b $COOKIE_FILE)
        
        status_code=${response: -3}
        body=${response%???}
        
        if [ "$status_code" -eq 200 ]; then
            print_success "  $endpoint 请求成功"
        else
            print_error "  $endpoint 请求失败"
        fi
    done
}

# 测试分类更新和删除
test_category_operations() {
    print_info "6. 测试分类更新和删除..."
    
    # 方法1: 直接创建一个新分类用于测试操作
    print_info "  创建专门用于测试的分类..."
    response=$(echo "name=测试操作分类" | \
        curl -s -X POST $BASE_URL/category \
        -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
        --data-binary @- \
        -b $COOKIE_FILE)
    
    echo "创建分类响应: $response"
    
    # 从响应中提取分类ID
    category_id=""
    
    # 尝试多种方式提取分类ID
    if [[ "$response" == *"category_id"* ]]; then
        # 从JSON响应中提取category_id
        category_id=$(echo "$response" | grep -o '"category_id":[^,}]*' | cut -d':' -f2 | tr -d ' "')
    fi
    
    # 如果还是没找到，尝试从分类列表中获取第一个分类
    if [ -z "$category_id" ]; then
        print_info "  从分类列表中获取分类ID..."
        list_response=$(curl -s -X GET $BASE_URL/categories -b $COOKIE_FILE)
        echo "分类列表: $list_response"
        
        # 尝试提取第一个分类的ID
        if [[ "$list_response" == *"id"* ]]; then
            category_id=$(echo "$list_response" | grep -o '"id":[^,}]*' | head -1 | cut -d':' -f2 | tr -d ' "')
        fi
    fi
    
    if [ -n "$category_id" ] && [ "$category_id" != "null" ]; then
        print_success "  找到分类ID: $category_id"
        
        # 更新分类
        print_info "  更新分类 ID: $category_id"
        response=$(echo "name=更新后的分类名称_$(date +%s)" | \
            curl -s -w "%{http_code}" -X PUT "$BASE_URL/category/$category_id" \
            -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
            --data-binary @- \
            -b $COOKIE_FILE)
        
        status_code=${response: -3}
        if [ "$status_code" -eq 200 ]; then
            print_success "  分类更新成功"
        else
            print_error "  分类更新失败: ${response%???}"
        fi
        
        # 删除分类
        print_info "  删除分类 ID: $category_id"
        response=$(curl -s -w "%{http_code}" -X DELETE "$BASE_URL/category/$category_id" \
            -b $COOKIE_FILE)
        
        status_code=${response: -3}
        if [ "$status_code" -eq 200 ]; then
            print_success "  分类删除成功"
        else
            print_error "  分类删除失败: ${response%???}"
        fi
    else
        print_error "  无法获取分类ID，跳过更新删除测试"
        echo "  原始响应: $response"
    fi
}

# 测试退出登录
test_logout() {
    print_info "7. 测试退出登录..."
    
    response=$(curl -s -w "%{http_code}" -X POST $BASE_URL/logout \
        -b $COOKIE_FILE)
    
    status_code=${response: -3}
    if [ "$status_code" -eq 200 ]; then
        print_success "退出登录成功"
    else
        print_error "退出登录失败"
    fi
}

# 主测试流程
main() {
    echo "开始测试 $(date)"
    echo ""
    
    check_server
    test_register
    test_login
    test_categories
    test_transactions
    test_statistics
    test_category_operations
    test_logout
    
    echo ""
    echo "=========================================="
    print_success "测试完成 $(date)"
    echo "=========================================="
    sleep 30
    
    # 清理
    rm -f $COOKIE_FILE
}

# 运行主函数
main

