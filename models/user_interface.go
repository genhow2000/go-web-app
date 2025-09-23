package models

import (
	"database/sql"
	"time"
)

// 統一的用戶接口，用於認證和權限管理
type UserInterface interface {
	GetID() int
	GetName() string
	GetEmail() string
	GetRole() string
	GetIsActive() bool
	GetLastLogin() *time.Time
	GetLoginCount() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// 確保所有用戶類型都實現了 UserInterface
var _ UserInterface = (*Customer)(nil)
var _ UserInterface = (*Merchant)(nil)
var _ UserInterface = (*Admin)(nil)

// Customer 實現 UserInterface
func (c *Customer) GetID() int                    { return c.ID }
func (c *Customer) GetName() string               { return c.Name }
func (c *Customer) GetEmail() string              { return c.Email }
func (c *Customer) GetRole() string               { return "customer" }
func (c *Customer) GetIsActive() bool             { return c.IsActive }
func (c *Customer) GetLastLogin() *time.Time      { return c.LastLogin }
func (c *Customer) GetLoginCount() int            { return c.LoginCount }
func (c *Customer) GetCreatedAt() time.Time       { return c.CreatedAt }
func (c *Customer) GetUpdatedAt() time.Time       { return c.UpdatedAt }

// Merchant 實現 UserInterface
func (m *Merchant) GetID() int                    { return m.ID }
func (m *Merchant) GetName() string               { return m.Name }
func (m *Merchant) GetEmail() string              { return m.Email }
func (m *Merchant) GetRole() string               { return "merchant" }
func (m *Merchant) GetIsActive() bool             { return m.IsActive }
func (m *Merchant) GetLastLogin() *time.Time      { return m.LastLogin }
func (m *Merchant) GetLoginCount() int            { return m.LoginCount }
func (m *Merchant) GetCreatedAt() time.Time       { return m.CreatedAt }
func (m *Merchant) GetUpdatedAt() time.Time       { return m.UpdatedAt }

// Admin 實現 UserInterface
func (a *Admin) GetID() int                       { return a.ID }
func (a *Admin) GetName() string                  { return a.Name }
func (a *Admin) GetEmail() string                 { return a.Email }
func (a *Admin) GetRole() string                  { return "admin" }
func (a *Admin) GetIsActive() bool                { return a.IsActive }
func (a *Admin) GetLastLogin() *time.Time         { return a.LastLogin }
func (a *Admin) GetLoginCount() int               { return a.LoginCount }
func (a *Admin) GetCreatedAt() time.Time          { return a.CreatedAt }
func (a *Admin) GetUpdatedAt() time.Time          { return a.UpdatedAt }

// 統一的用戶倉庫接口
type UserRepositoryInterface interface {
	GetByEmail(email string) (UserInterface, error)
	GetByID(id int) (UserInterface, error)
	GetByOAuthID(provider, oauthID string) (UserInterface, error)
	UpdateLoginInfo(id int) error
	LogLogin(userID int, ipAddress, userAgent string, success bool) error
}

// 統一的用戶倉庫，包含所有角色的倉庫
type UnifiedUserRepository struct {
	CustomerRepo *CustomerRepository
	MerchantRepo *MerchantRepository
	AdminRepo    *AdminRepository
}

func NewUnifiedUserRepository(db *sql.DB) *UnifiedUserRepository {
	return &UnifiedUserRepository{
		CustomerRepo: NewCustomerRepository(db),
		MerchantRepo: NewMerchantRepository(db),
		AdminRepo:    NewAdminRepository(db),
	}
}

func (r *UnifiedUserRepository) GetByEmail(email string) (UserInterface, error) {
	// 嘗試從所有表中查找用戶
	if customer, err := r.CustomerRepo.GetByEmail(email); err == nil {
		return customer, nil
	}
	if merchant, err := r.MerchantRepo.GetByEmail(email); err == nil {
		return merchant, nil
	}
	if admin, err := r.AdminRepo.GetByEmail(email); err == nil {
		return admin, nil
	}
	return nil, sql.ErrNoRows
}

func (r *UnifiedUserRepository) GetByID(id int) (UserInterface, error) {
	// 嘗試從所有表中查找用戶
	if customer, err := r.CustomerRepo.GetByID(id); err == nil {
		return customer, nil
	}
	if merchant, err := r.MerchantRepo.GetByID(id); err == nil {
		return merchant, nil
	}
	if admin, err := r.AdminRepo.GetByID(id); err == nil {
		return admin, nil
	}
	return nil, sql.ErrNoRows
}

func (r *UnifiedUserRepository) GetByOAuthID(provider, oauthID string) (UserInterface, error) {
	// 嘗試從所有表中查找OAuth用戶
	if customer, err := r.CustomerRepo.GetByOAuthID(provider, oauthID); err == nil {
		return customer, nil
	}
	if merchant, err := r.MerchantRepo.GetByOAuthID(provider, oauthID); err == nil {
		return merchant, nil
	}
	if admin, err := r.AdminRepo.GetByOAuthID(provider, oauthID); err == nil {
		return admin, nil
	}
	return nil, sql.ErrNoRows
}

func (r *UnifiedUserRepository) UpdateLoginInfo(id int) error {
	// 嘗試更新所有表中的登入信息
	if err := r.CustomerRepo.UpdateLoginInfo(id); err == nil {
		return nil
	}
	if err := r.MerchantRepo.UpdateLoginInfo(id); err == nil {
		return nil
	}
	if err := r.AdminRepo.UpdateLoginInfo(id); err == nil {
		return nil
	}
	return sql.ErrNoRows
}

func (r *UnifiedUserRepository) LogLogin(userID int, ipAddress, userAgent string, success bool) error {
	// 記錄登入日誌到統一的登入日誌表
	query := `INSERT INTO login_logs (user_type, user_id, ip_address, user_agent, success) VALUES (?, ?, ?, ?, ?)`
	
	// 需要先確定用戶類型
	user, err := r.GetByID(userID)
	if err != nil {
		return err
	}
	
	_, err = r.CustomerRepo.db.Exec(query, user.GetRole(), userID, ipAddress, userAgent, success)
	return err
}

// 獲取所有用戶統計
func (r *UnifiedUserRepository) GetUserStats() (map[string]interface{}, error) {
	customerCount, _ := r.CustomerRepo.Count()
	merchantCount, _ := r.MerchantRepo.Count()
	adminCount, _ := r.AdminRepo.Count()
	
	return map[string]interface{}{
		"customers": customerCount,
		"merchants": merchantCount,
		"admins":    adminCount,
		"total":     customerCount + merchantCount + adminCount,
	}, nil
}
