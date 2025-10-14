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

	// 插入模擬價格數據
	priceData := []struct {
		code          string
		price         float64
		openPrice     float64
		highPrice     float64
		lowPrice      float64
		closePrice    float64
		volume        int64
		amount        float64
		change        float64
		changePercent float64
	}{
		{"2330", 580.0, 575.0, 585.0, 570.0, 575.0, 25000000, 14500000000, 5.0, 0.87},
		{"2317", 105.5, 104.0, 106.0, 103.5, 104.0, 15000000, 1582500000, 1.5, 1.44},
		{"2454", 950.0, 940.0, 955.0, 935.0, 940.0, 8000000, 7600000000, 10.0, 1.06},
		{"6505", 85.2, 84.5, 86.0, 84.0, 84.5, 12000000, 1022400000, 0.7, 0.83},
		{"2881", 65.8, 65.0, 66.2, 64.8, 65.0, 18000000, 1184400000, 0.8, 1.23},
		{"2882", 58.5, 58.0, 59.0, 57.8, 58.0, 20000000, 1170000000, 0.5, 0.86},
		{"1101", 42.3, 42.0, 42.8, 41.8, 42.0, 10000000, 423000000, 0.3, 0.71},
		{"1216", 75.6, 75.0, 76.2, 74.8, 75.0, 8000000, 604800000, 0.6, 0.80},
		{"1303", 68.9, 68.5, 69.5, 68.2, 68.5, 12000000, 826800000, 0.4, 0.58},
		{"2002", 32.8, 32.5, 33.2, 32.3, 32.5, 15000000, 492000000, 0.3, 0.92},
		{"2412", 125.5, 125.0, 126.0, 124.8, 125.0, 5000000, 627500000, 0.5, 0.40},
		{"2891", 28.6, 28.4, 28.8, 28.2, 28.4, 25000000, 715000000, 0.2, 0.70},
		{"2886", 33.2, 33.0, 33.5, 32.8, 33.0, 18000000, 597600000, 0.2, 0.61},
		{"2884", 25.8, 25.6, 26.0, 25.4, 25.6, 20000000, 516000000, 0.2, 0.78},
		{"3711", 95.5, 94.8, 96.2, 94.5, 94.8, 10000000, 955000000, 0.7, 0.74},
		{"2308", 285.0, 283.0, 287.0, 282.0, 283.0, 3000000, 855000000, 2.0, 0.71},
		{"2382", 185.5, 184.0, 186.8, 183.5, 184.0, 5000000, 927500000, 1.5, 0.82},
		{"2474", 125.8, 125.0, 126.5, 124.5, 125.0, 2000000, 251600000, 0.8, 0.64},
		{"3231", 45.6, 45.2, 46.0, 45.0, 45.2, 8000000, 364800000, 0.4, 0.88},
		{"3008", 2150.0, 2140.0, 2160.0, 2135.0, 2140.0, 500000, 1075000000, 10.0, 0.47},
	}

	for _, price := range priceData {
		_, err := db.Exec(`
			INSERT OR REPLACE INTO stock_prices 
			(stock_code, price, open_price, high_price, low_price, close_price, volume, amount, change, change_percent, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, price.code, price.price, price.openPrice, price.highPrice, price.lowPrice, 
			price.closePrice, price.volume, price.amount, price.change, price.changePercent, time.Now())
		if err != nil {
			return fmt.Errorf("插入股票價格 %s 失敗: %v", price.code, err)
		}
	}

	log.Println("股票數據種子完成")
	return nil
}
