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

	// 用戶 seeder
	userSeeder := NewUserSeeder(sm.db)
	if err := userSeeder.Run(); err != nil {
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

	// 用戶 seeder
	userSeeder := NewUserSeeder(sm.db)
	if err := userSeeder.Clear(); err != nil {
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
	case "user":
		userSeeder := NewUserSeeder(sm.db)
		return userSeeder.Run()
	// 未來可以添加其他 seeder
	// case "product":
	//     productSeeder := NewProductSeeder(sm.db)
	//     return productSeeder.Run()
	default:
		log.Printf("未知的 seeder: %s", seederName)
		return nil
	}
}
