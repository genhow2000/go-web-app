package seeders

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// CustomerSeeder 客戶 seeder
type CustomerSeeder struct {
	db *sql.DB
}

// NewCustomerSeeder 創建客戶 seeder
func NewCustomerSeeder(db *sql.DB) *CustomerSeeder {
	return &CustomerSeeder{db: db}
}

// Run 執行客戶 seeder
func (s *CustomerSeeder) Run() error {
	// 檢查是否已經有測試客戶
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM customers WHERE email = 'customer@example.com'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("測試客戶已存在，跳過客戶 seeder")
		return nil
	}

	log.Println("開始執行客戶 seeder...")

	// 生成密碼 hash
	password := "111111"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO customers (name, email, password, phone, address, birth_date, gender, is_active, email_verified, last_login, login_count, profile_data, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	
	_, err = s.db.Exec(query, 
		"測試客戶", 
		"customer@example.com", 
		string(hashedPassword), 
		"0912345678", 
		"台北市信義區信義路五段7號", 
		"1990-01-01", 
		"男", 
		true,
		false, // email_verified
		nil,   // last_login
		0,     // login_count
		"{}",  // profile_data
	)
	
	if err != nil {
		return err
	}
	
	log.Println("創建客戶: 測試客戶 (customer@example.com)")
	log.Println("客戶 seeder 執行完成!")
	return nil
}

// Clear 清除測試客戶
func (s *CustomerSeeder) Clear() error {
	log.Println("清除測試客戶...")
	
	_, err := s.db.Exec("DELETE FROM customers WHERE email = 'customer@example.com'")
	if err != nil {
		return err
	}
	
	log.Println("測試客戶已清除")
	return nil
}
