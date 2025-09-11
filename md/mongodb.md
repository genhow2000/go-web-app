# MongoDB 文檔資料庫

## 概述

NoSQL 文檔資料庫，支援複雜資料結構，提供靈活的資料模型和強大的查詢功能。

## 🚀 功能特色

- **文檔存儲**：以 JSON 格式存儲資料
- **靈活模式**：無固定表結構，支援動態模式
- **強大查詢**：支援複雜查詢和聚合操作
- **水平擴展**：支援分片和集群部署
- **豐富索引**：支援多種索引類型

## 📊 資料模型

### 文檔結構

```json
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
  "name": "張三",
  "age": 25,
  "email": "zhangsan@example.com",
  "address": {
    "street": "中山路123號",
    "city": "台北市",
    "zip": "10001"
  },
  "hobbies": ["閱讀", "游泳", "程式設計"],
  "created_at": ISODate("2025-09-11T05:22:05Z")
}
```

### 集合 (Collection)

- 類似關聯資料庫的表
- 無固定結構
- 動態添加欄位
- 支援嵌套文檔

## 🛠️ 基本操作

### 插入文檔

```javascript
// 插入單一文檔
db.users.insertOne({
  name: "李四",
  age: 30,
  email: "lisi@example.com",
});

// 插入多個文檔
db.users.insertMany([
  { name: "王五", age: 28 },
  { name: "趙六", age: 35 },
]);
```

### 查詢文檔

```javascript
// 查詢所有文檔
db.users.find();

// 條件查詢
db.users.find({ age: { $gt: 25 } });

// 投影查詢
db.users.find({}, { name: 1, email: 1 });

// 排序查詢
db.users.find().sort({ age: -1 });
```

### 更新文檔

```javascript
// 更新單一文檔
db.users.updateOne({ name: "張三" }, { $set: { age: 26 } });

// 更新多個文檔
db.users.updateMany({ age: { $lt: 30 } }, { $set: { status: "young" } });
```

### 刪除文檔

```javascript
// 刪除單一文檔
db.users.deleteOne({ name: "張三" });

// 刪除多個文檔
db.users.deleteMany({ age: { $gt: 50 } });
```

## 🔍 高級查詢

### 聚合管道

```javascript
db.orders.aggregate([
  { $match: { status: "completed" } },
  {
    $group: {
      _id: "$customer_id",
      total: { $sum: "$amount" },
    },
  },
  { $sort: { total: -1 } },
  { $limit: 10 },
]);
```

### 索引優化

```javascript
// 創建單一索引
db.users.createIndex({ email: 1 });

// 創建複合索引
db.users.createIndex({ name: 1, age: -1 });

// 創建文字索引
db.articles.createIndex({ title: "text", content: "text" });
```

## 🔧 Go 整合

### 驅動程序

```go
import "go.mongodb.org/mongo-driver/mongo"

// 連接 MongoDB
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
```

### 基本操作

```go
// 插入文檔
user := User{Name: "張三", Age: 25}
result, err := collection.InsertOne(ctx, user)

// 查詢文檔
var user User
err := collection.FindOne(ctx, bson.M{"name": "張三"}).Decode(&user)

// 更新文檔
update := bson.M{"$set": bson.M{"age": 26}}
result, err := collection.UpdateOne(ctx, bson.M{"name": "張三"}, update)
```

## 📈 性能優化

### 索引策略

- 為常用查詢欄位創建索引
- 使用複合索引優化多欄位查詢
- 定期分析查詢性能
- 監控索引使用情況

### 查詢優化

- 使用投影減少網路傳輸
- 合理使用分頁查詢
- 避免全表掃描
- 使用聚合管道優化複雜查詢

## 🔒 安全特性

### 認證授權

```javascript
// 創建用戶
db.createUser({
  user: "app_user",
  pwd: "password",
  roles: ["readWrite"],
});
```

### 資料驗證

```javascript
// 創建帶驗證的集合
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      required: ["name", "email"],
      properties: {
        name: { type: "string" },
        email: {
          type: "string",
          pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
        },
      },
    },
  },
});
```

## 🚀 未來擴展

### 分片集群

- 水平分片
- 自動負載均衡
- 故障轉移
- 資料分佈優化

### 高級功能

- 地理空間查詢
- 全文搜尋
- 圖形查詢
- 時間序列資料

## 📊 使用場景

### 內容管理

- 文章和博客
- 產品目錄
- 用戶生成內容
- 多媒體資料

### 物聯網

- 感測器資料
- 設備狀態
- 時間序列資料
- 即時監控

### 電商平台

- 商品資訊
- 訂單管理
- 用戶行為
- 推薦系統

---

**這個 MongoDB 文檔資料庫展現了 NoSQL 資料庫的靈活性和強大查詢能力！**
