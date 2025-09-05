# 用戶管理系統使用指南

## 概述

本系統已擴展為完整的用戶管理系統，支持客戶和管理員兩種角色，提供完整的用戶註冊、登入、管理功能。

## 功能特性

### 🔐 認證系統

- 用戶註冊（支持角色選擇）
- 用戶登入（支持角色驗證）
- JWT Token 認證
- 密碼加密存儲
- 帳戶狀態管理（啟用/停用）

### 👥 角色管理

- **客戶 (Customer)**: 基本用戶權限
- **管理員 (Admin)**: 完整管理權限

### 🛡️ 權限控制

- 基於角色的訪問控制
- 管理員專用功能
- 中間件權限驗證

### 📊 管理後台

- 用戶統計儀表板
- 用戶列表管理
- 用戶創建/編輯/刪除
- 用戶狀態管理
- 角色管理

## 快速開始

### 1. 數據庫遷移

首先執行數據庫遷移以添加角色和狀態字段：

```bash
# 編譯遷移工具
go build -o migrate cmd/migrate/main.go

# 執行遷移
./migrate
```

### 2. 創建管理員帳戶

```bash
# 編譯管理員初始化工具
go build -o init-admin cmd/init-admin/main.go

# 創建管理員帳戶
./init-admin
```

### 3. 啟動服務器

```bash
go run main.go
```

## 使用指南

### 用戶註冊

1. 訪問 `/register` 頁面
2. 填寫用戶信息：
   - 姓名
   - 電子郵件
   - 密碼（至少 6 個字符）
   - 確認密碼
   - 註冊類型（客戶/管理員）
3. 點擊註冊按鈕

### 用戶登入

1. 訪問 `/login` 頁面
2. 輸入電子郵件和密碼
3. 系統會根據角色重定向到相應的儀表板

### 管理員功能

管理員登入後可以：

1. **訪問管理後台**: 點擊側邊欄的"管理後台"按鈕
2. **查看統計**: 在管理儀表板查看用戶統計
3. **管理用戶**:
   - 查看所有用戶列表
   - 創建新用戶
   - 編輯用戶信息
   - 更改用戶角色
   - 啟用/停用用戶
   - 刪除用戶

## API 端點

### 認證端點

- `POST /register` - 用戶註冊
- `POST /login` - 用戶登入
- `POST /logout` - 用戶登出

### 管理員 API 端點

- `GET /admin/api/users` - 獲取所有用戶
- `GET /admin/api/users/role/:role` - 根據角色獲取用戶
- `GET /admin/api/users/:id` - 獲取特定用戶
- `POST /admin/api/users` - 創建用戶
- `PUT /admin/api/users/:id` - 更新用戶
- `PUT /admin/api/users/:id/status` - 更新用戶狀態
- `PUT /admin/api/users/:id/role` - 更新用戶角色
- `DELETE /admin/api/users/:id` - 刪除用戶
- `GET /admin/api/stats` - 獲取用戶統計

## 數據庫結構

### users 表

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'customer' NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_role CHECK (role IN ('customer', 'admin'))
);
```

## 安全特性

1. **密碼加密**: 使用 bcrypt 進行密碼雜湊
2. **JWT 認證**: 安全的 token 認證機制
3. **角色驗證**: 基於角色的權限控制
4. **輸入驗證**: 所有輸入都經過驗證
5. **SQL 注入防護**: 使用參數化查詢

## 配置

### 環境變量

```bash
# 服務器配置
PORT=8080
HOST=0.0.0.0

# 數據庫配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=goapp
DB_SSLMODE=disable

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24
```

## 故障排除

### 常見問題

1. **數據庫連接失敗**

   - 檢查數據庫服務是否運行
   - 驗證數據庫配置信息

2. **遷移失敗**

   - 確保數據庫用戶有足夠權限
   - 檢查遷移文件語法

3. **認證失敗**

   - 檢查 JWT_SECRET 配置
   - 驗證 token 是否過期

4. **權限不足**
   - 確認用戶角色設置正確
   - 檢查中間件配置

## 開發指南

### 添加新角色

1. 更新 `models/user.go` 中的角色檢查
2. 更新 `services/auth_service.go` 中的角色驗證
3. 更新數據庫約束
4. 更新前端界面

### 添加新權限

1. 在 `middleware/auth.go` 中添加新的中間件
2. 在路由中應用相應的中間件
3. 更新前端權限檢查

## 技術棧

- **後端**: Go, Gin, PostgreSQL
- **前端**: HTML, CSS, JavaScript, Bootstrap
- **認證**: JWT, bcrypt
- **數據庫**: PostgreSQL
- **日誌**: Logrus

## 版本歷史

- **v2.0.0**: 添加完整的用戶管理系統
  - 角色管理（客戶/管理員）
  - 管理後台
  - 權限控制
  - 用戶管理 API

## 支持

如有問題或建議，請聯繫開發團隊。
