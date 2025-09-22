package models

import (
	"database/sql"
	"time"
)

// Cart 購物車模型
type Cart struct {
	ID         int         `json:"id" db:"id"`
	CustomerID int         `json:"customer_id" db:"customer_id"`
	Items      []CartItem  `json:"items,omitempty"`
	TotalPrice float64     `json:"total_price"`
	ItemCount  int         `json:"item_count"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}

// CartItem 購物車項目模型
type CartItem struct {
	ID        int      `json:"id" db:"id"`
	CartID    int      `json:"cart_id" db:"cart_id"`
	ProductID int      `json:"product_id" db:"product_id"`
	Quantity  int      `json:"quantity" db:"quantity"`
	Price     float64  `json:"price" db:"price"`
	Product   *Product `json:"product,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CartRepository 購物車數據庫操作
type CartRepository struct {
	db *sql.DB
}

// NewCartRepository 創建購物車倉庫
func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetCartByCustomerID 根據客戶ID獲取購物車
func (r *CartRepository) GetCartByCustomerID(customerID int) (*Cart, error) {
	cart := &Cart{
		CustomerID: customerID,
		Items:      []CartItem{},
		TotalPrice: 0,
		ItemCount:  0,
	}

	// 獲取購物車項目
	query := `
		SELECT sc.id, sc.product_id, sc.quantity, sc.created_at, sc.updated_at,
		       p.name, p.description, p.price, p.original_price, p.category, 
		       p.sub_category, p.brand, p.sku, p.stock, p.image_url, p.images, 
		       p.tags, p.is_active, p.is_featured, p.is_on_sale, p.merchant_id,
		       p.view_count, p.sales_count, p.rating, p.review_count, p.weight,
		       p.dimensions, p.created_at as product_created_at, p.updated_at as product_updated_at
		FROM shopping_cart sc
		JOIN products p ON sc.product_id = p.id
		WHERE sc.customer_id = ? AND p.is_active = 1
		ORDER BY sc.created_at DESC`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		item := CartItem{}
		product := &Product{}
		
		err := rows.Scan(
			&item.ID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&product.Name, &product.Description, &product.Price, &product.OriginalPrice,
			&product.Category, &product.SubCategory, &product.Brand, &product.SKU,
			&product.Stock, &product.ImageURL, &product.Images, &product.Tags,
			&product.IsActive, &product.IsFeatured, &product.IsOnSale, &product.MerchantID,
			&product.ViewCount, &product.SalesCount, &product.Rating, &product.ReviewCount,
			&product.Weight, &product.Dimensions, &product.CreatedAt, &product.UpdatedAt)
		
		if err != nil {
			return nil, err
		}

		product.ID = item.ProductID
		item.Product = product
		item.Price = product.Price
		item.CartID = cart.ID // 這裡會在後續優化中處理
		
		items = append(items, item)
	}

	cart.Items = items
	cart.ItemCount = len(items)
	cart.TotalPrice = r.calculateCartTotal(items)

	return cart, nil
}

// AddItemToCart 添加商品到購物車
func (r *CartRepository) AddItemToCart(customerID, productID, quantity int) error {
	// 檢查商品是否存在且可用
	var isActive bool
	var stock int
	err := r.db.QueryRow("SELECT is_active, stock FROM products WHERE id = ?", productID).Scan(&isActive, &stock)
	if err != nil {
		return err
	}
	if !isActive {
		return ErrProductNotAvailable
	}
	if stock < quantity {
		return ErrInsufficientStock
	}

	// 檢查購物車中是否已有該商品
	var existingQuantity int
	err = r.db.QueryRow("SELECT quantity FROM shopping_cart WHERE customer_id = ? AND product_id = ?", 
		customerID, productID).Scan(&existingQuantity)
	
	if err == sql.ErrNoRows {
		// 商品不在購物車中，直接添加
		query := `INSERT INTO shopping_cart (customer_id, product_id, quantity) VALUES (?, ?, ?)`
		_, err = r.db.Exec(query, customerID, productID, quantity)
		return err
	} else if err != nil {
		return err
	}

	// 商品已在購物車中，更新數量
	newQuantity := existingQuantity + quantity
	if newQuantity > stock {
		return ErrInsufficientStock
	}

	query := `UPDATE shopping_cart SET quantity = ?, updated_at = CURRENT_TIMESTAMP WHERE customer_id = ? AND product_id = ?`
	_, err = r.db.Exec(query, newQuantity, customerID, productID)
	return err
}

