package models

import (
	"database/sql"
	"time"
)

type Admin struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"` // 不序列化密碼
	AdminLevel   string     `json:"admin_level" db:"admin_level"` // normal, senior, super
	Department   *string    `json:"department,omitempty" db:"department"`
	Phone        *string    `json:"phone,omitempty" db:"phone"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty" db:"last_login"`
	LoginCount   int        `json:"login_count" db:"login_count"`
	AdminData    string     `json:"admin_data,omitempty" db:"admin_data"` // JSON 格式
	OAuthProvider *string   `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthID      *string    `json:"oauth_id,omitempty" db:"oauth_id"`
	OAuthData    *string    `json:"oauth_data,omitempty" db:"oauth_data"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(admin *Admin) error {
	query := `
		INSERT INTO admins (name, email, password, admin_level, department, phone, is_active, login_count, admin_data) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, admin.Name, admin.Email, admin.Password,
		admin.AdminLevel, admin.Department, admin.Phone, admin.IsActive, admin.LoginCount, admin.AdminData)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	admin.ID = int(id)
	
	// 獲取創建時間
	query = `SELECT created_at, updated_at FROM admins WHERE id = ?`
	err = r.db.QueryRow(query, admin.ID).Scan(&admin.CreatedAt, &admin.UpdatedAt)
	
	return err
}

func (r *AdminRepository) GetByEmail(email string) (*Admin, error) {
	admin := &Admin{}
	query := `SELECT id, name, email, password, admin_level, department, phone, is_active, last_login, login_count, admin_data, created_at, updated_at FROM admins WHERE email = ?`
	
	err := r.db.QueryRow(query, email).Scan(
		&admin.ID, &admin.Name, &admin.Email, &admin.Password,
		&admin.AdminLevel, &admin.Department, &admin.Phone, &admin.IsActive,
		&admin.LastLogin, &admin.LoginCount, &admin.AdminData, &admin.CreatedAt, &admin.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return admin, nil
}

func (r *AdminRepository) GetByID(id int) (*Admin, error) {
	admin := &Admin{}
	query := `SELECT id, name, email, password, admin_level, department, phone, is_active, last_login, login_count, admin_data, created_at, updated_at FROM admins WHERE id = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&admin.ID, &admin.Name, &admin.Email, &admin.Password,
		&admin.AdminLevel, &admin.Department, &admin.Phone, &admin.IsActive,
		&admin.LastLogin, &admin.LoginCount, &admin.AdminData, &admin.CreatedAt, &admin.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return admin, nil
}

func (r *AdminRepository) GetAll() ([]*Admin, error) {
	query := `SELECT id, name, email, admin_level, department, phone, is_active, last_login, login_count, admin_data, created_at, updated_at FROM admins ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.AdminLevel,
			&admin.Department, &admin.Phone, &admin.IsActive, &admin.LastLogin,
			&admin.LoginCount, &admin.AdminData, &admin.CreatedAt, &admin.UpdatedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *AdminRepository) Update(admin *Admin) error {
	query := `
		UPDATE admins 
		SET name = ?, email = ?, admin_level = ?, department = ?, phone = ?, 
		    is_active = ?, admin_data = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`
	
	_, err := r.db.Exec(query, admin.Name, admin.Email, admin.AdminLevel,
		admin.Department, admin.Phone, admin.IsActive, admin.AdminData, admin.ID)
	return err
}

func (r *AdminRepository) Delete(id int) error {
	query := `DELETE FROM admins WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AdminRepository) UpdateLoginInfo(id int) error {
	query := `UPDATE admins SET last_login = CURRENT_TIMESTAMP, login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AdminRepository) UpdateStatus(id int, isActive bool) error {
	query := `UPDATE admins SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, isActive, id)
	return err
}

func (r *AdminRepository) UpdateAdminLevel(id int, adminLevel string) error {
	query := `UPDATE admins SET admin_level = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, adminLevel, id)
	return err
}

func (r *AdminRepository) GetByAdminLevel(adminLevel string) ([]*Admin, error) {
	query := `SELECT id, name, email, admin_level, department, phone, is_active, last_login, login_count, admin_data, created_at, updated_at FROM admins WHERE admin_level = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, adminLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.AdminLevel,
			&admin.Department, &admin.Phone, &admin.IsActive, &admin.LastLogin,
			&admin.LoginCount, &admin.AdminData, &admin.CreatedAt, &admin.UpdatedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *AdminRepository) GetByDepartment(department string) ([]*Admin, error) {
	query := `SELECT id, name, email, admin_level, department, phone, is_active, last_login, login_count, admin_data, created_at, updated_at FROM admins WHERE department = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, department)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.AdminLevel,
			&admin.Department, &admin.Phone, &admin.IsActive, &admin.LastLogin,
			&admin.LoginCount, &admin.AdminData, &admin.CreatedAt, &admin.UpdatedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *AdminRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM admins`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// OAuth相關方法
func (r *AdminRepository) GetByOAuthID(provider, oauthID string) (*Admin, error) {
	admin := &Admin{}
	query := `SELECT id, name, email, password, admin_level, department, phone, is_active, last_login, login_count, admin_data, oauth_provider, oauth_id, oauth_data, created_at, updated_at FROM admins WHERE oauth_provider = ? AND oauth_id = ?`
	
	err := r.db.QueryRow(query, provider, oauthID).Scan(
		&admin.ID, &admin.Name, &admin.Email, &admin.Password,
		&admin.AdminLevel, &admin.Department, &admin.Phone, &admin.IsActive,
		&admin.LastLogin, &admin.LoginCount, &admin.AdminData, &admin.OAuthProvider, &admin.OAuthID, &admin.OAuthData,
		&admin.CreatedAt, &admin.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return admin, nil
}

func (r *AdminRepository) UpdateOAuthData(id int, provider, oauthID, oauthData string) error {
	query := `UPDATE admins SET oauth_provider = ?, oauth_id = ?, oauth_data = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, provider, oauthID, oauthData, id)
	return err
}

// 檢查管理員權限級別
func (r *AdminRepository) HasPermission(adminID int, requiredLevel string) bool {
	admin, err := r.GetByID(adminID)
	if err != nil {
		return false
	}
	
	// 權限級別：super > senior > normal
	levelHierarchy := map[string]int{
		"normal": 1,
		"senior": 2,
		"super":  3,
	}
	
	adminLevel, ok := levelHierarchy[admin.AdminLevel]
	requiredLevelInt, ok2 := levelHierarchy[requiredLevel]
	
	if !ok || !ok2 {
		return false
	}
	
	return adminLevel >= requiredLevelInt
}
