package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"go-simple-app/models"
)

// StockService 股票服務
type StockService struct {
	stockRepo models.StockRepository
	httpClient *http.Client
	ticker    *time.Ticker
	stopChan  chan bool
}

// NewStockService 創建股票服務實例
func NewStockService(stockRepo models.StockRepository) *StockService {
	return &StockService{
		stockRepo: stockRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		stopChan: make(chan bool),
	}
}

// GetStocksWithPagination 獲取股票列表（含分頁）
func (s *StockService) GetStocksWithPagination(filter models.StockFilter, page, limit int) (*models.StockListResponse, error) {
	// 計算分頁參數
	pagination := models.Pagination{
		CurrentPage: page,
		PerPage:     limit,
	}
	
	// 獲取總筆數
	totalCount, err := s.stockRepo.GetStockCount(filter)
	if err != nil {
		return nil, err
	}
	
	// 計算總頁數
	totalPages := (totalCount + limit - 1) / limit
	pagination.TotalCount = totalCount
	pagination.TotalPages = totalPages
	pagination.HasNext = page < totalPages
	pagination.HasPrev = page > 1
	
	// 獲取股票列表
	stocks, err := s.stockRepo.GetStocks(filter, pagination)
	if err != nil {
		return nil, err
	}
	
	return &models.StockListResponse{
		Stocks:     stocks,
		Pagination: pagination,
		TotalCount: totalCount,
	}, nil
}

// GetStockByCode 根據股票代碼獲取股票資訊
func (s *StockService) GetStockByCode(code string) (*models.StockWithPrice, error) {
	return s.stockRepo.GetStockByCode(code)
}

// GetStockCategories 獲取股票分類列表
func (s *StockService) GetStockCategories() ([]models.StockCategory, error) {
	return s.stockRepo.GetCategories()
}

// UpdateStockPricesFromAPI 從外部API更新股票價格
func (s *StockService) UpdateStockPricesFromAPI() error {
	// 這裡可以整合證交所或其他股票API
	// 目前先實作一個模擬的更新邏輯
	
	// 獲取所有活躍的股票代碼
	activeFilter := models.StockFilter{
		IsActive: &[]bool{true}[0],
	}
	
	// 獲取前100支股票進行價格更新
	pagination := models.Pagination{
		CurrentPage: 1,
		PerPage:     100,
	}
	
	stocks, err := s.stockRepo.GetStocks(activeFilter, pagination)
	if err != nil {
		return err
	}
	
	// 模擬價格更新（實際應該從API獲取）
	for _, stock := range stocks {
		if stock.Price != nil {
			// 模擬價格變動
			priceChange := (float64(time.Now().Unix()%100) - 50) / 1000 // -5% 到 +5% 的變動
			newPrice := stock.Price.Price * (1 + priceChange)
			
			// 更新價格資訊
			updatedPrice := *stock.Price
			updatedPrice.Price = newPrice
			updatedPrice.Change = newPrice - stock.Price.ClosePrice
			updatedPrice.ChangePercent = (updatedPrice.Change / stock.Price.ClosePrice) * 100
			updatedPrice.UpdatedAt = time.Now()
			
			// 保存到資料庫
			err := s.stockRepo.UpdateStockPrice(&updatedPrice)
			if err != nil {
				fmt.Printf("更新股票 %s 價格失敗: %v\n", stock.Code, err)
			}
		}
	}
	
	return nil
}

