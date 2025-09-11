package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // 不序列化密碼
	Role      string    `json:"role" db:"role"` // "customer" 或 "admin"
	IsActive  bool      `json:"is_active" db:"is_active"` // 帳戶是否啟用
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (name, email, password, role, is_active) 
		VALUES (?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.Role, user.IsActive)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	
	// 獲取創建時間
	query = `SELECT created_at, updated_at FROM users WHERE id = ?`
	err = r.db.QueryRow(query, user.ID).Scan(&user.CreatedAt, &user.UpdatedAt)
	
	return err
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at FROM users WHERE email = ?`
	
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at FROM users WHERE id = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *UserRepository) GetAll() ([]*User, error) {
	query := `SELECT id, name, email, role, is_active, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, role = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`
	
	_, err := r.db.Exec(query, user.Name, user.Email, user.Role, user.IsActive, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// 根據角色獲取用戶
func (r *UserRepository) GetByRole(role string) ([]*User, error) {
	query := `SELECT id, name, email, role, is_active, created_at, updated_at FROM users WHERE role = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// 更新用戶狀態
func (r *UserRepository) UpdateStatus(id int, isActive bool) error {
	query := `UPDATE users SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, isActive, id)
	return err
}

// 更新用戶角色
func (r *UserRepository) UpdateRole(id int, role string) error {
	query := `UPDATE users SET role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, role, id)
	return err
}
