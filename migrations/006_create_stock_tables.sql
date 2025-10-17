-- 創建股票相關資料表

-- 股票分類表
CREATE TABLE IF NOT EXISTS stock_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,           -- 分類名稱
    code VARCHAR(20) NOT NULL UNIQUE,     -- 分類代碼
    sort INTEGER DEFAULT 0,               -- 排序
    is_active BOOLEAN DEFAULT TRUE,       -- 是否啟用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 股票基本資訊表
CREATE TABLE IF NOT EXISTS stocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code VARCHAR(10) NOT NULL UNIQUE,     -- 股票代碼
    name VARCHAR(100) NOT NULL,           -- 股票名稱
    category VARCHAR(50),                 -- 產業分類
    market VARCHAR(10) NOT NULL,          -- 市場別 (TSE/OTC)
    is_active BOOLEAN DEFAULT TRUE,       -- 是否交易中
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 股票價格表
CREATE TABLE IF NOT EXISTS stock_prices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stock_code VARCHAR(10) NOT NULL,      -- 股票代碼
    price DECIMAL(10,2),                  -- 現價
    open_price DECIMAL(10,2),             -- 開盤價
    high_price DECIMAL(10,2),             -- 最高價
    low_price DECIMAL(10,2),              -- 最低價
    close_price DECIMAL(10,2),            -- 昨收價
    volume BIGINT DEFAULT 0,              -- 成交量
    amount DECIMAL(15,2) DEFAULT 0,       -- 成交金額
    change DECIMAL(10,2) DEFAULT 0,       -- 漲跌
    change_percent DECIMAL(8,4) DEFAULT 0, -- 漲跌幅
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (stock_code) REFERENCES stocks(code) ON DELETE CASCADE
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_stocks_code ON stocks(code);
CREATE INDEX IF NOT EXISTS idx_stocks_category ON stocks(category);
CREATE INDEX IF NOT EXISTS idx_stocks_market ON stocks(market);
CREATE INDEX IF NOT EXISTS idx_stocks_active ON stocks(is_active);

CREATE INDEX IF NOT EXISTS idx_stock_prices_code ON stock_prices(stock_code);
CREATE INDEX IF NOT EXISTS idx_stock_prices_updated ON stock_prices(updated_at);

CREATE INDEX IF NOT EXISTS idx_stock_categories_code ON stock_categories(code);
CREATE INDEX IF NOT EXISTS idx_stock_categories_active ON stock_categories(is_active);

-- 插入股票分類數據
INSERT OR IGNORE INTO stock_categories (name, code, sort) VALUES
('電子工業', 'ELECTRONICS', 1),
('金融保險', 'FINANCE', 2),
('傳產工業', 'INDUSTRY', 3),
('營建業', 'CONSTRUCTION', 4),
('航運業', 'TRANSPORTATION', 5),
('觀光業', 'TOURISM', 6),
('生技醫療', 'BIOTECH', 7),
('其他', 'OTHER', 99);

-- 插入一些測試股票數據
INSERT OR IGNORE INTO stocks (code, name, category, market) VALUES
('2330', '台積電', 'ELECTRONICS', 'TSE'),
('2317', '鴻海', 'ELECTRONICS', 'TSE'),
('2454', '聯發科', 'ELECTRONICS', 'TSE'),
('6505', '台塑化', 'INDUSTRY', 'TSE'),
('2881', '富邦金', 'FINANCE', 'TSE'),
('2882', '國泰金', 'FINANCE', 'TSE'),
('1101', '台泥', 'INDUSTRY', 'TSE'),
('1216', '統一', 'INDUSTRY', 'TSE'),
('1303', '南亞', 'INDUSTRY', 'TSE'),
('2002', '中鋼', 'INDUSTRY', 'TSE'),
('2412', '中華電', 'ELECTRONICS', 'TSE'),
('2408', '南亞科', 'ELECTRONICS', 'TSE'),
('2891', '中信金', 'FINANCE', 'TSE'),
('2886', '兆豐金', 'FINANCE', 'TSE'),
('2884', '玉山金', 'FINANCE', 'TSE'),
('3711', '日月光投控', 'ELECTRONICS', 'TSE'),
('2308', '台達電', 'ELECTRONICS', 'TSE'),
('2382', '廣達', 'ELECTRONICS', 'TSE'),
('2474', '可成', 'ELECTRONICS', 'TSE'),
('3231', '緯創', 'ELECTRONICS', 'TSE'),
('3008', '大立光', 'ELECTRONICS', 'TSE');