// UpdateCartItem 更新購物車商品數量
func (r *CartRepository) UpdateCartItem(customerID, productID, quantity int) error {
	if quantity <= 0 {
		return r.RemoveItemFromCart(customerID, productID)
	}

	// 檢查庫存
	var stock int
	err := r.db.QueryRow("SELECT stock FROM products WHERE id = ? AND is_active = 1", productID).Scan(&stock)
	if err != nil {
		return err
	}
	if stock < quantity {
		return ErrInsufficientStock
	}

	query := `UPDATE shopping_cart SET quantity = ?, updated_at = CURRENT_TIMESTAMP WHERE customer_id = ? AND product_id = ?`
	result, err := r.db.Exec(query, quantity, customerID, productID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCartItemNotFound
	}

	return nil
}

// RemoveItemFromCart 從購物車移除商品
func (r *CartRepository) RemoveItemFromCart(customerID, productID int) error {
	query := `DELETE FROM shopping_cart WHERE customer_id = ? AND product_id = ?`
	result, err := r.db.Exec(query, customerID, productID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCartItemNotFound
	}

	return nil
}

// ClearCart 清空購物車
func (r *CartRepository) ClearCart(customerID int) error {
	query := `DELETE FROM shopping_cart WHERE customer_id = ?`
	_, err := r.db.Exec(query, customerID)
	return err
}

// GetCartItemCount 獲取購物車商品數量
func (r *CartRepository) GetCartItemCount(customerID int) (int, error) {
	var count int
	query := `SELECT COALESCE(SUM(quantity), 0) FROM shopping_cart WHERE customer_id = ?`
	err := r.db.QueryRow(query, customerID).Scan(&count)
	return count, err
}

// GetCartItemByProductID 根據商品ID獲取購物車項目
func (r *CartRepository) GetCartItemByProductID(customerID, productID int) (*CartItem, error) {
	item := &CartItem{}
	product := &Product{}
	
	query := `
		SELECT sc.id, sc.product_id, sc.quantity, sc.created_at, sc.updated_at,
		       p.name, p.description, p.price, p.original_price, p.category, 
		       p.sub_category, p.brand, p.sku, p.stock, p.image_url, p.images, 
		       p.tags, p.is_active, p.is_featured, p.is_on_sale, p.merchant_id,
		       p.view_count, p.sales_count, p.rating, p.review_count, p.weight,
		       p.dimensions, p.created_at as product_created_at, p.updated_at as product_updated_at
		FROM shopping_cart sc
		JOIN products p ON sc.product_id = p.id
		WHERE sc.customer_id = ? AND sc.product_id = ?`

	err := r.db.QueryRow(query, customerID, productID).Scan(
		&item.ID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
		&product.Name, &product.Description, &product.Price, &product.OriginalPrice,
		&product.Category, &product.SubCategory, &product.Brand, &product.SKU,
		&product.Stock, &product.ImageURL, &product.Images, &product.Tags,
		&product.IsActive, &product.IsFeatured, &product.IsOnSale, &product.MerchantID,
		&product.ViewCount, &product.SalesCount, &product.Rating, &product.ReviewCount,
		&product.Weight, &product.Dimensions, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	product.ID = item.ProductID
	item.Product = product
	item.Price = product.Price

	return item, nil
}

// calculateCartTotal 計算購物車總價
func (r *CartRepository) calculateCartTotal(items []CartItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// 錯誤定義
var (
	ErrProductNotAvailable = &CartError{Code: "PRODUCT_NOT_AVAILABLE", Message: "商品不可用"}
	ErrInsufficientStock   = &CartError{Code: "INSUFFICIENT_STOCK", Message: "庫存不足"}
	ErrCartItemNotFound    = &CartError{Code: "CART_ITEM_NOT_FOUND", Message: "購物車項目不存在"}
)

// CartError 購物車錯誤
type CartError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *CartError) Error() string {
	return e.Message
}
