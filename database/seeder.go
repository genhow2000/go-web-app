package database

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Seeder 結構
type Seeder struct {
	db *sql.DB
}

// NewSeeder 創建新的 Seeder 實例
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// RunSeeders 執行所有 seeder
func (s *Seeder) RunSeeders() error {
	// 檢查是否已經執行過 seeder
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email IN ('admin@example.com', 'merchant@example.com', 'user@example.com')").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("測試用戶已存在，跳過 seeder")
		return nil
	}

	log.Println("開始執行 seeder...")

	// 執行用戶 seeder
	if err := s.seedUsers(); err != nil {
		return err
	}

	log.Println("Seeder 執行完成!")
	return nil
}

// seedUsers 創建測試用戶
func (s *Seeder) seedUsers() error {
	// 生成密碼 hash
	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []struct {
		Name     string
		Email    string
		Role     string
		IsActive bool
	}{
		{"系統管理員", "admin@example.com", "admin", true},
		{"商戶用戶", "merchant@example.com", "merchant", true},
		{"一般用戶", "user@example.com", "customer", true},
	}

	for _, user := range users {
		query := `
			INSERT INTO users (name, email, password, role, is_active, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
		
		_, err := s.db.Exec(query, user.Name, user.Email, string(hashedPassword), user.Role, user.IsActive)
		if err != nil {
			return err
		}
		
		log.Printf("創建用戶: %s (%s)", user.Name, user.Email)
	}

	return nil
}