// GetMarketStats 獲取市場統計資訊
func (s *StockService) GetMarketStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 獲取所有股票數據
	filter := models.StockFilter{}
	pagination := models.Pagination{
		CurrentPage: 1,
		PerPage:     1000,
	}
	stocks, err := s.stockRepo.GetStocks(filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("獲取股票列表失敗: %w", err)
	}

	// 計算統計數據
	totalStocks := len(stocks)
	var totalAmount float64
	var stocksWithPrice int
	var totalVolume int64
	var totalChange float64
	var positiveCount, negativeCount int

	for _, stock := range stocks {
		if stock.Price != nil {
			stocksWithPrice++
			totalAmount += stock.Price.Amount
			totalVolume += stock.Price.Volume
			totalChange += stock.Price.Change
			
			if stock.Price.Change > 0 {
				positiveCount++
			} else if stock.Price.Change < 0 {
				negativeCount++
			}
		}
	}

	// 獲取真實的台灣加權指數
	taiex, taiexChange, taiexChangePercent, err := s.getTaiwanIndex()
	if err != nil {
		// 如果獲取失敗，使用模擬數據
		if stocksWithPrice > 0 {
			avgChange := totalChange / float64(stocksWithPrice)
			taiex = 17500.0 + (avgChange * 100)
			taiexChange = avgChange * 100
			taiexChangePercent = (taiexChange / taiex) * 100
		} else {
			taiex = 17500.0
			taiexChange = 0
			taiexChangePercent = 0
		}
	}

	// 獲取上櫃指數
	otcIndex, otcChange, otcChangePercent, err := s.getOTCIndex()
	if err != nil {
		// 如果獲取失敗，使用模擬數據
		otcIndex = 220.0
		otcChange = 0
		otcChangePercent = 0
	}

	// 獲取各市場股票數量
	tseFilter := models.StockFilter{Market: "TSE"}
	otcFilter := models.StockFilter{Market: "OTC"}
	
	tseCount, err := s.stockRepo.GetStockCount(tseFilter)
	if err != nil {
		return nil, err
	}
	
	otcCount, err := s.stockRepo.GetStockCount(otcFilter)
	if err != nil {
		return nil, err
	}
	
	stats["taiex"] = taiex
	stats["taiexChange"] = taiexChange
	stats["taiexChangePercent"] = taiexChangePercent
	stats["otc_index"] = otcIndex
	stats["otc_change"] = otcChange
	stats["otc_change_percent"] = otcChangePercent
	stats["tse_count"] = tseCount
	stats["otc_count"] = otcCount
	stats["total_count"] = totalStocks
	stats["stocks_with_price"] = stocksWithPrice
	stats["our_stocks_amount"] = totalAmount / 100000000 // 我們44支股票的合計
	stats["total_volume"] = totalVolume
	stats["positive_count"] = positiveCount
	stats["negative_count"] = negativeCount
	stats["last_updated"] = time.Now().Format("2006-01-02 15:04:05")
	
	return stats, nil
}

// getTaiwanIndex 獲取台灣加權指數
func (s *StockService) getTaiwanIndex() (float64, float64, float64, error) {
	// 使用台灣證交所的加權指數API
	url := "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_t00.tw&json=1&delay=0"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, 0, err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://mis.twse.com.tw/")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, fmt.Errorf("API返回錯誤狀態碼: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, err
	}
	
	var apiResp TSEAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return 0, 0, 0, err
	}
	
	if apiResp.RTCode != "0000" || len(apiResp.MsgArray) == 0 {
		return 0, 0, 0, fmt.Errorf("API返回錯誤: %s", apiResp.RTMessage)
	}
	
	// 解析加權指數數據
	indexData := apiResp.MsgArray[0]
	tse := NewTSEAPIService()
	
	// 加權指數的字段映射
	index := tse.ParseFloat(indexData.Price) // 現價
	prevClose := tse.ParseFloat(indexData.ClosePrice) // 昨收
	change := index - prevClose
	changePercent := (change / prevClose) * 100
	
	return index, change, changePercent, nil
}

// getOTCIndex 獲取上櫃指數
func (s *StockService) getOTCIndex() (float64, float64, float64, error) {
	// 使用台灣證交所的上櫃指數API
	url := "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=otc_o00.tw&json=1&delay=0"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, 0, err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://mis.twse.com.tw/")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, fmt.Errorf("API返回錯誤狀態碼: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, err
	}
	
	var apiResp TSEAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return 0, 0, 0, err
	}
	
	if apiResp.RTCode != "0000" || len(apiResp.MsgArray) == 0 {
		return 0, 0, 0, fmt.Errorf("API返回錯誤: %s", apiResp.RTMessage)
	}
	
	// 解析上櫃指數數據
	indexData := apiResp.MsgArray[0]
	tse := NewTSEAPIService()
	
	// 上櫃指數的字段映射
	index := tse.ParseFloat(indexData.Price) // 現價
	prevClose := tse.ParseFloat(indexData.ClosePrice) // 昨收
	change := index - prevClose
	changePercent := (change / prevClose) * 100
	
	return index, change, changePercent, nil
}

