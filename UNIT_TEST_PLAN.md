# 單元測試規劃文檔

## 概述

本文檔記錄了整個系統的單元測試規劃，按照優先級和模組分類，確保系統的穩定性和安全性。

## 測試目標

- **核心業務邏輯覆蓋率**: 90%+
- **權限相關代碼覆蓋率**: 95%+
- **整體專案覆蓋率**: 80%+
- **安全模組覆蓋率**: 100%

## 測試框架選擇

```go
// 主要測試框架
github.com/stretchr/testify/assert    // 斷言庫
github.com/stretchr/testify/mock      // Mock 框架
github.com/stretchr/testify/suite     // 測試套件
github.com/DATA-DOG/go-sqlmock        // 資料庫 Mock
github.com/gin-gonic/gin              // HTTP 測試
```

## 測試目錄結構

```
tests/
├── unit/                    # 單元測試
│   ├── services/           # 服務層測試
│   ├── middleware/         # 中間件測試
│   ├── models/            # 模型層測試
│   ├── controllers/       # 控制器測試
│   └── utils/             # 工具函數測試
├── integration/           # 整合測試
├── fixtures/              # 測試數據
│   ├── test_data.json
│   ├── mock_responses.json
│   └── test_users.json
└── helpers/               # 測試輔助函數
    ├── test_db.go
    ├── mock_auth.go
    └── test_utils.go
```

---

## 🔥 第一優先級 - 核心安全模組

### 1. UnifiedAuthService 測試

**文件**: `tests/unit/services/unified_auth_service_test.go`

#### 測試案例

- [ ] **註冊功能測試**

  - [ ] 客戶註冊成功
  - [ ] 商戶註冊成功
  - [ ] 管理員註冊成功
  - [ ] 重複郵箱註冊失敗
  - [ ] 無效角色註冊失敗
  - [ ] 密碼長度驗證
  - [ ] 必填欄位驗證

- [ ] **登入功能測試**

  - [ ] 客戶登入成功
  - [ ] 商戶登入成功
  - [ ] 管理員登入成功
  - [ ] 錯誤密碼登入失敗
  - [ ] 不存在的用戶登入失敗
  - [ ] 停用帳戶登入失敗
  - [ ] 角色不匹配登入失敗

- [ ] **JWT Token 測試**

  - [ ] Token 生成正確性
  - [ ] Token 驗證成功
  - [ ] 過期 Token 驗證失敗
  - [ ] 無效 Token 驗證失敗
  - [ ] Token 內容正確性

- [ ] **密碼處理測試**
  - [ ] 密碼雜湊正確性
  - [ ] 密碼驗證正確性
  - [ ] 不同密碼雜湊不同

### 2. 權限中間件測試

**文件**: `tests/unit/middleware/auth_test.go`

#### 測試案例

- [ ] **AdminMiddleware 測試**

  - [ ] 管理員訪問允許
  - [ ] 非管理員訪問拒絕
  - [ ] 未認證用戶重定向
  - [ ] 無效用戶信息處理

- [ ] **CustomerMiddleware 測試**

  - [ ] 客戶訪問允許
  - [ ] 非客戶訪問拒絕
  - [ ] 未認證用戶重定向

- [ ] **MerchantMiddleware 測試**

  - [ ] 商戶訪問允許
  - [ ] 非商戶訪問拒絕
  - [ ] 未認證用戶重定向

- [ ] **MultiRoleMiddleware 測試**

  - [ ] 允許的角色訪問
  - [ ] 不允許的角色拒絕
  - [ ] 多角色配置正確性

- [ ] **UnifiedAuthMiddleware 測試**
  - [ ] Bearer Token 認證
  - [ ] Cookie Token 認證
  - [ ] Query Parameter Token 認證
  - [ ] 無效 Token 處理
  - [ ] 過期 Token 處理

---

## 🟡 第二優先級 - 資料層

### 3. 模型層測試

**文件**: `tests/unit/models/customer_test.go`
**文件**: `tests/unit/models/merchant_test.go`
**文件**: `tests/unit/models/admin_test.go`

#### 測試案例

- [ ] **Customer 模型測試**

  - [ ] 創建客戶成功
  - [ ] 查詢客戶成功
  - [ ] 更新客戶成功
  - [ ] 刪除客戶成功
  - [ ] 郵箱唯一性驗證
  - [ ] 必填欄位驗證

