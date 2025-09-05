# Go Web Application with Modern Dashboard

一個基於 Go 語言開發的現代化 Web 應用程序，具儀表板界面。

## 🚀 功能特色

- **現代化儀表板** - 類似專業界面設計
- **用戶認證系統** - JWT token 認證
- **PostgreSQL 數據庫** - 完整的數據持久化
- **日誌系統** - 使用 logrus 的專業日誌記錄
- **Docker 容器化** - 一鍵部署
- **響應式設計** - 支持桌面和移動設備
- **實時統計** - 動態數據展示

## 📋 技術棧

- **後端**: Go 1.21 + Gin 框架
- **數據庫**: PostgreSQL
- **前端**: HTML5 + CSS3 + JavaScript + Font Awesome
- **容器化**: Docker + Docker Compose
- **日誌**: Logrus
- **認證**: JWT

## 🛠️ 快速開始

### 使用 Docker（推薦）

1. 克隆倉庫

```bash
git clone https://github.com/genhow2000/go-web-app.git
cd go-web-app
```

2. 啟動服務

```bash
docker-compose up -d
```

3. 訪問應用

- 應用地址: http://localhost:8080
- 登入頁面: http://localhost:8080/login
- 註冊頁面: http://localhost:8080/register

### 本地開發

1. 安裝 Go 1.21+
2. 安裝 PostgreSQL
3. 設置環境變量

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=goapp
export JWT_SECRET=your-secret-key
```

4. 運行應用

```bash
go mod tidy
go run main.go
```

## 📁 項目結構

```
go-web-app/
├── config/          # 配置管理
├── controllers/     # 控制器層
├── database/        # 數據庫連接
├── logger/          # 日誌系統
├── middleware/      # 中間件
├── models/          # 數據模型
├── routes/          # 路由配置
├── services/        # 業務邏輯
├── templates/       # HTML模板
├── docker-compose.yml
├── Dockerfile
└── main.go
```

## 🔧 配置說明

### 環境變量

| 變量名      | 默認值          | 說明         |
| ----------- | --------------- | ------------ |
| PORT        | 8080            | 服務端口     |
| HOST        | 0.0.0.0         | 服務地址     |
| DB_HOST     | localhost       | 數據庫地址   |
| DB_PORT     | 5432            | 數據庫端口   |
| DB_USER     | postgres        | 數據庫用戶名 |
| DB_PASSWORD | password        | 數據庫密碼   |
| DB_NAME     | goapp           | 數據庫名稱   |
| JWT_SECRET  | your-secret-key | JWT 密鑰     |

## 📊 API 接口

### 認證接口

- `POST /register` - 用戶註冊
- `POST /login` - 用戶登入
- `POST /logout` - 用戶登出

### 用戶管理

- `GET /users` - 獲取所有用戶
- `GET /users/:id` - 獲取特定用戶
- `POST /users` - 創建用戶
- `PUT /users/:id` - 更新用戶
- `DELETE /users/:id` - 刪除用戶

### 系統接口

- `GET /health` - 健康檢查
- `GET /dashboard` - 儀表板頁面

## 🎨 界面預覽

### 登入頁面

- 現代化漸層背景設計
- 響應式表單布局
- 實時錯誤提示

### 儀表板

- 側邊欄導航菜單
- 統計卡片展示
- 實時數據更新
- 用戶列表管理

## 📝 日誌系統

應用使用 logrus 提供專業的日誌記錄：

- **文件日誌**: `logs/YYYY-MM-DD.log`
- **控制台日誌**: 實時輸出
- **日誌級別**: Debug, Info, Warn, Error, Fatal
- **JSON 格式**: 便於日誌分析

查看日誌：

```bash
# Docker環境
docker logs go-simple-app

# 本地文件
tail -f logs/2025-09-03.log
```

## 🐳 Docker 部署

### 構建鏡像

```bash
docker build -t go-web-app .
```

### 運行容器

```bash
docker-compose up -d
```

### 查看日誌

```bash
docker-compose logs -f
```

## 🔒 安全特性

- JWT token 認證
- 密碼 bcrypt 加密
- CORS 跨域保護
- SQL 注入防護
- XSS 防護

## 📈 性能優化

- 數據庫連接池
- 中間件緩存
- 靜態資源優化
- 響應式設計

## 🤝 貢獻指南

1. Fork 本倉庫
2. 創建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

## 📄 許可證

本項目採用 MIT 許可證 - 查看 [LICENSE](LICENSE) 文件了解詳情

## 👨‍💻 作者

**genhow2000**

- GitHub: [@genhow2000](https://github.com/genhow2000)

## 🙏 致謝

- [Gin Web Framework](https://gin-gonic.com/)
- [Logrus](https://github.com/sirupsen/logrus)
- [Font Awesome](https://fontawesome.com/)
- [PostgreSQL](https://www.postgresql.org/)

---

⭐ 如果這個項目對你有幫助，請給個 Star！