// getMarketTotalAmount 獲取市場總成交金額
func (s *StockService) getMarketTotalAmount() (float64, error) {
	// 使用台灣證交所的市場統計API
	url := "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_t00.tw&json=1&delay=0"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://mis.twse.com.tw/")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API返回錯誤狀態碼: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	
	var apiResp TSEAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return 0, err
	}
	
	if apiResp.RTCode != "0000" || len(apiResp.MsgArray) == 0 {
		return 0, fmt.Errorf("API返回錯誤: %s", apiResp.RTMessage)
	}
	
	// 解析市場總成交金額
	indexData := apiResp.MsgArray[0]
	tse := NewTSEAPIService()
	
	// 從加權指數數據中獲取總成交金額
	totalAmount := tse.ParseFloat(indexData.Amount) // 總成交金額（單位：千元）
	
	// 轉換為億元
	return totalAmount / 100000000, nil
}

// SearchStocks 搜尋股票
func (s *StockService) SearchStocks(keyword string, page, limit int) (*models.StockListResponse, error) {
	filter := models.StockFilter{
		Search: keyword,
	}
	
	return s.GetStocksWithPagination(filter, page, limit)
}

// GetStocksByCategory 根據分類獲取股票
func (s *StockService) GetStocksByCategory(category string, page, limit int) (*models.StockListResponse, error) {
	filter := models.StockFilter{
		Category: category,
	}
	
	return s.GetStocksWithPagination(filter, page, limit)
}

// GetTopGainers 獲取漲幅榜
func (s *StockService) GetTopGainers(limit int) ([]models.StockWithPrice, error) {
	filter := models.StockFilter{
		SortBy:    "change_percent",
		SortOrder: "desc",
	}
	
	pagination := models.Pagination{
		CurrentPage: 1,
		PerPage:     limit,
	}
	
	stocks, err := s.stockRepo.GetStocks(filter, pagination)
	if err != nil {
		return nil, err
	}
	
	// 只返回有價格資訊且漲幅為正的股票
	var gainers []models.StockWithPrice
	for _, stock := range stocks {
		if stock.Price != nil && stock.Price.ChangePercent > 0 {
			gainers = append(gainers, stock)
		}
	}
	
	return gainers, nil
}

// GetTopLosers 獲取跌幅榜
func (s *StockService) GetTopLosers(limit int) ([]models.StockWithPrice, error) {
	filter := models.StockFilter{
		SortBy:    "change_percent",
		SortOrder: "asc",
	}
	
	pagination := models.Pagination{
		CurrentPage: 1,
		PerPage:     limit,
	}
	
	stocks, err := s.stockRepo.GetStocks(filter, pagination)
	if err != nil {
		return nil, err
	}
	
	// 只返回有價格資訊且跌幅為負的股票
	var losers []models.StockWithPrice
	for _, stock := range stocks {
		if stock.Price != nil && stock.Price.ChangePercent < 0 {
			losers = append(losers, stock)
		}
	}
	
	return losers, nil
}

// GetTopVolume 獲取成交量榜
func (s *StockService) GetTopVolume(limit int) ([]models.StockWithPrice, error) {
	filter := models.StockFilter{
		SortBy:    "volume",
		SortOrder: "desc",
	}
	
	pagination := models.Pagination{
		CurrentPage: 1,
		PerPage:     limit,
	}
	
	stocks, err := s.stockRepo.GetStocks(filter, pagination)
	if err != nil {
		return nil, err
	}
	
	// 只返回有價格資訊且成交量大於0的股票
	var volumeStocks []models.StockWithPrice
	for _, stock := range stocks {
		if stock.Price != nil && stock.Price.Volume > 0 {
			volumeStocks = append(volumeStocks, stock)
		}
	}
	
	return volumeStocks, nil
}

// UpdateStockPricesFromTSE 從台灣證交所更新股票價格
func (s *StockService) UpdateStockPricesFromTSE() error {
	return s.UpdateStockPricesFromTSEWithForce(false)
}