-- 插入一些測試價格數據
INSERT OR IGNORE INTO stock_prices (stock_code, price, open_price, high_price, low_price, close_price, volume, amount, change, change_percent) VALUES
('2330', 580.0, 575.0, 585.0, 570.0, 575.0, 25000000, 14500000000, 5.0, 0.87),
('2317', 105.5, 104.0, 106.0, 103.5, 104.0, 15000000, 1582500000, 1.5, 1.44),
('2454', 950.0, 940.0, 955.0, 935.0, 940.0, 8000000, 7600000000, 10.0, 1.06),
('6505', 85.2, 84.5, 86.0, 84.0, 84.5, 12000000, 1022400000, 0.7, 0.83),
('2881', 65.8, 65.0, 66.2, 64.8, 65.0, 18000000, 1184400000, 0.8, 1.23),
('2882', 58.5, 58.0, 59.0, 57.8, 58.0, 20000000, 1170000000, 0.5, 0.86),
('1101', 42.3, 42.0, 42.8, 41.8, 42.0, 10000000, 423000000, 0.3, 0.71),
('1216', 75.6, 75.0, 76.2, 74.8, 75.0, 8000000, 604800000, 0.6, 0.80),
('1303', 68.9, 68.5, 69.5, 68.2, 68.5, 12000000, 826800000, 0.4, 0.58),
('2002', 32.8, 32.5, 33.2, 32.3, 32.5, 15000000, 492000000, 0.3, 0.92),
('2412', 125.5, 125.0, 126.0, 124.8, 125.0, 5000000, 627500000, 0.5, 0.40),
('2891', 28.6, 28.4, 28.8, 28.2, 28.4, 25000000, 715000000, 0.2, 0.70),
('2886', 33.2, 33.0, 33.5, 32.8, 33.0, 18000000, 597600000, 0.2, 0.61),
('2884', 25.8, 25.6, 26.0, 25.4, 25.6, 20000000, 516000000, 0.2, 0.78),
('3711', 95.5, 94.8, 96.2, 94.5, 94.8, 10000000, 955000000, 0.7, 0.74),
('2308', 285.0, 283.0, 287.0, 282.0, 283.0, 3000000, 855000000, 2.0, 0.71),
('2382', 185.5, 184.0, 186.8, 183.5, 184.0, 5000000, 927500000, 1.5, 0.82),
('2474', 125.8, 125.0, 126.5, 124.5, 125.0, 2000000, 251600000, 0.8, 0.64),
('3231', 45.6, 45.2, 46.0, 45.0, 45.2, 8000000, 364800000, 0.4, 0.88),
('3008', 2150.0, 2140.0, 2160.0, 2135.0, 2140.0, 500000, 1075000000, 10.0, 0.47);

-- 更新 stocks 表的 updated_at 欄位觸發器
CREATE TRIGGER IF NOT EXISTS update_stocks_updated_at 
    AFTER UPDATE ON stocks
    FOR EACH ROW
BEGIN
    UPDATE stocks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- 更新 stock_prices 表的 updated_at 欄位觸發器
CREATE TRIGGER IF NOT EXISTS update_stock_prices_updated_at 
    AFTER UPDATE ON stock_prices
    FOR EACH ROW
BEGIN
    UPDATE stock_prices SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
