# Redis 快取系統

## 概述

高效能記憶體快取，提升系統性能，支援多種資料結構和快取策略。

## 🚀 功能特色

- **高效能**：記憶體級別的高速存取
- **多種資料結構**：支援字串、列表、集合、有序集合等
- **持久化**：支援 RDB 和 AOF 持久化
- **分佈式**：支援主從複製和集群
- **豐富功能**：支援發布訂閱、事務等

## 📊 資料結構

### 字串 (String)

```redis
SET user:1:name "張三"
GET user:1:name
```

### 列表 (List)

```redis
LPUSH tasks "任務1"
RPUSH tasks "任務2"
LRANGE tasks 0 -1
```

### 集合 (Set)

```redis
SADD tags "Go" "Redis" "Web"
SMEMBERS tags
```

### 有序集合 (Sorted Set)

```redis
ZADD leaderboard 100 "玩家1"
ZADD leaderboard 200 "玩家2"
ZRANGE leaderboard 0 -1 WITHSCORES
```

### 哈希 (Hash)

```redis
HSET user:1 name "張三" age 25
HGET user:1 name
HGETALL user:1
```

## 🛠️ 快取策略

### 讀取策略

- **Cache-Aside**：應用程序管理快取
- **Read-Through**：快取自動讀取資料庫
- **Refresh-Ahead**：預先刷新快取

### 寫入策略

- **Write-Through**：同時寫入快取和資料庫
- **Write-Behind**：先寫快取，異步寫資料庫
- **Write-Around**：直接寫資料庫，失效快取

## 🔧 使用場景

### 會話存儲

```go
// 存儲用戶會話
redis.Set("session:"+sessionID, userData, 24*time.Hour)

// 獲取用戶會話
userData := redis.Get("session:" + sessionID)
```

### 資料庫查詢快取

```go
// 檢查快取
cached := redis.Get("user:" + userID)
if cached != nil {
    return cached
}

// 查詢資料庫
user := db.GetUser(userID)

// 存入快取
redis.Set("user:"+userID, user, 1*time.Hour)
```

### 計數器

```go
// 增加計數
redis.Incr("page:views:" + pageID)

// 獲取計數
views := redis.Get("page:views:" + pageID)
```

### 排行榜

```go
// 更新分數
redis.ZAdd("leaderboard", score, userID)

// 獲取排行榜
topUsers := redis.ZRevRange("leaderboard", 0, 9)
```

## ⚡ 性能優化

### 連接池

```go
// 配置連接池
pool := &redis.Pool{
    MaxIdle:     10,
    MaxActive:   100,
    IdleTimeout: 300 * time.Second,
}
```

### 管道操作

```go
// 批量操作
pipe := redis.Pipeline()
pipe.Set("key1", "value1", 0)
pipe.Set("key2", "value2", 0)
pipe.Exec()
```

### 批量操作

```go
// 批量獲取
keys := []string{"key1", "key2", "key3"}
values := redis.MGet(keys...)
```

## 🔒 安全特性

### 認證

```redis
# 設置密碼
CONFIG SET requirepass "your_password"

# 認證
AUTH your_password
```

### 網路安全

- 綁定特定 IP
- 防火牆保護
- SSL/TLS 加密
- 訪問控制列表

## 📈 監控與維護

### 性能監控

```redis
# 查看統計信息
INFO stats

# 查看記憶體使用
INFO memory

# 查看連接信息
INFO clients
```

### 持久化配置

```redis
# RDB 快照
SAVE
BGSAVE

# AOF 日誌
CONFIG SET appendonly yes
```

## 🚀 未來擴展

### 集群模式

- 主從複製
- 分片集群
- 故障轉移
- 負載均衡

### 高級功能

- 發布訂閱
- 事務支持
- Lua 腳本
- 模組擴展

## 📊 使用範例

### Go 客戶端

```go
import "github.com/go-redis/redis/v8"

// 創建客戶端
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// 基本操作
err := rdb.Set(ctx, "key", "value", 0).Err()
val, err := rdb.Get(ctx, "key").Result()
```

### 快取模式

```go
func GetUser(userID string) (*User, error) {
    // 1. 檢查快取
    cached, err := redis.Get("user:" + userID)
    if err == nil {
        return cached, nil
    }

    // 2. 查詢資料庫
    user, err := db.GetUser(userID)
    if err != nil {
        return nil, err
    }

    // 3. 存入快取
    redis.Set("user:"+userID, user, 1*time.Hour)

    return user, nil
}
```

---

**這個 Redis 快取系統展現了高效能快取設計和記憶體資料庫的最佳實踐！**
