-- 初始資料庫結構
-- 執行時間: 2024-01-01

-- 創建基本的系統表
CREATE TABLE IF NOT EXISTS system_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入系統資訊
INSERT OR IGNORE INTO system_info (key, value, description) VALUES 
('version', '1.0.0', '系統版本'),
('database_version', '1.0.0', '資料庫版本'),
('last_migration', '001_initial_schema', '最後執行的 migration');

-- 創建系統資訊索引
CREATE INDEX IF NOT EXISTS idx_system_info_key ON system_info(key);
