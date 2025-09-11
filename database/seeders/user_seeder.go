package seeders

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// UserSeeder 用戶 seeder
type UserSeeder struct {
	db *sql.DB
}

// NewUserSeeder 創建用戶 seeder
func NewUserSeeder(db *sql.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// Run 執行用戶 seeder
func (s *UserSeeder) Run() error {
	// 檢查是否已經有測試用戶
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email IN ('admin@example.com', 'merchant@example.com', 'user@example.com')").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("測試用戶已存在，跳過用戶 seeder")
		return nil
	}

	log.Println("開始執行用戶 seeder...")

	// 生成密碼 hash
	password := "111111"
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

	log.Println("用戶 seeder 執行完成!")
	return nil
}

// Clear 清除測試用戶
func (s *UserSeeder) Clear() error {
	log.Println("清除測試用戶...")
	
	_, err := s.db.Exec("DELETE FROM users WHERE email IN ('admin@example.com', 'merchant@example.com', 'user@example.com')")
	if err != nil {
		return err
	}
	
	log.Println("測試用戶已清除")
	return nil
}
