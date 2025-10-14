-- 修復股票價格表的重複記錄問題

-- 先清理重複記錄，只保留每個股票代碼的最新記錄
DELETE FROM stock_prices 
WHERE id NOT IN (
    SELECT MAX(id) 
    FROM stock_prices 
    GROUP BY stock_code
);

-- 添加 UNIQUE 約束到 stock_code
CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_prices_code_unique ON stock_prices(stock_code);
