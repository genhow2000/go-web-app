package services

import (
	"errors"
	"go-simple-app/models"
)

type UserService struct {
	userRepo *models.UserRepository
}

func NewUserService(userRepo *models.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	// 檢查用戶是否已存在
	existingUser, _ := s.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("該電子郵件已被註冊")
	}

	// 如果沒有密碼，設置默認密碼
	if user.Password == "" {
		user.Password = "default123"
	}

	return s.userRepo.Create(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) GetUserCount() (int, error) {
	return s.userRepo.Count()
}
