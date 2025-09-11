package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// DBAuthMiddleware 資料庫管理專用認證中間件
func DBAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查是否有資料庫管理認證 cookie
		dbToken, err := c.Cookie("db_auth_token")
		if err != nil || dbToken == "" {
			// 重定向到資料庫管理登入頁面
			c.Redirect(http.StatusFound, "/admin/db/login")
			c.Abort()
			return
		}

		// 驗證 token (這裡使用簡單的環境變數驗證)
		expectedToken := os.Getenv("DB_AUTH_TOKEN")
		if expectedToken == "" {
			expectedToken = "system" // 預設密碼
		}

		if dbToken != expectedToken {
			// Token 無效，重定向到登入頁面
			c.Redirect(http.StatusFound, "/admin/db/login")
			c.Abort()
			return
		}

		// 認證通過，繼續處理請求
		c.Next()
	}
}

