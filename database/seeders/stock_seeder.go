package seeders

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// SeedStockData 種子股票數據
func SeedStockData(db *sql.DB) error {
	log.Println("開始種子股票數據...")

	// 插入股票分類
	categories := []struct {
		name string
		code string
		sort int
	}{
		{"電子工業", "ELECTRONICS", 1},
		{"金融保險", "FINANCE", 2},
		{"傳產工業", "INDUSTRY", 3},
		{"營建業", "CONSTRUCTION", 4},
		{"航運業", "TRANSPORTATION", 5},
		{"觀光業", "TOURISM", 6},
		{"生技醫療", "BIOTECH", 7},
		{"其他", "OTHER", 99},
	}

	for _, cat := range categories {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO stock_categories (name, code, sort, is_active, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, cat.name, cat.code, cat.sort, true, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("插入股票分類失敗: %v", err)
		}
	}

	// 插入股票基本資訊
	stocks := []struct {
		code     string
		name     string
		category string
		market   string
	}{
		// 電子股
		{"2330", "台積電", "ELECTRONICS", "TSE"},
		{"2317", "鴻海", "ELECTRONICS", "TSE"},
		{"2454", "聯發科", "ELECTRONICS", "TSE"},
		{"3711", "日月光投控", "ELECTRONICS", "TSE"},
		{"2308", "台達電", "ELECTRONICS", "TSE"},
		{"2382", "廣達", "ELECTRONICS", "TSE"},
		{"2474", "可成", "ELECTRONICS", "TSE"},
		{"3231", "緯創", "ELECTRONICS", "TSE"},
		{"3008", "大立光", "ELECTRONICS", "TSE"},
		{"2412", "中華電", "ELECTRONICS", "TSE"},
		{"2377", "微星", "ELECTRONICS", "TSE"},
		{"2357", "華碩", "ELECTRONICS", "TSE"},
		{"2376", "技嘉", "ELECTRONICS", "TSE"},
		{"2324", "仁寶", "ELECTRONICS", "TSE"},
		{"2382", "廣達", "ELECTRONICS", "TSE"},
		
		// 金融股
		{"2881", "富邦金", "FINANCE", "TSE"},
		{"2882", "國泰金", "FINANCE", "TSE"},
		{"2891", "中信金", "FINANCE", "TSE"},
		{"2886", "兆豐金", "FINANCE", "TSE"},
		{"2884", "玉山金", "FINANCE", "TSE"},
		{"2880", "華南金", "FINANCE", "TSE"},
		{"2885", "元大金", "FINANCE", "TSE"},
		{"2883", "開發金", "FINANCE", "TSE"},
		{"2887", "台新金", "FINANCE", "TSE"},
		{"2888", "新光金", "FINANCE", "TSE"},
		
		// 傳產股
		{"6505", "台塑化", "INDUSTRY", "TSE"},
		{"1101", "台泥", "INDUSTRY", "TSE"},
		{"1216", "統一", "INDUSTRY", "TSE"},
		{"1303", "南亞", "INDUSTRY", "TSE"},
		{"2002", "中鋼", "INDUSTRY", "TSE"},
		{"1326", "台化", "INDUSTRY", "TSE"},
		{"1402", "遠東新", "INDUSTRY", "TSE"},
		{"2201", "裕隆", "INDUSTRY", "TSE"},
		{"2207", "和泰車", "INDUSTRY", "TSE"},
		{"2912", "統一超", "INDUSTRY", "TSE"},
		
		// 航運股
		{"2603", "長榮", "TRANSPORTATION", "TSE"},
		{"2609", "陽明", "TRANSPORTATION", "TSE"},
		{"2615", "萬海", "TRANSPORTATION", "TSE"},
		{"2618", "長榮航", "TRANSPORTATION", "TSE"},
		{"2610", "華航", "TRANSPORTATION", "TSE"},
		
		// 生技股
		{"4743", "合一", "BIOTECH", "TSE"},
		{"6547", "高端疫苗", "BIOTECH", "TSE"},
		{"4142", "國光生", "BIOTECH", "TSE"},
		{"4162", "智擎", "BIOTECH", "TSE"},
		{"4174", "浩鼎", "BIOTECH", "TSE"},
	}

	for _, stock := range stocks {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO stocks (code, name, category, market, is_active, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, stock.code, stock.name, stock.category, stock.market, true, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("插入股票 %s 失敗: %v", stock.code, err)
		}
	}

	// 清空所有現有的模擬價格數據
	_, err := db.Exec("DELETE FROM stock_prices")
	if err != nil {
		return fmt.Errorf("清空模擬價格數據失敗: %v", err)
	}

	log.Println("股票基本資訊種子完成，價格數據將從真實API獲取")
	return nil
}
