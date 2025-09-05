package database

import (
	"database/sql"
	"fmt"
	"go-simple-app/config"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(cfg *config.DatabaseConfig) error {
	var err error
	
	// 構建連接字串
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	// 連接資料庫
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("無法連接資料庫: %w", err)
	}

	// 測試連接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("資料庫連接測試失敗: %w", err)
	}

	log.Println("資料庫連接成功!")
	
	// 創建表
	if err := createTables(); err != nil {
		return fmt.Errorf("創建表失敗: %w", err)
	}

	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("資料表創建成功!")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
