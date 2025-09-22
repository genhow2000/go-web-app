# GitHub Secrets 檢查清單

為了讓 GitHub Actions 自動部署正常工作，您需要在 GitHub 倉庫中設置以下 Secrets：

## 必需的 Secrets

### 1. GCP_SA_KEY
- **描述**: Google Cloud 服務帳戶的 JSON 金鑰
- **如何獲取**:
  1. 前往 Google Cloud Console > IAM & Admin > Service Accounts
  2. 找到或創建一個服務帳戶
  3. 點擊 "Actions" > "Manage keys" > "Add key" > "Create new key"
  4. 選擇 JSON 格式並下載
  5. 將整個 JSON 內容複製到 GitHub Secrets

### 2. JWT_SECRET
- **描述**: JWT 簽名密鑰
- **建議值**: 使用一個強隨機字符串，例如：
  ```
  your-super-secret-jwt-key-here-make-it-long-and-random
  ```

### 3. GROQ_API_KEY
- **描述**: Groq API 金鑰
- **如何獲取**: 前往 https://console.groq.com/ 獲取 API 金鑰

### 4. GEMINI_API_KEY
- **描述**: Google Gemini API 金鑰
- **如何獲取**: 前往 https://makersuite.google.com/app/apikey 獲取 API 金鑰

## 如何設置 GitHub Secrets

1. 前往您的 GitHub 倉庫：https://github.com/genhow2000/go-web-app
2. 點擊 "Settings" 標籤
3. 在左側選單中點擊 "Secrets and variables" > "Actions"
4. 點擊 "New repository secret"
5. 輸入 Secret 名稱和值
6. 點擊 "Add secret"

## 檢查 Secrets 是否設置正確

設置完成後，您可以：
1. 推送代碼到 main 分支觸發部署
2. 前往 "Actions" 標籤查看部署狀態
3. 如果部署失敗，檢查 Actions 日誌中的錯誤信息

## 當前配置狀態

- ✅ GitHub Actions 工作流程已配置
- ✅ Docker 構建配置正確
- ✅ Cloud Run 部署配置正確
- ❓ 需要檢查 GitHub Secrets 是否已設置
