package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// Stock 股票基本資訊模型
type Stock struct {
	ID           int       `json:"id" db:"id"`
	Code         string    `json:"code" db:"code"`                   // 股票代碼
	Name         string    `json:"name" db:"name"`                   // 股票名稱
	Category     string    `json:"category" db:"category"`           // 產業分類
	Market       string    `json:"market" db:"market"`               // 市場別 (TSE/OTC)
	IsActive     bool      `json:"is_active" db:"is_active"`         // 是否交易中
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// StockPrice 股票價格資訊模型
type StockPrice struct {
	ID           int       `json:"id" db:"id"`
	StockCode    string    `json:"stock_code" db:"stock_code"`       // 股票代碼
	Price        float64   `json:"price" db:"price"`                 // 現價
	OpenPrice    float64   `json:"open_price" db:"open_price"`       // 開盤價
	HighPrice    float64   `json:"high_price" db:"high_price"`       // 最高價
	LowPrice     float64   `json:"low_price" db:"low_price"`         // 最低價
	ClosePrice   float64   `json:"close_price" db:"close_price"`     // 昨收價
	Volume       int64     `json:"volume" db:"volume"`               // 成交量
	Amount       float64   `json:"amount" db:"amount"`               // 成交金額
	Change       float64   `json:"change" db:"change"`               // 漲跌
	ChangePercent float64  `json:"change_percent" db:"change_percent"` // 漲跌幅
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// StockListResponse 股票列表回應模型
type StockListResponse struct {
	Stocks      []StockWithPrice `json:"stocks"`
	Pagination  Pagination       `json:"pagination"`
	TotalCount  int              `json:"total_count"`
}

// StockWithPrice 包含價格資訊的股票模型
type StockWithPrice struct {
	Stock
	Price        *StockPrice `json:"price,omitempty"`
}

// Pagination 分頁模型
type Pagination struct {
	CurrentPage int `json:"current_page"`  // 當前頁碼
	PerPage     int `json:"per_page"`      // 每頁筆數
	TotalPages  int `json:"total_pages"`   // 總頁數
	TotalCount  int `json:"total_count"`   // 總筆數
	HasNext     bool `json:"has_next"`     // 是否有下一頁
	HasPrev     bool `json:"has_prev"`     // 是否有上一頁
}

// StockFilter 股票篩選條件
type StockFilter struct {
	Category    string `json:"category"`     // 產業分類
	Market      string `json:"market"`       // 市場別
	Search      string `json:"search"`       // 搜尋關鍵字
	SortBy      string `json:"sort_by"`      // 排序欄位
	SortOrder   string `json:"sort_order"`   // 排序方向 (asc/desc)
	MinPrice    float64 `json:"min_price"`   // 最低價格
	MaxPrice    float64 `json:"max_price"`   // 最高價格
	IsActive    *bool   `json:"is_active"`   // 是否交易中
}

// StockCategory 股票分類模型
type StockCategory struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`         // 分類名稱
	Code     string `json:"code" db:"code"`         // 分類代碼
	Sort     int    `json:"sort" db:"sort"`         // 排序
	IsActive bool   `json:"is_active" db:"is_active"` // 是否啟用
}

// StockRepository 股票資料庫操作介面
type StockRepository interface {
	// 基本CRUD操作
	GetStocks(filter StockFilter, pagination Pagination) ([]StockWithPrice, error)
	GetStockByCode(code string) (*StockWithPrice, error)
	GetStockByID(id int) (*StockWithPrice, error)
	CreateStock(stock *Stock) error
	UpdateStock(stock *Stock) error
	DeleteStock(id int) error
	
	// 價格相關操作
	UpdateStockPrice(price *StockPrice) error
	GetStockPrice(code string) (*StockPrice, error)
	GetStockPrices(codes []string) ([]StockPrice, error)
	
	// 分類相關操作
	GetCategories() ([]StockCategory, error)
	GetCategoryByCode(code string) (*StockCategory, error)
	
	// 統計相關操作
	GetStockCount(filter StockFilter) (int, error)
	GetMarketStats() (map[string]interface{}, error)
}

// StockRepositoryImpl 股票資料庫操作實作
type StockRepositoryImpl struct {
	db *sql.DB
}

// NewStockRepository 創建股票資料庫操作實例
func NewStockRepository(db *sql.DB) StockRepository {
	return &StockRepositoryImpl{db: db}
}

// GetStocks 獲取股票列表（含分頁和篩選）
func (r *StockRepositoryImpl) GetStocks(filter StockFilter, pagination Pagination) ([]StockWithPrice, error) {
	query := `
		SELECT s.id, s.code, s.name, s.category, s.market, s.is_active, s.created_at, s.updated_at,
		       sp.price, sp.open_price, sp.high_price, sp.low_price, sp.close_price,
		       sp.volume, sp.amount, sp.change, sp.change_percent, sp.updated_at as price_updated_at
		FROM stocks s
		LEFT JOIN stock_prices sp ON s.code = sp.stock_code
		WHERE 1=1
	`
	
	args := []interface{}{}
	argIndex := 1
	
	// 添加篩選條件
	if filter.Category != "" {
		query += " AND s.category = $" + strconv.Itoa(argIndex)
		args = append(args, filter.Category)
		argIndex++
	}
	
	if filter.Market != "" {
		query += " AND s.market = $" + strconv.Itoa(argIndex)
		args = append(args, filter.Market)
		argIndex++
	}
	
	if filter.Search != "" {
		query += " AND (s.code LIKE $" + strconv.Itoa(argIndex) + " OR s.name LIKE $" + strconv.Itoa(argIndex+1) + ")"
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}
	
	if filter.IsActive != nil {
		query += " AND s.is_active = $" + strconv.Itoa(argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	// 添加排序
	if filter.SortBy != "" {
		order := "ASC"
		if filter.SortOrder == "desc" {
			order = "DESC"
		}
		
		switch filter.SortBy {
		case "price":
			query += " ORDER BY sp.price " + order
		case "change_percent":
			query += " ORDER BY sp.change_percent " + order
		case "volume":
			query += " ORDER BY sp.volume " + order
		case "name":
			query += " ORDER BY s.name " + order
		case "code":
			query += " ORDER BY s.code " + order
		default:
			query += " ORDER BY s.code ASC"
		}
	} else {
		query += " ORDER BY s.code ASC"
	}
	
	// 添加分頁
	offset := (pagination.CurrentPage - 1) * pagination.PerPage
	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, pagination.PerPage, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var stocks []StockWithPrice
	for rows.Next() {
		var stock StockWithPrice
		var price StockPrice
		var priceUpdatedAt *time.Time
		
		var priceValue, openPrice, highPrice, lowPrice, closePrice, amount, change, changePercent sql.NullFloat64
		var volume sql.NullInt64
		
		err := rows.Scan(
			&stock.ID, &stock.Code, &stock.Name, &stock.Category, &stock.Market, &stock.IsActive,
			&stock.CreatedAt, &stock.UpdatedAt,
			&priceValue, &openPrice, &highPrice, &lowPrice, &closePrice,
			&volume, &amount, &change, &changePercent, &priceUpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if priceUpdatedAt != nil {
			price.UpdatedAt = *priceUpdatedAt
			
			// 處理NULL值
			if priceValue.Valid {
				price.Price = priceValue.Float64
			}
			if openPrice.Valid {
				price.OpenPrice = openPrice.Float64
			}
			if highPrice.Valid {
				price.HighPrice = highPrice.Float64
			}
			if lowPrice.Valid {
				price.LowPrice = lowPrice.Float64
			}
			if closePrice.Valid {
				price.ClosePrice = closePrice.Float64
			}
			if volume.Valid {
				price.Volume = volume.Int64
			}
			if amount.Valid {
				price.Amount = amount.Float64
			}
			if change.Valid {
				price.Change = change.Float64
			}
			if changePercent.Valid {
				price.ChangePercent = changePercent.Float64
			}
			
			stock.Price = &price
		}
		
		stocks = append(stocks, stock)
	}
	
	return stocks, nil
}

// GetStockByCode 根據股票代碼獲取股票資訊
func (r *StockRepositoryImpl) GetStockByCode(code string) (*StockWithPrice, error) {
	query := `
		SELECT s.id, s.code, s.name, s.category, s.market, s.is_active, s.created_at, s.updated_at,
		       sp.price, sp.open_price, sp.high_price, sp.low_price, sp.close_price,
		       sp.volume, sp.amount, sp.change, sp.change_percent, sp.updated_at as price_updated_at
		FROM stocks s
		LEFT JOIN stock_prices sp ON s.code = sp.stock_code
		WHERE s.code = $1
	`
	
	var stock StockWithPrice
	var price StockPrice
	var priceUpdatedAt *time.Time
	
	var priceValue, openPrice, highPrice, lowPrice, closePrice, amount, change, changePercent sql.NullFloat64
	var volume sql.NullInt64
	
	err := r.db.QueryRow(query, code).Scan(
		&stock.ID, &stock.Code, &stock.Name, &stock.Category, &stock.Market, &stock.IsActive,
		&stock.CreatedAt, &stock.UpdatedAt,
		&priceValue, &openPrice, &highPrice, &lowPrice, &closePrice,
		&volume, &amount, &change, &changePercent, &priceUpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	if priceUpdatedAt != nil {
		price.UpdatedAt = *priceUpdatedAt
		
		// 處理NULL值
		if priceValue.Valid {
			price.Price = priceValue.Float64
		}
		if openPrice.Valid {
			price.OpenPrice = openPrice.Float64
		}
		if highPrice.Valid {
			price.HighPrice = highPrice.Float64
		}
		if lowPrice.Valid {
			price.LowPrice = lowPrice.Float64
		}
		if closePrice.Valid {
			price.ClosePrice = closePrice.Float64
		}
		if volume.Valid {
			price.Volume = volume.Int64
		}
		if amount.Valid {
			price.Amount = amount.Float64
		}
		if change.Valid {
			price.Change = change.Float64
		}
		if changePercent.Valid {
			price.ChangePercent = changePercent.Float64
		}
		
		stock.Price = &price
	}
	
	return &stock, nil
}

// GetStockCount 獲取股票總數（用於分頁計算）
func (r *StockRepositoryImpl) GetStockCount(filter StockFilter) (int, error) {
	query := "SELECT COUNT(*) FROM stocks s WHERE 1=1"
	args := []interface{}{}
	argIndex := 1
	
	// 添加篩選條件
	if filter.Category != "" {
		query += " AND s.category = $" + strconv.Itoa(argIndex)
		args = append(args, filter.Category)
		argIndex++
	}
	
	if filter.Market != "" {
		query += " AND s.market = $" + strconv.Itoa(argIndex)
		args = append(args, filter.Market)
		argIndex++
	}
	
	if filter.Search != "" {
		query += " AND (s.code LIKE $" + strconv.Itoa(argIndex) + " OR s.name LIKE $" + strconv.Itoa(argIndex+1) + ")"
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}
	
	if filter.IsActive != nil {
		query += " AND s.is_active = $" + strconv.Itoa(argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}
	
	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// GetCategories 獲取股票分類列表
func (r *StockRepositoryImpl) GetCategories() ([]StockCategory, error) {
	query := `
		SELECT id, name, code, sort, is_active 
		FROM stock_categories 
		WHERE is_active = true 
		ORDER BY sort ASC, name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []StockCategory
	for rows.Next() {
		var category StockCategory
		err := rows.Scan(&category.ID, &category.Name, &category.Code, &category.Sort, &category.IsActive)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	
	return categories, nil
}

// 其他方法的實作...
func (r *StockRepositoryImpl) GetStockByID(id int) (*StockWithPrice, error) {
	// 實作根據ID獲取股票
	return nil, nil
}

func (r *StockRepositoryImpl) CreateStock(stock *Stock) error {
	// 實作創建股票
	return nil
}

func (r *StockRepositoryImpl) UpdateStock(stock *Stock) error {
	// 實作更新股票
	return nil
}

func (r *StockRepositoryImpl) DeleteStock(id int) error {
	// 實作刪除股票
	return nil
}

func (r *StockRepositoryImpl) UpdateStockPrice(price *StockPrice) error {
	// 檢查股票是否存在
	var stockID int
	err := r.db.QueryRow("SELECT id FROM stocks WHERE code = ?", price.StockCode).Scan(&stockID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("股票代碼 %s 不存在", price.StockCode)
		}
		return fmt.Errorf("查詢股票失敗: %w", err)
	}
	
	// 先刪除現有的價格記錄，然後插入新的
	_, err = r.db.Exec("DELETE FROM stock_prices WHERE stock_code = ?", price.StockCode)
	if err != nil {
		return fmt.Errorf("刪除舊價格記錄失敗: %w", err)
	}
	
	// 插入新的價格記錄
	query := `
		INSERT INTO stock_prices (stock_code, price, open_price, high_price, low_price, close_price, volume, amount, change, change_percent, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = r.db.Exec(query, 
		price.StockCode, price.Price, price.OpenPrice, price.HighPrice, 
		price.LowPrice, price.ClosePrice, price.Volume, price.Amount, 
		price.Change, price.ChangePercent, price.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("插入股票價格失敗: %w", err)
	}
	
	return nil
}

func (r *StockRepositoryImpl) GetStockPrice(code string) (*StockPrice, error) {
	// 實作獲取股票價格
	return nil, nil
}

func (r *StockRepositoryImpl) GetStockPrices(codes []string) ([]StockPrice, error) {
	// 實作獲取多個股票價格
	return nil, nil
}

func (r *StockRepositoryImpl) GetCategoryByCode(code string) (*StockCategory, error) {
	// 實作根據代碼獲取分類
	return nil, nil
}

func (r *StockRepositoryImpl) GetMarketStats() (map[string]interface{}, error) {
	// 實作獲取市場統計
	return nil, nil
}
