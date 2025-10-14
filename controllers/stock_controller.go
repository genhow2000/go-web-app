package controllers

import (
	"net/http"
	"strconv"
	"go-simple-app/models"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

// StockController 股票控制器
type StockController struct {
	stockService *services.StockService
}

// NewStockController 創建股票控制器
func NewStockController(stockService *services.StockService) *StockController {
	return &StockController{
		stockService: stockService,
	}
}

// GetStocks 獲取股票列表
func (sc *StockController) GetStocks(c *gin.Context) {
	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")
	market := c.Query("market")
	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	
	// 限制每頁最大筆數
	if limit > 100 {
		limit = 100
	}
	
	// 構建篩選條件
	filter := models.StockFilter{
		Category:  category,
		Market:    market,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	
	// 獲取股票列表
	result, err := sc.stockService.GetStocksWithPagination(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取股票列表失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetStock 獲取單一股票資訊
func (sc *StockController) GetStock(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "股票代碼不能為空",
		})
		return
	}
	
	stock, err := sc.stockService.GetStockByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取股票資訊失敗",
			"message": err.Error(),
		})
		return
	}
	
	if stock == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "股票不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stock,
	})
}

// UpdateStockPrices 更新股票價格（從台灣證交所）
func (sc *StockController) UpdateStockPrices(c *gin.Context) {
	err := sc.stockService.UpdateStockPricesFromTSE()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "股票價格更新成功",
	})
}

// ForceUpdateStockPrices 強制更新所有股票價格（手動更新）
func (sc *StockController) ForceUpdateStockPrices(c *gin.Context) {
	err := sc.stockService.UpdateStockPricesFromTSEWithForce(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "所有股票價格強制更新成功",
	})
}

// GetStockCategories 獲取股票分類
func (sc *StockController) GetStockCategories(c *gin.Context) {
	categories, err := sc.stockService.GetStockCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取股票分類失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}

// SearchStocks 搜尋股票
func (sc *StockController) SearchStocks(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "搜尋關鍵字不能為空",
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if limit > 100 {
		limit = 100
	}
	
	result, err := sc.stockService.SearchStocks(keyword, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜尋股票失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetStocksByCategory 根據分類獲取股票
func (sc *StockController) GetStocksByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "分類不能為空",
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if limit > 100 {
		limit = 100
	}
	
	result, err := sc.stockService.GetStocksByCategory(category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取分類股票失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetTopGainers 獲取漲幅榜
func (sc *StockController) GetTopGainers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if limit > 50 {
		limit = 50
	}
	
	stocks, err := sc.stockService.GetTopGainers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取漲幅榜失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stocks,
	})
}

// GetTopLosers 獲取跌幅榜
func (sc *StockController) GetTopLosers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if limit > 50 {
		limit = 50
	}
	
	stocks, err := sc.stockService.GetTopLosers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取跌幅榜失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stocks,
	})
}

// GetTopVolume 獲取成交量榜
func (sc *StockController) GetTopVolume(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if limit > 50 {
		limit = 50
	}
	
	stocks, err := sc.stockService.GetTopVolume(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取成交量榜失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stocks,
	})
}

// GetMarketStats 獲取市場統計
func (sc *StockController) GetMarketStats(c *gin.Context) {
	stats, err := sc.stockService.GetMarketStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取市場統計失敗",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}


// ShowStockListPage 顯示股票列表頁面
func (sc *StockController) ShowStockListPage(c *gin.Context) {
	// 這裡可以返回HTML頁面或重定向到Vue.js路由
	c.HTML(http.StatusOK, "stock_list.html", gin.H{
		"title": "台股列表",
	})
}

// ShowStockDetailPage 顯示股票詳情頁面
func (sc *StockController) ShowStockDetailPage(c *gin.Context) {
	code := c.Param("code")
	
	stock, err := sc.stockService.GetStockByCode(code)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "錯誤",
			"message": "獲取股票資訊失敗",
		})
		return
	}
	
	if stock == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "股票不存在",
			"message": "找不到指定的股票",
		})
		return
	}
	
	c.HTML(http.StatusOK, "stock_detail.html", gin.H{
		"title": stock.Name + " (" + stock.Code + ")",
		"stock": stock,
	})
}
