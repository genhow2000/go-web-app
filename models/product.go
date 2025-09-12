package models

import (
	"database/sql"
	"time"
)

// Product 商品模型
type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	OriginalPrice *float64 `json:"original_price,omitempty" db:"original_price"`
	Category    string    `json:"category" db:"category"`
	SubCategory *string   `json:"sub_category,omitempty" db:"sub_category"`
	Brand       *string   `json:"brand,omitempty" db:"brand"`
	SKU         *string   `json:"sku,omitempty" db:"sku"`
	Stock       int       `json:"stock" db:"stock"`
	ImageURL    *string   `json:"image_url,omitempty" db:"image_url"`
	Images      *string   `json:"images,omitempty" db:"images"` // JSON 格式存儲多張圖片
	Tags        *string   `json:"tags,omitempty" db:"tags"`     // JSON 格式存儲標籤
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsFeatured  bool      `json:"is_featured" db:"is_featured"`   // 是否為精選商品
	IsOnSale    bool      `json:"is_on_sale" db:"is_on_sale"`     // 是否在促銷
	MerchantID  int       `json:"merchant_id" db:"merchant_id"`
	ViewCount   int       `json:"view_count" db:"view_count"`
	SalesCount  int       `json:"sales_count" db:"sales_count"`
	Rating      float64   `json:"rating" db:"rating"`
	ReviewCount int       `json:"review_count" db:"review_count"`
	Weight      *float64  `json:"weight,omitempty" db:"weight"`
	Dimensions  *string   `json:"dimensions,omitempty" db:"dimensions"` // JSON 格式存儲尺寸
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductRepository 商品數據庫操作
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository 創建商品倉庫
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create 創建商品
func (r *ProductRepository) Create(product *Product) error {
	query := `
		INSERT INTO products (name, description, price, original_price, category, sub_category, 
		                     brand, sku, stock, image_url, images, tags, is_active, is_featured, 
		                     is_on_sale, merchant_id, view_count, sales_count, rating, review_count, 
		                     weight, dimensions) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, product.Name, product.Description, product.Price,
		product.OriginalPrice, product.Category, product.SubCategory, product.Brand,
		product.SKU, product.Stock, product.ImageURL, product.Images, product.Tags,
		product.IsActive, product.IsFeatured, product.IsOnSale, product.MerchantID,
		product.ViewCount, product.SalesCount, product.Rating, product.ReviewCount,
		product.Weight, product.Dimensions)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)
	
	// 獲取創建時間
	query = `SELECT created_at, updated_at FROM products WHERE id = ?`
	err = r.db.QueryRow(query, product.ID).Scan(&product.CreatedAt, &product.UpdatedAt)
	
	return err
}

// GetByID 根據ID獲取商品
func (r *ProductRepository) GetByID(id int) (*Product, error) {
	product := &Product{}
	query := `SELECT id, name, description, price, original_price, category, sub_category, 
	          brand, sku, stock, image_url, images, tags, is_active, is_featured, is_on_sale, 
	          merchant_id, view_count, sales_count, rating, review_count, weight, dimensions, 
	          created_at, updated_at FROM products WHERE id = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price, &product.OriginalPrice,
		&product.Category, &product.SubCategory, &product.Brand, &product.SKU, &product.Stock,
		&product.ImageURL, &product.Images, &product.Tags, &product.IsActive, &product.IsFeatured,
		&product.IsOnSale, &product.MerchantID, &product.ViewCount, &product.SalesCount,
		&product.Rating, &product.ReviewCount, &product.Weight, &product.Dimensions,
		&product.CreatedAt, &product.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return product, nil
}

