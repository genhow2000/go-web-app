# 資料庫工具說明

## Migration 系統

您的 Go 專案已經實現了完整的 Migration 系統！

### 功能特點

1. **自動執行**：應用程式啟動時會自動執行所有未執行的 migration
2. **版本控制**：使用檔案名前的數字來控制執行順序（如：`001_add_role_and_status.sql`）
3. **重複執行保護**：已執行的 migration 不會重複執行
4. **記錄追蹤**：在 `migrations` 表中記錄所有已執行的 migration

### 現有 Migration 檔案

- `migrations/001_add_role_and_status.sql` - 添加角色和狀態欄位

### 如何添加新的 Migration

1. 在 `migrations/` 目錄下創建新的 SQL 檔案
2. 檔案名格式：`XXX_description.sql`（XXX 是三位數版本號）
3. 例如：`002_add_products_table.sql`

## Seeder 系統

### 功能特點

1. **自動執行**：應用程式啟動時會自動執行 seeder
2. **重複執行保護**：如果測試用戶已存在，會跳過執行
3. **測試數據**：自動創建管理員、商戶和一般用戶

### 現有測試用戶

- **管理員**：`admin@example.com` / `admin123`
- **商戶**：`merchant@example.com` / `admin123`
- **一般用戶**：`user@example.com` / `admin123`

## 命令工具

### 使用 Makefile（推薦）

```bash
# 建構應用程式
make build

# 運行應用程式
make run

# 執行 migration
make migrate

# 執行 seeder
make seed

# Docker 相關命令
make docker-build
make docker-run
make docker-logs
make docker-stop

# 在 Docker 容器中執行命令
make docker-migrate
make docker-seed
make docker-shell
```

### 直接使用 Docker 命令

```bash
# 執行 migration
docker-compose exec go-app ./migrate

# 執行 seeder
docker-compose exec go-app ./seed

# 進入容器
docker-compose exec go-app /bin/sh
```

### 本地執行（需要 Go 環境）

```bash
# 執行 migration
go run cmd/migrate/main.go

# 執行 seeder
go run cmd/seed/main.go
```

## 資料庫結構

### Migration 表

```sql
CREATE TABLE migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 用戶表

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

## 最佳實踐

1. **Migration 命名**：使用描述性的檔案名，如 `001_add_user_roles.sql`
2. **版本號**：使用三位數版本號，確保執行順序
3. **可逆性**：盡量讓 migration 可逆（雖然目前系統不支援自動回滾）
4. **測試**：在開發環境中測試 migration 後再部署到生產環境
5. **備份**：執行重要 migration 前先備份資料庫

## 與其他框架的比較

| 功能      | 您的系統 | Laravel | Rails | Django |
| --------- | -------- | ------- | ----- | ------ |
| Migration | ✅       | ✅      | ✅    | ✅     |
| Seeder    | ✅       | ✅      | ✅    | ✅     |
| 版本控制  | ✅       | ✅      | ✅    | ✅     |
| 自動執行  | ✅       | ✅      | ✅    | ✅     |
| 回滾功能  | ❌       | ✅      | ✅    | ✅     |

您的 Go 專案已經有了完整的 Migration 和 Seeder 系統，功能與主流框架相當！
