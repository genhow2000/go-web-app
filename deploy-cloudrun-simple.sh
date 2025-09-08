#!/bin/bash

# Cloud Run 簡化版部署腳本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置變數
PROJECT_ID="fleet-day-383710"
SERVICE_NAME="go-simple-app"
REGION="asia-east1"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo -e "${BLUE}🚀 開始部署簡化版 Go 應用到 Cloud Run...${NC}"

# 1. 設置 GCP 專案
echo -e "${YELLOW}🔧 設置 GCP 專案...${NC}"
gcloud config set project $PROJECT_ID

# 2. 啟用必要的 API
echo -e "${YELLOW}📡 啟用必要的 API...${NC}"
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# 3. 構建 Docker 映像
echo -e "${YELLOW}🐳 構建 Docker 映像...${NC}"
docker build -f Dockerfile.simple -t $IMAGE_NAME:latest .

# 4. 推送映像到 Container Registry
echo -e "${YELLOW}📤 推送映像到 Container Registry...${NC}"
docker push $IMAGE_NAME:latest

# 5. 部署到 Cloud Run
echo -e "${YELLOW}🚀 部署到 Cloud Run...${NC}"
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --min-instances 0

# 6. 獲取服務 URL
echo -e "${YELLOW}🔍 獲取服務 URL...${NC}"
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format='value(status.url)')

echo -e "${GREEN}🎉 部署完成！${NC}"
echo -e "${BLUE}📋 服務資訊:${NC}"
echo -e "  • 服務名稱: $SERVICE_NAME"
echo -e "  • 服務 URL: $SERVICE_URL"
echo -e "  • 區域: $REGION"
echo -e "  • 記憶體: 512Mi"
echo -e "  • CPU: 1"

echo -e "${YELLOW}💡 測試命令:${NC}"
echo -e "  curl $SERVICE_URL"
echo -e "  curl $SERVICE_URL/health"
echo -e "  curl $SERVICE_URL/api/status"

echo -e "${GREEN}💰 費用預估: $0/月 (免費層級)${NC}"
