#!/bin/bash

# 開發環境啟動腳本
# 同時啟動 Vite 開發服務器和 nginx 反向代理

echo "啟動開發環境..."

# 檢查 nginx 是否安裝
if ! command -v nginx &> /dev/null; then
    echo "請先安裝 nginx:"
    echo "brew install nginx"
    exit 1
fi

# 啟動 Vite 開發服務器（後台）
echo "啟動 Vite 開發服務器..."
cd frontend
npm run dev &
VITE_PID=$!

# 等待 Vite 啟動
sleep 3

# 啟動 nginx 反向代理
echo "啟動 nginx 反向代理..."
nginx -c $(pwd)/nginx.conf

echo "開發環境已啟動！"
echo "前端: http://localhost:8080 (通過 nginx 代理到 Vite)"
echo "後端 API: http://localhost:8080/api (直接訪問 Go 後端)"
echo ""
echo "按 Ctrl+C 停止所有服務"

# 等待用戶中斷
wait $VITE_PID
