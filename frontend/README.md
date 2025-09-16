# Vue.js 前端應用

這是阿和商城的Vue.js前端應用，使用現代化的前端技術棧構建。

## 技術棧

- **Vue 3** - 漸進式JavaScript框架
- **Vite** - 快速的前端構建工具
- **Vue Router** - 官方路由管理器
- **Pinia** - 現代化的狀態管理庫
- **Element Plus** - 基於Vue 3的組件庫
- **Axios** - HTTP客戶端

## 項目結構

```
frontend/
├── src/
│   ├── components/          # 可重用組件
│   │   ├── common/         # 通用組件
│   │   ├── product/        # 商品相關組件
│   │   ├── chat/           # 聊天相關組件
│   │   └── tech/           # 技術展示組件
│   ├── views/              # 頁面組件
│   │   ├── auth/           # 認證頁面
│   │   ├── dashboard/      # 儀表板頁面
│   │   └── merchant/       # 商戶管理頁面
│   ├── stores/             # Pinia狀態管理
│   ├── services/           # API服務
│   ├── router/             # 路由配置
│   └── assets/             # 靜態資源
├── public/                 # 公共資源
├── package.json           # 依賴配置
└── vite.config.js         # Vite配置
```

## 安裝和運行

### 1. 安裝依賴

```bash
cd frontend
npm install
```

### 2. 開發模式運行

```bash
npm run dev
```

前端應用將在 `http://localhost:3000` 運行

### 3. 構建生產版本

```bash
npm run build
```

構建文件將輸出到 `../static/dist` 目錄

## 功能特性

### 🏠 首頁
- 響應式設計
- 商品展示
- 分類導航
- AI聊天助手
- 搜尋功能

### 🔐 認證系統
- 客戶登入/註冊
- 商戶登入/註冊
- 管理員登入
- JWT Token認證
- 路由守衛

### 📊 儀表板
- 客戶儀表板
- 商戶儀表板
- 管理員儀表板
- 統計數據展示
- 快速操作

### 🛍️ 商品管理
- 商品列表
- 商品創建/編輯
- 商品狀態管理
- 庫存管理
- 分類管理

### 🤖 AI聊天
- 實時聊天界面
- 多AI提供商支持
- 對話歷史
- 打字指示器

### 🚀 技術展示
- 系統狀態監控
- 技術棧展示
- 功能介紹
- 實時數據更新

## API集成

前端通過Axios與後端API進行通信：

- **基礎URL**: `http://localhost:8080`
- **認證**: JWT Token自動添加
- **錯誤處理**: 統一的錯誤處理機制
- **請求攔截**: 自動添加認證頭
- **響應攔截**: 處理認證錯誤

## 狀態管理

使用Pinia進行狀態管理：

- **authStore**: 用戶認證狀態
- **userStore**: 用戶信息
- **productStore**: 商品數據（可擴展）

## 路由配置

- **首頁**: `/`
- **技術展示**: `/tech-showcase`
- **登入頁面**: `/customer/login`, `/merchant/login`, `/admin/login`
- **儀表板**: `/customer/dashboard`, `/merchant/dashboard`, `/admin/dashboard`
- **商品管理**: `/merchant/products`

## 開發指南

### 添加新頁面

1. 在 `src/views/` 創建頁面組件
2. 在 `src/router/index.js` 添加路由
3. 更新導航菜單

### 添加新組件

1. 在 `src/components/` 創建組件
2. 在需要的地方導入使用
3. 遵循Vue 3 Composition API

### 添加新API

1. 在 `src/services/api.js` 添加API方法
2. 在組件中調用API
3. 處理加載狀態和錯誤

## 部署

### 開發環境

```bash
npm run dev
```

### 生產環境

```bash
npm run build
```

構建後的文件會輸出到 `../static/dist`，可以通過Go後端提供靜態文件服務。

## 注意事項

1. 確保後端API服務正在運行
2. 檢查API端點是否正確配置
3. 注意CORS設置
4. 確保JWT Token正確處理

## 故障排除

### 常見問題

1. **API請求失敗**: 檢查後端服務是否運行
2. **認證問題**: 檢查JWT Token是否正確
3. **路由問題**: 檢查路由配置是否正確
4. **樣式問題**: 檢查CSS導入是否正確

### 調試技巧

1. 使用瀏覽器開發者工具
2. 檢查網絡請求
3. 查看控制台錯誤
4. 使用Vue DevTools
