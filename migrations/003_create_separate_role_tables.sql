-- 創建三張獨立的角色表
-- 執行時間: 2024-01-02

-- 1. 客戶表 (customers)
CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    birth_date DATE,
    gender VARCHAR(10),
    is_active BOOLEAN DEFAULT 1,
    email_verified BOOLEAN DEFAULT 0,
    last_login DATETIME,
    login_count INTEGER DEFAULT 0,
    profile_data TEXT, -- JSON 格式存儲額外資料
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. 商戶表 (merchants)
CREATE TABLE IF NOT EXISTS merchants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    business_name VARCHAR(200),
    business_license VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    business_type VARCHAR(50),
    is_active BOOLEAN DEFAULT 1,
    is_verified BOOLEAN DEFAULT 0, -- 商戶認證狀態
    last_login DATETIME,
    login_count INTEGER DEFAULT 0,
    business_data TEXT, -- JSON 格式存儲商戶專用資料
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. 管理員表 (admins)
CREATE TABLE IF NOT EXISTS admins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    admin_level VARCHAR(20) DEFAULT 'normal', -- normal, senior, super
    department VARCHAR(100),
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT 1,
    last_login DATETIME,
    login_count INTEGER DEFAULT 0,
    admin_data TEXT, -- JSON 格式存儲管理員專用資料
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 創建索引
-- 客戶表索引
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers(is_active);
CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at);

-- 商戶表索引
CREATE INDEX IF NOT EXISTS idx_merchants_email ON merchants(email);
CREATE INDEX IF NOT EXISTS idx_merchants_is_active ON merchants(is_active);
CREATE INDEX IF NOT EXISTS idx_merchants_is_verified ON merchants(is_verified);
CREATE INDEX IF NOT EXISTS idx_merchants_business_type ON merchants(business_type);
CREATE INDEX IF NOT EXISTS idx_merchants_created_at ON merchants(created_at);

-- 管理員表索引
CREATE INDEX IF NOT EXISTS idx_admins_email ON admins(email);
CREATE INDEX IF NOT EXISTS idx_admins_is_active ON admins(is_active);
CREATE INDEX IF NOT EXISTS idx_admins_admin_level ON admins(admin_level);
CREATE INDEX IF NOT EXISTS idx_admins_department ON admins(department);
CREATE INDEX IF NOT EXISTS idx_admins_created_at ON admins(created_at);

-- 創建登入日誌表（統一記錄所有角色的登入）
CREATE TABLE IF NOT EXISTS login_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_type VARCHAR(20) NOT NULL, -- 'customer', 'merchant', 'admin'
    user_id INTEGER NOT NULL,
    login_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    success BOOLEAN DEFAULT 1,
    login_method VARCHAR(20) DEFAULT 'password', -- password, oauth, etc.
    FOREIGN KEY (user_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES merchants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES admins(id) ON DELETE CASCADE
);

-- 登入日誌索引
CREATE INDEX IF NOT EXISTS idx_login_logs_user_type ON login_logs(user_type);
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_login_time ON login_logs(login_time);
CREATE INDEX IF NOT EXISTS idx_login_logs_success ON login_logs(success);

-- 創建角色權限表
CREATE TABLE IF NOT EXISTS role_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    permissions TEXT, -- JSON 格式存儲權限
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入默認角色權限
INSERT OR IGNORE INTO role_permissions (role, description, permissions) VALUES 
('customer', '客戶角色', '["view_own_profile", "update_own_profile", "create_chat", "view_own_chats", "place_orders", "view_own_orders"]'),
('merchant', '商戶角色', '["view_own_profile", "update_own_profile", "create_chat", "view_own_chats", "manage_products", "view_orders", "manage_orders", "view_analytics", "manage_business"]'),
('admin', '管理員角色', '["view_all_profiles", "update_all_profiles", "delete_users", "manage_roles", "view_system_stats", "manage_system", "manage_all_orders", "manage_all_products", "view_analytics"]');

-- 創建角色權限索引
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role);

-- 從舊的 users 表遷移數據（如果存在）
-- 注意：這需要在應用層處理，因為需要根據角色分配到不同的表

-- 創建視圖來統一查詢所有用戶（可選）
CREATE VIEW IF NOT EXISTS all_users AS
SELECT 'customer' as user_type, id, name, email, is_active, last_login, login_count, created_at, updated_at FROM customers
UNION ALL
SELECT 'merchant' as user_type, id, name, email, is_active, last_login, login_count, created_at, updated_at FROM merchants
UNION ALL
SELECT 'admin' as user_type, id, name, email, is_active, last_login, login_count, created_at, updated_at FROM admins;
