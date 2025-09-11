package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go-simple-app/database/seeders"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  ./seed run          - 執行所有 seeder")
		fmt.Println("  ./seed clear        - 清除所有測試數據")
		fmt.Println("  ./seed user         - 只執行用戶 seeder")
		fmt.Println("  ./seed clear user   - 只清除用戶測試數據")
		os.Exit(1)
	}

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

	// 創建 seeder 管理器
	seederManager := seeders.NewSeederManager(db)

	// 根據命令執行相應操作
	command := os.Args[1]
	switch command {
	case "run":
		if err := seederManager.RunAll(); err != nil {
			log.Fatal("執行 seeder 失敗:", err)
		}
		fmt.Println("所有 seeder 執行完成!")
		
	case "clear":
		if err := seederManager.ClearAll(); err != nil {
			log.Fatal("清除測試數據失敗:", err)
		}
		fmt.Println("所有測試數據已清除!")
		
	case "user":
		if err := seederManager.RunSpecific("user"); err != nil {
			log.Fatal("執行用戶 seeder 失敗:", err)
		}
		fmt.Println("用戶 seeder 執行完成!")
		
	default:
		fmt.Printf("未知命令: %s\n", command)
		os.Exit(1)
	}
}