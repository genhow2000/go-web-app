# Seeder 系統說明

## 目錄結構

```
database/seeders/
├── README.md           # 說明文檔
├── manager.go          # Seeder 管理器
├── user_seeder.go      # 用戶 seeder
└── product_seeder.go   # 產品 seeder (示例)
```

## 使用方式

### 1. 執行所有 seeder

```bash
# 本地執行
make seed

# Docker 執行
make docker-seed
# 或
docker-compose exec go-app ./seed run
```

### 2. 清除所有測試數據

```bash
# 本地執行
make seed-clear

# Docker 執行
make docker-seed-clear
# 或
docker-compose exec go-app ./seed clear
```

### 3. 執行特定 seeder

```bash
# 只執行用戶 seeder
make seed-user
# 或
docker-compose exec go-app ./seed user
```

## 如何添加新的 Seeder

### 1. 創建新的 seeder 檔案

### 2. 在 manager.go 中註冊新的 seeder

### 3. 更新 Makefile (可選)

```makefile
# 只執行訂單 seeder
seed-order:
	go run cmd/seed/main.go order

# Docker 執行訂單 seeder
docker-seed-order:
	docker-compose exec go-app ./seed order
```

## 最佳實踐

1. **命名規範**：使用 `{entity}_seeder.go` 格式
2. **重複執行保護**：檢查是否已存在測試數據
3. **錯誤處理**：處理表不存在的情況
4. **日誌記錄**：記錄執行過程
5. **清除功能**：提供清除測試數據的功能

## 優勢

✅ **模組化**：每個 seeder 獨立管理  
✅ **可擴展**：容易添加新的 seeder  
✅ **可控制**：可以執行特定或所有 seeder  
✅ **安全**：重複執行保護  
✅ **靈活**：支援清除和重新執行
