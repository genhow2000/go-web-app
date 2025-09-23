package models

import (
	"database/sql"
	"time"
)

type Merchant struct {
	ID             int        `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Email          string     `json:"email" db:"email"`
	Password       string     `json:"-" db:"password"` // 不序列化密碼
	BusinessName   *string    `json:"business_name,omitempty" db:"business_name"`
	BusinessLicense *string   `json:"business_license,omitempty" db:"business_license"`
	Phone          *string    `json:"phone,omitempty" db:"phone"`
	Address        *string    `json:"address,omitempty" db:"address"`
	BusinessType   *string    `json:"business_type,omitempty" db:"business_type"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	IsVerified     bool       `json:"is_verified" db:"is_verified"` // 商戶認證狀態
	LastLogin      *time.Time `json:"last_login,omitempty" db:"last_login"`
	LoginCount     int        `json:"login_count" db:"login_count"`
	BusinessData   string     `json:"business_data,omitempty" db:"business_data"` // JSON 格式
	OAuthProvider  *string    `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthID        *string    `json:"oauth_id,omitempty" db:"oauth_id"`
	OAuthData      *string    `json:"oauth_data,omitempty" db:"oauth_data"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type MerchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(merchant *Merchant) error {
	query := `
		INSERT INTO merchants (name, email, password, business_name, business_license, phone, address, business_type, is_active, is_verified, login_count, business_data) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, merchant.Name, merchant.Email, merchant.Password,
		merchant.BusinessName, merchant.BusinessLicense, merchant.Phone, merchant.Address,
		merchant.BusinessType, merchant.IsActive, merchant.IsVerified, merchant.LoginCount, merchant.BusinessData)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	merchant.ID = int(id)
	
	// 獲取創建時間
	query = `SELECT created_at, updated_at FROM merchants WHERE id = ?`
	err = r.db.QueryRow(query, merchant.ID).Scan(&merchant.CreatedAt, &merchant.UpdatedAt)
	
	return err
}

func (r *MerchantRepository) GetByEmail(email string) (*Merchant, error) {
	merchant := &Merchant{}
	query := `SELECT id, name, email, password, business_name, business_license, phone, address, business_type, is_active, is_verified, last_login, login_count, business_data, created_at, updated_at FROM merchants WHERE email = ?`
	
	err := r.db.QueryRow(query, email).Scan(
		&merchant.ID, &merchant.Name, &merchant.Email, &merchant.Password,
		&merchant.BusinessName, &merchant.BusinessLicense, &merchant.Phone, &merchant.Address,
		&merchant.BusinessType, &merchant.IsActive, &merchant.IsVerified, &merchant.LastLogin,
		&merchant.LoginCount, &merchant.BusinessData, &merchant.CreatedAt, &merchant.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return merchant, nil
}

func (r *MerchantRepository) GetByID(id int) (*Merchant, error) {
	merchant := &Merchant{}
	query := `SELECT id, name, email, password, business_name, business_license, phone, address, business_type, is_active, is_verified, last_login, login_count, business_data, created_at, updated_at FROM merchants WHERE id = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&merchant.ID, &merchant.Name, &merchant.Email, &merchant.Password,
		&merchant.BusinessName, &merchant.BusinessLicense, &merchant.Phone, &merchant.Address,
		&merchant.BusinessType, &merchant.IsActive, &merchant.IsVerified, &merchant.LastLogin,
		&merchant.LoginCount, &merchant.BusinessData, &merchant.CreatedAt, &merchant.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return merchant, nil
}

func (r *MerchantRepository) GetAll() ([]*Merchant, error) {
	query := `SELECT id, name, email, business_name, business_license, phone, address, business_type, is_active, is_verified, last_login, login_count, business_data, created_at, updated_at FROM merchants ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []*Merchant
	for rows.Next() {
		merchant := &Merchant{}
		err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Email, &merchant.BusinessName,
			&merchant.BusinessLicense, &merchant.Phone, &merchant.Address, &merchant.BusinessType,
			&merchant.IsActive, &merchant.IsVerified, &merchant.LastLogin, &merchant.LoginCount,
			&merchant.BusinessData, &merchant.CreatedAt, &merchant.UpdatedAt)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (r *MerchantRepository) Update(merchant *Merchant) error {
	query := `
		UPDATE merchants 
		SET name = ?, email = ?, business_name = ?, business_license = ?, phone = ?, address = ?, 
		    business_type = ?, is_active = ?, is_verified = ?, business_data = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`
	
	_, err := r.db.Exec(query, merchant.Name, merchant.Email, merchant.BusinessName,
		merchant.BusinessLicense, merchant.Phone, merchant.Address, merchant.BusinessType,
		merchant.IsActive, merchant.IsVerified, merchant.BusinessData, merchant.ID)
	return err
}

func (r *MerchantRepository) Delete(id int) error {
	query := `DELETE FROM merchants WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MerchantRepository) UpdateLoginInfo(id int) error {
	query := `UPDATE merchants SET last_login = CURRENT_TIMESTAMP, login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MerchantRepository) UpdateStatus(id int, isActive bool) error {
	query := `UPDATE merchants SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, isActive, id)
	return err
}

func (r *MerchantRepository) UpdateVerification(id int, isVerified bool) error {
	query := `UPDATE merchants SET is_verified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, isVerified, id)
	return err
}

func (r *MerchantRepository) GetByBusinessType(businessType string) ([]*Merchant, error) {
	query := `SELECT id, name, email, business_name, business_license, phone, address, business_type, is_active, is_verified, last_login, login_count, business_data, created_at, updated_at FROM merchants WHERE business_type = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, businessType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []*Merchant
	for rows.Next() {
		merchant := &Merchant{}
		err := rows.Scan(&merchant.ID, &merchant.Name, &merchant.Email, &merchant.BusinessName,
			&merchant.BusinessLicense, &merchant.Phone, &merchant.Address, &merchant.BusinessType,
			&merchant.IsActive, &merchant.IsVerified, &merchant.LastLogin, &merchant.LoginCount,
			&merchant.BusinessData, &merchant.CreatedAt, &merchant.UpdatedAt)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (r *MerchantRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM merchants`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// OAuth相關方法
func (r *MerchantRepository) GetByOAuthID(provider, oauthID string) (*Merchant, error) {
	merchant := &Merchant{}
	query := `SELECT id, name, email, password, business_name, business_license, phone, address, business_type, is_active, is_verified, last_login, login_count, business_data, oauth_provider, oauth_id, oauth_data, created_at, updated_at FROM merchants WHERE oauth_provider = ? AND oauth_id = ?`
	
	err := r.db.QueryRow(query, provider, oauthID).Scan(
		&merchant.ID, &merchant.Name, &merchant.Email, &merchant.Password,
		&merchant.BusinessName, &merchant.BusinessLicense, &merchant.Phone, &merchant.Address,
		&merchant.BusinessType, &merchant.IsActive, &merchant.IsVerified, &merchant.LastLogin,
		&merchant.LoginCount, &merchant.BusinessData, &merchant.OAuthProvider, &merchant.OAuthID, &merchant.OAuthData,
		&merchant.CreatedAt, &merchant.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return merchant, nil
}

func (r *MerchantRepository) UpdateOAuthData(id int, provider, oauthID, oauthData string) error {
	query := `UPDATE merchants SET oauth_provider = ?, oauth_id = ?, oauth_data = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, provider, oauthID, oauthData, id)
	return err
}
