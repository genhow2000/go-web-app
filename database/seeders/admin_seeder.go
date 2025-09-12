package seeders

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// AdminSeeder 管理員 seeder
type AdminSeeder struct {
	db *sql.DB
}

// NewAdminSeeder 創建管理員 seeder
func NewAdminSeeder(db *sql.DB) *AdminSeeder {
	return &AdminSeeder{db: db}
}

// Run 執行管理員 seeder
func (s *AdminSeeder) Run() error {
	// 檢查是否已經有測試管理員
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM admins WHERE email = 'admin@example.com'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("測試管理員已存在，跳過管理員 seeder")
		return nil
	}

	log.Println("開始執行管理員 seeder...")

	// 生成密碼 hash
	password := "111111"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO admins (name, email, password, phone, admin_level, department, is_active, last_login, login_count, admin_data, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	
	_, err = s.db.Exec(query, 
		"系統管理員", 
		"admin@example.com", 
		string(hashedPassword), 
		"0911111111", 
		"super", 
		"IT部門", 
		true,
		nil,   // last_login
		0,     // login_count
		"{}",  // admin_data
	)
	
	if err != nil {
		return err
	}
	
	log.Println("創建管理員: 系統管理員 (admin@example.com)")
	log.Println("管理員 seeder 執行完成!")
	return nil
}

// Clear 清除測試管理員
func (s *AdminSeeder) Clear() error {
	log.Println("清除測試管理員...")
	
	_, err := s.db.Exec("DELETE FROM admins WHERE email = 'admin@example.com'")
	if err != nil {
		return err
	}
	
	log.Println("測試管理員已清除")
	return nil
}
