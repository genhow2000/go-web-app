package services

import (
	"database/sql"
	"go-simple-app/models"
)

type MallService struct {
	productRepo *models.ProductRepository
}

func NewMallService(db *sql.DB) *MallService {
	return &MallService{
		productRepo: models.NewProductRepository(db),
	}
}

// GetHomepageData 獲取首頁數據
func (s *MallService) GetHomepageData() (map[string]interface{}, error) {
	// 獲取精選商品
	featuredProducts, err := s.productRepo.GetFeatured(8)
	if err != nil {
		return nil, err
	}

	// 獲取分類
	categories, err := s.productRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	// 獲取最新商品
	latestProducts, err := s.productRepo.GetAll(6, 0)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"featured_products": featuredProducts,
		"categories":        categories,
		"latest_products":   latestProducts,
	}, nil
}

// GetProductStats 獲取商品統計
func (s *MallService) GetProductStats() (map[string]interface{}, error) {
	// 這裡可以添加更複雜的統計邏輯
	// 例如：總商品數、分類統計、熱銷商品等
	
	return map[string]interface{}{
		"total_products": 0,
		"categories":     []string{},
		"featured_count": 0,
	}, nil
}

// SearchProductsWithFilters 帶篩選條件的商品搜尋
func (s *MallService) SearchProductsWithFilters(keyword string, category string, minPrice, maxPrice float64, limit, offset int) ([]*models.Product, error) {
	// 這裡可以實現更複雜的搜尋邏輯
	// 例如：價格範圍篩選、分類篩選、排序等
	
	if keyword != "" {
		return s.productRepo.Search(keyword, limit, offset)
	}
	
	if category != "" {
		return s.productRepo.GetByCategory(category, limit, offset)
	}
	
	return s.productRepo.GetAll(limit, offset)
}

// GetRecommendedProducts 獲取推薦商品
func (s *MallService) GetRecommendedProducts(userID int, limit int) ([]*models.Product, error) {
	// 這裡可以實現基於用戶行為的推薦算法
	// 目前簡單返回精選商品
	return s.productRepo.GetFeatured(limit)
}

// GetProductCategories 獲取商品分類（包含子分類）
func (s *MallService) GetProductCategories() ([]map[string]interface{}, error) {
	categories, err := s.productRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, category := range categories {
		// 獲取該分類的商品數量
		products, err := s.productRepo.GetByCategory(category, 1, 0)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"name":        category,
			"product_count": len(products),
			"icon":        getCategoryIcon(category),
		})
	}

	return result, nil
}

// getCategoryIcon 根據分類名稱獲取圖標
func getCategoryIcon(category string) string {
	iconMap := map[string]string{
		"電子產品": "📱",
		"服飾":    "👕",
		"家居":    "🏠",
		"美妝":    "💄",
		"運動":    "⚽",
		"食品":    "🍎",
		"圖書":    "📚",
		"其他":    "📦",
	}

	if icon, exists := iconMap[category]; exists {
		return icon
	}
	return "📦"
}
