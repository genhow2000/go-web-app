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

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Migration struct {
	Version int
	Name    string
	Path    string
}

func main() {
	// 獲取資料庫路徑
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/tmp/app.db"
	}

	// 確保資料庫目錄存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("無法創建資料庫目錄: %v", err)
	}

	// 連接資料庫
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}
	defer db.Close()

	// 創建遷移表
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("創建遷移表失敗: %v", err)
	}

	// 獲取遷移檔案
	migrations, err := getMigrations()
	if err != nil {
		log.Fatalf("獲取遷移檔案失敗: %v", err)
	}

	// 執行遷移
	if err := runMigrations(db, migrations); err != nil {
		log.Fatalf("執行遷移失敗: %v", err)
	}

	fmt.Println("資料庫遷移完成!")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func getMigrations() ([]Migration, error) {
	migrationsDir := "../../migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
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

func runMigrations(db *sql.DB, migrations []Migration) error {
	for _, migration := range migrations {
		// 檢查是否已執行
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", migration.Version).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			fmt.Printf("跳過已執行的遷移: %s\n", migration.Name)
			continue
		}

		// 讀取遷移檔案
		content, err := ioutil.ReadFile(migration.Path)
		if err != nil {
			return fmt.Errorf("讀取遷移檔案 %s 失敗: %v", migration.Path, err)
		}

		// 執行遷移
		fmt.Printf("執行遷移: %s\n", migration.Name)
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("執行遷移 %s 失敗: %v", migration.Name, err)
		}

		// 記錄遷移
		_, err = db.Exec("INSERT INTO migrations (version, name) VALUES (?, ?)", migration.Version, migration.Name)
		if err != nil {
			return fmt.Errorf("記錄遷移 %s 失敗: %v", migration.Name, err)
		}

		fmt.Printf("遷移 %s 執行成功\n", migration.Name)
	}

	return nil
}
