#!/bin/bash

# 環境變數設置腳本
# 請在運行前設置您的 API 金鑰

echo "設置 Go 應用程式環境變數..."

# 檢查是否已設置 API 金鑰
if [ -z "$GROQ_API_KEY" ]; then
    echo "請設置 GROQ_API_KEY 環境變數"
    echo "export GROQ_API_KEY=\"your-groq-api-key\""
    exit 1
fi

if [ -z "$GEMINI_API_KEY" ]; then
    echo "請設置 GEMINI_API_KEY 環境變數"
    echo "export GEMINI_API_KEY=\"your-gemini-api-key\""
    exit 1
fi

# 設置其他環境變數
export JWT_SECRET="your-secure-jwt-secret-key-$(date +%s)"
export AI_PRIMARY_PROVIDER="groq"
export AI_FALLBACK_PROVIDER="gemini"
export DB_PATH="/app/data/app.db"
export MONGODB_URI="mongodb+srv://genhow2000_db_user:PKwCLU4qXzrpAjSY@cluster0.atzcxpw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
export MONGODB_DATABASE="chatbot"

echo "環境變數設置完成！"
echo "GROQ_API_KEY: ${GROQ_API_KEY:0:10}..."
echo "GEMINI_API_KEY: ${GEMINI_API_KEY:0:10}..."
echo "JWT_SECRET: ${JWT_SECRET:0:10}..."
