#!/bin/bash

echo "🐳 開始驗證 Docker 系統狀態..."
echo "=================================================="

# 檢查 Docker 是否安裝
echo "1. 檢查 Docker..."
if command -v docker &> /dev/null; then
    echo "✅ Docker 已安裝"
    docker --version
else
    echo "❌ Docker 未安裝"
    exit 1
fi

# 檢查 Docker Compose 是否安裝
echo -e "\n2. 檢查 Docker Compose..."
if command -v docker-compose &> /dev/null; then
    echo "✅ Docker Compose 已安裝"
    docker-compose --version
else
    echo "❌ Docker Compose 未安裝"
    exit 1
fi

# 檢查必要文件
echo -e "\n3. 檢查必要文件..."
files=(
    "Dockerfile"
    "docker-compose.yml"
    "go.mod"
    "main.go"
    "models/customer.go"
    "models/merchant.go"
    "models/admin.go"
    "models/user_interface.go"
    "services/unified_auth_service.go"
    "controllers/unified_auth_controller.go"
    "migrations/003_create_separate_role_tables.sql"
    "templates/customer_login.html"
    "templates/customer_dashboard.html"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 不存在"
    fi
done

# 檢查 Docker 構建
echo -e "\n4. 檢查 Docker 構建..."
if docker build -t go-simple-app . > /dev/null 2>&1; then
    echo "✅ Docker 構建成功"
else
    echo "❌ Docker 構建失敗"
    echo "請運行 'docker build -t go-simple-app .' 查看詳細錯誤"
fi

# 檢查容器狀態
echo -e "\n5. 檢查容器狀態..."
if docker ps -a | grep -q go-simple-app; then
    echo "✅ 容器 go-simple-app 存在"
    docker ps -a | grep go-simple-app
else
    echo "ℹ️  容器 go-simple-app 不存在（正常，尚未啟動）"
fi

echo -e "\n=================================================="
echo "✅ Docker 系統驗證完成！"
echo ""
echo "📋 下一步操作："
echo "1. 啟動服務: docker-compose up -d"
echo "2. 查看日誌: docker-compose logs -f"
echo "3. 停止服務: docker-compose down"
echo "4. 重新構建: docker-compose up --build -d"
echo ""
echo "🌐 服務啟動後可訪問："
echo "- 客戶登入: http://localhost:8080/auth/customer/login"
echo "- 商戶登入: http://localhost:8080/merchant/login"
echo "- 管理員登入: http://localhost:8080/admin/login"
echo "- 註冊頁面: http://localhost:8080/auth/register"
echo "- 健康檢查: http://localhost:8080/health"
