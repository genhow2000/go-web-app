# AI 智能聊天系統

## 🤖 系統概述

本系統整合了多個 AI 服務提供商，提供智能對話功能，支援匿名用戶和認證用戶使用。

## ✨ 主要功能

### 1. 多 AI 提供商支援

- **主要提供商**：Groq API (Llama 3.1 8B)
- **備用提供商**：Google Gemini 2.0 Flash
- **模擬服務**：本地模擬 AI（當 API 不可用時）

### 2. 智能服務切換

- 自動檢測 API 可用性
- 主要服務失敗時自動切換到備用服務
- 支援錯誤處理和重試機制

### 3. 用戶支援

- **匿名用戶**：無需註冊即可使用
- **認證用戶**：支援對話歷史記錄
- **多角色支援**：客戶、商戶、管理員

## 🛠️ 技術架構

### 後端架構

```
services/
├── ai_interface.go      # AI 服務接口定義
├── ai_manager.go        # AI 服務管理器
├── groq_service.go      # Groq API 服務
├── gemini_service.go    # Gemini API 服務
├── simulation_service.go # 模擬 AI 服務
└── chat_service.go      # 聊天業務邏輯
```

### 資料庫設計

```sql
-- 對話記錄表
CREATE TABLE conversations (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    title TEXT NOT NULL,
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 消息記錄表
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    role TEXT NOT NULL, -- 'user' or 'assistant'
    content TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id)
);
```

## 🔧 配置說明

### 環境變數

```bash
# AI 服務配置
GROQ_API_KEY=your_groq_api_key
GEMINI_API_KEY=your_gemini_api_key
AI_PRIMARY_PROVIDER=groq
AI_FALLBACK_PROVIDER=gemini

# 服務限制
GROQ_DAILY_LIMIT=10000
GEMINI_DAILY_LIMIT=1500
```

### API 端點

```
POST /api/chat/conversations     # 創建新對話
POST /api/chat/send             # 發送消息
GET  /api/chat/conversations    # 獲取用戶對話列表（需認證）
GET  /api/chat/conversations/:id # 獲取特定對話（需認證）
DELETE /api/chat/conversations/:id # 刪除對話（需認證）
```

## 📊 使用統計

### 限制配置

- **匿名用戶**：每分鐘 5 次請求，每日 50 次
- **認證用戶**：每分鐘 10 次請求，每日 100 次
- **管理員**：無限制

### 監控功能

- 實時使用統計
- API 錯誤率監控
- 服務切換記錄
- 用戶行為分析

## 🚀 部署指南

### 1. 本地開發

```bash
# 設置環境變數
export GROQ_API_KEY="your-groq-key"
export GEMINI_API_KEY="your-gemini-key"

# 啟動服務
go run main.go
```

### 2. Docker 部署

```bash
# 使用 docker-compose
docker-compose up -d
```

### 3. 雲端部署

```bash
# 更新 Cloud Run 環境變數
gcloud run services update go-app \
  --region=asia-east1 \
  --set-env-vars="GROQ_API_KEY=your-key,GEMINI_API_KEY=your-key"
```

## 🔒 安全考量

### API 金鑰管理

- 使用環境變數存儲敏感信息
- 支援 Google Cloud Secret Manager
- 定期輪換 API 金鑰

### 速率限制

- 防止 API 濫用
- 支援不同用戶等級的限制
- 自動封鎖異常行為

### 資料隱私

- 支援匿名使用
- 可選的對話歷史記錄
- 資料加密存儲

## 📈 性能優化

### 快取策略

- 常用回應快取
- 用戶會話快取
- API 回應快取

### 並發處理

- Goroutine 池管理
- 請求隊列機制
- 超時控制

### 監控指標

- 回應時間
- 錯誤率
- 吞吐量
- 資源使用率

## 🧪 測試功能

### 單元測試

```bash
go test ./services/...
```

### 整合測試

```bash
# 測試 AI 服務切換
curl -X POST /api/chat/send \
  -H "Content-Type: application/json" \
  -d '{"message":"測試消息","conversation_id":"test123"}'
```

### 壓力測試

```bash
# 使用 Apache Bench 進行壓力測試
ab -n 1000 -c 10 -p test.json -T application/json http://localhost:8080/api/chat/send
```

## 🔮 未來規劃

### 功能增強

- 支援更多 AI 提供商
- 自定義 AI 模型
- 語音對話功能
- 多語言支援

### 性能優化

- 流式回應
- 更智能的快取策略
- 負載均衡

### 監控增強

- 詳細的日誌記錄
- 實時監控面板
- 告警機制
