package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"go-simple-app/config"
	"go-simple-app/models"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	config       *config.OAuthConfig
	authService  *UnifiedAuthService
}

type LineUserInfo struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

func NewOAuthService(config *config.OAuthConfig, authService *UnifiedAuthService) *OAuthService {
	return &OAuthService{
		config:      config,
		authService: authService,
	}
}

func (s *OAuthService) GetLineAuthURL(state string) string {
	// 添加調試日誌
	fmt.Printf("DEBUG: LINE OAuth Config - ClientID: %s, RedirectURL: %s\n", s.config.LINE.ClientID, s.config.LINE.RedirectURL)
	
	oauth2Config := &oauth2.Config{
		ClientID:     s.config.LINE.ClientID,
		ClientSecret: s.config.LINE.ClientSecret,
		RedirectURL:  s.config.LINE.RedirectURL,
		Scopes:       s.config.LINE.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
	}
	
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("DEBUG: Generated Auth URL: %s\n", authURL)
	
	return authURL
}

func (s *OAuthService) HandleLineCallback(ctx context.Context, code string) (models.UserInterface, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     s.config.LINE.ClientID,
		ClientSecret: s.config.LINE.ClientSecret,
		RedirectURL:  s.config.LINE.RedirectURL,
		Scopes:       s.config.LINE.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
	}
	
	// 用code換取token
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %v", err)
	}
	
	// 獲取用戶資料
	userInfo, err := s.getLineUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	
	// 查找或創建用戶
	return s.findOrCreateUser(userInfo, token.AccessToken)
}

func (s *OAuthService) getLineUserInfo(ctx context.Context, accessToken string) (*LineUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LINE API returned status: %d", resp.StatusCode)
	}
	
	var userInfo LineUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	
	return &userInfo, nil
}

func (s *OAuthService) findOrCreateUser(lineUser *LineUserInfo, accessToken string) (models.UserInterface, error) {
	// 查找現有用戶
	user, err := s.authService.userRepo.GetByOAuthID("line", lineUser.UserID)
	if err == nil && user != nil {
		// 更新OAuth資料
		s.updateOAuthData(user, lineUser, accessToken)
		return user, nil
	}
	
	// 創建新用戶（默認為customer角色）
	oauthData, _ := json.Marshal(map[string]interface{}{
		"line_user_id":    lineUser.UserID,
		"display_name":    lineUser.DisplayName,
		"picture_url":     lineUser.PictureURL,
		"status_message":  lineUser.StatusMessage,
		"access_token":    accessToken,
	})
	
	registerReq := &UnifiedRegisterRequest{
		Name:          lineUser.DisplayName,
		Email:         fmt.Sprintf("line_%s@line.local", lineUser.UserID), // 虛擬email
		Password:      "", // OAuth用戶不需要密碼
		Role:          "customer",
		OAuthProvider: &[]string{"line"}[0],
		OAuthID:       &lineUser.UserID,
		OAuthData:     &[]string{string(oauthData)}[0],
	}
	
	return s.authService.RegisterWithOAuth(registerReq)
}

func (s *OAuthService) updateOAuthData(user models.UserInterface, lineUser *LineUserInfo, accessToken string) {
	oauthData, _ := json.Marshal(map[string]interface{}{
		"line_user_id":    lineUser.UserID,
		"display_name":    lineUser.DisplayName,
		"picture_url":     lineUser.PictureURL,
		"status_message":  lineUser.StatusMessage,
		"access_token":    accessToken,
	})
	
	// 根據用戶類型更新OAuth資料
	switch user.GetRole() {
	case "customer":
		s.authService.userRepo.CustomerRepo.UpdateOAuthData(user.GetID(), "line", lineUser.UserID, string(oauthData))
	case "merchant":
		s.authService.userRepo.MerchantRepo.UpdateOAuthData(user.GetID(), "line", lineUser.UserID, string(oauthData))
	case "admin":
		s.authService.userRepo.AdminRepo.UpdateOAuthData(user.GetID(), "line", lineUser.UserID, string(oauthData))
	}
}

// GenerateToken 生成JWT token
func (s *OAuthService) GenerateToken(user models.UserInterface) (string, error) {
	return s.authService.GenerateToken(user)
}
