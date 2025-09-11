package database

import (
	"database/sql"
	"go-simple-app/database/seeders"
)

// Seeder 結構 (保持向後兼容)
type Seeder struct {
	db *sql.DB
}

// NewSeeder 創建新的 Seeder 實例
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// RunSeeders 執行所有 seeder (使用新的 seeder 管理器)
func (s *Seeder) RunSeeders() error {
	seederManager := seeders.NewSeederManager(s.db)
	return seederManager.RunAll()
}
