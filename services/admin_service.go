package services

import (
	"errors"
	"go-simple-app/logger"
	"go-simple-app/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	userRepo *models.UserRepository
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=customer admin"`
	IsActive bool   `json:"is_active"`
}

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=customer admin"`
}

func NewAdminService(userRepo *models.UserRepository) *AdminService {
	return &AdminService{
		userRepo: userRepo,
	}
}

// 獲取所有用戶
func (s *AdminService) GetAllUsers() ([]*models.User, error) {
	logger.Info("管理員獲取所有用戶列表")
	return s.userRepo.GetAll()
}

// 根據角色獲取用戶
func (s *AdminService) GetUsersByRole(role string) ([]*models.User, error) {
	logger.Info("管理員根據角色獲取用戶列表", logrus.Fields{
		"role": role,
	})
	return s.userRepo.GetByRole(role)
}

// 獲取用戶詳情
func (s *AdminService) GetUserByID(id int) (*models.User, error) {
	logger.Info("管理員獲取用戶詳情", logrus.Fields{
		"user_id": id,
	})
	return s.userRepo.GetByID(id)
}

// 創建用戶
func (s *AdminService) CreateUser(req *UserCreateRequest) (*models.User, error) {
	logger.Info("管理員創建用戶", logrus.Fields{
		"email": req.Email,
		"name":  req.Name,
		"role":  req.Role,
	})

	// 檢查用戶是否已存在
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		logger.Warn("創建用戶失敗：電子郵件已被註冊", logrus.Fields{
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
		Role:     req.Role,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("用戶創建失敗", err, logrus.Fields{
			"email": req.Email,
			"name":  req.Name,
		})
		return nil, errors.New("用戶創建失敗")
	}

	// 清除密碼
	user.Password = ""
	
	logger.Info("管理員創建用戶成功", logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
	})
	
	return user, nil
}

// 更新用戶
func (s *AdminService) UpdateUser(id int, req *UserUpdateRequest) (*models.User, error) {
	logger.Info("管理員更新用戶", logrus.Fields{
		"user_id": id,
		"email":   req.Email,
		"name":    req.Name,
		"role":    req.Role,
	})

	// 檢查用戶是否存在
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Warn("更新用戶失敗：用戶不存在", logrus.Fields{
			"user_id": id,
		})
		return nil, errors.New("用戶不存在")
	}

	// 檢查郵箱是否被其他用戶使用
	if user.Email != req.Email {
		existingUser, _ := s.userRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			logger.Warn("更新用戶失敗：電子郵件已被其他用戶使用", logrus.Fields{
				"email": req.Email,
				"user_id": id,
			})
			return nil, errors.New("該電子郵件已被其他用戶使用")
		}
	}

	// 更新用戶信息
	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	user.IsActive = req.IsActive

	if err := s.userRepo.Update(user); err != nil {
		logger.Error("用戶更新失敗", err, logrus.Fields{
			"user_id": id,
		})
		return nil, errors.New("用戶更新失敗")
	}

	// 清除密碼
	user.Password = ""
	
	logger.Info("管理員更新用戶成功", logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
	})
	
	return user, nil
}

// 更新用戶狀態
func (s *AdminService) UpdateUserStatus(id int, isActive bool) error {
	logger.Info("管理員更新用戶狀態", logrus.Fields{
		"user_id":   id,
		"is_active": isActive,
	})

	// 檢查用戶是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Warn("更新用戶狀態失敗：用戶不存在", logrus.Fields{
			"user_id": id,
		})
		return errors.New("用戶不存在")
	}

	if err := s.userRepo.UpdateStatus(id, isActive); err != nil {
		logger.Error("用戶狀態更新失敗", err, logrus.Fields{
			"user_id": id,
		})
		return errors.New("用戶狀態更新失敗")
	}

	logger.Info("管理員更新用戶狀態成功", logrus.Fields{
		"user_id":   id,
		"is_active": isActive,
	})
	
	return nil
}

// 更新用戶角色
func (s *AdminService) UpdateUserRole(id int, role string) error {
	logger.Info("管理員更新用戶角色", logrus.Fields{
		"user_id": id,
		"role":    role,
	})

	// 驗證角色
	if role != "customer" && role != "admin" {
		return errors.New("無效的角色")
	}

	// 檢查用戶是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Warn("更新用戶角色失敗：用戶不存在", logrus.Fields{
			"user_id": id,
		})
		return errors.New("用戶不存在")
	}

	if err := s.userRepo.UpdateRole(id, role); err != nil {
		logger.Error("用戶角色更新失敗", err, logrus.Fields{
			"user_id": id,
		})
		return errors.New("用戶角色更新失敗")
	}

	logger.Info("管理員更新用戶角色成功", logrus.Fields{
		"user_id": id,
		"role":    role,
	})
	
	return nil
}

// 刪除用戶
func (s *AdminService) DeleteUser(id int) error {
	logger.Info("管理員刪除用戶", logrus.Fields{
		"user_id": id,
	})

	// 檢查用戶是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Warn("刪除用戶失敗：用戶不存在", logrus.Fields{
			"user_id": id,
		})
		return errors.New("用戶不存在")
	}

	if err := s.userRepo.Delete(id); err != nil {
		logger.Error("用戶刪除失敗", err, logrus.Fields{
			"user_id": id,
		})
		return errors.New("用戶刪除失敗")
	}

	logger.Info("管理員刪除用戶成功", logrus.Fields{
		"user_id": id,
	})
	
	return nil
}

// 獲取用戶統計
func (s *AdminService) GetUserStats() (map[string]interface{}, error) {
	logger.Info("管理員獲取用戶統計")

	// 獲取總用戶數
	totalUsers, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}

	// 獲取客戶數
	customers, err := s.userRepo.GetByRole("customer")
	if err != nil {
		return nil, err
	}

	// 獲取管理員數
	admins, err := s.userRepo.GetByRole("admin")
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users": totalUsers,
		"customers":   len(customers),
		"admins":      len(admins),
	}

	logger.Info("管理員獲取用戶統計成功", logrus.Fields{
		"total_users": totalUsers,
		"customers":   len(customers),
		"admins":      len(admins),
	})

	return stats, nil
}
