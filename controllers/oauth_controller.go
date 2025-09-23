package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"go-simple-app/services"
	"go-simple-app/logger"
	
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OAuthController struct {
	oauthService *services.OAuthService
}

func NewOAuthController(oauthService *services.OAuthService) *OAuthController {
	return &OAuthController{
		oauthService: oauthService,
	}
}

// LINE登入重定向
func (c *OAuthController) LineLogin(ctx *gin.Context) {
	// 生成state參數防止CSRF攻擊
	state := generateRandomState()
	
	// 將state存儲到session或cookie
	ctx.SetCookie("oauth_state", state, 600, "/", "", false, false)
	
	// 重定向到LINE授權頁面
	authURL := c.oauthService.GetLineAuthURL(state)
	ctx.Redirect(http.StatusFound, authURL)
}

// LINE登入回調
func (c *OAuthController) LineCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	
	// 驗證state參數
	cookieState, _ := ctx.Cookie("oauth_state")
	if state != cookieState {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state parameter",
		})
		return
	}
	
	// 清除state cookie
	ctx.SetCookie("oauth_state", "", -1, "/", "", false, false)
	
	// 處理OAuth回調
	user, err := c.oauthService.HandleLineCallback(ctx.Request.Context(), code)
	if err != nil {
		logger.Error("OAuth回調處理失敗", err, logrus.Fields{
			"code": code,
			"state": state,
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "OAuth login failed",
		})
		return
	}
	
	logger.Info("OAuth回調處理成功", logrus.Fields{
		"user_id": user.GetID(),
		"role": user.GetRole(),
		"email": user.GetEmail(),
	})
	
	// 生成JWT token
	token, err := c.oauthService.GenerateToken(user)
	if err != nil {
		logger.Error("JWT token生成失敗", err, logrus.Fields{
			"user_id": user.GetID(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token generation failed",
		})
		return
	}
	
	// 設置認證cookie
	ctx.SetCookie("auth_token", token, 3600*24, "/", "", false, false)
	
	logger.Info("OAuth登入完成，重定向到儀表板", logrus.Fields{
		"user_id": user.GetID(),
		"redirect_url": "/customer/dashboard",
	})
	
	// 重定向到儀表板，並在URL中傳遞token
	// 使用環境變數或默認值
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	ctx.Redirect(http.StatusFound, frontendURL+"/customer/dashboard?token="+token)
}

func generateRandomState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
