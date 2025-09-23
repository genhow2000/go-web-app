# 腳本文件說明

本資料夾包含所有開發和部署相關的腳本文件。

## 開發環境腳本

### `start-dev.sh`
啟動開發環境（推薦使用）
```bash
./sh/start-dev.sh
```
- 啟動 Go 後端服務（Docker 容器）
- 啟動 Vue 前端服務（本地 Vite）
- 前端修改會自動熱重載

### `stop-dev.sh`
停止開發環境
```bash
./sh/stop-dev.sh
```

### `test-dev-env.sh`
測試開發環境配置
```bash
./sh/test-dev-env.sh
```

## 部署腳本

### `deploy-gcp.sh`
部署到 Google Cloud Platform
```bash
./sh/deploy-gcp.sh
```

## 系統驗證腳本

### `setup-budget-minimal.sh`
設置最小預算配置
```bash
./sh/setup-budget-minimal.sh
```

### `verify_docker.sh`
驗證 Docker 環境
```bash
./sh/verify_docker.sh
```

### `verify_system.sh`
驗證系統配置
```bash
./sh/verify_system.sh
```

## 優化腳本

### `optimize-scripts.sh`
優化腳本文件
```bash
./sh/optimize-scripts.sh
```

## 使用建議

1. **日常開發**：使用 `start-dev.sh` 和 `stop-dev.sh`
2. **環境測試**：使用 `test-dev-env.sh`
3. **部署**：使用 `deploy-gcp.sh`
4. **系統驗證**：使用 `verify_*.sh` 腳本

## 注意事項

- 所有腳本都需要在專案根目錄執行
- 確保 Docker 和 Node.js 已安裝
- 開發環境需要本地安裝 Node.js
- 生產環境使用 Docker 構建
