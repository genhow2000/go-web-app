package seeders

import (
	"database/sql"
	"log"
)

// ProductSeeder 產品 seeder (示例)
type ProductSeeder struct {
	db *sql.DB
}

// NewProductSeeder 創建產品 seeder
func NewProductSeeder(db *sql.DB) *ProductSeeder {
	return &ProductSeeder{db: db}
}

// Run 執行產品 seeder
func (s *ProductSeeder) Run() error {
	// 檢查是否已經有測試產品
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM products WHERE name LIKE '測試%'").Scan(&count)
	if err != nil {
		// 如果 products 表不存在，跳過
		log.Println("products 表不存在，跳過產品 seeder")
		return nil
	}

	if count > 0 {
		log.Println("測試產品已存在，跳過產品 seeder")
		return nil
	}

	log.Println("開始執行產品 seeder...")

	// 創建測試商品
	products := []struct {
		Name         string
		Description  string
		Price        float64
		OriginalPrice *float64
		Category     string
		SubCategory  *string
		Brand        *string
		Stock        int
		ImageURL     *string
		IsFeatured   bool
		IsOnSale     bool
		MerchantID   int
	}{
		{
			Name:         "測試智能手機",
			Description:  "最新款智能手機，配備先進的攝影系統和強大的處理器",
			Price:        29999.00,
			OriginalPrice: func() *float64 { v := 32999.00; return &v }(),
			Category:     "電子產品",
			SubCategory:  func() *string { v := "手機"; return &v }(),
			Brand:        func() *string { v := "TechBrand"; return &v }(),
			Stock:        50,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=智能手機"; return &v }(),
			IsFeatured:   true,
			IsOnSale:     true,
			MerchantID:   1,
		},
		{
			Name:         "測試筆記本電腦",
			Description:  "高性能筆記本電腦，適合辦公和娛樂使用",
			Price:        45999.00,
			OriginalPrice: nil,
			Category:     "電子產品",
			SubCategory:  func() *string { v := "筆記本"; return &v }(),
			Brand:        func() *string { v := "LaptopPro"; return &v }(),
			Stock:        30,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=筆記本電腦"; return &v }(),
			IsFeatured:   true,
			IsOnSale:     false,
			MerchantID:   1,
		},
		{
			Name:         "測試時尚T恤",
			Description:  "舒適的純棉T恤，多種顏色可選",
			Price:        599.00,
			OriginalPrice: func() *float64 { v := 799.00; return &v }(),
			Category:     "服飾",
			SubCategory:  func() *string { v := "男裝"; return &v }(),
			Brand:        func() *string { v := "FashionStyle"; return &v }(),
			Stock:        100,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=時尚T恤"; return &v }(),
			IsFeatured:   false,
			IsOnSale:     true,
			MerchantID:   1,
		},
		{
			Name:         "測試運動鞋",
			Description:  "專業運動鞋，提供優異的緩震和支撐",
			Price:        2999.00,
			OriginalPrice: nil,
			Category:     "運動",
			SubCategory:  func() *string { v := "鞋類"; return &v }(),
			Brand:        func() *string { v := "SportMax"; return &v }(),
			Stock:        80,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=運動鞋"; return &v }(),
			IsFeatured:   true,
			IsOnSale:     false,
			MerchantID:   1,
		},
		{
			Name:         "測試護膚套裝",
			Description:  "天然護膚套裝，包含潔面乳、爽膚水和面霜",
			Price:        1299.00,
			OriginalPrice: func() *float64 { v := 1599.00; return &v }(),
			Category:     "美妝",
			SubCategory:  func() *string { v := "護膚"; return &v }(),
			Brand:        func() *string { v := "BeautyCare"; return &v }(),
			Stock:        60,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=護膚套裝"; return &v }(),
			IsFeatured:   false,
			IsOnSale:     true,
			MerchantID:   1,
		},
		{
			Name:         "測試咖啡機",
			Description:  "全自動咖啡機，可製作多種咖啡飲品",
			Price:        8999.00,
			OriginalPrice: nil,
			Category:     "家居",
			SubCategory:  func() *string { v := "廚房用品"; return &v }(),
			Brand:        func() *string { v := "CoffeeMaster"; return &v }(),
			Stock:        25,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=咖啡機"; return &v }(),
			IsFeatured:   true,
			IsOnSale:     false,
			MerchantID:   1,
		},
		{
			Name:         "測試有機零食",
			Description:  "健康有機零食，無添加防腐劑",
			Price:        299.00,
			OriginalPrice: func() *float64 { v := 399.00; return &v }(),
			Category:     "食品",
			SubCategory:  func() *string { v := "零食"; return &v }(),
			Brand:        func() *string { v := "OrganicFood"; return &v }(),
			Stock:        200,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=有機零食"; return &v }(),
			IsFeatured:   false,
			IsOnSale:     true,
			MerchantID:   1,
		},
		{
			Name:         "測試編程書籍",
			Description:  "Go語言編程入門書籍，適合初學者",
			Price:        599.00,
			OriginalPrice: nil,
			Category:     "圖書",
			SubCategory:  func() *string { v := "技術書籍"; return &v }(),
			Brand:        func() *string { v := "TechBooks"; return &v }(),
			Stock:        150,
			ImageURL:     func() *string { v := "https://via.placeholder.com/300x300?text=編程書籍"; return &v }(),
			IsFeatured:   true,
			IsOnSale:     false,
			MerchantID:   1,
		},
	}

	for _, product := range products {
		query := `INSERT INTO products (name, description, price, original_price, category, sub_category, 
		                               brand, stock, image_url, is_featured, is_on_sale, merchant_id, 
		                               view_count, sales_count, rating, review_count, created_at, updated_at) 
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
		_, err := s.db.Exec(query, product.Name, product.Description, product.Price, product.OriginalPrice,
			product.Category, product.SubCategory, product.Brand, product.Stock, product.ImageURL,
			product.IsFeatured, product.IsOnSale, product.MerchantID, 0, 0, 0.0, 0)
		if err != nil {
			return err
		}
		log.Printf("創建產品: %s", product.Name)
	}

	log.Println("產品 seeder 執行完成!")
	return nil
}

// Clear 清除測試產品
func (s *ProductSeeder) Clear() error {
	log.Println("清除測試產品...")
	
	_, err := s.db.Exec("DELETE FROM products WHERE name LIKE '測試%'")
	if err != nil {
		// 如果 products 表不存在，跳過
		log.Println("products 表不存在，跳過清除")
		return nil
	}
	
	log.Println("測試產品已清除")
	return nil
}