// GetAll 獲取所有商品
func (r *ProductRepository) GetAll(limit, offset int) ([]*Product, error) {
	query := `SELECT id, name, description, price, original_price, category, sub_category, 
	          brand, sku, stock, image_url, images, tags, is_active, is_featured, is_on_sale, 
	          merchant_id, view_count, sales_count, rating, review_count, weight, dimensions, 
	          created_at, updated_at FROM products WHERE is_active = 1 
	          ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price, &product.OriginalPrice,
			&product.Category, &product.SubCategory, &product.Brand, &product.SKU, &product.Stock,
			&product.ImageURL, &product.Images, &product.Tags, &product.IsActive, &product.IsFeatured,
			&product.IsOnSale, &product.MerchantID, &product.ViewCount, &product.SalesCount,
			&product.Rating, &product.ReviewCount, &product.Weight, &product.Dimensions,
			&product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	
	return products, nil
}

// GetFeatured 獲取精選商品
func (r *ProductRepository) GetFeatured(limit int) ([]*Product, error) {
	query := `SELECT id, name, description, price, original_price, category, sub_category, 
	          brand, sku, stock, image_url, images, tags, is_active, is_featured, is_on_sale, 
	          merchant_id, view_count, sales_count, rating, review_count, weight, dimensions, 
	          created_at, updated_at FROM products WHERE is_active = 1 AND is_featured = 1 
	          ORDER BY created_at DESC LIMIT ?`
	
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price, &product.OriginalPrice,
			&product.Category, &product.SubCategory, &product.Brand, &product.SKU, &product.Stock,
			&product.ImageURL, &product.Images, &product.Tags, &product.IsActive, &product.IsFeatured,
			&product.IsOnSale, &product.MerchantID, &product.ViewCount, &product.SalesCount,
			&product.Rating, &product.ReviewCount, &product.Weight, &product.Dimensions,
			&product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	
	return products, nil
}

// GetByCategory 根據分類獲取商品
func (r *ProductRepository) GetByCategory(category string, limit, offset int) ([]*Product, error) {
	query := `SELECT id, name, description, price, original_price, category, sub_category, 
	          brand, sku, stock, image_url, images, tags, is_active, is_featured, is_on_sale, 
	          merchant_id, view_count, sales_count, rating, review_count, weight, dimensions, 
	          created_at, updated_at FROM products WHERE is_active = 1 AND category = ? 
	          ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := r.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price, &product.OriginalPrice,
			&product.Category, &product.SubCategory, &product.Brand, &product.SKU, &product.Stock,
			&product.ImageURL, &product.Images, &product.Tags, &product.IsActive, &product.IsFeatured,
			&product.IsOnSale, &product.MerchantID, &product.ViewCount, &product.SalesCount,
			&product.Rating, &product.ReviewCount, &product.Weight, &product.Dimensions,
			&product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	
	return products, nil
}

// Search 搜索商品
func (r *ProductRepository) Search(keyword string, limit, offset int) ([]*Product, error) {
	query := `SELECT id, name, description, price, original_price, category, sub_category, 
	          brand, sku, stock, image_url, images, tags, is_active, is_featured, is_on_sale, 
	          merchant_id, view_count, sales_count, rating, review_count, weight, dimensions, 
	          created_at, updated_at FROM products WHERE is_active = 1 AND 
	          (name LIKE ? OR description LIKE ? OR category LIKE ? OR brand LIKE ?) 
	          ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	searchTerm := "%" + keyword + "%"
	rows, err := r.db.Query(query, searchTerm, searchTerm, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price, &product.OriginalPrice,
			&product.Category, &product.SubCategory, &product.Brand, &product.SKU, &product.Stock,
			&product.ImageURL, &product.Images, &product.Tags, &product.IsActive, &product.IsFeatured,
			&product.IsOnSale, &product.MerchantID, &product.ViewCount, &product.SalesCount,
			&product.Rating, &product.ReviewCount, &product.Weight, &product.Dimensions,
			&product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	
	return products, nil
}

// GetCategories 獲取所有分類
func (r *ProductRepository) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM products WHERE is_active = 1 ORDER BY category`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	
	return categories, nil
}

// Update 更新商品
func (r *ProductRepository) Update(product *Product) error {
	query := `
		UPDATE products SET name = ?, description = ?, price = ?, original_price = ?, 
		                   category = ?, sub_category = ?, brand = ?, sku = ?, stock = ?, 
		                   image_url = ?, images = ?, tags = ?, is_active = ?, is_featured = ?, 
		                   is_on_sale = ?, view_count = ?, sales_count = ?, rating = ?, 
		                   review_count = ?, weight = ?, dimensions = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`
	
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price,
		product.OriginalPrice, product.Category, product.SubCategory, product.Brand,
		product.SKU, product.Stock, product.ImageURL, product.Images, product.Tags,
		product.IsActive, product.IsFeatured, product.IsOnSale, product.ViewCount,
		product.SalesCount, product.Rating, product.ReviewCount, product.Weight,
		product.Dimensions, product.ID)
	
	return err
}

// Delete 刪除商品
func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// IncrementViewCount 增加瀏覽次數
func (r *ProductRepository) IncrementViewCount(id int) error {
	query := `UPDATE products SET view_count = view_count + 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// IncrementSalesCount 增加銷售次數
func (r *ProductRepository) IncrementSalesCount(id int, quantity int) error {
	query := `UPDATE products SET sales_count = sales_count + ?, stock = stock - ? WHERE id = ?`
	_, err := r.db.Exec(query, quantity, quantity, id)
	return err
}
