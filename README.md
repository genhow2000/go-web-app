# Go 簡單應用

這是一個完整的 Go 語言 Web 應用，包含 PostgreSQL 資料庫，使用 Docker 進行容器化部署。

## 功能

- 提供基本的 HTTP 服務器
- 首頁 API 端點 (`/`)
- 健康檢查端點 (`/health`)
- 用戶管理 API (`/users`)
- PostgreSQL 資料庫整合
- 支援環境變數配置

## 文件結構

```
go/
├── main.go              # 主應用程序（包含資料庫操作）
├── go.mod              # Go 模組文件
├── go.sum              # 依賴校驗文件
├── Dockerfile          # Docker 構建文件
├── docker-compose.yml  # Docker Compose 配置（包含 PostgreSQL）
└── README.md           # 說明文件
```

## 技術棧

- **後端**: Go 1.21
- **資料庫**: PostgreSQL 15
- **容器化**: Docker & Docker Compose
- **資料庫驅動**: github.com/lib/pq

## 快速開始

### 使用 Docker Compose（推薦）

1. 進入 go 目錄：

   ```bash
   cd go
   ```

2. 構建並啟動應用：

   ```bash
   docker-compose up --build
   ```

3. 訪問應用：

   - 首頁：http://localhost:8080
   - 健康檢查：http://localhost:8080/health
   - 用戶列表：http://localhost:8080/users

4. 測試用戶 API：

   ```bash
   # 創建新用戶
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name": "張三", "email": "zhang@example.com"}'

   # 獲取所有用戶
   curl http://localhost:8080/users
   ```

### 使用 Docker

1. 構建映像：

   ```bash
   docker build -t go-simple-app .
   ```

2. 運行容器：
   ```bash
   docker run -p 8080:8080 go-simple-app
   ```

### 本地開發

1. 確保已安裝 Go 1.21 或更高版本

2. 運行應用：
   ```bash
   go run main.go
   ```

## API 端點

### GET /

返回歡迎訊息和應用狀態。

**響應示例：**

```json
{
  "message": "歡迎來到 Go 服務器！",
  "status": "running",
  "version": "1.0.0"
}
```

### GET /health

健康檢查端點，包含資料庫連接狀態。

**響應示例：**

```json
{
  "status": "healthy",
  "service": "go-simple-app",
  "database": "connected"
}
```

### GET /users

獲取所有用戶列表。

**響應示例：**

```json
{
  "users": [
    {
      "id": 1,
      "name": "張三",
      "email": "zhang@example.com",
      "created_at": "2024-01-01T10:00:00Z"
    }
  ],
  "count": 1
}
```

### POST /users

創建新用戶。

**請求體：**

```json
{
  "name": "張三",
  "email": "zhang@example.com"
}
```

**響應示例：**

```json
{
  "message": "用戶創建成功",
  "user": {
    "id": 1,
    "name": "張三",
    "email": "zhang@example.com"
  }
}
```

## 環境變數

### 應用配置

- `PORT`: 服務器端口（默認：8080）

### 資料庫配置

- `DB_HOST`: 資料庫主機（默認：localhost）
- `DB_PORT`: 資料庫端口（默認：5432）
- `DB_USER`: 資料庫用戶名（默認：postgres）
- `DB_PASSWORD`: 資料庫密碼（默認：password）
- `DB_NAME`: 資料庫名稱（默認：goapp）

## 停止應用

使用 Docker Compose：

```bash
docker-compose down
```

使用 Docker：

```bash
docker stop go-simple-app
docker stop go-postgres
```

## 資料庫管理

### 連接資料庫

```bash
# 使用 Docker 連接 PostgreSQL
docker exec -it go-postgres psql -U postgres -d goapp
```

### 查看資料表

```sql
-- 查看所有表
\dt

-- 查看用戶表結構
\d users

-- 查看用戶數據
SELECT * FROM users;
```

### 重置資料庫

```bash
# 停止並刪除所有容器和數據
docker-compose down -v

# 重新啟動
docker-compose up --build
```
