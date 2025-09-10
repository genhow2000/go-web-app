package database

import (
	"database/sql"
	"fmt"
	"go-simple-app/config"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
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

// InitSQLite 初始化 SQLite 資料庫
func InitSQLite() error {
	var err error
	
	// 從環境變數獲取資料庫路徑
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// 統一使用 /tmp 目錄，本地和雲端都一樣
		dbDir := "/tmp"
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("無法創建資料庫目錄: %w", err)
		}
		dbPath = filepath.Join(dbDir, "app.db")
	}
	
	// 連接 SQLite 資料庫 (使用純Go驅動)
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("無法連接資料庫: %w", err)
	}

	// 測試連接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("資料庫連接測試失敗: %w", err)
	}

	log.Println("SQLite 資料庫連接成功!")
	
	// 創建表
	if err := createTablesSQLite(); err != nil {
		return fmt.Errorf("創建表失敗: %w", err)
	}

	return nil
}

// createTablesSQLite 創建 SQLite 表
func createTablesSQLite() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) DEFAULT 'user',
		status VARCHAR(50) DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("SQLite 資料表創建成功!")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
