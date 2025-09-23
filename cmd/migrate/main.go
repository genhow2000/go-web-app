package main

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

	_ "modernc.org/sqlite"
)

func main() {
	// 從環境變數獲取資料庫路徑
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/tmp/app.db"
	}

	// 連接資料庫
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("無法連接資料庫:", err)
	}
	defer db.Close()

	// 測試連接
	if err = db.Ping(); err != nil {
		log.Fatal("資料庫連接測試失敗:", err)
	}

	fmt.Println("資料庫連接成功!")

	// 執行遷移
	if err := runMigrations(db); err != nil {
		log.Fatal("遷移執行失敗:", err)
	}

	fmt.Println("Migration 工具準備就緒")
}

// runMigrations 執行所有遷移
func runMigrations(db *sql.DB) error {
	// 刪除舊的遷移記錄表（如果存在）
	db.Exec("DROP TABLE IF EXISTS migrations")
	
	// 創建遷移記錄表
	createMigrationTable := `
		CREATE TABLE migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version INTEGER NOT NULL,
			name TEXT NOT NULL,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(version)
		)
	`
	
	if _, err := db.Exec(createMigrationTable); err != nil {
		return fmt.Errorf("創建遷移記錄表失敗: %v", err)
	}

	// 讀取遷移文件目錄
	migrationsDir := "migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("讀取遷移目錄失敗: %v", err)
	}

	// 按文件名排序
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// 執行每個遷移文件
	for _, filename := range migrationFiles {
		// 從文件名提取版本號
		versionStr := strings.Split(filename, "_")[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			fmt.Printf("跳過無效的遷移文件名: %s\n", filename)
			continue
		}

		// 檢查是否已經執行過
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", version).Scan(&count)
		if err != nil {
			return fmt.Errorf("檢查遷移狀態失敗: %v", err)
		}

		if count > 0 {
			fmt.Printf("跳過已執行的遷移: %s\n", filename)
			continue
		}

		// 讀取遷移文件內容
		filePath := filepath.Join(migrationsDir, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("讀取遷移文件失敗 %s: %v", filename, err)
		}

		// 執行遷移
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("執行遷移失敗 %s: %v", filename, err)
		}

		// 記錄遷移執行
		if _, err := db.Exec("INSERT INTO migrations (version, name) VALUES (?, ?)", version, filename); err != nil {
			return fmt.Errorf("記錄遷移執行失敗 %s: %v", filename, err)
		}

		fmt.Printf("成功執行遷移: %s\n", filename)
	}

	return nil
}