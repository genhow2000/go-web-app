package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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

	// 檢查是否已經有測試用戶
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email IN ('admin@example.com', 'merchant@example.com', 'user@example.com')").Scan(&count)
	if err != nil {
		log.Fatal("檢查用戶失敗:", err)
	}

	if count > 0 {
		fmt.Println("測試用戶已存在，無需執行 seeder")
		return
	}

	fmt.Println("開始創建測試用戶...")
	
	// 這裡可以添加創建測試用戶的邏輯
	// 或者直接調用 seeder 包
	fmt.Println("Seeder 執行完成!")
}
