package controllers

import (
	"net/http"
	"strconv"
	"go-simple-app/models"

	"github.com/gin-gonic/gin"
)

type MerchantProductController struct {
	productRepo *models.ProductRepository
}

func NewMerchantProductController(productRepo *models.ProductRepository) *MerchantProductController {
	return &MerchantProductController{
		productRepo: productRepo,
	}
}

// GetMerchantProducts 獲取商戶的商品列表
func (c *MerchantProductController) GetMerchantProducts(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	// 解析查詢參數
	status := ctx.Query("status")
	search := ctx.Query("search")
	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// 獲取商品列表
	products, err := c.productRepo.GetByMerchantID(merchantID, limit, offset, status, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取商品列表失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
	})
}

// GetMerchantProductStats 獲取商戶商品統計
func (c *MerchantProductController) GetMerchantProductStats(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	// 獲取統計數據
	stats, err := c.productRepo.GetMerchantStats(merchantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取統計數據失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// CreateMerchantProduct 創建商戶商品
func (c *MerchantProductController) CreateMerchantProduct(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤: " + err.Error(),
		})
		return
	}

	// 設置商戶ID
	product.MerchantID = merchantID

	// 創建商品
	if err := c.productRepo.Create(&product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "創建商品失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品創建成功",
		"product": product,
	})
}

// GetMerchantProduct 獲取單個商戶商品
func (c *MerchantProductController) GetMerchantProduct(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	// 獲取商品ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	// 獲取商品
	product, err := c.productRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "商品不存在",
		})
		return
	}

	// 檢查商品是否屬於該商戶
	if product.MerchantID != merchantID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "無權限訪問此商品",
		})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// UpdateMerchantProduct 更新商戶商品
func (c *MerchantProductController) UpdateMerchantProduct(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	// 獲取商品ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	// 先檢查商品是否存在且屬於該商戶
	existingProduct, err := c.productRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "商品不存在",
		})
		return
	}

	if existingProduct.MerchantID != merchantID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "無權限修改此商品",
		})
		return
	}

	// 綁定更新數據
	var updateData models.Product
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤: " + err.Error(),
		})
		return
	}

	// 更新商品數據
	updateData.ID = id
	updateData.MerchantID = merchantID // 確保商戶ID不變

	if err := c.productRepo.Update(&updateData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新商品失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品更新成功",
		"product": updateData,
	})
}

// ToggleMerchantProductStatus 切換商品上架狀態
func (c *MerchantProductController) ToggleMerchantProductStatus(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	// 獲取商品ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	// 先檢查商品是否存在且屬於該商戶
	product, err := c.productRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "商品不存在",
		})
		return
	}

	if product.MerchantID != int(merchantID) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "無權限修改此商品",
		})
		return
	}

	// 切換狀態
	product.IsActive = !product.IsActive

	if err := c.productRepo.Update(product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "切換商品狀態失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品狀態更新成功",
		"is_active": product.IsActive,
	})
}

// DeleteMerchantProduct 刪除商戶商品
func (c *MerchantProductController) DeleteMerchantProduct(ctx *gin.Context) {
	// 從JWT token中獲取商戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}
	
	// 將user轉換為UserInterface以獲取ID
	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}
	
	merchantID := userInterface.GetID()

	// 獲取商品ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	// 先檢查商品是否存在且屬於該商戶
	product, err := c.productRepo.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "商品不存在",
		})
		return
	}

	if product.MerchantID != int(merchantID) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "無權限刪除此商品",
		})
		return
	}

	// 刪除商品
	if err := c.productRepo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "刪除商品失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品刪除成功",
	})
}
