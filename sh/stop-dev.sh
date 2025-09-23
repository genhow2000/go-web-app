#!/bin/bash

# 停止開發環境腳本
echo "🛑 停止開發環境..."

# 停止 Docker 容器
echo "🐳 停止 Docker 容器..."
docker-compose -f docker-compose.dev.yml down

echo "✅ 開發環境已停止"
echo ""
echo "💡 提示:"
echo "   - 前端服務會自動停止（如果正在運行）"
echo "   - 後端 Docker 容器已停止"
echo "   - 使用 './start-dev.sh' 重新啟動開發環境"
