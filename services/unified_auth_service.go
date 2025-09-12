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

type UnifiedAuthService struct {
	userRepo *models.UnifiedUserRepository
	config   *config.JWTConfig
}

type UnifiedLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // 可選，用於指定登入類型
}

type UnifiedRegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=customer merchant admin"`
	// 角色專用字段
	Phone         *string `json:"phone,omitempty"`
	Address       *string `json:"address,omitempty"`
	BirthDate     *string `json:"birth_date,omitempty"`
	Gender        *string `json:"gender,omitempty"`
	BusinessName  *string `json:"business_name,omitempty"`
	BusinessType  *string `json:"business_type,omitempty"`
	Department    *string `json:"department,omitempty"`
	AdminLevel    *string `json:"admin_level,omitempty"`
}

type UnifiedLoginResponse struct {
	Token string           `json:"token"`
	User  models.UserInterface `json:"user"`
}

func NewUnifiedAuthService(userRepo *models.UnifiedUserRepository, config *config.JWTConfig) *UnifiedAuthService {
	return &UnifiedAuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *UnifiedAuthService) Register(req *UnifiedRegisterRequest) (models.UserInterface, error) {
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

	// 雜湊密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密碼雜湊失敗", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("密碼處理失敗")
	}

	// 根據角色創建不同的用戶
	switch req.Role {
	case "customer":
		return s.createCustomer(req, string(hashedPassword))
	case "merchant":
		return s.createMerchant(req, string(hashedPassword))
	case "admin":
		return s.createAdmin(req, string(hashedPassword))
	default:
		return nil, errors.New("無效的角色")
	}
}

