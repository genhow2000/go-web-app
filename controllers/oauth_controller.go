package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"go-simple-app/services"
	"go-simple-app/logger"
	
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 內存存儲state參數，避免cookie問題
var stateStore = struct {
	sync.RWMutex
	states map[string]time.Time
}{
	states: make(map[string]time.Time),
}

// 清理過期的state
func cleanupExpiredStates() {
	stateStore.Lock()
	defer stateStore.Unlock()
	
	now := time.Now()
	for state, expiry := range stateStore.states {
		if now.After(expiry) {
			delete(stateStore.states, state)
		}
	}
}

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
	
	// 判斷是否為 HTTPS，用來決定是否加上 Secure 屬性
	isSecure := ctx.Request.TLS != nil || ctx.GetHeader("X-Forwarded-Proto") == "https"

	// 將state存儲到cookie中，避免多實例問題
	// 不設定 Domain，讓瀏覽器自動處理；使用 SameSite=Lax 來支援跨站重定向
	baseCookie := fmt.Sprintf("oauth_state=%s; Path=/; Max-Age=600; SameSite=Lax", state)
	if isSecure {
		// 線上 HTTPS 環境：加上 Secure 與 HttpOnly
		ctx.Header("Set-Cookie", baseCookie+"; Secure; HttpOnly")
	} else {
		// 本地開發 HTTP 環境：不能加 Secure，否則瀏覽器不會帶 cookie 回來，導致 state 驗證失敗
		ctx.Header("Set-Cookie", baseCookie+"; HttpOnly")
	}
	
	// 重定向到LINE授權頁面
	authURL := c.oauthService.GetLineAuthURL(state)
	ctx.Redirect(http.StatusFound, authURL)
}

// LINE登入回調
func (c *OAuthController) LineCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	
	// 從cookie中獲取state參數
	cookieState, err := ctx.Cookie("oauth_state")
	if err != nil {
		logger.Error("無法獲取state cookie", err, logrus.Fields{
			"code": code,
			"state": state,
			"cookies": ctx.Request.Header.Get("Cookie"),
			"user_agent": ctx.GetHeader("User-Agent"),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state parameter - cookie not found",
		})
		return
	}
	
	// 驗證state參數
	if state != cookieState {
		logger.Error("State參數不匹配", nil, logrus.Fields{
			"code": code,
			"state": state,
			"cookieState": cookieState,
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state parameter - mismatch",
		})
		return
	}
	
	// 清除state cookie
	ctx.Header("Set-Cookie", "oauth_state=; Path=/; Max-Age=0; Secure; HttpOnly; SameSite=Lax")
	
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
	// 不設定 Domain，讓瀏覽器自動處理
	ctx.SetCookie("auth_token", token, 3600*24, "/", "", true, true)
	
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
