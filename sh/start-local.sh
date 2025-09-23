#!/bin/bash

# 本地開發腳本
echo "🏠 啟動本地開發環境..."

# 設置本地環境變數
export LINE_CLIENT_ID="2008159551"
export LINE_CLIENT_SECRET="2cca495d6b53e8b2a2d684ee87113f01"
export LINE_REDIRECT_URL="http://localhost:8080/auth/line/callback"
export BASE_URL="http://localhost:8080"

echo "🔧 本地 LINE OAuth 配置:"
echo "  - LINE_CLIENT_ID: $LINE_CLIENT_ID"
echo "  - LINE_REDIRECT_URL: $LINE_REDIRECT_URL"
echo "  - BASE_URL: $BASE_URL"

# 啟動 Docker Compose
echo "🐳 啟動 Docker 容器..."
docker-compose up -d --build

echo "✅ 本地開發環境已啟動！"
echo "🌐 應用程式網址: http://localhost:8080"
echo "🔗 LINE 登入網址: http://localhost:8080/customer/login"
echo ""
echo "📝 請確保在 LINE 開發者控制台添加以下 Callback URL:"
echo "   http://localhost:8080/auth/line/callback"
