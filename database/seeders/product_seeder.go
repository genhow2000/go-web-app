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

	// 這裡可以添加創建測試產品的邏輯
	// 例如：
	// products := []struct {
	//     Name        string
	//     Description string
	//     Price       float64
	//     Category    string
	// }{
	//     {"測試商品1", "這是一個測試商品", 99.99, "電子產品"},
	//     {"測試商品2", "這是另一個測試商品", 199.99, "服飾"},
	// }
	//
	// for _, product := range products {
	//     query := `INSERT INTO products (name, description, price, category, created_at, updated_at) 
	//               VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	//     _, err := s.db.Exec(query, product.Name, product.Description, product.Price, product.Category)
	//     if err != nil {
	//         return err
	//     }
	//     log.Printf("創建產品: %s", product.Name)
	// }

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
