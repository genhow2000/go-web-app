# API 服務

## 概述

RESTful API 設計，支援前後端分離，提供完整的 API 服務和文檔。

## 🚀 功能特色

- **RESTful 設計**：標準的 REST API 設計
- **前後端分離**：完全的前後端分離架構
- **統一響應格式**：標準化的 API 響應格式
- **錯誤處理**：完整的錯誤處理機制
- **API 文檔**：完整的 API 文檔

## 📋 API 設計原則

### RESTful 規範

- 使用標準 HTTP 方法
- 資源導向的 URL 設計
- 統一的狀態碼使用
- 無狀態的 API 設計

### 響應格式

```json
{
  "status": "success",
  "message": "操作成功",
  "data": {
    "id": 1,
    "name": "用戶名稱"
  },
  "timestamp": "2025-09-11T05:22:05Z"
}
```

## 🛠️ 核心 API

### 認證 API

- `POST /auth/login` - 用戶登入
- `POST /auth/register` - 用戶註冊
- `POST /auth/logout` - 用戶登出

### 用戶管理 API

- `GET /users` - 獲取用戶列表
- `GET /users/:id` - 獲取特定用戶
- `POST /users` - 創建用戶
- `PUT /users/:id` - 更新用戶
- `DELETE /users/:id` - 刪除用戶

### 管理員 API

- `GET /admin/api/users` - 管理員獲取用戶列表
- `POST /admin/api/users` - 管理員創建用戶
- `PUT /admin/api/users/:id` - 管理員更新用戶
- `DELETE /admin/api/users/:id` - 管理員刪除用戶

### 系統 API

- `GET /health` - 健康檢查
- `GET /api/stats` - 系統統計

## 🔐 認證機制

### JWT Token 認證

```http
Authorization: Bearer <jwt_token>
```

### 認證流程

1. 用戶登入獲取 token
2. 後續請求攜帶 token
3. 服務器驗證 token
4. 返回請求結果

## 📊 響應狀態碼

### 成功響應

- `200 OK` - 請求成功
- `201 Created` - 資源創建成功
- `204 No Content` - 請求成功無內容

### 客戶端錯誤

- `400 Bad Request` - 請求參數錯誤
- `401 Unauthorized` - 未認證
- `403 Forbidden` - 權限不足
- `404 Not Found` - 資源不存在

### 服務器錯誤

- `500 Internal Server Error` - 服務器內部錯誤
- `502 Bad Gateway` - 網關錯誤
- `503 Service Unavailable` - 服務不可用

## 🛡️ 安全特性

### 輸入驗證

- 參數類型驗證
- 數據格式驗證
- 長度限制驗證
- 特殊字符過濾

### 安全防護

- SQL 注入防護
- XSS 攻擊防護
- CSRF 防護
- 速率限制

## 📈 性能優化

### 緩存機制

- 響應數據緩存
- 查詢結果緩存
- 靜態資源緩存
- 分佈式緩存

### 分頁查詢

- 支援分頁參數
- 總數統計
- 排序功能
- 篩選功能

## 🔧 中間件

### 認證中間件

```go
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 驗證 JWT token
        // 設置用戶信息到上下文
    }
}
```

### 日誌中間件

```go
func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 記錄請求日誌
        // 記錄響應時間
    }
}
```

## 📝 API 文檔

### 自動生成

- 基於代碼註釋生成
- 實時更新文檔
- 互動式 API 測試
- 多格式文檔導出

### 文檔內容

- API 端點說明
- 請求參數描述
- 響應格式示例
- 錯誤碼說明

## 🚀 未來擴展

- 支援 GraphQL API
- 支援 WebSocket 實時通信
- 支援 API 版本控制
- 支援 API 限流
- 支援 API 分析

## 📊 使用範例

### 用戶登入

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"111111"}'
```

### 獲取用戶列表

```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer <jwt_token>"
```

### 創建用戶

```bash
curl -X POST http://localhost:8080/users \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"新用戶","email":"new@example.com","password":"123456"}'
```

---

**這個 API 服務展現了現代 Web 應用的 API 設計理念和 RESTful 架構的最佳實踐！**
