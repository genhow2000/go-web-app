#!/bin/bash

# 測試開發環境配置
echo "🧪 測試開發環境配置..."

# 檢查必要文件是否存在
echo "📁 檢查文件結構..."

files=(
    "Dockerfile.dev"
    "docker-compose.dev.yml"
    "start-dev.sh"
    "stop-dev.sh"
    "deploy-gcp.sh"
    "frontend/package.json"
    "frontend/vite.config.js"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file"
    else
        echo "❌ $file (缺失)"
    fi
done

echo ""
echo "🔧 檢查環境依賴..."

# 檢查 Node.js
if command -v node &> /dev/null; then
    echo "✅ Node.js $(node --version)"
else
    echo "❌ Node.js (未安裝)"
fi

# 檢查 npm
if command -v npm &> /dev/null; then
    echo "✅ npm $(npm --version)"
else
    echo "❌ npm (未安裝)"
fi

# 檢查 Docker
if command -v docker &> /dev/null; then
    echo "✅ Docker $(docker --version)"
else
    echo "❌ Docker (未安裝)"
fi

# 檢查 docker-compose
if command -v docker-compose &> /dev/null; then
    echo "✅ docker-compose $(docker-compose --version)"
else
    echo "❌ docker-compose (未安裝)"
fi

echo ""
echo "📋 測試完成！"
echo ""
echo "💡 下一步："
echo "   1. 確保所有依賴都已安裝"
echo "   2. 運行 './start-dev.sh' 啟動開發環境"
echo "   3. 訪問 http://localhost:3000 查看前端"
echo "   4. 訪問 http://localhost:8080 查看後端 API"
