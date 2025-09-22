package routes

import (
	"go-simple-app/controllers"
	"go-simple-app/middleware"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

// SetupCartRoutes 設置購物車路由
func SetupCartRoutes(router *gin.Engine, cartService *services.CartService, unifiedAuthService *services.UnifiedAuthService) {
	// 創建購物車控制器
	cartController := controllers.NewCartController(cartService)

	// 購物車API路由組（需要客戶端認證）
	cartAPI := router.Group("/api/cart")
	cartAPI.Use(middleware.UnifiedAuthMiddleware(unifiedAuthService)) // 使用統一認證中間件
	cartAPI.Use(middleware.CustomerMiddleware()) // 只允許客戶端訪問
	{
		// 獲取購物車
		cartAPI.GET("", cartController.GetCart)
		
		// 獲取購物車摘要
		cartAPI.GET("/summary", cartController.GetCartSummary)
		
		// 獲取購物車商品數量
		cartAPI.GET("/count", cartController.GetCartItemCount)
		
		// 清空購物車
		cartAPI.DELETE("", cartController.ClearCart)
		
		// 購物車項目管理
		cartItems := cartAPI.Group("/items")
		{
			// 添加商品到購物車
			cartItems.POST("", cartController.AddToCart)
			
			// 更新購物車商品數量
			cartItems.PUT("/:productId", cartController.UpdateCartItem)
			
			// 從購物車移除商品
			cartItems.DELETE("/:productId", cartController.RemoveFromCart)
		}
	}
}
