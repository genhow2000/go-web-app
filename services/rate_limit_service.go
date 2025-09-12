package services

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"
)

// RateLimitService 限制服務
type RateLimitService struct {
	// 每分鐘請求限制
	minuteLimits map[string][]time.Time
	// 每日請求限制
	dailyLimits map[string]int
	// 鎖定機制
	mutex sync.RWMutex
}

// 限制配置
const (
	// 匿名用戶每分鐘最多 5 次請求
	AnonymousMinuteLimit = 5
	// 匿名用戶每天最多 50 次請求
	AnonymousDailyLimit = 50
	// 清理過期記錄的間隔
	CleanupInterval = 5 * time.Minute
)

// NewRateLimitService 創建限制服務
func NewRateLimitService() *RateLimitService {
	service := &RateLimitService{
		minuteLimits: make(map[string][]time.Time),
		dailyLimits:  make(map[string]int),
	}
	
	// 啟動清理協程
	go service.startCleanup()
	
	return service
}

// CheckRateLimit 檢查請求頻率限制
func (rls *RateLimitService) CheckRateLimit(identifier string, isAnonymous bool) (bool, string) {
	rls.mutex.Lock()
	defer rls.mutex.Unlock()
	
	now := time.Now()
	
	// 生成限制鍵（匿名用戶使用IP+日期，認證用戶使用用戶ID）
	limitKey := rls.generateLimitKey(identifier, isAnonymous, now)
	
	// 檢查每分鐘限制
	if !rls.checkMinuteLimit(limitKey, now) {
		return false, "請求過於頻繁，請稍後再試（每分鐘最多5次請求）"
	}
	
	// 檢查每日限制
	if !rls.checkDailyLimit(limitKey, now) {
		return false, "今日使用次數已達上限，請註冊會員獲得更多使用次數"
	}
	
	// 記錄請求
	rls.recordRequest(limitKey, now)
	
	return true, ""
}

// checkMinuteLimit 檢查每分鐘限制
func (rls *RateLimitService) checkMinuteLimit(limitKey string, now time.Time) bool {
	requests := rls.minuteLimits[limitKey]
	
	// 清理過期的請求記錄（1分鐘前）
	cutoff := now.Add(-time.Minute)
	validRequests := []time.Time{}
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	rls.minuteLimits[limitKey] = validRequests
	
	// 檢查是否超過限制
	return len(validRequests) < AnonymousMinuteLimit
}

// checkDailyLimit 檢查每日限制
func (rls *RateLimitService) checkDailyLimit(limitKey string, now time.Time) bool {
	// 生成每日限制鍵（只包含日期）
	dailyKey := rls.generateDailyKey(limitKey, now)
	
	// 檢查今日請求次數
	return rls.dailyLimits[dailyKey] < AnonymousDailyLimit
}

// recordRequest 記錄請求
func (rls *RateLimitService) recordRequest(limitKey string, now time.Time) {
	// 記錄到每分鐘限制
	rls.minuteLimits[limitKey] = append(rls.minuteLimits[limitKey], now)
	
	// 記錄到每日限制
	dailyKey := rls.generateDailyKey(limitKey, now)
	rls.dailyLimits[dailyKey]++
}

// generateLimitKey 生成限制鍵
func (rls *RateLimitService) generateLimitKey(identifier string, isAnonymous bool, now time.Time) string {
	if isAnonymous {
		// 匿名用戶：IP + 日期
		return fmt.Sprintf("anon_%s_%s", identifier, now.Format("2006-01-02"))
	}
	// 認證用戶：用戶ID + 日期
	return fmt.Sprintf("user_%s_%s", identifier, now.Format("2006-01-02"))
}

// generateDailyKey 生成每日限制鍵
func (rls *RateLimitService) generateDailyKey(limitKey string, now time.Time) string {
	// 提取基礎標識符
	if len(limitKey) > 5 {
		base := limitKey[5:] // 移除 "anon_" 或 "user_" 前綴
		return fmt.Sprintf("daily_%s", base)
	}
	return limitKey
}

// startCleanup 啟動清理協程
func (rls *RateLimitService) startCleanup() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		rls.cleanup()
	}
}

// cleanup 清理過期記錄
func (rls *RateLimitService) cleanup() {
	rls.mutex.Lock()
	defer rls.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-24 * time.Hour) // 清理24小時前的記錄
	
	// 清理過期的每分鐘限制記錄
	for key, requests := range rls.minuteLimits {
		validRequests := []time.Time{}
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		if len(validRequests) == 0 {
			delete(rls.minuteLimits, key)
		} else {
			rls.minuteLimits[key] = validRequests
		}
	}
	
	// 清理過期的每日限制記錄
	for key := range rls.dailyLimits {
		// 檢查是否為昨天的記錄
		if len(key) > 6 && key[:6] == "daily_" {
			dateStr := key[6:]
			if date, err := time.Parse("2006-01-02", dateStr); err == nil {
				if date.Before(cutoff) {
					delete(rls.dailyLimits, key)
				}
			}
		}
	}
}

// GetUsageStats 獲取使用統計
func (rls *RateLimitService) GetUsageStats(identifier string, isAnonymous bool) map[string]interface{} {
	rls.mutex.RLock()
	defer rls.mutex.RUnlock()
	
	now := time.Now()
	limitKey := rls.generateLimitKey(identifier, isAnonymous, now)
	dailyKey := rls.generateDailyKey(limitKey, now)
	
	// 計算每分鐘請求次數
	minuteRequests := rls.minuteLimits[limitKey]
	cutoff := now.Add(-time.Minute)
	minuteCount := 0
	for _, reqTime := range minuteRequests {
		if reqTime.After(cutoff) {
			minuteCount++
		}
	}
	
	// 獲取每日請求次數
	dailyCount := rls.dailyLimits[dailyKey]
	
	return map[string]interface{}{
		"minute_requests": minuteCount,
		"minute_limit":    AnonymousMinuteLimit,
		"daily_requests":  dailyCount,
		"daily_limit":     AnonymousDailyLimit,
		"is_anonymous":    isAnonymous,
	}
}

// GenerateAnonymousIdentifier 為匿名用戶生成標識符
func (rls *RateLimitService) GenerateAnonymousIdentifier(ip, userAgent string) string {
	// 使用IP和User-Agent生成匿名用戶標識符
	hash := md5.Sum([]byte(ip + userAgent + time.Now().Format("2006-01-02")))
	return fmt.Sprintf("%x", hash)[:16] // 取前16位
}
