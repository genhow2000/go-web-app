package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

type Migration struct {
	Version int
	Name    string
	Path    string
}


// Init 初始化 SQLite 資料庫
func Init() error {
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
	if err := createTables(); err != nil {
		return fmt.Errorf("創建表失敗: %w", err)
	}

	// 執行資料庫遷移
	if err := runMigrations(); err != nil {
		return fmt.Errorf("執行遷移失敗: %w", err)
	}

	// 執行 seeder
	seeder := NewSeeder(DB)
	if err := seeder.RunSeeders(); err != nil {
		return fmt.Errorf("執行 seeder 失敗: %w", err)
	}

	return nil
}

// createTables 創建 SQLite 表
func createTables() error {
	// 基本表創建已移至遷移文件中
	// 這裡不再創建任何表，完全依賴遷移文件
	log.Println("SQLite 資料表創建完成（由遷移文件處理）!")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// runMigrations 執行資料庫遷移
func runMigrations() error {
	// 創建遷移表
	if err := createMigrationsTable(); err != nil {
		return fmt.Errorf("創建遷移表失敗: %w", err)
	}

	// 獲取遷移檔案
	migrations, err := getMigrations()
	if err != nil {
		return fmt.Errorf("獲取遷移檔案失敗: %w", err)
	}

	// 執行遷移
	for _, migration := range migrations {
		// 檢查是否已執行
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", migration.Version).Scan(&count)
		if err != nil {
			return fmt.Errorf("檢查遷移狀態失敗: %w", err)
		}

		if count > 0 {
			log.Printf("跳過已執行的遷移: %s", migration.Name)
			continue
		}

		// 讀取遷移檔案
		content, err := ioutil.ReadFile(migration.Path)
		if err != nil {
			return fmt.Errorf("讀取遷移檔案 %s 失敗: %w", migration.Path, err)
		}

		// 執行遷移
		log.Printf("執行遷移: %s", migration.Name)
		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("執行遷移 %s 失敗: %w", migration.Name, err)
		}

		// 記錄遷移
		_, err = DB.Exec("INSERT INTO migrations (version, name) VALUES (?, ?)", migration.Version, migration.Name)
		if err != nil {
			return fmt.Errorf("記錄遷移 %s 失敗: %w", migration.Name, err)
		}

		log.Printf("遷移 %s 執行成功", migration.Name)
	}

	return nil
}

// createMigrationsTable 創建遷移記錄表
func createMigrationsTable() error {
	// 先刪除舊的遷移記錄表（如果存在）
	DB.Exec("DROP TABLE IF EXISTS migrations")
	
	// 創建新的遷移記錄表
	query := `
	CREATE TABLE migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version INTEGER NOT NULL,
		name TEXT NOT NULL,
		executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(version)
	)`
	_, err := DB.Exec(query)
	return err
}

// getMigrations 獲取所有遷移檔案
func getMigrations() ([]Migration, error) {
	migrationsDir := "migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		// 如果 migrations 目錄不存在，返回空切片
		if os.IsNotExist(err) {
			return []Migration{}, nil
		}
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			// 從檔案名提取版本號 (例如: 001_add_role_and_status.sql -> 1)
			versionStr := strings.Split(file.Name(), "_")[0]
			version, err := strconv.Atoi(versionStr)
			if err != nil {
				continue
			}

			migrations = append(migrations, Migration{
				Version: version,
				Name:    strings.TrimSuffix(file.Name(), ".sql"),
				Path:    filepath.Join(migrationsDir, file.Name()),
			})
		}
	}

	// 按版本號排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}
