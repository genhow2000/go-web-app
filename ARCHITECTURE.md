# Go 應用架構說明

## 🏗️ 架構概述

這個 Go 應用現在使用了類似 Laravel 的 MVC 架構，並整合了現代 Go 語言的最佳實踐。

## 📁 目錄結構

```
go/
├── config/                 # 配置管理
│   └── config.go          # 應用配置
├── controllers/           # 控制器層
│   ├── auth_controller.go # 認證控制器
│   └── user_controller.go # 用戶控制器
├── database/              # 資料庫層
│   └── database.go       # 資料庫連接管理
├── middleware/            # 中間件
│   ├── auth.go           # 認證中間件
│   └── cors.go           # CORS 中間件
├── models/                # 模型層
│   └── user.go           # 用戶模型和 Repository
├── routes/                # 路由管理
│   └── routes.go         # 路由配置
├── services/              # 服務層
│   ├── auth_service.go   # 認證服務
│   └── user_service.go   # 用戶服務
├── templates/             # HTML 模板
│   ├── login.html        # 登入頁面
│   ├── register.html     # 註冊頁面
│   └── dashboard.html    # 儀表板
├── main.go               # 主程序入口
├── go.mod                # Go 模組文件
├── Dockerfile            # Docker 配置
└── docker-compose.yml    # Docker Compose 配置
```

## 🔧 技術棧

### 核心框架

- **Gin**: 高性能 HTTP Web 框架
- **PostgreSQL**: 關係型資料庫
- **JWT**: JSON Web Token 認證
- **bcrypt**: 密碼雜湊

### 依賴管理

- `github.com/gin-gonic/gin` - Web 框架
- `github.com/lib/pq` - PostgreSQL 驅動
- `github.com/golang-jwt/jwt/v5` - JWT 處理
- `golang.org/x/crypto` - 密碼雜湊
- `github.com/joho/godotenv` - 環境變數管理

## 🏛️ 架構模式

### 1. MVC 模式

- **Model**: `models/` - 資料模型和 Repository 模式
- **View**: `templates/` - HTML 模板
- **Controller**: `controllers/` - 請求處理邏輯

### 2. 服務層模式

- **Service Layer**: `services/` - 業務邏輯處理
- 分離控制器和業務邏輯
- 提供可重用的業務功能

### 3. Repository 模式

- 封裝資料庫操作
- 提供統一的資料存取介面
- 便於測試和維護

### 4. 中間件模式

- **認證中間件**: 保護需要認證的路由
- **CORS 中間件**: 處理跨域請求
- **日誌中間件**: 請求日誌記錄

## 🚀 主要功能

### 認證系統

- 用戶註冊/登入
- JWT Token 認證
- 密碼 bcrypt 雜湊
- 會話管理

### 用戶管理

- CRUD 操作
- 用戶列表
- 用戶詳情
- 受保護的路由

### API 端點

#### 公開端點

- `GET /` - 首頁
- `GET /health` - 健康檢查
- `GET /login` - 登入頁面
- `GET /register` - 註冊頁面
- `POST /login` - 登入 API
- `POST /register` - 註冊 API

#### 受保護端點

- `GET /dashboard` - 儀表板
- `GET /users` - 用戶列表
- `GET /users/:id` - 用戶詳情
- `POST /users` - 創建用戶
- `PUT /users/:id` - 更新用戶
- `DELETE /users/:id` - 刪除用戶

## 🔒 安全特性

1. **密碼安全**: 使用 bcrypt 雜湊密碼
2. **JWT 認證**: 安全的 token 認證機制
3. **CORS 支援**: 跨域請求處理
4. **輸入驗證**: Gin 的內建驗證
5. **SQL 注入防護**: 使用參數化查詢

## 🐳 Docker 部署

### 多階段構建

- 構建階段: 使用 `golang:1.21-alpine`
- 運行階段: 使用 `alpine:latest`
- 最小化映像大小

### 服務配置

- Go 應用: 端口 8080
- PostgreSQL: 端口 5432
- 資料持久化: Docker volumes

## 📊 與 Laravel 的對比

| 功能   | Laravel      | Go + Gin   |
| ------ | ------------ | ---------- |
| 框架   | Laravel      | Gin        |
| 路由   | Route::get() | r.GET()    |
| 控制器 | Controller   | Controller |
| 模型   | Eloquent     | Repository |
| 服務   | Service      | Service    |
| 中間件 | Middleware   | Middleware |
| 認證   | Auth         | JWT        |
| 資料庫 | Eloquent ORM | 原生 SQL   |
| 模板   | Blade        | HTML + JS  |

## 🚀 啟動方式

```bash
# 使用 Docker Compose
docker-compose up --build

# 訪問應用
http://localhost:8080
```

## 🔄 微服務架構建議

如果要進一步發展為微服務架構，可以考慮：

1. **Go-kit**: 微服務工具包
2. **gRPC**: 高效能的 RPC 框架
3. **Consul/Eureka**: 服務發現
4. **Kubernetes**: 容器編排
5. **Istio**: 服務網格

## 📈 性能優勢

- **高並發**: Go 的 goroutine 支援
- **低延遲**: 編譯型語言的性能
- **小記憶體**: 高效的記憶體使用
- **快速啟動**: 容器啟動時間短

這個架構提供了類似 Laravel 的開發體驗，同時發揮了 Go 語言的性能優勢。
