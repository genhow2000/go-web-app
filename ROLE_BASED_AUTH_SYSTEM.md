# 基於角色的權限系統設計

## 概述

本系統採用三張獨立表的方式來管理三種不同角色的用戶，實現了完整的權限分離和數據隔離。

## 資料庫設計

### 三張獨立表

1. **customers** - 客戶表

   - 基本字段：id, name, email, password, phone, address
   - 客戶專用字段：birth_date, gender, email_verified
   - 系統字段：is_active, last_login, login_count, profile_data, created_at, updated_at

2. **merchants** - 商戶表

   - 基本字段：id, name, email, password, phone, address
   - 商戶專用字段：business_name, business_license, business_type, is_verified
   - 系統字段：is_active, last_login, login_count, business_data, created_at, updated_at

3. **admins** - 管理員表
   - 基本字段：id, name, email, password, phone
   - 管理員專用字段：admin_level (normal/senior/super), department
   - 系統字段：is_active, last_login, login_count, admin_data, created_at, updated_at

### 輔助表

4. **login_logs** - 統一登入日誌表

   - 記錄所有角色的登入活動
   - 字段：user_type, user_id, login_time, ip_address, user_agent, success

5. **role_permissions** - 角色權限表
   - 定義每種角色的權限
   - 字段：role, description, permissions (JSON)

## 模型層設計

### 統一接口

- `UserInterface` - 所有用戶類型的統一接口
- `UnifiedUserRepository` - 統一用戶倉庫，支持跨表查詢

### 角色專用模型

- `Customer` + `CustomerRepository` - 客戶模型
- `Merchant` + `MerchantRepository` - 商戶模型
- `Admin` + `AdminRepository` - 管理員模型

## 認證服務

### UnifiedAuthService

- 統一的認證服務，支持三種角色的註冊和登入
- 自動根據角色選擇對應的數據表
- 統一的 JWT token 生成和驗證

## 權限中間件

### 角色專用中間件

- `CustomerMiddleware()` - 只允許 customer 角色
- `MerchantMiddleware()` - 只允許 merchant 角色
- `AdminMiddleware()` - 只允許 admin 角色
- `MultiRoleMiddleware(allowedRoles...)` - 允許多個角色

## 路由設計

### 登入路由

- `/auth/customer/login` - 客戶登入
- `/merchant/login` - 商戶登入
- `/admin/login` - 管理員登入

### 受保護路由

- `/customer/*` - 客戶專用路由（需要 customer 權限）
- `/merchant/*` - 商戶專用路由（需要 merchant 權限）
- `/admin/*` - 管理員專用路由（需要 admin 權限）

### 共享路由

- `/api/chat/*` - 聊天功能（所有角色都可使用）

## 頁面設計

### 客戶頁面

- `customer_login.html` - 客戶登入頁面
- `customer_dashboard.html` - 客戶儀表板

### 商戶頁面

- `merchant_login.html` - 商戶登入頁面
- `merchant_dashboard.html` - 商戶儀表板

### 管理員頁面

- `admin_login.html` - 管理員登入頁面
- `admin_dashboard.html` - 管理員儀表板

## 權限級別

### 客戶權限

- 查看個人資料
- 更新個人資料
- 創建和查看對話
- 下單和查看訂單

### 商戶權限

- 客戶所有權限
- 管理商品
- 查看和管理訂單
- 查看分析報告
- 管理商戶信息

### 管理員權限

- 所有權限
- 管理所有用戶
- 系統管理
- 查看系統統計
- 管理角色和權限

## 安全特性

1. **數據隔離** - 每種角色的數據完全分離
2. **權限分離** - 每種角色有獨立的權限檢查
3. **登入日誌** - 完整的登入活動記錄
4. **JWT 認證** - 安全的 token 認證機制
5. **密碼加密** - bcrypt 密碼雜湊

## 使用方式

### 註冊用戶

```json
POST /auth/register
{
  "name": "用戶名",
  "email": "user@example.com",
  "password": "password123",
  "role": "customer", // 或 "merchant" 或 "admin"
  "phone": "0912345678", // 可選
  "address": "地址", // 可選
  // 角色專用字段...
}
```

### 登入

```json
POST /auth/customer/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

### 訪問受保護資源

- 在請求頭中添加：`Authorization: Bearer <token>`
- 或使用 cookie：`auth_token=<token>`

## 遷移文件

1. `001_add_role_and_status.sql` - 添加角色和狀態字段
2. `002_add_merchant_role.sql` - 添加 merchant 角色支持
3. `003_create_separate_role_tables.sql` - 創建三張獨立表

## 優勢

1. **清晰的業務邏輯** - 每種角色有獨立的業務流程
2. **更好的安全性** - 數據和權限完全分離
3. **靈活的擴展性** - 可以為每種角色添加專用字段
4. **統一的認證** - 使用統一的認證服務和接口
5. **完整的日誌** - 詳細的登入和操作日誌

這個設計確保了每種角色都有獨立的數據存儲和權限管理，同時保持了代碼的統一性和可維護性。
