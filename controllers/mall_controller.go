package controllers

import (
	"net/http"
	"strconv"
	"go-simple-app/models"

	"github.com/gin-gonic/gin"
)

type MallController struct {
	productRepo *models.ProductRepository
}

func NewMallController(productRepo *models.ProductRepository) *MallController {
	return &MallController{
		productRepo: productRepo,
	}
}

// ShowHomepage 顯示商城首頁
func (c *MallController) ShowHomepage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "homepage.html", gin.H{
		"title": "阿和商城 - 精選商品，優質服務",
	})
}

// ShowTechShowcase 顯示技術展示頁面
func (c *MallController) ShowTechShowcase(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "tech-showcase.html", gin.H{
		"version":           "2.0.0",
		"uptime":            "24小時",
		"migration_version": "004",
	})
}

// GetCategories 獲取商品分類
func (c *MallController) GetCategories(ctx *gin.Context) {
	categories, err := c.productRepo.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取分類失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// GetFeaturedProducts 獲取精選商品
func (c *MallController) GetFeaturedProducts(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "8")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 8
	}

	products, err := c.productRepo.GetFeatured(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取精選商品失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProducts 獲取商品列表
func (c *MallController) GetProducts(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "12")
	offsetStr := ctx.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 12
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	products, err := c.productRepo.GetAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取商品失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProductsByCategory 根據分類獲取商品
func (c *MallController) GetProductsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "分類參數不能為空",
		})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "12")
	offsetStr := ctx.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 12
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	products, err := c.productRepo.GetByCategory(category, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取分類商品失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// SearchProducts 搜尋商品
func (c *MallController) SearchProducts(ctx *gin.Context) {
	keyword := ctx.Query("q")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "搜尋關鍵字不能為空",
		})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "12")
	offsetStr := ctx.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 12
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	products, err := c.productRepo.Search(keyword, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜尋商品失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProduct 獲取單個商品詳情
func (c *MallController) GetProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	product, err := c.productRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "商品不存在",
		})
		return
	}

	// 增加瀏覽次數
	go c.productRepo.IncrementViewCount(id)

	ctx.JSON(http.StatusOK, product)
}

// ShowProductPage 顯示商品詳情頁面
func (c *MallController) ShowProductPage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title":   "錯誤",
			"message": "無效的商品ID",
		})
		return
	}

	product, err := c.productRepo.GetByID(id)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "商品不存在",
			"message": "您要查看的商品不存在或已下架",
		})
		return
	}

	// 增加瀏覽次數
	go c.productRepo.IncrementViewCount(id)

	ctx.HTML(http.StatusOK, "product_detail.html", gin.H{
		"title":   product.Name,
		"product": product,
	})
}

// ShowCategoryPage 顯示分類頁面
func (c *MallController) ShowCategoryPage(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title":   "錯誤",
			"message": "分類參數不能為空",
		})
		return
	}

	ctx.HTML(http.StatusOK, "category.html", gin.H{
		"title":    category + " - 商品分類",
		"category": category,
	})
}

// ShowSearchPage 顯示搜尋結果頁面
func (c *MallController) ShowSearchPage(ctx *gin.Context) {
	keyword := ctx.Query("q")
	if keyword == "" {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title":   "錯誤",
			"message": "搜尋關鍵字不能為空",
		})
		return
	}

	ctx.HTML(http.StatusOK, "search.html", gin.H{
		"title":   "搜尋結果 - " + keyword,
		"keyword": keyword,
	})
}
