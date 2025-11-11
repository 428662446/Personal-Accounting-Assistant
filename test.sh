#!/bin/bash
#!/bin/bash
echo "=== 测试开始 ==="

# 注册用户（使用正确的编码）
echo "username=testuser&password=testpass123" | \
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
  --data-binary @-

echo -e "\n2. 登录用户"
echo "username=testuser&password=testpass123" | \
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
  --data-binary @- \
  -c cookies.txt -b cookies.txt

echo -e "\n3. 创建类别（中文测试）"
echo "name=餐饮" | \
curl -X POST http://localhost:8080/category \
  -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
  --data-binary @- \
  -b cookies.txt

echo -e "\n4. 再创建一个类别"
echo "name=交通出行" | \
curl -X POST http://localhost:8080/category \
  -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" \
  --data-binary @- \
  -b cookies.txt

echo -e "\n5. 获取类别"
curl -X GET http://localhost:8080/categories -b cookies.txt

echo -e "\n=== 测试结束 ==="

echo -e "\n=== 测试结束 ==="

sleep 30

