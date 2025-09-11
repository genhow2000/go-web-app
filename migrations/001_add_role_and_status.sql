-- 添加角色和狀態字段到用戶表
-- 執行時間: 2024-01-01

-- 添加角色字段 (SQLite 不支援 NOT NULL 與 DEFAULT 同時使用，所以分步執行)
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'customer';

-- 添加狀態字段 (SQLite 使用 INTEGER 代替 BOOLEAN)
ALTER TABLE users ADD COLUMN is_active INTEGER DEFAULT 1;

-- 更新現有用戶的默認值（如果有的話）
UPDATE users SET role = 'customer' WHERE role IS NULL;
UPDATE users SET is_active = 1 WHERE is_active IS NULL;

-- 創建角色索引以提高查詢性能
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 創建狀態索引
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- 創建一個默認管理員用戶（可選，根據需要調整）
-- INSERT INTO users (name, email, password, role, is_active) 
-- VALUES ('系統管理員', 'admin@example.com', '$2a$10$example_hash', 'admin', 1);
