package seeders

import (
	"database/sql"
	"log"
)

// UpdateProductImagesSeeder 更新商品圖片 seeder
type UpdateProductImagesSeeder struct {
	db *sql.DB
}

// NewUpdateProductImagesSeeder 創建更新商品圖片 seeder
func NewUpdateProductImagesSeeder(db *sql.DB) *UpdateProductImagesSeeder {
	return &UpdateProductImagesSeeder{db: db}
}

// Run 執行更新商品圖片 seeder
func (s *UpdateProductImagesSeeder) Run() error {
	log.Println("開始執行更新商品圖片 seeder...")

	// 定義商品圖片更新映射
	imageUpdates := map[string]string{
		"Pexels測試商品": "https://images.pexels.com/photos/788946/pexels-photo-788946.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試智能手機":   "https://images.pexels.com/photos/788946/pexels-photo-788946.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試筆記本電腦":  "https://images.pexels.com/photos/1181671/pexels-photo-1181671.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試時尚T恤":   "https://images.pexels.com/photos/1043474/pexels-photo-1043474.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試運動鞋":    "https://images.pexels.com/photos/1598505/pexels-photo-1598505.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試護膚套裝":   "https://images.pexels.com/photos/3373736/pexels-photo-3373736.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試咖啡機":    "https://images.pexels.com/photos/302899/pexels-photo-302899.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試有機零食":   "https://images.pexels.com/photos/143133/pexels-photo-143133.jpeg?auto=compress&cs=tinysrgb&w=400",
		"測試編程書籍":   "https://images.pexels.com/photos/159711/books-bookstore-book-reading-159711.jpeg?auto=compress&cs=tinysrgb&w=400",
	}

	// 更新每個商品的圖片
	for productName, imageURL := range imageUpdates {
		query := `UPDATE products SET image_url = ? WHERE name = ?`
		result, err := s.db.Exec(query, imageURL, productName)
		if err != nil {
			log.Printf("更新商品 %s 圖片失敗: %v", productName, err)
			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("獲取更新行數失敗: %v", err)
			continue
		}

		if rowsAffected > 0 {
			log.Printf("成功更新商品圖片: %s -> %s", productName, imageURL)
		} else {
			log.Printf("商品不存在: %s", productName)
		}
	}

	// 為沒有特定名稱的商品添加通用圖片
	genericImages := map[string]string{
		"電子產品": "https://images.pexels.com/photos/788946/pexels-photo-788946.jpeg?auto=compress&cs=tinysrgb&w=400",
		"服飾":   "https://images.pexels.com/photos/1043474/pexels-photo-1043474.jpeg?auto=compress&cs=tinysrgb&w=400",
		"運動":   "https://images.pexels.com/photos/1598505/pexels-photo-1598505.jpeg?auto=compress&cs=tinysrgb&w=400",
		"美妝":   "https://images.pexels.com/photos/3373736/pexels-photo-3373736.jpeg?auto=compress&cs=tinysrgb&w=400",
		"家居":   "https://images.pexels.com/photos/302899/pexels-photo-302899.jpeg?auto=compress&cs=tinysrgb&w=400",
		"食品":   "https://images.pexels.com/photos/143133/pexels-photo-143133.jpeg?auto=compress&cs=tinysrgb&w=400",
		"圖書":   "https://images.pexels.com/photos/159711/books-bookstore-book-reading-159711.jpeg?auto=compress&cs=tinysrgb&w=400",
	}

	// 為沒有圖片的商品按分類添加圖片
	for category, imageURL := range genericImages {
		query := `UPDATE products SET image_url = ? WHERE category = ? AND (image_url IS NULL OR image_url = '')`
		result, err := s.db.Exec(query, imageURL, category)
		if err != nil {
			log.Printf("更新分類 %s 商品圖片失敗: %v", category, err)
			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("獲取更新行數失敗: %v", err)
			continue
		}

		if rowsAffected > 0 {
			log.Printf("為分類 %s 的 %d 個商品添加了圖片", category, rowsAffected)
		}
	}

	log.Println("更新商品圖片 seeder 執行完成!")
	return nil
}

// Clear 清除商品圖片
func (s *UpdateProductImagesSeeder) Clear() error {
	log.Println("清除商品圖片...")
	
	_, err := s.db.Exec("UPDATE products SET image_url = NULL")
	if err != nil {
		log.Printf("清除商品圖片失敗: %v", err)
		return err
	}
	
	log.Println("商品圖片已清除")
	return nil
}
