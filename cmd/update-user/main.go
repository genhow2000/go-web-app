package main

import (
	"fmt"
	"log"
	"os"
	"go-simple-app/database"
	"go-simple-app/models"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: update-user <user_id> <new_role>")
		fmt.Println("例如: update-user 2 admin")
		os.Exit(1)
	}

	userID := os.Args[1]
	newRole := os.Args[2]

	// 初始化SQLite資料庫
	if err := database.Init(); err != nil {
		log.Fatal("SQLite資料庫初始化失敗:", err)
	}
	defer database.Close()

	// 初始化 Repository
	userRepo := models.NewUserRepository(database.DB)

	// 更新用戶角色
	err := userRepo.UpdateRole(2, newRole)
	if err != nil {
		log.Fatal("更新用戶角色失敗:", err)
	}

	fmt.Printf("用戶 ID %s 的角色已更新為 %s\n", userID, newRole)
}
