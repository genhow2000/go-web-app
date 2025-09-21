package controllers

import (
	"strings"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ImageProxyController struct {
	client *http.Client
}

func NewImageProxyController() *ImageProxyController {
	return &ImageProxyController{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyImage 代理外部圖片
func (c *ImageProxyController) ProxyImage(ctx *gin.Context) {
	// 獲取圖片URL參數
	imageURL := ctx.Query("url")
	if imageURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少圖片URL參數",
		})
		return
	}

	// 驗證URL格式
	if !isValidImageURL(imageURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的圖片URL",
		})
		return
	}

	// 創建請求
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無法創建請求: " + err.Error(),
		})
		return
	}

	// 設置請求頭
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// 發送請求
	resp, err := c.client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "無法獲取圖片: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 檢查響應狀態
	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "圖片服務返回錯誤狀態: " + resp.Status,
		})
		return
	}

	// 設置響應頭
	ctx.Header("Content-Type", resp.Header.Get("Content-Type"))
	ctx.Header("Cache-Control", "public, max-age=3600")
	ctx.Header("Access-Control-Allow-Origin", "*")

	// 複製圖片數據
	_, err = io.Copy(ctx.Writer, resp.Body)
	if err != nil {
		// 如果複製失敗，返回錯誤
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "圖片傳輸失敗: " + err.Error(),
		})
		return
	}
}

// GenerateExternalImage 生成外部圖片URL
func (c *ImageProxyController) GenerateExternalImage(ctx *gin.Context) {
	// 獲取參數
	service := ctx.DefaultQuery("service", "placeholder")

	// 根據服務選擇不同的圖片URL
	var imageURL string
	switch service {
	case "placeholder":
		// 使用httpbin.org的圖片服務
		imageURL = "https://httpbin.org/image/png"
	case "picsum":
		imageURL = "https://picsum.photos/300/300"
	case "httpbin":
		imageURL = "https://httpbin.org/image/png"
	case "unsplash":
		imageURL = "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=300&h=300&fit=crop"
	default:
		// 使用httpbin.org的圖片服務
		imageURL = "https://httpbin.org/image/png"
	}

	// 返回代理URL
	proxyURL := "/api/image/proxy?url=" + imageURL
	ctx.JSON(http.StatusOK, gin.H{
		"original_url": imageURL,
		"proxy_url":    proxyURL,
		"service":      service,
	})
}

// isValidImageURL 驗證圖片URL
func isValidImageURL(url string) bool {
	// 簡單的URL驗證
	if len(url) < 10 {
		return false
	}
	
	// 檢查是否以http或https開頭
	if url[:4] != "http" {
		return false
	}
	
	// 檢查是否包含圖片相關的域名或路徑
	validDomains := []string{
		"via.placeholder.com",
		"picsum.photos",
		"httpbin.org",
		"images.unsplash.com",
		"loremflickr.com",
		"placeimg.com",
		"images.pexels.com",
	}
	
	for _, domain := range validDomains {
		if contains(url, domain) {
			return true
		}
	}
	
	return false
}

// contains 檢查字符串是否包含子字符串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
