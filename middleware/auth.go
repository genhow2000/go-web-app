package middleware

import (
	"net/http"
	"strings"
	"go-simple-app/services"
	"go-simple-app/models"

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
			} else {
				// 從 query parameter 獲取 token (用於頁面訪問)
				tokenString = c.Query("token")
			}
		}

		if tokenString == "" {
			// 如果是頁面請求，重定向到登入頁面
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/merchant/login")
				c.Abort()
				return
			}
			
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供認證 token",
			})
			c.Abort()
			return
		}

		// 驗證 token
		user, err := authService.ValidateToken(tokenString)
		if err != nil {
			// 如果是頁面請求，重定向到登入頁面
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/merchant/login")
				c.Abort()
				return
			}
			
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

// 管理員權限中間件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/merchant/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "未認證的用戶",
				})
			}
			c.Abort()
			return
		}

		// 檢查用戶角色
		if userObj, ok := user.(*models.User); ok {
			if userObj.Role != "admin" {
				if c.Request.Header.Get("Accept") == "text/html" {
					c.HTML(http.StatusForbidden, "error.html", gin.H{
						"error": "權限不足，需要管理員權限",
					})
				} else {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "權限不足，需要管理員權限",
					})
				}
				c.Abort()
				return
			}
		} else {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/merchant/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "無效的用戶信息",
				})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// 客戶權限中間件
func CustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/merchant/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "未認證的用戶",
				})
			}
			c.Abort()
			return
		}

		// 檢查用戶角色
		if userObj, ok := user.(*models.User); ok {
			if userObj.Role != "customer" && userObj.Role != "admin" {
				if c.Request.Header.Get("Accept") == "text/html" {
					c.HTML(http.StatusForbidden, "error.html", gin.H{
						"error": "權限不足",
					})
				} else {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "權限不足",
					})
				}
				c.Abort()
				return
			}
		} else {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/merchant/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "無效的用戶信息",
				})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
