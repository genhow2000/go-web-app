package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"go-simple-app/config"
)

func main() {
	// 載入配置
	cfg := config.Load()

	// 連接數據庫
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("無法連接數據庫:", err)
	}
	defer db.Close()

	// 測試連接
	if err := db.Ping(); err != nil {
		log.Fatal("數據庫連接失敗:", err)
	}

	fmt.Println("數據庫連接成功")

	// 創建遷移表
	if err := createMigrationsTable(db); err != nil {
		log.Fatal("創建遷移表失敗:", err)
	}

	// 執行遷移
	if err := runMigrations(db); err != nil {
		log.Fatal("執行遷移失敗:", err)
	}

	fmt.Println("所有遷移執行完成")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func runMigrations(db *sql.DB) error {
	// 讀取遷移文件
	migrationDir := "migrations"
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("讀取遷移目錄失敗: %v", err)
	}

	// 過濾並排序SQL文件
	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// 獲取已執行的遷移
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("獲取已執行遷移失敗: %v", err)
	}

	// 執行未執行的遷移
	for _, filename := range sqlFiles {
		if _, exists := executedMigrations[filename]; exists {
			fmt.Printf("跳過已執行的遷移: %s\n", filename)
			continue
		}

		fmt.Printf("執行遷移: %s\n", filename)
		if err := executeMigration(db, filepath.Join(migrationDir, filename), filename); err != nil {
			return fmt.Errorf("執行遷移 %s 失敗: %v", filename, err)
		}
	}

	return nil
}

func getExecutedMigrations(db *sql.DB) (map[string]bool, error) {
	query := "SELECT filename FROM migrations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		executed[filename] = true
	}

	return executed, nil
}

func executeMigration(db *sql.DB, filepath, filename string) error {
	// 讀取SQL文件
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// 開始事務
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 執行SQL
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// 記錄遷移
	recordQuery := "INSERT INTO migrations (filename) VALUES ($1)"
	if _, err := tx.Exec(recordQuery, filename); err != nil {
		return err
	}

	// 提交事務
	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("遷移 %s 執行成功\n", filename)
	return nil
}
