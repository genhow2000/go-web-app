#!/bin/bash

echo "🚀 測試Vue.js Docker構建..."

# 檢查Docker是否運行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未運行，請先啟動Docker"
    exit 1
fi

echo "✅ Docker正在運行"

# 構建Docker映像
echo "🔨 構建Docker映像..."
docker build -t go-vue-app .

if [ $? -eq 0 ]; then
    echo "✅ Docker映像構建成功"
else
    echo "❌ Docker映像構建失敗"
    exit 1
fi

# 檢查構建的文件
echo "📁 檢查構建文件..."
docker run --rm go-vue-app ls -la /root/static/dist/

if [ $? -eq 0 ]; then
    echo "✅ Vue.js構建文件存在"
else
    echo "❌ Vue.js構建文件不存在"
    exit 1
fi

echo "🎉 測試完成！Vue.js已成功集成到Docker中"
echo "💡 現在可以運行: docker-compose up"
