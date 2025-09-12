package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-simple-app/config"
)

// AIManager AI服务管理器
type AIManager struct {
	config   config.AIConfig
	services map[string]AIService
}

// NewAIManager 创建AI管理器
func NewAIManager(cfg config.AIConfig) *AIManager {
	manager := &AIManager{
		config:   cfg,
		services: make(map[string]AIService),
	}
	
	manager.initializeServices()
	return manager
}

// initializeServices 初始化所有AI服务
func (m *AIManager) initializeServices() {
	// 初始化Hugging Face服务
	if m.config.HuggingFace.APIURL != "" {
		m.services["huggingface"] = NewHuggingFaceService(m.config.HuggingFace)
		log.Printf("Initialized Hugging Face service: %s", m.config.HuggingFace.Model)
	}

	// 初始化Gemini服务
	if m.config.Gemini.APIKey != "" && m.config.Gemini.APIURL != "" {
		m.services["gemini"] = NewGeminiService(m.config.Gemini)
		log.Printf("Initialized Gemini service: %s", m.config.Gemini.Model)
	}

	// 初始化Groq服务
	if m.config.Groq.APIKey != "" && m.config.Groq.APIURL != "" {
		m.services["groq"] = NewGroqService(m.config.Groq)
		log.Printf("Initialized Groq service: %s", m.config.Groq.Model)
	}

	// 初始化模拟服务
	m.services["simulation"] = NewSimulationService()
	log.Printf("Initialized Simulation service")
}

// GenerateResponse 生成AI回复
func (m *AIManager) GenerateResponse(ctx context.Context, message, conversationID string) (string, error) {
	// 尝试主要服务
	primaryService := m.getPrimaryService(ctx)
	if primaryService != nil {
		response, err := primaryService.GenerateResponse(ctx, message, conversationID)
		if err == nil {
			log.Printf("Generated response using %s API", primaryService.GetServiceName())
			return response, nil
		}
		
		// 处理AI错误
		m.handleAIError(err, primaryService.GetServiceName())
		
		// 尝试备用服务
		backupService := m.getBackupService(ctx)
		if backupService != nil {
			log.Printf("Trying backup service: %s", backupService.GetServiceName())
			response, err := backupService.GenerateResponse(ctx, message, conversationID)
			if err == nil {
				log.Printf("Generated response using %s API", backupService.GetServiceName())
				return response, nil
			}
			
			// 处理备用服务错误
			m.handleAIError(err, backupService.GetServiceName())
		}
	}

	// 最后使用模拟服务
	if simulationService, exists := m.services["simulation"]; exists {
		log.Printf("Using simulation service as fallback")
		return simulationService.GenerateResponse(ctx, message, conversationID)
	}

	return "", fmt.Errorf("no available service")
}

// getPrimaryService 获取主要服务
func (m *AIManager) getPrimaryService(ctx context.Context) AIService {
	serviceName := string(m.config.PrimaryProvider)
	if service, exists := m.services[serviceName]; exists && service.IsAvailable(ctx) {
		return service
	}
	return nil
}

// getBackupService 获取备用服务
func (m *AIManager) getBackupService(ctx context.Context) AIService {
	// 如果主要服务是Groq，尝试Gemini
	if m.config.PrimaryProvider == "groq" {
		if geminiService, exists := m.services["gemini"]; exists && geminiService.IsAvailable(ctx) {
			return geminiService
		}
	}

	// 如果主要服务是Gemini，尝试Groq
	if m.config.PrimaryProvider == "gemini" {
		if groqService, exists := m.services["groq"]; exists && groqService.IsAvailable(ctx) {
			return groqService
		}
	}

	// 最后尝试模拟服务
	if simulationService, exists := m.services["simulation"]; exists {
		return simulationService
	}

	return nil
}

// handleAIError 处理AI错误
func (m *AIManager) handleAIError(err error, serviceName string) {
	log.Printf("AI Error from %s: %v", serviceName, err)
	
	// 这里可以添加更复杂的错误处理逻辑
	// 比如根据错误类型决定是否切换服务
	if aiErr, ok := err.(*AIError); ok {
		if aiErr.IsQuotaExceededError() {
			log.Printf("Quota exceeded for %s, switching to backup", serviceName)
		} else if aiErr.IsRateLimitedError() {
			log.Printf("Rate limited for %s, switching to backup", serviceName)
		}
	}
}

// GetServiceStats 获取所有服务统计
func (m *AIManager) GetServiceStats() map[string]map[string]interface{} {
	stats := make(map[string]map[string]interface{})
	
	for name, service := range m.services {
		stats[name] = service.GetUsageStats()
	}
	
	return stats
}
