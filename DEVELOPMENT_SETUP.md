# 開發環境設定

## Docker 環境配置
- **運行方式**: Docker Compose
- **開發環境文件**: `docker-compose.dev.yml`
- **容器名稱**: `go-simple-app-dev`
- **端口**: 8080:8080
- **Dockerfile**: `Dockerfile.dev`

## 常用命令
```bash
# 啟動開發環境
docker-compose -f docker-compose.dev.yml up -d

# 查看日誌
docker-compose -f docker-compose.dev.yml logs -f

# 重啟服務
docker-compose -f docker-compose.dev.yml restart

# 進入容器
docker exec -it go-simple-app-dev sh

# 停止服務
docker-compose -f docker-compose.dev.yml down
```

## 環境變數
- MongoDB: 已配置
- AI 服務: Groq (主要) + Gemini (備用)
- LINE OAuth: 已配置
- 前端映射: 已配置即時更新

## 注意事項
- 前端代碼在 `frontend/` 目錄
- 構建後的文件會自動映射到容器
- 日誌文件在 `logs/` 目錄
