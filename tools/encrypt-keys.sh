#!/bin/sh

# API 金鑰加密工具
# 用於加密 API 金鑰配置檔案

echo "🔐 API 金鑰加密工具"
echo "=================="

# 檢查原始配置檔案
if [ ! -f "config/api-keys.conf" ]; then
    echo "❌ 找不到 config/api-keys.conf 檔案"
    exit 1
fi

echo "✅ 找到 API 金鑰配置檔案"

# 設定密碼
PASSWORD="go-app-2024"

# 加密檔案
echo "🔒 加密 API 金鑰..."
openssl enc -aes-256-cbc -salt -in config/api-keys.conf -out config/encrypted-keys.enc -pass pass:"$PASSWORD"

if [ $? -eq 0 ]; then
    echo "✅ API 金鑰已加密並儲存到 config/encrypted-keys.enc"
    echo "💡 現在可以安全地提交到 Git"
else
    echo "❌ 加密失敗"
    exit 1
fi
