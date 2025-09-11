package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-simple-app/database"
	"go-simple-app/logger"
)

type DBController struct{}

func NewDBController() *DBController {
	return &DBController{}
}

// ShowDBManager 顯示資料庫管理頁面
func (dc *DBController) ShowDBManager(c *gin.Context) {
	c.HTML(http.StatusOK, "db_manager.html", gin.H{
		"title": "資料庫管理",
	})
}

// GetTables 獲取所有資料表
func (dc *DBController) GetTables(c *gin.Context) {
	query := `
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		logger.Error("查詢資料表失敗", err, logrus.Fields{})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢資料表失敗"})
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			logger.Error("掃描資料表名稱失敗", err, logrus.Fields{})
			continue
		}
		tables = append(tables, tableName)
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}

// GetTableData 獲取資料表資料
func (dc *DBController) GetTableData(c *gin.Context) {
	tableName := c.Param("table")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	
	offset := (page - 1) * limit

	// 獲取總數
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var total int
	err := database.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		logger.Error("查詢資料表總數失敗", err, logrus.Fields{"table": tableName})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢資料表總數失敗"})
		return
	}

	// 獲取欄位資訊
	columnsQuery := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := database.DB.Query(columnsQuery)
	if err != nil {
		logger.Error("查詢資料表欄位失敗", err, logrus.Fields{"table": tableName})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢資料表欄位失敗"})
		return
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue sql.NullString
		
		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			continue
		}
		columns = append(columns, name)
	}

	// 獲取資料
	dataQuery := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", tableName, limit, offset)
	rows, err = database.DB.Query(dataQuery)
	if err != nil {
		logger.Error("查詢資料表資料失敗", err, logrus.Fields{"table": tableName})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢資料表資料失敗"})
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				row[col] = nil
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"columns": columns,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// ExecuteQuery 執行 SQL 查詢
func (dc *DBController) ExecuteQuery(c *gin.Context) {
	var request struct {
		Query string `json:"query"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	// 只允許 SELECT 查詢
	if len(request.Query) < 6 || request.Query[:6] != "SELECT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只允許 SELECT 查詢"})
		return
	}

	rows, err := database.DB.Query(request.Query)
	if err != nil {
		logger.Error("執行查詢失敗", err, logrus.Fields{"query": request.Query})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢執行失敗: " + err.Error()})
		return
	}
	defer rows.Close()

	// 獲取欄位名稱
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取欄位名稱失敗"})
		return
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				row[col] = nil
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"columns": columns,
		"data":    data,
		"count":   len(data),
	})
}

// GetDBStats 獲取資料庫統計資訊
func (dc *DBController) GetDBStats(c *gin.Context) {
	query := `
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		logger.Error("查詢資料表失敗", err, logrus.Fields{})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢資料表失敗"})
		return
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}

		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
		var count int
		err := database.DB.QueryRow(countQuery).Scan(&count)
		if err != nil {
			continue
		}

		stats = append(stats, map[string]interface{}{
			"table": tableName,
			"count": count,
		})
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// ShowDBLogin 顯示資料庫管理登入頁面
func (dc *DBController) ShowDBLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "db_login.html", gin.H{
		"title": "資料庫管理登入",
	})
}

// DBLogin 資料庫管理登入
func (dc *DBController) DBLogin(c *gin.Context) {
	var request struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	// 檢查密碼 (這裡使用環境變數或預設值)
	expectedPassword := os.Getenv("DB_AUTH_TOKEN")
	if expectedPassword == "" {
		expectedPassword = "system" // 預設密碼
	}

	if request.Password != expectedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密碼錯誤"})
		return
	}

	// 設置認證 cookie
	c.SetCookie("db_auth_token", expectedPassword, 3600*24, "/admin/db", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{"message": "登入成功"})
}

// DBLogout 資料庫管理登出
func (dc *DBController) DBLogout(c *gin.Context) {
	// 清除認證 cookie
	c.SetCookie("db_auth_token", "", -1, "/admin/db", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

