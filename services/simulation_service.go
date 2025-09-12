package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// SimulationService 模拟AI服务
type SimulationService struct {
	usageStats AIUsageStats
}

// NewSimulationService 创建模拟服务
func NewSimulationService() *SimulationService {
	return &SimulationService{
		usageStats: AIUsageStats{
			Provider:    "simulation",
			DailyUsage:  0,
			DailyLimit:  999999,
			LastReset:   time.Now(),
			IsExhausted: false,
			ErrorCount:  0,
			LastError:   "",
			LastUsed:    time.Time{},
		},
	}
}

// GenerateResponse 生成模拟回复
func (s *SimulationService) GenerateResponse(ctx context.Context, message, conversationID string) (string, error) {
	// 更新使用统计
	s.usageStats.DailyUsage++
	s.usageStats.LastUsed = time.Now()

	// 生成模拟回复
	response := s.generateSimulatedResponse(message)
	return response, nil
}

// GetServiceName 获取服务名称
func (s *SimulationService) GetServiceName() string {
	return "Simulation Service"
}

// IsAvailable 检查服务是否可用
func (s *SimulationService) IsAvailable(ctx context.Context) bool {
	return !s.usageStats.IsExhausted
}

// GetUsageStats 获取使用统计
func (s *SimulationService) GetUsageStats() map[string]interface{} {
	return map[string]interface{}{
		"provider":        s.usageStats.Provider,
		"daily_usage":     s.usageStats.DailyUsage,
		"daily_limit":     s.usageStats.DailyLimit,
		"last_reset":      s.usageStats.LastReset,
		"is_exhausted":    s.usageStats.IsExhausted,
		"error_count":     s.usageStats.ErrorCount,
		"last_error":      s.usageStats.LastError,
		"last_used":       s.usageStats.LastUsed,
		"usage_percentage": float64(s.usageStats.DailyUsage) / float64(s.usageStats.DailyLimit) * 100,
	}
}

// generateSimulatedResponse 生成模拟回复
func (s *SimulationService) generateSimulatedResponse(message string) string {
	message = strings.ToLower(message)
	
	// 关键词匹配回复
	responses := map[string][]string{
		"你好": {
			"您好！我是阿和商城的AI購物助手，很高興為您服務！有什麼可以幫助您的嗎？",
			"你好！歡迎來到阿和商城，我是您的專屬購物助手，有什麼需要幫助的嗎？",
			"您好！很高興見到您，我是AI助手，隨時為您提供購物建議和幫助！",
		},
		"產品": {
			"我們有很多優質商品可以滿足您的需求，讓我為您介紹一下。",
			"我們提供各種類型的產品，包括電子產品、服裝、家居用品等，您對哪類產品感興趣？",
			"我們有豐富的產品選擇，從時尚服飾到智能家電，應有盡有！",
		},
		"優惠": {
			"我們經常推出各種優惠活動，包括折扣、滿減、贈品等，請關注我們的促銷信息！",
			"目前有多個優惠活動正在進行中，您可以查看我們的優惠頁面了解更多詳情。",
			"我們會定期推出限時優惠，建議您關注我們的官方通知！",
		},
		"配送": {
			"我們提供快速配送服務，一般情況下1-3個工作日內送達。",
			"配送時間根據您的地區而定，我們會盡快為您安排發貨。",
			"我們與多家物流公司合作，確保您的訂單安全快速送達。",
		},
		"退換": {
			"我們提供7天無理由退換貨服務，讓您購物更放心。",
			"如果商品有問題，請聯繫我們的客服，我們會為您處理退換貨事宜。",
			"我們有完善的售後服務，確保您的購物體驗。",
		},
		"推薦": {
			"根據您的需求，我推薦以下熱銷商品：智能手機、筆記本電腦、時尚服飾等。",
			"我們有很多熱門商品，包括最新款的電子產品和時尚單品。",
			"讓我為您推薦一些性價比很高的商品，相信您會喜歡的！",
		},
	}
	
	// 查找匹配的關鍵詞
	for keyword, responseList := range responses {
		if strings.Contains(message, keyword) {
			rand.Seed(time.Now().UnixNano())
			return responseList[rand.Intn(len(responseList))]
		}
	}
	
	// 默認回復
	defaultResponses := []string{
		"感謝您的詢問！我是阿和商城的AI助手，很樂意為您提供幫助。",
		"您好！我是購物助手，有什麼關於商品的問題都可以問我。",
		"很高興為您服務！請告訴我您需要什麼幫助，我會盡力協助您。",
		"我是阿和商城的智能助手，隨時為您提供購物建議和產品信息。",
		"歡迎來到阿和商城！我是您的專屬購物助手，有什麼可以幫助您的嗎？",
	}
	
	rand.Seed(time.Now().UnixNano())
	return defaultResponses[rand.Intn(len(defaultResponses))]
}