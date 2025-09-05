package services

import (
	"errors"
	"go-simple-app/config"
	"go-simple-app/logger"
	"go-simple-app/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *models.UserRepository
	config   *config.JWTConfig
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *models.User `json:"user"`
}

func NewAuthService(userRepo *models.UserRepository, config *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	logger.Info("開始用戶註冊流程", logrus.Fields{
		"email": req.Email,
		"name":  req.Name,
	})

	// 檢查用戶是否已存在
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		logger.Warn("註冊失敗：電子郵件已被註冊", logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("該電子郵件已被註冊")
	}

	// 雜湊密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密碼雜湊失敗", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("密碼處理失敗")
	}

	// 創建用戶
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("用戶創建失敗", err, logrus.Fields{
			"email": req.Email,
			"name":  req.Name,
		})
		return nil, errors.New("註冊失敗")
	}

	// 清除密碼
	user.Password = ""
	
	logger.Info("用戶註冊成功", logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	})
	
	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	logger.Info("開始用戶登入流程", logrus.Fields{
		"email": req.Email,
	})

	// 查找用戶
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		logger.Warn("登入失敗：用戶不存在", logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		})
		return nil, errors.New("用戶不存在或密碼錯誤")
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn("登入失敗：密碼錯誤", logrus.Fields{
			"email": req.Email,
			"user_id": user.ID,
		})
		return nil, errors.New("用戶不存在或密碼錯誤")
	}

	// 生成 JWT token
	token, err := s.generateToken(user)
	if err != nil {
		logger.Error("JWT token生成失敗", err, logrus.Fields{
			"email": req.Email,
			"user_id": user.ID,
		})
		return nil, errors.New("token 生成失敗")
	}

	// 清除密碼
	user.Password = ""

	logger.Info("用戶登入成功", logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	})

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(s.config.ExpiresIn) * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無效的 token 簽名方法")
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(float64)
		if !ok {
			return nil, errors.New("無效的 token")
		}

		user, err := s.userRepo.GetByID(int(userID))
		if err != nil {
			return nil, errors.New("用戶不存在")
		}

		user.Password = ""
		return user, nil
	}

	return nil, errors.New("無效的 token")
}
