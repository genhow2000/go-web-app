package middleware

import (
	"net/http"
	"go-simple-app/models"

	"github.com/gin-gonic/gin"
)

// 管理員權限中間件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/admin/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "未認證的用戶",
				})
			}
			c.Abort()
			return
		}

		// 檢查用戶角色
		if userObj, ok := user.(models.UserInterface); ok {
			if userObj.GetRole() != "admin" {
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
				c.Redirect(http.StatusFound, "/admin/login")
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

// 客戶權限中間件 - 只允許 customer 角色
func CustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/customer/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "未認證的用戶",
				})
			}
			c.Abort()
			return
		}

		// 檢查用戶角色 - 只允許 customer
		if userObj, ok := user.(models.UserInterface); ok {
			if userObj.GetRole() != "customer" {
				if c.Request.Header.Get("Accept") == "text/html" {
					c.HTML(http.StatusForbidden, "error.html", gin.H{
						"error": "權限不足，需要客戶權限",
					})
				} else {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "權限不足，需要客戶權限",
					})
				}
				c.Abort()
				return
			}
		} else {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/customer/login")
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

// 商戶權限中間件 - 只允許 merchant 角色
func MerchantMiddleware() gin.HandlerFunc {
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

		// 檢查用戶角色 - 只允許 merchant
		if userObj, ok := user.(models.UserInterface); ok {
			if userObj.GetRole() != "merchant" {
				if c.Request.Header.Get("Accept") == "text/html" {
					c.HTML(http.StatusForbidden, "error.html", gin.H{
						"error": "權限不足，需要商戶權限",
					})
				} else {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "權限不足，需要商戶權限",
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

// 多角色權限中間件 - 允許指定的多個角色
func MultiRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/customer/login")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "未認證的用戶",
				})
			}
			c.Abort()
			return
		}

		// 檢查用戶角色是否在允許的角色列表中
		if userObj, ok := user.(models.UserInterface); ok {
			roleAllowed := false
			userRole := userObj.GetRole()
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					roleAllowed = true
					break
				}
			}
			
			if !roleAllowed {
				if c.Request.Header.Get("Accept") == "text/html" {
					c.HTML(http.StatusForbidden, "error.html", gin.H{
						"error": "權限不足，需要適當的權限",
					})
				} else {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "權限不足，需要適當的權限",
					})
				}
				c.Abort()
				return
			}
		} else {
			if c.Request.Header.Get("Accept") == "text/html" {
				c.Redirect(http.StatusFound, "/customer/login")
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