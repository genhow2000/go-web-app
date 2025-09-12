# 安全部署指南

## 🔐 API 金鑰安全設置

### 1. 本地開發環境

```bash
# 設置環境變數（不要提交到 Git）
export GROQ_API_KEY="your-groq-api-key"
export GEMINI_API_KEY="your-gemini-api-key"

# 或使用腳本
./setup-env.sh
```

### 2. 雲端部署環境

#### 方法 1：使用 Google Cloud Secret Manager（推薦）

```bash
# 創建 Secret
gcloud secrets create groq-api-key --data-file=- <<< "your-groq-api-key"
gcloud secrets create gemini-api-key --data-file=- <<< "your-gemini-api-key"

# 更新 Cloud Run 服務
gcloud run services update go-app \
  --region=asia-east1 \
  --set-secrets="GROQ_API_KEY=groq-api-key:latest,GEMINI_API_KEY=gemini-api-key:latest"
```

#### 方法 2：直接設置環境變數

```bash
gcloud run services update go-app \
  --region=asia-east1 \
  --set-env-vars="GROQ_API_KEY=your-groq-key,GEMINI_API_KEY=your-gemini-key"
```

### 3. CI/CD 安全設置

在 Google Cloud Build 中設置替換變數：

```bash
gcloud builds submit --substitutions=_GROQ_API_KEY="your-groq-key",_GEMINI_API_KEY="your-gemini-key"
```

## ⚠️ 安全注意事項

1. **永遠不要**將 API 金鑰直接寫在代碼中
2. **永遠不要**將 API 金鑰提交到 Git 倉庫
3. 使用環境變數或 Secret Manager 存儲敏感信息
4. 定期輪換 API 金鑰
5. 監控 API 使用情況

## 📁 文件說明

- `.env` - 本地環境變數（已加入 .gitignore）
- `setup-env.sh` - 環境變數設置腳本
- `cloudbuild.yaml` - CI/CD 配置（不包含敏感信息）
