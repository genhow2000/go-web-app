#!/bin/sh

# 解密 API 金鑰並啟動應用程式
# 一鍵執行，自動解密並啟動

echo "🔓 解密 API 金鑰並啟動應用程式"
echo "=============================="

# 切換到專案根目錄
cd "$(dirname "$0")/.."

# 檢查加密檔案
if [ ! -f "config/encrypted-keys.enc" ]; then
    echo "❌ 找不到 config/encrypted-keys.enc 檔案"
    echo "💡 請確保專案檔案完整"
    exit 1
fi

echo "✅ 找到加密的 API 金鑰檔案"

# 設定密碼
PASSWORD="go-app-2024"

# 解密檔案
echo "🔓 解密 API 金鑰..."
openssl enc -aes-256-cbc -d -salt -in config/encrypted-keys.enc -out .env -pass pass:"$PASSWORD"

if [ $? -ne 0 ]; then
    echo "❌ 解密失敗，請檢查密碼或檔案完整性"
    exit 1
fi

echo "✅ API 金鑰已解密並寫入 .env 檔案"

# 檢查 Docker
if ! command -v docker >/dev/null 2>&1; then
    echo "❌ Docker 未安裝，請先安裝 Docker"
    exit 1
fi

if ! command -v docker-compose >/dev/null 2>&1; then
    echo "❌ Docker Compose 未安裝，請先安裝 Docker Compose"
    exit 1
fi

echo "✅ Docker 環境檢查通過"

# 停止現有容器
echo "🛑 停止現有容器..."
docker-compose down 2>/dev/null

# 啟動新容器
echo "🏗️  構建和啟動容器..."
docker-compose up --build -d

echo ""
echo "🎉 專案啟動完成！"
echo ""
echo "📋 專案資訊:"
echo "  - 應用程式: http://localhost:8080"
echo "  - 健康檢查: http://localhost:8080/health"
echo "  - 管理介面: http://localhost:8080/db-manager"
echo ""
echo "💡 常用命令:"
echo "  - 查看日誌: docker-compose logs -f"
echo "  - 停止專案: docker-compose down"
echo "  - 重啟專案: docker-compose restart"

# 清理臨時檔案（可選）
# rm -f .env