- [ ] **Merchant 模型測試**

  - [ ] 創建商戶成功
  - [ ] 商戶認證狀態管理
  - [ ] 商戶專用欄位驗證

- [ ] **Admin 模型測試**

  - [ ] 創建管理員成功
  - [ ] 管理員層級驗證
  - [ ] 權限級別檢查

- [ ] **Repository 測試**
  - [ ] 資料庫連接測試
  - [ ] 查詢語法正確性
  - [ ] 錯誤處理測試

### 4. 控制器測試

**文件**: `tests/unit/controllers/unified_auth_controller_test.go`

#### 測試案例

- [ ] **註冊端點測試**

  - [ ] POST /auth/register 成功
  - [ ] 請求驗證測試
  - [ ] 錯誤響應測試

- [ ] **登入端點測試**

  - [ ] POST /auth/customer/login 成功
  - [ ] POST /merchant/login 成功
  - [ ] POST /admin/login 成功
  - [ ] 錯誤響應測試

- [ ] **權限端點測試**
  - [ ] 受保護端點訪問
  - [ ] 權限檢查正確性

---

## 🟢 第三優先級 - 輔助功能

### 5. 管理員服務測試

**文件**: `tests/unit/services/unified_admin_service_test.go`

#### 測試案例

- [ ] **用戶管理測試**

  - [ ] 創建用戶
  - [ ] 更新用戶
  - [ ] 刪除用戶
  - [ ] 查詢用戶列表

- [ ] **權限管理測試**
  - [ ] 角色權限分配
  - [ ] 權限檢查

### 6. 聊天服務測試

**文件**: `tests/unit/services/chat_service_test.go`

#### 測試案例

- [ ] **聊天功能測試**
  - [ ] 創建聊天
  - [ ] 發送訊息
  - [ ] 查詢聊天記錄

---

## 🔧 測試輔助工具

### 1. 測試數據管理

**文件**: `tests/fixtures/test_data.json`

```json
{
  "users": {
    "customer": {
      "valid": [...],
      "invalid": [...]
    },
    "merchant": {
      "valid": [...],
      "invalid": [...]
    },
    "admin": {
      "valid": [...],
      "invalid": [...]
    }
  },
  "tokens": {
    "valid": "...",
    "expired": "...",
    "invalid": "..."
  }
}
```

### 2. Mock 服務

**文件**: `tests/helpers/mock_auth.go`

- 模擬認證服務
- 模擬資料庫操作
- 模擬外部 API 調用

### 3. 測試資料庫

**文件**: `tests/helpers/test_db.go`

- 測試專用資料庫設置
- 測試數據清理
- 資料庫連接管理

---

## 📋 實施計劃

### 階段一：核心安全模組 (第 1-2 週)

1. 設置測試框架和基礎設施
2. 實施 UnifiedAuthService 測試
3. 實施權限中間件測試
4. 達到 90% 覆蓋率

### 階段二：資料層測試 (第 3-4 週)

1. 實施模型層測試
2. 實施控制器測試
3. 添加整合測試
4. 達到 85% 覆蓋率

### 階段三：輔助功能測試 (第 5-6 週)

1. 實施管理員服務測試
2. 實施聊天服務測試
3. 完善測試文檔
4. 達到 80% 覆蓋率

### 階段四：優化和維護 (第 7-8 週)

1. 性能測試
2. 壓力測試
3. 測試覆蓋率優化
4. CI/CD 整合

---

## 🎯 成功指標

- [ ] 所有核心業務邏輯有完整測試覆蓋
- [ ] 權限相關代碼 100% 測試覆蓋
- [ ] 所有測試通過率 100%
- [ ] 測試執行時間 < 30 秒
- [ ] CI/CD 流程整合完成
- [ ] 測試文檔完整

---

## 📝 測試執行命令

```bash
# 執行所有測試
make test

# 執行特定模組測試
go test ./tests/unit/services/...

# 執行測試並生成覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 執行測試並顯示詳細輸出
go test -v ./...

# 執行測試並生成基準測試
go test -bench=. ./...
```

---

## 🔄 持續改進

- 每週檢查測試覆蓋率
- 定期更新測試案例
- 根據新功能添加測試
- 優化測試執行速度
- 改進測試可讀性

---

**最後更新**: 2024-01-02
**負責人**: 開發團隊
**狀態**: 規劃中
