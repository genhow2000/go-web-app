package seeders

import (
	"database/sql"
	"log"
)

// SeederManager seeder 管理器
type SeederManager struct {
	db *sql.DB
}

// NewSeederManager 創建 seeder 管理器
func NewSeederManager(db *sql.DB) *SeederManager {
	return &SeederManager{db: db}
}

// RunAll 執行所有 seeder
func (sm *SeederManager) RunAll() error {
	log.Println("開始執行所有 seeder...")

	// 客戶 seeder
	customerSeeder := NewCustomerSeeder(sm.db)
	if err := customerSeeder.Run(); err != nil {
		return err
	}

	// 商戶 seeder
	merchantSeeder := NewMerchantSeeder(sm.db)
	if err := merchantSeeder.Run(); err != nil {
		return err
	}

	// 管理員 seeder
	adminSeeder := NewAdminSeeder(sm.db)
	if err := adminSeeder.Run(); err != nil {
		return err
	}

	// 未來可以在這裡添加其他 seeder
	// productSeeder := NewProductSeeder(sm.db)
	// if err := productSeeder.Run(); err != nil {
	//     return err
	// }

	log.Println("所有 seeder 執行完成!")
	return nil
}

// ClearAll 清除所有測試數據
func (sm *SeederManager) ClearAll() error {
	log.Println("開始清除所有測試數據...")

	// 客戶 seeder
	customerSeeder := NewCustomerSeeder(sm.db)
	if err := customerSeeder.Clear(); err != nil {
		return err
	}

	// 商戶 seeder
	merchantSeeder := NewMerchantSeeder(sm.db)
	if err := merchantSeeder.Clear(); err != nil {
		return err
	}

	// 管理員 seeder
	adminSeeder := NewAdminSeeder(sm.db)
	if err := adminSeeder.Clear(); err != nil {
		return err
	}

	// 未來可以在這裡添加其他 seeder 的清除
	// productSeeder := NewProductSeeder(sm.db)
	// if err := productSeeder.Clear(); err != nil {
	//     return err
	// }

	log.Println("所有測試數據已清除!")
	return nil
}

// RunSpecific 執行特定的 seeder
func (sm *SeederManager) RunSpecific(seederName string) error {
	switch seederName {
	case "customer":
		customerSeeder := NewCustomerSeeder(sm.db)
		return customerSeeder.Run()
	case "merchant":
		merchantSeeder := NewMerchantSeeder(sm.db)
		return merchantSeeder.Run()
	case "admin":
		adminSeeder := NewAdminSeeder(sm.db)
		return adminSeeder.Run()
	// 未來可以添加其他 seeder
	// case "product":
	//     productSeeder := NewProductSeeder(sm.db)
	//     return productSeeder.Run()
	default:
		log.Printf("未知的 seeder: %s", seederName)
		return nil
	}
}
