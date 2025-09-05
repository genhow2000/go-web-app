-- 添加角色和狀態字段到用戶表
-- 執行時間: 2024-01-01

-- 添加角色字段
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'customer' NOT NULL;

-- 添加狀態字段
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true NOT NULL;

-- 添加角色檢查約束
ALTER TABLE users ADD CONSTRAINT check_role CHECK (role IN ('customer', 'admin'));

-- 創建角色索引以提高查詢性能
CREATE INDEX idx_users_role ON users(role);

-- 創建狀態索引
CREATE INDEX idx_users_is_active ON users(is_active);

-- 更新現有用戶的默認值（如果有的話）
UPDATE users SET role = 'customer', is_active = true WHERE role IS NULL OR is_active IS NULL;

-- 創建一個默認管理員用戶（可選，根據需要調整）
-- INSERT INTO users (name, email, password, role, is_active) 
-- VALUES ('系統管理員', 'admin@example.com', '$2a$10$example_hash', 'admin', true);
