# 用戶管理系統

## 概述

基於 Go + Gin 框架開發的完整用戶管理系統，支援用戶註冊、登入、權限管理等功能。

## 🚀 功能特色

- **用戶註冊**：支援新用戶註冊功能
- **用戶登入**：安全的用戶認證系統
- **權限管理**：基於角色的權限控制
- **用戶管理**：完整的 CRUD 操作

## 👥 用戶角色

- **admin**：系統管理員，擁有所有權限
- **merchant**：商戶用戶，擁有商戶相關權限
- **customer**：一般用戶，基本權限

## 🔐 安全認證

### JWT Token 認證

```go
// 生成 JWT Token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(secretKey))
```

### 密碼加密

```go
// 使用 bcrypt 加密密碼
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

## 📊 資料庫設計

### 用戶表結構

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    is_active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🛠️ API 接口

### 認證接口

- `POST /register` - 用戶註冊
- `POST /login` - 用戶登入
- `POST /logout` - 用戶登出

### 用戶管理接口

- `GET /users` - 獲取所有用戶
- `GET /users/:id` - 獲取特定用戶
- `POST /users` - 創建用戶
- `PUT /users/:id` - 更新用戶
- `DELETE /users/:id` - 刪除用戶

## 🎨 分離式登入

### 商戶登入

- 路徑：`/merchant/login`
- 角色：merchant
- 儀表板：`/merchant/dashboard`

### 管理員登入

- 路徑：`/admin/login`
- 角色：admin
- 儀表板：`/admin/dashboard`

## 🔒 安全特性

- **JWT 認證**：安全的 token 認證機制
- **密碼加密**：使用 bcrypt 加密存儲
- **角色權限**：基於角色的訪問控制
- **會話管理**：安全的會話管理
- **輸入驗證**：完整的輸入驗證

## 📈 性能優化

- **資料庫索引**：email 字段唯一索引
- **連接池**：資料庫連接池管理
- **緩存機制**：用戶信息緩存
- **分頁查詢**：支援分頁查詢

## 🎯 技術亮點

- **分離式架構**：商戶和管理員分離登入
- **中間件設計**：認證和權限中間件
- **Repository 模式**：資料訪問層抽象
- **Service 層**：業務邏輯層分離

## 🚀 未來擴展

- 支援 OAuth 登入
- 支援多因素認證
- 支援用戶群組管理
- 支援權限細粒度控制
- 支援用戶行為分析

## 📝 使用範例

### 用戶註冊

```json
POST /register
{
    "name": "測試用戶",
    "email": "test@example.com",
    "password": "123456"
}
```

### 用戶登入

```json
POST /login
{
    "email": "test@example.com",
    "password": "123456"
}
```

### 獲取用戶列表

```
GET /users
Authorization: Bearer <jwt_token>
```

---

**這個用戶管理系統展現了完整的前後端分離架構設計和安全的認證機制實現！**
