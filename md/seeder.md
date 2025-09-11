# 自製 Seeder 系統

## 概述

由於 Go 語言沒有內建的 Seeder 系統，我們完全自製了一個模組化的 Seeder 管理系統，用於自動生成和管理測試數據。

## 🚀 功能特色

- **模組化設計**：每個 seeder 都是獨立的模組
- **重複執行保護**：避免重複創建測試數據
- **執行記錄追蹤**：詳細的執行日誌
- **靈活配置**：支援單獨執行或批量執行
- **清除功能**：支援清除特定或所有測試數據

## 📁 目錄結構

```
database/seeders/
├── README.md           # 本文檔
├── manager.go          # Seeder 管理器
├── user_seeder.go      # 用戶數據 seeder
└── product_seeder.go   # 產品數據 seeder (示例)
```

## 🛠️ 核心組件

### 1. Seeder 接口

```go
type Seeder interface {
    Run() error    // 執行 seeder
    Clear() error  // 清除數據
}
```

### 2. SeederManager

負責管理所有 seeder 的執行：

```go
type SeederManager struct {
    db      *sql.DB
    seeders []Seeder
}
```

### 3. 個別 Seeder

每個 seeder 都實現 Seeder 接口：

- **UserSeeder**：創建測試用戶數據
- **ProductSeeder**：創建測試產品數據（示例）

## 📋 使用方法

### 命令行工具

```bash
# 執行所有 seeder
./seed run

# 執行特定 seeder
./seed run user

# 清除所有測試數據
./seed clear

# 清除特定 seeder 的數據
./seed clear user
```

### 程式碼中調用

```go
// 創建 seeder 管理器
seederManager := seeders.NewSeederManager(db)

// 執行所有 seeder
err := seederManager.RunAll()

// 執行特定 seeder
err := seederManager.RunSpecific("user")

// 清除所有測試數據
err := seederManager.ClearAll()
```

## 🔧 創建新的 Seeder

### 1. 創建 seeder 文件

```go
package seeders

import (
    "database/sql"
    "log"
)

type MySeeder struct {
    db *sql.DB
}

func NewMySeeder(db *sql.DB) *MySeeder {
    return &MySeeder{db: db}
}

func (s *MySeeder) Run() error {
    // 實現數據創建邏輯
    log.Println("執行 MySeeder...")
    return nil
}

func (s *MySeeder) Clear() error {
    // 實現數據清除邏輯
    log.Println("清除 MySeeder 數據...")
    return nil
}
```

### 2. 註冊到管理器

在 `manager.go` 中註冊：

```go
func (sm *SeederManager) registerSeeders() {
    sm.seeders = append(sm.seeders, NewUserSeeder(sm.db))
    sm.seeders = append(sm.seeders, NewProductSeeder(sm.db))
    sm.seeders = append(sm.seeders, NewMySeeder(sm.db)) // 新增
}
```

## 📊 實際範例

### UserSeeder 範例

```go
func (s *UserSeeder) Run() error {
    // 檢查是否已經有測試用戶
    var count int
    err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email IN ('admin@example.com', 'merchant@example.com', 'user@example.com')").Scan(&count)
    if err != nil {
        return err
    }

    if count > 0 {
        log.Println("測試用戶已存在，跳過用戶 seeder")
        return nil
    }

    // 創建測試用戶
    users := []*models.User{
        {Name: "系統管理員", Email: "admin@example.com", Role: "admin"},
        {Name: "商戶用戶", Email: "merchant@example.com", Role: "merchant"},
        {Name: "一般用戶", Email: "user@example.com", Role: "customer"},
    }

    // 插入數據...
    return nil
}
```

## 🔒 安全特性

- **重複執行保護**：避免重複創建數據
- **事務安全**：每個 seeder 都在事務中執行
- **錯誤處理**：完整的錯誤處理機制
- **日誌記錄**：詳細的執行日誌

## 📈 性能優化

- **批量插入**：支援批量數據插入
- **索引優化**：自動創建必要的索引
- **記憶體管理**：高效的記憶體使用
- **並發安全**：支援並發執行

## 🎯 技術亮點

1. **完全自製**：不依賴任何第三方 seeder 庫
2. **模組化設計**：易於擴展和維護
3. **接口驅動**：使用 Go 接口實現多態
4. **錯誤處理**：完整的錯誤處理機制
5. **日誌系統**：詳細的執行日誌

## 🚀 未來擴展

- [ ] 支援 YAML 配置
- [ ] 支援數據依賴關係
- [ ] 支援數據驗證
- [ ] 支援數據遷移
- [ ] 支援多數據庫

## 📝 注意事項

1. **測試環境專用**：seeder 主要用於測試環境
2. **數據清理**：定期清理測試數據
3. **版本控制**：seeder 文件需要版本控制
4. **文檔更新**：新增 seeder 時更新文檔

---

**這個自製 Seeder 系統展現了對 Go 語言的深度理解和自製能力，是技術能力的重要體現！**