// UpdateStockPricesFromTSEWithForce 從台灣證交所更新股票價格（可強制更新）
func (s *StockService) UpdateStockPricesFromTSEWithForce(forceUpdate bool) error {
	if !forceUpdate && !s.isTradingTime(time.Now()) {
		fmt.Println("非交易時間，跳過更新")
		return nil
	}
	
	fmt.Println("開始從台灣證交所更新股票價格...")
	
	// 獲取所有股票代碼
	filter := models.StockFilter{}
	pagination := models.Pagination{
		CurrentPage: 1,
		PerPage:     1000,
	}
	stocks, err := s.stockRepo.GetStocks(filter, pagination)
	if err != nil {
		return fmt.Errorf("獲取股票列表失敗: %w", err)
	}
	
	fmt.Printf("找到 %d 支股票需要更新\n", len(stocks))
	
	if len(stocks) == 0 {
		return nil
	}
	
	// 提取股票代碼
	var codes []string
	for _, stock := range stocks {
		codes = append(codes, stock.Code)
	}
	
	// 分批處理（每次最多20個股票代碼）
	batchSize := 20
	tseAPI := NewTSEAPIService()
	
	for i := 0; i < len(codes); i += batchSize {
		end := i + batchSize
		if end > len(codes) {
			end = len(codes)
		}
		
		batchCodes := codes[i:end]
		tseData, err := tseAPI.FetchStockData(batchCodes)
		if err != nil {
			fmt.Printf("獲取股票數據失敗 (批次 %d): %v\n", i/batchSize+1, err)
			continue
		}
		
		// 更新每個股票的價格
		for _, data := range tseData {
			stockPrice := ConvertTSEToStockPrice(data)
			fmt.Printf("股票 %s: 原始價格=%s, 解析後價格=%.2f\n", data.Code, data.Price, stockPrice.Price)
			if stockPrice.Price > 0 { // 只更新有價格的股票
				err := s.stockRepo.UpdateStockPrice(stockPrice)
				if err != nil {
					fmt.Printf("更新股票 %s 價格失敗: %v\n", data.Code, err)
				} else {
					fmt.Printf("成功更新股票 %s 價格: %.2f\n", data.Code, stockPrice.Price)
				}
			} else {
				fmt.Printf("跳過股票 %s (價格為0或無效)\n", data.Code)
			}
		}
		
		// 避免請求過於頻繁
		time.Sleep(1 * time.Second)
	}
	
	return nil
}

// TSEAPIService 台灣證交所API服務
type TSEAPIService struct {
	baseURL string
	client  *http.Client
}

// NewTSEAPIService 創建證交所API服務
func NewTSEAPIService() *TSEAPIService {
	return &TSEAPIService{
		baseURL: "https://mis.twse.com.tw/stock/api",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TSEStockData 證交所股票數據結構
type TSEStockData struct {
	Code         string  `json:"c"`   // 股票代碼
	Name         string  `json:"n"`   // 股票名稱
	Price        string  `json:"z"`   // 現價
	OpenPrice    string  `json:"o"`   // 開盤價
	HighPrice    string  `json:"h"`   // 最高價
	LowPrice     string  `json:"l"`   // 最低價
	ClosePrice   string  `json:"y"`   // 昨收價
	Volume       string  `json:"v"`   // 成交量
	Amount       string  `json:"a"`   // 成交金額
	Change       string  `json:"ch"`  // 漲跌（實際上是股票代碼）
	ChangePercent string `json:"%"`   // 漲跌幅（實際上是時間）
	Time         string  `json:"t"`   // 時間
}

// TSEAPIResponse 證交所API回應結構
type TSEAPIResponse struct {
	MsgArray []TSEStockData `json:"msgArray"`
	UserDelay int          `json:"userDelay"`
	RTMessage string       `json:"rtmessage"`
	Referer   string       `json:"referer"`
	QueryTime interface{}  `json:"queryTime"` // 可能是字符串或對象
	RTCode    string       `json:"rtcode"`
}

// FetchStockData 從證交所API獲取股票數據
func (t *TSEAPIService) FetchStockData(codes []string) ([]TSEStockData, error) {
	if len(codes) == 0 {
		return []TSEStockData{}, nil
	}
	
	// 構建查詢參數
	var exChList []string
	for _, code := range codes {
		// 判斷是上市還是上櫃
		if len(code) == 4 && code[0] >= '1' && code[0] <= '9' {
			exChList = append(exChList, fmt.Sprintf("tse_%s.tw", code))
		} else {
			exChList = append(exChList, fmt.Sprintf("otc_%s.tw", code))
		}
	}
	
	exCh := strings.Join(exChList, "|")
	url := fmt.Sprintf("%s/getStockInfo.jsp?ex_ch=%s&json=1&delay=0", t.baseURL, exCh)
	
	// 發送HTTP請求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("創建請求失敗: %w", err)
	}
	
	// 設置請求頭
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://mis.twse.com.tw/")
	
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("請求失敗: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回錯誤狀態碼: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取響應失敗: %w", err)
	}
	
	var apiResp TSEAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析JSON失敗: %w", err)
	}
	
	if apiResp.RTCode != "0000" {
		return nil, fmt.Errorf("API返回錯誤: %s", apiResp.RTMessage)
	}
	
	return apiResp.MsgArray, nil
}

