package models

import (
	"database/sql"
	"time"
)

type Customer struct {
	ID            int        `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"-" db:"password"` // 不序列化密碼
	Phone         *string    `json:"phone,omitempty" db:"phone"`
	Address       *string    `json:"address,omitempty" db:"address"`
	BirthDate     *time.Time `json:"birth_date,omitempty" db:"birth_date"`
	Gender        *string    `json:"gender,omitempty" db:"gender"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	LastLogin     *time.Time `json:"last_login,omitempty" db:"last_login"`
	LoginCount    int        `json:"login_count" db:"login_count"`
	ProfileData   string     `json:"profile_data,omitempty" db:"profile_data"` // JSON 格式
	OAuthProvider *string    `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthID       *string    `json:"oauth_id,omitempty" db:"oauth_id"`
	OAuthData     *string    `json:"oauth_data,omitempty" db:"oauth_data"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *Customer) error {
	query := `
		INSERT INTO customers (name, email, password, phone, address, birth_date, gender, is_active, email_verified, login_count, profile_data) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, customer.Name, customer.Email, customer.Password, 
		customer.Phone, customer.Address, customer.BirthDate, customer.Gender, 
		customer.IsActive, customer.EmailVerified, customer.LoginCount, customer.ProfileData)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	customer.ID = int(id)
	
	// 獲取創建時間
	query = `SELECT created_at, updated_at FROM customers WHERE id = ?`
	err = r.db.QueryRow(query, customer.ID).Scan(&customer.CreatedAt, &customer.UpdatedAt)
	
	return err
}

func (r *CustomerRepository) GetByEmail(email string) (*Customer, error) {
	customer := &Customer{}
	query := `SELECT id, name, email, password, phone, address, birth_date, gender, is_active, email_verified, last_login, login_count, profile_data, oauth_provider, oauth_id, oauth_data, created_at, updated_at FROM customers WHERE email = ?`
	
	err := r.db.QueryRow(query, email).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Password,
		&customer.Phone, &customer.Address, &customer.BirthDate, &customer.Gender,
		&customer.IsActive, &customer.EmailVerified, &customer.LastLogin, &customer.LoginCount,
		&customer.ProfileData, &customer.OAuthProvider, &customer.OAuthID, &customer.OAuthData,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

func (r *CustomerRepository) GetByID(id int) (*Customer, error) {
	customer := &Customer{}
	query := `SELECT id, name, email, password, phone, address, birth_date, gender, is_active, email_verified, last_login, login_count, profile_data, oauth_provider, oauth_id, oauth_data, created_at, updated_at FROM customers WHERE id = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Password,
		&customer.Phone, &customer.Address, &customer.BirthDate, &customer.Gender,
		&customer.IsActive, &customer.EmailVerified, &customer.LastLogin, &customer.LoginCount,
		&customer.ProfileData, &customer.OAuthProvider, &customer.OAuthID, &customer.OAuthData,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

func (r *CustomerRepository) GetAll() ([]*Customer, error) {
	query := `SELECT id, name, email, phone, address, birth_date, gender, is_active, email_verified, last_login, login_count, profile_data, oauth_provider, oauth_id, oauth_data, created_at, updated_at FROM customers ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		customer := &Customer{}
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, 
			&customer.BirthDate, &customer.Gender, &customer.IsActive, &customer.EmailVerified, 
			&customer.LastLogin, &customer.LoginCount, &customer.ProfileData, &customer.OAuthProvider, 
			&customer.OAuthID, &customer.OAuthData, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (r *CustomerRepository) Update(customer *Customer) error {
	query := `
		UPDATE customers 
		SET name = ?, email = ?, phone = ?, address = ?, birth_date = ?, gender = ?, 
		    is_active = ?, email_verified = ?, profile_data = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`
	
	_, err := r.db.Exec(query, customer.Name, customer.Email, customer.Phone, customer.Address,
		customer.BirthDate, customer.Gender, customer.IsActive, customer.EmailVerified,
		customer.ProfileData, customer.ID)
	return err
}

func (r *CustomerRepository) Delete(id int) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CustomerRepository) UpdateLoginInfo(id int) error {
	query := `UPDATE customers SET last_login = CURRENT_TIMESTAMP, login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CustomerRepository) UpdateStatus(id int, isActive bool) error {
	query := `UPDATE customers SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, isActive, id)
	return err
}

func (r *CustomerRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM customers`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// OAuth相關方法
func (r *CustomerRepository) GetByOAuthID(provider, oauthID string) (*Customer, error) {
	customer := &Customer{}
	query := `SELECT id, name, email, password, phone, address, birth_date, gender, is_active, email_verified, last_login, login_count, profile_data, oauth_provider, oauth_id, oauth_data, created_at, updated_at FROM customers WHERE oauth_provider = ? AND oauth_id = ?`
	
	err := r.db.QueryRow(query, provider, oauthID).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Password,
		&customer.Phone, &customer.Address, &customer.BirthDate, &customer.Gender,
		&customer.IsActive, &customer.EmailVerified, &customer.LastLogin, &customer.LoginCount,
		&customer.ProfileData, &customer.OAuthProvider, &customer.OAuthID, &customer.OAuthData,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

func (r *CustomerRepository) UpdateOAuthData(id int, provider, oauthID, oauthData string) error {
	query := `UPDATE customers SET oauth_provider = ?, oauth_id = ?, oauth_data = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, provider, oauthID, oauthData, id)
	return err
}