func (s *UnifiedAuthService) createCustomer(req *UnifiedRegisterRequest, hashedPassword string) (models.UserInterface, error) {
	customer := &models.Customer{
		Name:          req.Name,
		Email:         req.Email,
		Password:      hashedPassword,
		Phone:         req.Phone,
		Address:       req.Address,
		IsActive:      true,
		EmailVerified: false,
		LoginCount:    0,
		ProfileData:   "{}",
	}

	// 解析生日
	if req.BirthDate != nil && *req.BirthDate != "" {
		if birthDate, err := time.Parse("2006-01-02", *req.BirthDate); err == nil {
			customer.BirthDate = &birthDate
		}
	}

	if req.Gender != nil {
		customer.Gender = req.Gender
	}

	if err := s.userRepo.CustomerRepo.Create(customer); err != nil {
		logger.Error("客戶創建失敗", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("註冊失敗")
	}

	logger.Info("客戶註冊成功", logrus.Fields{
		"user_id": customer.ID,
		"email":   customer.Email,
	})
	
	return customer, nil
}

func (s *UnifiedAuthService) createMerchant(req *UnifiedRegisterRequest, hashedPassword string) (models.UserInterface, error) {
	merchant := &models.Merchant{
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashedPassword,
		Phone:        req.Phone,
		Address:      req.Address,
		BusinessType: req.BusinessType,
		IsActive:     true,
		IsVerified:   false,
		LoginCount:   0,
		BusinessData: "{}",
	}

	if req.BusinessName != nil {
		merchant.BusinessName = req.BusinessName
	}

	if err := s.userRepo.MerchantRepo.Create(merchant); err != nil {
		logger.Error("商戶創建失敗", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("註冊失敗")
	}

	logger.Info("商戶註冊成功", logrus.Fields{
		"user_id": merchant.ID,
		"email":   merchant.Email,
	})
	
	return merchant, nil
}

func (s *UnifiedAuthService) createAdmin(req *UnifiedRegisterRequest, hashedPassword string) (models.UserInterface, error) {
	adminLevel := "normal"
	if req.AdminLevel != nil {
		adminLevel = *req.AdminLevel
	}

	admin := &models.Admin{
		Name:        req.Name,
		Email:       req.Email,
		Password:    hashedPassword,
		AdminLevel:  adminLevel,
		Department:  req.Department,
		Phone:       req.Phone,
		IsActive:    true,
		LoginCount:  0,
		AdminData:   "{}",
	}

	if err := s.userRepo.AdminRepo.Create(admin); err != nil {
		logger.Error("管理員創建失敗", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("註冊失敗")
	}

	logger.Info("管理員註冊成功", logrus.Fields{
		"user_id": admin.ID,
		"email":   admin.Email,
	})
	
	return admin, nil
}

func (s *UnifiedAuthService) Login(req *UnifiedLoginRequest) (*UnifiedLoginResponse, error) {
	logger.Info("開始用戶登入流程", logrus.Fields{
		"email": req.Email,
		"role":  req.Role,
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

	// 如果指定了角色，檢查角色是否匹配
	if req.Role != "" && user.GetRole() != req.Role {
		logger.Warn("登入失敗：角色不匹配", logrus.Fields{
			"email":      req.Email,
			"user_role":  user.GetRole(),
			"login_role": req.Role,
		})
		return nil, errors.New("角色不匹配")
	}

	// 檢查帳戶是否啟用
	if !user.GetIsActive() {
		logger.Warn("登入失敗：帳戶已停用", logrus.Fields{
			"email":   req.Email,
			"user_id": user.GetID(),
		})
		return nil, errors.New("帳戶已停用，請聯繫管理員")
	}

	// 獲取用戶密碼進行驗證
	var userPassword string
	switch u := user.(type) {
	case *models.Customer:
		userPassword = u.Password
	case *models.Merchant:
		userPassword = u.Password
	case *models.Admin:
		userPassword = u.Password
	default:
		return nil, errors.New("無效的用戶類型")
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(req.Password)); err != nil {
		logger.Warn("登入失敗：密碼錯誤", logrus.Fields{
			"email":   req.Email,
			"user_id": user.GetID(),
		})
		return nil, errors.New("用戶不存在或密碼錯誤")
	}

	// 更新登入信息
	if err := s.userRepo.UpdateLoginInfo(user.GetID()); err != nil {
		logger.Warn("更新登入信息失敗", logrus.Fields{
			"email":   req.Email,
			"user_id": user.GetID(),
			"error":   err.Error(),
		})
		// 不影響登入流程，只記錄警告
	}

	// 生成 JWT token
	token, err := s.generateToken(user)
	if err != nil {
		logger.Error("JWT token生成失敗", err, logrus.Fields{
			"email":   req.Email,
			"user_id": user.GetID(),
		})
		return nil, errors.New("token 生成失敗")
	}

	logger.Info("用戶登入成功", logrus.Fields{
		"user_id": user.GetID(),
		"email":   user.GetEmail(),
		"name":    user.GetName(),
		"role":    user.GetRole(),
	})

	return &UnifiedLoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UnifiedAuthService) generateToken(user models.UserInterface) (string, error) {
	claims := jwt.MapClaims{
		"id":        user.GetID(),
		"name":      user.GetName(),
		"email":     user.GetEmail(),
		"role":      user.GetRole(),
		"is_active": user.GetIsActive(),
		"exp":       time.Now().Add(time.Duration(s.config.ExpiresIn) * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

func (s *UnifiedAuthService) ValidateToken(tokenString string) (models.UserInterface, error) {
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

		role, ok := claims["role"].(string)
		if !ok {
			return nil, errors.New("無效的 token")
		}

		// 根據角色查找用戶
		var user models.UserInterface
		var err error
		switch role {
		case "customer":
			user, err = s.userRepo.CustomerRepo.GetByID(int(userID))
		case "merchant":
			user, err = s.userRepo.MerchantRepo.GetByID(int(userID))
		case "admin":
			user, err = s.userRepo.AdminRepo.GetByID(int(userID))
		default:
			return nil, errors.New("無效的角色")
		}

		if err != nil {
			return nil, errors.New("用戶不存在")
		}

		return user, nil
	}

	return nil, errors.New("無效的 token")
}

// 獲取用戶統計數據
func (s *UnifiedAuthService) GetUserStats() (map[string]interface{}, error) {
	return s.userRepo.GetUserStats()
}

// 檢查用戶是否為指定角色
func (s *UnifiedAuthService) IsRole(user models.UserInterface, role string) bool {
	return user.GetRole() == role
}

// 檢查用戶是否為管理員
func (s *UnifiedAuthService) IsAdmin(user models.UserInterface) bool {
	return s.IsRole(user, "admin")
}

// 檢查用戶是否為商戶
func (s *UnifiedAuthService) IsMerchant(user models.UserInterface) bool {
	return s.IsRole(user, "merchant")
}

// 檢查用戶是否為客戶
func (s *UnifiedAuthService) IsCustomer(user models.UserInterface) bool {
	return s.IsRole(user, "customer")
}