// ParseFloat 安全解析浮點數
func (t *TSEAPIService) ParseFloat(s string) float64 {
	if s == "" || s == "--" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// ParseInt64 安全解析整數
func (t *TSEAPIService) ParseInt64(s string) int64 {
	if s == "" || s == "--" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// ConvertTSEToStockPrice 將證交所數據轉換為內部股票價格模型
func ConvertTSEToStockPrice(tseData TSEStockData) *models.StockPrice {
	tse := NewTSEAPIService()
	
	// 計算漲跌點數和漲跌幅
	currentPrice := tse.ParseFloat(tseData.Price)
	prevClose := tse.ParseFloat(tseData.ClosePrice)
	change := currentPrice - prevClose
	var changePercent float64
	if prevClose > 0 {
		changePercent = (change / prevClose) * 100
	}
	
	return &models.StockPrice{
		StockCode:     tseData.Code,
		Price:         currentPrice,
		OpenPrice:     tse.ParseFloat(tseData.OpenPrice),
		HighPrice:     tse.ParseFloat(tseData.HighPrice),
		LowPrice:      tse.ParseFloat(tseData.LowPrice),
		ClosePrice:    prevClose,
		Volume:        tse.ParseInt64(tseData.Volume),
		Amount:        tse.ParseFloat(tseData.Amount),
		Change:        change,
		ChangePercent: changePercent,
		UpdatedAt:     time.Now(),
	}
}

// StartAutoUpdate 開始自動更新股票價格（每5秒）
func (s *StockService) StartAutoUpdate() {
	s.ticker = time.NewTicker(5 * time.Second)
	
	go func() {
		for {
			select {
			case <-s.ticker.C:
				// 檢查是否在交易時間內（週一至週五 9:00-13:30）
				now := time.Now()
				if s.isTradingTime(now) {
					fmt.Println("開始自動更新股票價格...")
					err := s.UpdateStockPricesFromTSE()
					if err != nil {
						fmt.Printf("自動更新股票價格失敗: %v\n", err)
					} else {
						fmt.Println("股票價格更新完成")
					}
				} else {
					fmt.Println("非交易時間，跳過更新")
				}
			case <-s.stopChan:
				fmt.Println("停止自動更新股票價格")
				return
			}
		}
	}()
	
	fmt.Println("股票價格自動更新已啟動（每5秒，僅交易時間）")
}

// isTradingTime 檢查是否在交易時間內
func (s *StockService) isTradingTime(t time.Time) bool {
	// 轉換為台灣時間 (UTC+8)
	taiwanTime := t.In(time.FixedZone("CST", 8*60*60))
	
	// 檢查是否為週一至週五
	weekday := taiwanTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	// 檢查是否在交易時間內
	hour := taiwanTime.Hour()
	minute := taiwanTime.Minute()
	
	// 台灣股市交易時間：9:00-13:30（連續交易）
	if hour >= 9 && hour < 13 {
		return true
	}
	if hour == 13 && minute <= 30 {
		return true
	}
	
	return false
}

// StopAutoUpdate 停止自動更新
func (s *StockService) StopAutoUpdate() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	select {
	case s.stopChan <- true:
	default:
	}
}
