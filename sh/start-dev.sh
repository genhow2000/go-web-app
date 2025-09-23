#!/bin/bash

# 開發環境啟動腳本
echo "🚀 啟動開發環境..."

# 檢查是否已安裝 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 錯誤: 未安裝 Node.js"
    echo "請先安裝 Node.js: https://nodejs.org/"
    exit 1
fi

# 檢查是否已安裝 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ 錯誤: 未安裝 Docker"
    echo "請先安裝 Docker: https://www.docker.com/"
    exit 1
fi

# 檢查是否已安裝 docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ 錯誤: 未安裝 docker-compose"
    echo "請先安裝 docker-compose"
    exit 1
fi

echo "✅ 環境檢查通過"

# 啟動 Go 後端服務
echo "🐹 啟動 Go 後端服務..."
docker-compose -f docker-compose.dev.yml up -d go-app

# 等待後端服務啟動
echo "⏳ 等待後端服務啟動..."
sleep 5

# 檢查後端服務是否正常
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ 後端服務啟動成功"
else
    echo "⚠️  後端服務可能還在啟動中，請稍等..."
fi

# 進入前端目錄並啟動前端開發服務器
echo "🎨 啟動 Vue 前端開發服務器..."
cd frontend

# 檢查是否已安裝依賴
if [ ! -d "node_modules" ]; then
    echo "📦 安裝前端依賴..."
    npm install
fi

echo "🌟 前端開發服務器將在 http://localhost:3000 啟動"
echo "🔗 後端 API 服務在 http://localhost:8080"
echo ""
echo "💡 提示:"
echo "   - 修改前端代碼會自動熱重載"
echo "   - 修改後端代碼需要重新啟動後端服務"
echo "   - 按 Ctrl+C 停止前端服務"
echo "   - 使用 'docker-compose -f docker-compose.dev.yml down' 停止後端服務"
echo ""

# 啟動前端開發服務器
npm run dev