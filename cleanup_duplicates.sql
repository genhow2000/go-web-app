-- 清理重複的股票價格記錄
-- 只保留每個股票代碼的最新記錄

-- 創建臨時表，只保留每個股票代碼的最新記錄
CREATE TEMP TABLE temp_latest_prices AS
SELECT 
    stock_code,
    MAX(updated_at) as latest_updated_at
FROM stock_prices
GROUP BY stock_code;

-- 刪除所有舊的價格記錄
DELETE FROM stock_prices;

-- 重新插入最新的價格記錄
INSERT INTO stock_prices (stock_id, stock_code, price, open_price, high_price, low_price, close_price, volume, amount, change, change_percent, updated_at)
SELECT 
    sp.stock_id,
    sp.stock_code,
    sp.price,
    sp.open_price,
    sp.high_price,
    sp.low_price,
    sp.close_price,
    sp.volume,
    sp.amount,
    sp.change,
    sp.change_percent,
    sp.updated_at
FROM stock_prices sp
JOIN temp_latest_prices tlp ON sp.stock_code = tlp.stock_code AND sp.updated_at = tlp.latest_updated_at;

-- 添加 UNIQUE 約束到 stock_code
CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_prices_code_unique ON stock_prices(stock_code);
