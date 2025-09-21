package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImageController struct{}

func NewImageController() *ImageController {
	return &ImageController{}
}

// GenerateProductImage 生成商品圖片
func (c *ImageController) GenerateProductImage(ctx *gin.Context) {
	// 獲取參數
	width := ctx.DefaultQuery("w", "300")
	height := ctx.DefaultQuery("h", "300")
	text := ctx.Query("text")
	
	// 轉換為整數
	w, err := strconv.Atoi(width)
	if err != nil {
		w = 300
	}
	
	h, err := strconv.Atoi(height)
	if err != nil {
		h = 300
	}
	
	// 如果沒有提供文字，使用預設
	if text == "" {
		text = "商品圖片"
	}
	
	// 生成SVG圖片
	svg := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
		<rect width="%d" height="%d" fill="#f8f9fa" stroke="#dee2e6" stroke-width="2"/>
		<rect x="10" y="10" width="%d" height="%d" fill="#e9ecef" stroke="#adb5bd" stroke-width="1"/>
		<text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" font-family="Arial, sans-serif" font-size="%d" fill="#495057">
			%s
		</text>
		<text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" font-family="Arial, sans-serif" font-size="%d" fill="#6c757d">
			%dx%d
		</text>
	</svg>`, 
		w, h, w, h, w-20, h-20, w/2, h/2-10, w/20, text, w/2, h/2+20, w/25, w, h)
	
	// 設置響應頭
	ctx.Header("Content-Type", "image/svg+xml")
	ctx.Header("Cache-Control", "public, max-age=3600")
	
	// 返回SVG
	ctx.String(http.StatusOK, svg)
}

// GeneratePlaceholderImage 生成佔位符圖片
func (c *ImageController) GeneratePlaceholderImage(ctx *gin.Context) {
	// 獲取參數
	width := ctx.DefaultQuery("w", "300")
	height := ctx.DefaultQuery("h", "300")
	text := ctx.Query("text")
	
	// 轉換為整數
	w, err := strconv.Atoi(width)
	if err != nil {
		w = 300
	}
	
	h, err := strconv.Atoi(height)
	if err != nil {
		h = 300
	}
	
	// 如果沒有提供文字，使用預設
	if text == "" {
		text = "圖片載入中"
	}
	
	// 生成簡單的佔位符SVG
	svg := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
		<rect width="%d" height="%d" fill="#f8f9fa"/>
		<rect x="%d" y="%d" width="%d" height="%d" fill="#e9ecef" rx="8"/>
		<text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" font-family="Arial, sans-serif" font-size="%d" fill="#6c757d">
			%s
		</text>
	</svg>`, 
		w, h, w, h, w/4, h/4, w/2, h/2, w/2, h/2, w/15, text)
	
	// 設置響應頭
	ctx.Header("Content-Type", "image/svg+xml")
	ctx.Header("Cache-Control", "public, max-age=3600")
	
	// 返回SVG
	ctx.String(http.StatusOK, svg)
}
