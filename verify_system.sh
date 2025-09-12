#!/bin/bash

echo "🔍 開始驗證系統狀態..."
echo "=================================================="

# 檢查 Go 模組
echo "1. 檢查 Go 模組..."
if go mod tidy; then
    echo "✅ Go 模組檢查通過"
else
    echo "❌ Go 模組檢查失敗"
    exit 1
fi

# 檢查編譯
echo -e "\n2. 檢查編譯..."
if go build -o app .; then
    echo "✅ 編譯成功"
    rm -f app  # 清理編譯文件
else
    echo "❌ 編譯失敗"
    exit 1
fi

# 檢查遷移文件
echo -e "\n3. 檢查遷移文件..."
migration_files=(
    "migrations/001_add_role_and_status.sql"
    "migrations/002_add_merchant_role.sql" 
    "migrations/003_create_separate_role_tables.sql"
)

for file in "${migration_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 不存在"
    fi
done

# 檢查模型文件
echo -e "\n4. 檢查模型文件..."
model_files=(
    "models/customer.go"
    "models/merchant.go"
    "models/admin.go"
    "models/user_interface.go"
)

for file in "${model_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 不存在"
    fi
done

# 檢查服務文件
echo -e "\n5. 檢查服務文件..."
service_files=(
    "services/unified_auth_service.go"
)

for file in "${service_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 不存在"
    fi
done

# 檢查控制器文件
echo -e "\n6. 檢查控制器文件..."
controller_files=(
    "controllers/unified_auth_controller.go"
)

for file in "${controller_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 不存在"
    fi
done

# 檢查模板文件
echo -e "\n7. 檢查模板文件..."
template_files=(
    "templates/customer_login.html"
    "templates/customer_dashboard.html"
)

for file in "${template_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 不存在"
    fi
done

# 檢查語法錯誤
echo -e "\n8. 檢查語法錯誤..."
if go vet ./...; then
    echo "✅ 語法檢查通過"
else
    echo "❌ 語法檢查發現問題"
fi

echo -e "\n=================================================="
echo "✅ 系統驗證完成！"
echo ""
echo "📋 下一步操作："
echo "1. 啟動服務器: go run main.go"
echo "2. 訪問以下 URL 測試："
echo "   - 客戶登入: http://localhost:8080/auth/customer/login"
echo "   - 商戶登入: http://localhost:8080/merchant/login"
echo "   - 管理員登入: http://localhost:8080/admin/login"
echo "   - 註冊頁面: http://localhost:8080/auth/register"
echo "3. 測試 API 端點："
echo "   - 健康檢查: curl http://localhost:8080/health"
echo "   - 系統統計: curl http://localhost:8080/api/stats"
