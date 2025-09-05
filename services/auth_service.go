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
	Role     string `json:"role"` // 可選，默認為 "customer"
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
		"role":  req.Role,
	})

	// 檢查用戶是否已存在
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		logger.Warn("註冊失敗：電子郵件已被註冊", logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("該電子郵件已被註冊")
	}

	// 設置默認角色
	role := req.Role
	if role == "" {
		role = "customer"
	}

	// 驗證角色
	if role != "customer" && role != "admin" {
		logger.Warn("註冊失敗：無效的角色", logrus.Fields{
			"email": req.Email,
			"role":  req.Role,
		})
		return nil, errors.New("無效的角色")
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
		Role:     role,
		IsActive: true, // 默認啟用
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
		"role":    user.Role,
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

	// 檢查帳戶是否啟用
	if !user.IsActive {
		logger.Warn("登入失敗：帳戶已停用", logrus.Fields{
			"email": req.Email,
			"user_id": user.ID,
		})
		return nil, errors.New("帳戶已停用，請聯繫管理員")
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
		"role":    user.Role,
	})

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"role":     user.Role,
		"is_active": user.IsActive,
		"exp":      time.Now().Add(time.Duration(s.config.ExpiresIn) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
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

// 檢查用戶是否為管理員
func (s *AuthService) IsAdmin(user *models.User) bool {
	return user.Role == "admin"
}

// 檢查用戶是否為客戶
func (s *AuthService) IsCustomer(user *models.User) bool {
	return user.Role == "customer"
}

// 創建管理員用戶（僅供系統初始化使用）
func (s *AuthService) CreateAdmin(name, email, password string) (*models.User, error) {
	logger.Info("創建管理員用戶", logrus.Fields{
		"email": email,
		"name":  name,
	})

	// 檢查是否已存在
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("該電子郵件已被註冊")
	}

	// 雜湊密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密碼處理失敗")
	}

	// 創建管理員用戶
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "admin",
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 清除密碼
	user.Password = ""
	
	logger.Info("管理員用戶創建成功", logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	})
	
	return user, nil
}
