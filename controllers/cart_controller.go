package controllers

import (
	"net/http"
	"strconv"
	"go-simple-app/models"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

// CartController 購物車控制器
type CartController struct {
	cartService *services.CartService
}

// NewCartController 創建購物車控制器
func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{
		cartService: cartService,
	}
}

// GetCart 獲取購物車
// @Summary 獲取購物車
// @Description 獲取當前用戶的購物車內容
// @Tags 購物車
// @Accept json
// @Produce json
// @Success 200 {object} models.Cart
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart [get]
func (c *CartController) GetCart(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
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

	customerID := userInterface.GetID()

	cart, err := c.cartService.GetCart(customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取購物車失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// AddToCart 添加商品到購物車
// @Summary 添加商品到購物車
// @Description 將指定商品添加到購物車
// @Tags 購物車
// @Accept json
// @Produce json
// @Param request body AddToCartRequest true "添加商品請求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart/items [post]
func (c *CartController) AddToCart(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}

	customerID := userInterface.GetID()

	var req AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤: " + err.Error(),
		})
		return
	}

	// 參數驗證
	if req.ProductID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "商品ID必須大於0",
		})
		return
	}

	if req.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "數量必須大於0",
		})
		return
	}

	err := c.cartService.AddToCart(customerID, req.ProductID, req.Quantity)
	if err != nil {
		// 檢查是否為購物車錯誤
		if cartErr, ok := err.(*models.CartError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": cartErr.Message,
				"code":  cartErr.Code,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "添加商品到購物車失敗: " + err.Error(),
		})
		return
	}

	// 獲取更新後的購物車項目數量
	itemCount, err := c.cartService.GetCartItemCount(customerID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "商品已添加到購物車",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "商品已添加到購物車",
		"item_count": itemCount,
	})
}

// UpdateCartItem 更新購物車商品數量
// @Summary 更新購物車商品數量
// @Description 更新購物車中指定商品的數量
// @Tags 購物車
// @Accept json
// @Produce json
// @Param productId path int true "商品ID"
// @Param request body UpdateCartItemRequest true "更新商品請求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart/items/{productId} [put]
func (c *CartController) UpdateCartItem(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}

	customerID := userInterface.GetID()

	// 獲取商品ID
	productIDStr := ctx.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	var req UpdateCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "請求參數錯誤: " + err.Error(),
		})
		return
	}

	// 參數驗證
	if req.Quantity < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "數量不能為負數",
		})
		return
	}

	err = c.cartService.UpdateCartItem(customerID, productID, req.Quantity)
	if err != nil {
		// 檢查是否為購物車錯誤
		if cartErr, ok := err.(*models.CartError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": cartErr.Message,
				"code":  cartErr.Code,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新購物車商品失敗: " + err.Error(),
		})
		return
	}

	// 獲取更新後的購物車項目數量
	itemCount, err := c.cartService.GetCartItemCount(customerID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "購物車已更新",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "購物車已更新",
		"item_count": itemCount,
	})
}

// RemoveFromCart 從購物車移除商品
// @Summary 從購物車移除商品
// @Description 從購物車中移除指定商品
// @Tags 購物車
// @Accept json
// @Produce json
// @Param productId path int true "商品ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart/items/{productId} [delete]
func (c *CartController) RemoveFromCart(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}

	customerID := userInterface.GetID()

	// 獲取商品ID
	productIDStr := ctx.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的商品ID",
		})
		return
	}

	err = c.cartService.RemoveFromCart(customerID, productID)
	if err != nil {
		// 檢查是否為購物車錯誤
		if cartErr, ok := err.(*models.CartError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": cartErr.Message,
				"code":  cartErr.Code,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "從購物車移除商品失敗: " + err.Error(),
		})
		return
	}

	// 獲取更新後的購物車項目數量
	itemCount, err := c.cartService.GetCartItemCount(customerID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "商品已從購物車移除",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "商品已從購物車移除",
		"item_count": itemCount,
	})
}

// ClearCart 清空購物車
// @Summary 清空購物車
// @Description 清空當前用戶的購物車
// @Tags 購物車
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart [delete]
func (c *CartController) ClearCart(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}

	customerID := userInterface.GetID()

	err := c.cartService.ClearCart(customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "清空購物車失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "購物車已清空",
	})
}

// GetCartSummary 獲取購物車摘要
// @Summary 獲取購物車摘要
// @Description 獲取購物車的摘要信息，包括驗證錯誤
// @Tags 購物車
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart/summary [get]
func (c *CartController) GetCartSummary(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}

	customerID := userInterface.GetID()

	summary, err := c.cartService.GetCartSummary(customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取購物車摘要失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetCartItemCount 獲取購物車商品數量
// @Summary 獲取購物車商品數量
// @Description 獲取購物車中商品的總數量
// @Tags 購物車
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart/count [get]
func (c *CartController) GetCartItemCount(ctx *gin.Context) {
	// 從JWT token中獲取用戶ID
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權訪問",
		})
		return
	}

	userInterface, ok := user.(models.UserInterface)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "用戶信息格式錯誤",
		})
		return
	}

	customerID := userInterface.GetID()

	itemCount, err := c.cartService.GetCartItemCount(customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "獲取購物車商品數量失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"item_count": itemCount,
	})
}

// 請求結構體
type AddToCartRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}
