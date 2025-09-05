package middleware

import (
	"net/http"
	"strings"
	"go-simple-app/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Header 獲取 token
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			// 檢查 Bearer token 格式
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				tokenString = authHeader
			}
		} else {
			// 從 cookie 獲取 token
			if cookie, err := c.Cookie("auth_token"); err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供認證 token",
			})
			c.Abort()
			return
		}

		// 驗證 token
		user, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無效的 token",
			})
			c.Abort()
			return
		}

		// 將用戶信息存儲到 context
		c.Set("user", user)
		c.Next()
	}
}

func OptionalAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Header 獲取 token
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				tokenString = authHeader
			}
		} else {
			if cookie, err := c.Cookie("auth_token"); err == nil {
				tokenString = cookie
			}
		}

		if tokenString != "" {
			// 嘗試驗證 token
			if user, err := authService.ValidateToken(tokenString); err == nil {
				c.Set("user", user)
			}
		}

		c.Next()
	}
}
