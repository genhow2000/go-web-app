package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"go-simple-app/config"
	"go-simple-app/database"
	"go-simple-app/models"
	"go-simple-app/services"
)

func main() {
	// 載入配置
	cfg := config.Load()

	// 初始化資料庫
	if err := database.Init(); err != nil {
		log.Fatal("資料庫初始化失敗:", err)
	}
	defer database.Close()

	// 初始化服務
	userRepo := models.NewUserRepository(database.DB)
	authService := services.NewAuthService(userRepo, &cfg.JWT)

	fmt.Println("=== 創建管理員帳戶 ===")
	fmt.Println()

	// 獲取管理員信息
	name := getInput("管理員姓名", "系統管理員")
	email := getInput("管理員郵箱", "admin@example.com")
	password := getPassword("管理員密碼")

	// 創建管理員
	admin, err := authService.CreateAdmin(name, email, password)
	if err != nil {
		log.Fatal("創建管理員失敗:", err)
	}

	fmt.Println()
	fmt.Println("✅ 管理員創建成功！")
	fmt.Printf("姓名: %s\n", admin.Name)
	fmt.Printf("郵箱: %s\n", admin.Email)
	fmt.Printf("角色: %s\n", admin.Role)
	fmt.Println()
	fmt.Println("請記住這些信息，您可以使用此帳戶登入管理後台。")
}

func getInput(prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return defaultValue
	}
	return input
}

func getPassword(prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal("讀取密碼失敗:", err)
		}
		fmt.Println()
		
		passwordStr := string(password)
		if len(passwordStr) < 6 {
			fmt.Println("密碼至少需要6個字符，請重新輸入。")
			continue
		}
		
		fmt.Printf("確認密碼: ")
		confirmPassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal("讀取確認密碼失敗:", err)
		}
		fmt.Println()
		
		if passwordStr != string(confirmPassword) {
			fmt.Println("兩次輸入的密碼不一致，請重新輸入。")
			continue
		}
		
		return passwordStr
	}
}
