# 雲端資料庫解決方案

## 🚨 目前問題

你的 Go 專案目前使用 SQLite 存儲在 `/tmp/app.db`，這在 Cloud Run 環境中會導致：

1. **資料遺失**：每次容器重啟都會清空 `/tmp` 目錄
2. **無法持久化**：Cloud Run 是無狀態的，不適合存儲資料
3. **無法擴展**：多個實例無法共享資料

## 💡 推薦解決方案

### 方案 1：Google Cloud SQL (推薦)

**優點：**

- 完全託管，無需維護
- 自動備份和恢復
- 高可用性和擴展性
- 與 Cloud Run 完美整合

**實施步驟：**

1. **創建 Cloud SQL 實例：**

```bash
gcloud sql instances create go-app-db \
  --database-version=POSTGRES_13 \
  --tier=db-f1-micro \
  --region=asia-east1 \
  --storage-type=SSD \
  --storage-size=10GB
```

2. **創建資料庫：**

```bash
gcloud sql databases create goapp --instance=go-app-db
```

3. **創建用戶：**

```bash
gcloud sql users create goapp \
  --instance=go-app-db \
  --password=your-secure-password
```

4. **更新應用配置：**

```go
// 在 config/config.go 中
Database: DatabaseConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnv("DB_PORT", "5432"),
    User:     getEnv("DB_USER", "goapp"),
    Password: getEnv("DB_PASSWORD", ""),
    DBName:   getEnv("DB_NAME", "goapp"),
    SSLMode:  getEnv("DB_SSLMODE", "require"),
},
```

### 方案 2：使用 Cloud Storage + SQLite

**優點：**

- 保持現有 SQLite 代碼
- 成本較低
- 簡單實施

**缺點：**

- 不支援並發寫入
- 性能較差
- 需要額外的同步邏輯

**實施步驟：**

1. **創建 Cloud Storage Bucket：**

```bash
gsutil mb gs://your-project-go-app-db
```

2. **修改資料庫初始化邏輯：**

```go
// 在應用啟動時從 Cloud Storage 下載資料庫
// 在應用關閉時上傳資料庫到 Cloud Storage
```

### 方案 3：使用 Cloud Firestore

**優點：**

- NoSQL，靈活性高
- 自動擴展
- 與 Google Cloud 生態系統整合好

**缺點：**

- 需要重寫資料模型
- 學習成本較高
- 查詢語法不同

## 🔧 推薦實施：Cloud SQL + PostgreSQL

### 1. 更新資料庫驅動

在 `go.mod` 中添加：

```go
require (
    github.com/lib/pq v1.10.9
)
```

### 2. 更新資料庫連接

```go
// database/database.go
import (
    _ "github.com/lib/pq"
)

func Init() error {
    // 從環境變數構建連接字串
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSLMODE"),
    )

    DB, err = sql.Open("postgres", connStr)
    // ... 其餘邏輯
}
```

### 3. 更新 Cloud Build 配置

```yaml
# cloudbuild.yaml
steps:
  # ... 構建步驟 ...

  # 部署到 Cloud Run
  - name: "gcr.io/google.com/cloudsdktool/cloud-sdk"
    entrypoint: gcloud
    args:
      - "run"
      - "deploy"
      - "go-app"
      - "--image"
      - "gcr.io/fleet-day-383710/go-app:latest"
      - "--platform"
      - "managed"
      - "--region"
      - "asia-east1"
      - "--allow-unauthenticated"
      - "--port"
      - "8080"
      - "--set-env-vars"
      - "DB_HOST=your-cloud-sql-ip,DB_PORT=5432,DB_USER=goapp,DB_PASSWORD=your-password,DB_NAME=goapp,DB_SSLMODE=require"
      - "--add-cloudsql-instances"
      - "your-project:asia-east1:go-app-db"
```

## 📊 成本比較

| 方案                   | 月成本 | 複雜度 | 性能 | 推薦度     |
| ---------------------- | ------ | ------ | ---- | ---------- |
| Cloud SQL (PostgreSQL) | $25-50 | 中     | 高   | ⭐⭐⭐⭐⭐ |
| Cloud Storage + SQLite | $5-10  | 高     | 中   | ⭐⭐       |
| Cloud Firestore        | $10-30 | 高     | 高   | ⭐⭐⭐     |

## 🚀 快速開始

1. **選擇 Cloud SQL 方案**
2. **創建 Cloud SQL 實例**
3. **更新應用代碼**
4. **測試本地連接**
5. **部署到 Cloud Run**

這樣就能解決雲端資料庫持久化的問題了！
