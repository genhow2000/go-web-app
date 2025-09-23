#!/bin/bash

# GCP 部署腳本
echo "🚀 部署到 GCP..."

# 檢查是否已安裝 gcloud
if ! command -v gcloud &> /dev/null; then
    echo "❌ 錯誤: 未安裝 Google Cloud CLI"
    echo "請先安裝: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# 檢查是否已登入
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "❌ 錯誤: 未登入 Google Cloud"
    echo "請先執行: gcloud auth login"
    exit 1
fi

# 獲取專案 ID
PROJECT_ID=$(gcloud config get-value project)
if [ -z "$PROJECT_ID" ]; then
    echo "❌ 錯誤: 未設置 Google Cloud 專案"
    echo "請先執行: gcloud config set project YOUR_PROJECT_ID"
    exit 1
fi

echo "📋 專案 ID: $PROJECT_ID"

# 構建並推送 Docker 映像
echo "🐳 構建 Docker 映像..."
docker build -t gcr.io/$PROJECT_ID/go-web-app .

echo "📤 推送到 Container Registry..."
docker push gcr.io/$PROJECT_ID/go-web-app

# 部署到 Cloud Run
echo "☁️  部署到 Cloud Run..."
gcloud run deploy go-web-app \
  --image gcr.io/$PROJECT_ID/go-web-app \
  --platform managed \
  --region asia-east1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --max-instances 10

echo "✅ 部署完成！"
echo ""
echo "🌐 您的應用已部署到 Cloud Run"
echo "🔗 查看部署狀態: https://console.cloud.google.com/run"
