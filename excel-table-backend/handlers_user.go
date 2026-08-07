package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/3zheng/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var validUserRoles = map[string]bool{
	"admin":         true,
	"export_input":  true,
	"export_review": true,
	"import_input":  true,
	"import_review": true,
}

type UserListItem struct {
	ID              int      `json:"id"`
	Username        string   `json:"username"`
	UserRole        string   `json:"user_role"`
	CanEditProducts bool     `json:"can_edit_products"`
	CreatedAt       string   `json:"created_at"`
	Regions         []string `json:"regions"`
}

func requireAdmin(c *gin.Context) bool {
	if c.GetString("user_role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return false
	}
	return true
}

func isImportRole(role string) bool {
	return role == "import_input" || role == "import_review"
}

func syncUserRegions(tx *sql.Tx, userID int, role string, regions []string) error {
	if _, err := tx.Exec("DELETE FROM user_regions WHERE user_id = ?", userID); err != nil {
		return err
	}
	if !isImportRole(role) {
		return nil
	}
	seen := make(map[string]bool)
	for _, region := range regions {
		region = strings.TrimSpace(region)
		if region == "" || seen[region] {
			continue
		}
		seen[region] = true
		if _, err := tx.Exec("INSERT INTO user_regions (user_id, region) VALUES (?, ?)", userID, region); err != nil {
			return err
		}
	}
	return nil
}

func loadUserRegions(userID int) ([]string, error) {
	return getUserRegions(userID)
}

func getUserRoleByID(userID int) (string, error) {
	var role string
	err := db.QueryRow("SELECT user_role FROM users WHERE id = ?", userID).Scan(&role)
	return role, err
}

func handleUserList(c *gin.Context) {
	currentRole := c.GetString("user_role")
	currentID := c.GetInt("user_id")

	var rows *sql.Rows
	var err error

	if currentRole == "admin" {
		rows, err = db.Query(`
			SELECT id, username, user_role, can_edit_products, created_at
			FROM users ORDER BY id
		`)
	} else {
		rows, err = db.Query(`
			SELECT id, username, user_role, can_edit_products, created_at
			FROM users WHERE id = ?
		`, currentID)
	}

	if err != nil {
		logger.Error("handleUserList 查询失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []UserListItem
	for rows.Next() {
		var u UserListItem
		if err := rows.Scan(&u.ID, &u.Username, &u.UserRole, &u.CanEditProducts, &u.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		regions, err := loadUserRegions(u.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if regions == nil {
			regions = []string{}
		}
		u.Regions = regions
		result = append(result, u)
	}
	if result == nil {
		result = []UserListItem{}
	}
	c.JSON(http.StatusOK, result)
}

func handleUserCreate(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	var req struct {
		Username        string   `json:"username"`
		Password        string   `json:"password"`
		UserRole        string   `json:"user_role"`
		CanEditProducts bool     `json:"can_edit_products"`
		Regions         []string `json:"regions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}
	if !validUserRoles[req.UserRole] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
		return
	}
	// 禁止创建管理员
	if req.UserRole == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "不允许创建管理员账号"})
		return
	}
	if isImportRole(req.UserRole) && len(req.Regions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "进口类用户至少分配一个地区"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO users (username, password_hash, user_role, can_edit_products)
		VALUES (?, ?, ?, ?)
	`, req.Username, string(hash), req.UserRole, req.CanEditProducts)
	if err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID64, _ := res.LastInsertId()
	userID := int(userID64)
	if err := syncUserRegions(tx, userID, req.UserRole, req.Regions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	regions, _ := loadUserRegions(userID)
	if regions == nil {
		regions = []string{}
	}
	c.JSON(http.StatusCreated, UserListItem{
		ID:              userID,
		Username:        req.Username,
		UserRole:        req.UserRole,
		CanEditProducts: req.CanEditProducts,
		Regions:         regions,
	})
}

func handleUserUpdate(c *gin.Context) {
	currentRole := c.GetString("user_role")
	currentID := c.GetInt("user_id")

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	targetRole, err := getUserRoleByID(userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ========== 权限校验 ==========
	if currentRole == "admin" {
		// 管理员不能修改其他管理员
		if targetRole == "admin" && userID != currentID {
			c.JSON(http.StatusForbidden, gin.H{"error": "不能修改其他管理员账号"})
			return
		}
	} else {
		// 非管理员只能修改自己
		if userID != currentID {
			c.JSON(http.StatusForbidden, gin.H{"error": "只能修改自己的信息"})
			return
		}
	}

	var req struct {
		Username        string   `json:"username"`
		Password        string   `json:"password"`
		UserRole        string   `json:"user_role"`
		CanEditProducts bool     `json:"can_edit_products"`
		Regions         []string `json:"regions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// ---------- 非管理员：只允许改密码 ----------
	if currentRole != "admin" {
		if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入新密码"})
			return
		}
		// 强制保持原有信息
		var oldUsername string
		var oldCanEdit bool
		_ = db.QueryRow("SELECT username, can_edit_products FROM users WHERE id = ?", userID).
			Scan(&oldUsername, &oldCanEdit)
		req.Username = oldUsername
		req.UserRole = targetRole
		req.CanEditProducts = oldCanEdit
	} else {
		// ---------- 管理员 ----------
		req.Username = strings.TrimSpace(req.Username)
		if req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
			return
		}

		// 编辑自己时：角色强制锁定为 admin
		if userID == currentID {
			req.UserRole = "admin"
		} else {
			// 编辑别人时不允许设为管理员
			if req.UserRole == "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "不允许将其他用户设置为管理员"})
				return
			}
			if !validUserRoles[req.UserRole] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
				return
			}
			if isImportRole(req.UserRole) && len(req.Regions) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "进口类用户至少分配一个地区"})
				return
			}
		}
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		_, err = tx.Exec(`
			UPDATE users SET username = ?, password_hash = ?, user_role = ?, can_edit_products = ?
			WHERE id = ?
		`, req.Username, string(hash), req.UserRole, req.CanEditProducts, userID)
	} else {
		_, err = tx.Exec(`
			UPDATE users SET username = ?, user_role = ?, can_edit_products = ?
			WHERE id = ?
		`, req.Username, req.UserRole, req.CanEditProducts, userID)
	}
	if err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只有管理员修改非自己时才同步地区；自己或非管理员不碰地区
	if currentRole == "admin" && userID != currentID {
		if err := syncUserRegions(tx, userID, req.UserRole, req.Regions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var createdAt string
	_ = db.QueryRow("SELECT created_at FROM users WHERE id = ?", userID).Scan(&createdAt)
	regions, _ := loadUserRegions(userID)
	if regions == nil {
		regions = []string{}
	}
	c.JSON(http.StatusOK, UserListItem{
		ID:              userID,
		Username:        req.Username,
		UserRole:        req.UserRole,
		CanEditProducts: req.CanEditProducts,
		CreatedAt:       createdAt,
		Regions:         regions,
	})
}

func handleUserDelete(c *gin.Context) {
	if !requireAdmin(c) {
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	if userID == c.GetInt("user_id") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录用户"})
		return
	}

	targetRole, err := getUserRoleByID(userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if targetRole == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能删除管理员账号"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM user_regions WHERE user_id = ?", userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res, err := tx.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}