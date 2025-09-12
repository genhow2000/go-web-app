package seeders

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// MerchantSeeder 商戶 seeder
type MerchantSeeder struct {
	db *sql.DB
}

// NewMerchantSeeder 創建商戶 seeder
func NewMerchantSeeder(db *sql.DB) *MerchantSeeder {
	return &MerchantSeeder{db: db}
}

// Run 執行商戶 seeder
func (s *MerchantSeeder) Run() error {
	// 檢查是否已經有測試商戶
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM merchants WHERE email = 'merchant@example.com'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("測試商戶已存在，跳過商戶 seeder")
		return nil
	}

	log.Println("開始執行商戶 seeder...")

	// 生成密碼 hash
	password := "111111"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO merchants (name, email, password, phone, business_name, business_type, business_license, is_active, is_verified, last_login, login_count, business_data, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	
	_, err = s.db.Exec(query, 
		"測試商戶", 
		"merchant@example.com", 
		string(hashedPassword), 
		"0987654321", 
		"測試商店", 
		"零售", 
		"12345678", 
		true,
		false, // is_verified
		nil,   // last_login
		0,     // login_count
		"{}",  // business_data
	)
	
	if err != nil {
		return err
	}
	
	log.Println("創建商戶: 測試商戶 (merchant@example.com)")
	log.Println("商戶 seeder 執行完成!")
	return nil
}

// Clear 清除測試商戶
func (s *MerchantSeeder) Clear() error {
	log.Println("清除測試商戶...")
	
	_, err := s.db.Exec("DELETE FROM merchants WHERE email = 'merchant@example.com'")
	if err != nil {
		return err
	}
	
	log.Println("測試商戶已清除")
	return nil
}
