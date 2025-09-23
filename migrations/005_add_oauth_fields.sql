-- 為每個角色表添加OAuth相關字段（如果不存在）
-- 使用 PRAGMA table_info 檢查欄位是否存在
-- 注意：SQLite 不支援 IF NOT EXISTS 用於 ALTER TABLE ADD COLUMN
-- 所以我們需要先檢查欄位是否存在

-- 為 customers 表添加 OAuth 欄位
-- 如果欄位已存在，這些語句會失敗但不會影響其他操作
ALTER TABLE customers ADD COLUMN oauth_provider VARCHAR(20) DEFAULT NULL;
ALTER TABLE customers ADD COLUMN oauth_id VARCHAR(100) DEFAULT NULL;
ALTER TABLE customers ADD COLUMN oauth_data TEXT DEFAULT NULL;

-- 為 merchants 表添加 OAuth 欄位
ALTER TABLE merchants ADD COLUMN oauth_provider VARCHAR(20) DEFAULT NULL;
ALTER TABLE merchants ADD COLUMN oauth_id VARCHAR(100) DEFAULT NULL;
ALTER TABLE merchants ADD COLUMN oauth_data TEXT DEFAULT NULL;

-- 為 admins 表添加 OAuth 欄位
ALTER TABLE admins ADD COLUMN oauth_provider VARCHAR(20) DEFAULT NULL;
ALTER TABLE admins ADD COLUMN oauth_id VARCHAR(100) DEFAULT NULL;
ALTER TABLE admins ADD COLUMN oauth_data TEXT DEFAULT NULL;

-- 創建OAuth登入日誌表
CREATE TABLE oauth_login_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_type VARCHAR(20) NOT NULL,
    user_id INTEGER NOT NULL,
    oauth_provider VARCHAR(20) NOT NULL,
    oauth_id VARCHAR(100) NOT NULL,
    login_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    success BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 為OAuth查詢添加索引
CREATE INDEX idx_customers_oauth ON customers(oauth_provider, oauth_id);
CREATE INDEX idx_merchants_oauth ON merchants(oauth_provider, oauth_id);
CREATE INDEX idx_admins_oauth ON admins(oauth_provider, oauth_id);
