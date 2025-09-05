package middleware

import (
	"go-simple-app/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RequestLogger 記錄HTTP請求的中間件
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 記錄請求詳情
		logger.Info("HTTP請求", logrus.Fields{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"user_agent": param.Request.UserAgent(),
			"timestamp":  param.TimeStamp.Format("2006-01-02 15:04:05"),
		})
		
		return ""
	})
}

// ErrorLogger 記錄錯誤的中間件
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 記錄錯誤
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("請求處理錯誤", err, logrus.Fields{
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"client_ip": c.ClientIP(),
					"status":    c.Writer.Status(),
				})
			}
		}
	}
}
