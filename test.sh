# 1. 注册
curl -X POST http://localhost:8080/register \
  -F "username=test" \
  -F "password=123"

# 2. 登录（保存Cookie）
curl -X POST http://localhost:8080/login \
  -F "username=test" \
  -F "password=123" \
  -c cookies.txt

# 3. 记录交易（使用Cookie）
curl -X POST http://localhost:8080/transaction \
  -b cookies.txt \
  -F "type=income" \
  -F "amount=100" \
  -F "category=salary" \
  -F "note=test"

# 4. 获取交易
curl -X GET http://localhost:8080/transactions -b cookies.txt

# 5. 退出登录
curl -X POST http://localhost:8080/logout -b cookies.txt
sleep 30 