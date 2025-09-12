package services

import (
	"errors"
	"go-simple-app/models"
)

// UserCreateRequest 創建用戶請求
type UserCreateRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
	Role     string  `json:"role" binding:"required,oneof=customer merchant admin"`
}

// UserUpdateRequest 更新用戶請求
type UserUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UnifiedAdminService struct {
	userRepo *models.UnifiedUserRepository
}

func NewUnifiedAdminService(userRepo *models.UnifiedUserRepository) *UnifiedAdminService {
	return &UnifiedAdminService{
		userRepo: userRepo,
	}
}

// 獲取所有用戶
func (s *UnifiedAdminService) GetAllUsers() ([]models.UserInterface, error) {
	var allUsers []models.UserInterface
	
	// 獲取所有客戶
	customers, err := s.userRepo.CustomerRepo.GetAll()
	if err == nil {
		for _, customer := range customers {
			allUsers = append(allUsers, customer)
		}
	}
	
	// 獲取所有商戶
	merchants, err := s.userRepo.MerchantRepo.GetAll()
	if err == nil {
		for _, merchant := range merchants {
			allUsers = append(allUsers, merchant)
		}
	}
	
	// 獲取所有管理員
	admins, err := s.userRepo.AdminRepo.GetAll()
	if err == nil {
		for _, admin := range admins {
			allUsers = append(allUsers, admin)
		}
	}
	
	return allUsers, nil
}

// 根據ID獲取用戶
func (s *UnifiedAdminService) GetUserByID(id int) (models.UserInterface, error) {
	// 嘗試從所有表中查找用戶
	if customer, err := s.userRepo.CustomerRepo.GetByID(id); err == nil {
		return customer, nil
	}
	if merchant, err := s.userRepo.MerchantRepo.GetByID(id); err == nil {
		return merchant, nil
	}
	if admin, err := s.userRepo.AdminRepo.GetByID(id); err == nil {
		return admin, nil
	}
	return nil, errors.New("用戶不存在")
}

// 根據角色獲取用戶
func (s *UnifiedAdminService) GetUsersByRole(role string) ([]models.UserInterface, error) {
	switch role {
	case "customer":
		customers, err := s.userRepo.CustomerRepo.GetAll()
		if err != nil {
			return nil, err
		}
		var result []models.UserInterface
		for _, customer := range customers {
			result = append(result, customer)
		}
		return result, nil
	case "merchant":
		merchants, err := s.userRepo.MerchantRepo.GetAll()
		if err != nil {
			return nil, err
		}
		var result []models.UserInterface
		for _, merchant := range merchants {
			result = append(result, merchant)
		}
		return result, nil
	case "admin":
		admins, err := s.userRepo.AdminRepo.GetAll()
		if err != nil {
			return nil, err
		}
		var result []models.UserInterface
		for _, admin := range admins {
			result = append(result, admin)
		}
		return result, nil
	default:
		return nil, errors.New("無效的角色")
	}
}

// 創建用戶
func (s *UnifiedAdminService) CreateUser(req *UserCreateRequest) (models.UserInterface, error) {
	// 這裡需要 UnifiedAuthService 來創建用戶
	// 暫時返回錯誤，需要重構
	return nil, errors.New("創建用戶功能需要 UnifiedAuthService")
}

// 更新用戶
func (s *UnifiedAdminService) UpdateUser(id int, req *UserUpdateRequest) (models.UserInterface, error) {
	// 獲取現有用戶
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	
	// 根據角色更新
	switch user.GetRole() {
	case "customer":
		// 更新客戶信息
		customer := user.(*models.Customer)
		if req.Name != nil {
			customer.Name = *req.Name
		}
		if req.Email != nil {
			customer.Email = *req.Email
		}
		if req.Phone != nil {
			customer.Phone = req.Phone
		}
		if req.Address != nil {
			customer.Address = req.Address
		}
		if req.IsActive != nil {
			customer.IsActive = *req.IsActive
		}
		err = s.userRepo.CustomerRepo.Update(customer)
		return customer, err
	case "merchant":
		// 更新商戶信息
		merchant := user.(*models.Merchant)
		if req.Name != nil {
			merchant.Name = *req.Name
		}
		if req.Email != nil {
			merchant.Email = *req.Email
		}
		if req.Phone != nil {
			merchant.Phone = req.Phone
		}
		if req.Address != nil {
			merchant.Address = req.Address
		}
		if req.IsActive != nil {
			merchant.IsActive = *req.IsActive
		}
		err = s.userRepo.MerchantRepo.Update(merchant)
		return merchant, err
	case "admin":
		// 更新管理員信息
		admin := user.(*models.Admin)
		if req.Name != nil {
			admin.Name = *req.Name
		}
		if req.Email != nil {
			admin.Email = *req.Email
		}
		if req.Phone != nil {
			admin.Phone = req.Phone
		}
		if req.IsActive != nil {
			admin.IsActive = *req.IsActive
		}
		err = s.userRepo.AdminRepo.Update(admin)
		return admin, err
	default:
		return nil, errors.New("無效的角色")
	}
}

// 更新用戶狀態
func (s *UnifiedAdminService) UpdateUserStatus(id int, isActive bool) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}
	
	switch user.GetRole() {
	case "customer":
		return s.userRepo.CustomerRepo.UpdateStatus(id, isActive)
	case "merchant":
		return s.userRepo.MerchantRepo.UpdateStatus(id, isActive)
	case "admin":
		return s.userRepo.AdminRepo.UpdateStatus(id, isActive)
	default:
		return errors.New("無效的角色")
	}
}

// 更新用戶角色
func (s *UnifiedAdminService) UpdateUserRole(id int, role string) error {
	// 角色更新需要刪除舊用戶並創建新用戶
	// 這是一個複雜的操作，暫時返回錯誤
	return errors.New("角色更新功能尚未實現")
}

// 刪除用戶
func (s *UnifiedAdminService) DeleteUser(id int) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}
	
	switch user.GetRole() {
	case "customer":
		return s.userRepo.CustomerRepo.Delete(id)
	case "merchant":
		return s.userRepo.MerchantRepo.Delete(id)
	case "admin":
		return s.userRepo.AdminRepo.Delete(id)
	default:
		return errors.New("無效的角色")
	}
}

// 獲取用戶統計
func (s *UnifiedAdminService) GetUserStats() (map[string]interface{}, error) {
	return s.userRepo.GetUserStats()
}
