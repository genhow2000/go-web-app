package services

import (
	"database/sql"
	"go-simple-app/models"
)

// CartService 購物車業務邏輯服務
type CartService struct {
	cartRepo    *models.CartRepository
	productRepo *models.ProductRepository
}

// NewCartService 創建購物車服務
func NewCartService(db *sql.DB) *CartService {
	return &CartService{
		cartRepo:    models.NewCartRepository(db),
		productRepo: models.NewProductRepository(db),
	}
}

// AddToCart 添加商品到購物車
func (s *CartService) AddToCart(customerID, productID, quantity int) error {
	// 參數驗證
	if customerID <= 0 {
		return &models.CartError{Code: "INVALID_CUSTOMER_ID", Message: "無效的客戶ID"}
	}
	if productID <= 0 {
		return &models.CartError{Code: "INVALID_PRODUCT_ID", Message: "無效的商品ID"}
	}
	if quantity <= 0 {
		return &models.CartError{Code: "INVALID_QUANTITY", Message: "數量必須大於0"}
	}

	// 檢查商品是否存在
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.CartError{Code: "PRODUCT_NOT_FOUND", Message: "商品不存在"}
		}
		return err
	}

	// 檢查商品是否可用
	if !product.IsActive {
		return &models.CartError{Code: "PRODUCT_NOT_AVAILABLE", Message: "商品已下架"}
	}

	// 檢查庫存
	if product.Stock < quantity {
		return &models.CartError{Code: "INSUFFICIENT_STOCK", Message: "庫存不足"}
	}

	// 添加到購物車
	return s.cartRepo.AddItemToCart(customerID, productID, quantity)
}

// UpdateCartItem 更新購物車商品數量
func (s *CartService) UpdateCartItem(customerID, productID, quantity int) error {
	// 參數驗證
	if customerID <= 0 {
		return &models.CartError{Code: "INVALID_CUSTOMER_ID", Message: "無效的客戶ID"}
	}
	if productID <= 0 {
		return &models.CartError{Code: "INVALID_PRODUCT_ID", Message: "無效的商品ID"}
	}
	if quantity < 0 {
		return &models.CartError{Code: "INVALID_QUANTITY", Message: "數量不能為負數"}
	}

	// 如果數量為0，移除商品
	if quantity == 0 {
		return s.RemoveFromCart(customerID, productID)
	}

	// 檢查商品是否仍然可用
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.CartError{Code: "PRODUCT_NOT_FOUND", Message: "商品不存在"}
		}
		return err
	}

	if !product.IsActive {
		return &models.CartError{Code: "PRODUCT_NOT_AVAILABLE", Message: "商品已下架"}
	}

	// 檢查庫存
	if product.Stock < quantity {
		return &models.CartError{Code: "INSUFFICIENT_STOCK", Message: "庫存不足"}
	}

	// 更新購物車項目
	return s.cartRepo.UpdateCartItem(customerID, productID, quantity)
}

// RemoveFromCart 從購物車移除商品
func (s *CartService) RemoveFromCart(customerID, productID int) error {
	// 參數驗證
	if customerID <= 0 {
		return &models.CartError{Code: "INVALID_CUSTOMER_ID", Message: "無效的客戶ID"}
	}
	if productID <= 0 {
		return &models.CartError{Code: "INVALID_PRODUCT_ID", Message: "無效的商品ID"}
	}

	return s.cartRepo.RemoveItemFromCart(customerID, productID)
}

// GetCart 獲取購物車
func (s *CartService) GetCart(customerID int) (*models.Cart, error) {
	// 參數驗證
	if customerID <= 0 {
		return nil, &models.CartError{Code: "INVALID_CUSTOMER_ID", Message: "無效的客戶ID"}
	}

	cart, err := s.cartRepo.GetCartByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	// 重新計算總價（確保價格是最新的）
	s.calculateCartTotal(cart)

	return cart, nil
}

// ClearCart 清空購物車
func (s *CartService) ClearCart(customerID int) error {
	// 參數驗證
	if customerID <= 0 {
		return &models.CartError{Code: "INVALID_CUSTOMER_ID", Message: "無效的客戶ID"}
	}

	return s.cartRepo.ClearCart(customerID)
}

// GetCartItemCount 獲取購物車商品數量
func (s *CartService) GetCartItemCount(customerID int) (int, error) {
	// 參數驗證
	if customerID <= 0 {
		return 0, &models.CartError{Code: "INVALID_CUSTOMER_ID", Message: "無效的客戶ID"}
	}

	return s.cartRepo.GetCartItemCount(customerID)
}

// CalculateCartTotal 計算購物車總價
func (s *CartService) CalculateCartTotal(cart *models.Cart) float64 {
	return s.calculateCartTotal(cart)
}

// calculateCartTotal 內部方法：計算購物車總價
func (s *CartService) calculateCartTotal(cart *models.Cart) float64 {
	total := 0.0
	for _, item := range cart.Items {
		// 使用商品的最新價格
		if item.Product != nil {
			total += item.Product.Price * float64(item.Quantity)
		} else {
			total += item.Price * float64(item.Quantity)
		}
	}
	cart.TotalPrice = total
	return total
}

// ValidateCartItems 驗證購物車項目
func (s *CartService) ValidateCartItems(cart *models.Cart) []string {
	var errors []string

	for _, item := range cart.Items {
		// 檢查商品是否仍然存在
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			errors = append(errors, "商品 "+item.Product.Name+" 已不存在")
			continue
		}

		// 檢查商品是否仍然可用
		if !product.IsActive {
			errors = append(errors, "商品 "+product.Name+" 已下架")
			continue
		}

		// 檢查庫存
		if product.Stock < item.Quantity {
			errors = append(errors, "商品 "+product.Name+" 庫存不足，當前庫存："+string(rune(product.Stock)))
		}

		// 檢查價格是否變更
		if product.Price != item.Price {
			errors = append(errors, "商品 "+product.Name+" 價格已變更")
		}
	}

	return errors
}

// SyncCartWithProductPrices 同步購物車商品價格
func (s *CartService) SyncCartWithProductPrices(cart *models.Cart) error {
	for i, item := range cart.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			continue // 跳過不存在的商品
		}

		// 更新價格
		cart.Items[i].Price = product.Price
		if cart.Items[i].Product != nil {
			cart.Items[i].Product.Price = product.Price
		}
	}

	// 重新計算總價
	s.calculateCartTotal(cart)
	return nil
}

// GetCartSummary 獲取購物車摘要
func (s *CartService) GetCartSummary(customerID int) (map[string]interface{}, error) {
	cart, err := s.GetCart(customerID)
	if err != nil {
		return nil, err
	}

	// 驗證購物車項目
	validationErrors := s.ValidateCartItems(cart)

	// 同步價格
	s.SyncCartWithProductPrices(cart)

	summary := map[string]interface{}{
		"item_count":        cart.ItemCount,
		"total_price":       cart.TotalPrice,
		"items":             cart.Items,
		"validation_errors": validationErrors,
		"has_errors":        len(validationErrors) > 0,
	}

	return summary, nil
}
