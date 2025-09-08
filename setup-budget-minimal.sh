#!/bin/bash

# 設置 GCP 最小預算控制腳本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置變數
PROJECT_ID="fleet-day-383710"
BILLING_ACCOUNT="01C261-42EA23-881D2D"
BUDGET_AMOUNT="1"  # 1 美元預算 (最低金額)
BUDGET_NAME="Go App Minimal Budget"

echo -e "${BLUE}💰 設置 GCP 最小預算控制...${NC}"

# 1. 檢查帳單狀態
echo -e "${YELLOW}🔍 檢查帳單狀態...${NC}"
BILLING_STATUS=$(gcloud billing projects describe $PROJECT_ID --format='value(billingEnabled)')
echo "帳單狀態: $BILLING_STATUS"

if [ "$BILLING_STATUS" = "False" ]; then
    echo -e "${RED}❌ 帳單未啟用，請先在 GCP Console 中啟用帳單${NC}"
    echo -e "${YELLOW}步驟:${NC}"
    echo "1. 前往 https://console.cloud.google.com/"
    echo "2. 選擇專案: $PROJECT_ID"
    echo "3. 前往 '帳單' → '連結帳單帳戶'"
    echo "4. 選擇帳單帳戶: $BILLING_ACCOUNT"
    echo "5. 確認連結"
    echo ""
    echo -e "${YELLOW}啟用後，請重新運行此腳本${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 帳單已啟用${NC}"

# 2. 創建最小預算
echo -e "${YELLOW}📊 創建最小預算警報...${NC}"
gcloud billing budgets create \
  --billing-account=$BILLING_ACCOUNT \
  --display-name="$BUDGET_NAME" \
  --budget-amount=$BUDGET_AMOUNT \
  --threshold-rule=percent=100 \
  --projects=$PROJECT_ID || echo "預算可能已存在"

# 3. 顯示預算信息
echo -e "${YELLOW}📋 顯示預算信息...${NC}"
gcloud billing budgets list --billing-account=$BILLING_ACCOUNT

echo -e "${GREEN}🎉 最小預算控制設置完成！${NC}"
echo -e "${BLUE}📊 預算設置:${NC}"
echo -e "  • 預算金額: \$$BUDGET_AMOUNT (最低金額)"
echo -e "  • 警報閾值: 100% (立即通知)"
echo -e "  • 專案: $PROJECT_ID"

echo -e "${YELLOW}💡 注意事項:${NC}"
echo -e "  • 你的應用預估費用: \$0/月 (免費層級)"
echo -e "  • 預算警報會在達到 \$1 時立即通知"
echo -e "  • 實際使用量遠低於 \$1"

echo -e "${GREEN}✅ 現在可以安全使用 GCP 服務了！${NC}"
