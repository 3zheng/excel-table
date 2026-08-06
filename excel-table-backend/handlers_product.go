package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/3zheng/logger"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ModelNo      string  `json:"model_no"`
	Brand        *string `json:"brand"`
	Commodities  *string `json:"commodities"`
	Descriptions *string `json:"descriptions"`
	Unit         *string `json:"unit"`
	UpdatedAt    string  `json:"updated_at"`
	UpdatedBy    *string `json:"updated_by"`
	// 价格字段（来自 product_prices，按地区）
	ExportRefPrice *float64 `json:"export_ref_price"`
	ImportRefCost  *float64 `json:"import_ref_cost"`
}

// 获取产品列表（含当前地区价格）
func handleProductList(c *gin.Context) {
	userRole := c.GetString("user_role")
	userID := c.GetInt("user_id")
	region := c.Query("region")

	// 非管理员只能看自己地区的产品
	if userRole != "admin" {
		regions, err := getUserRegions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 如果没指定地区或指定的地区不在权限内，返回空
		if region == "" {
			if len(regions) > 0 {
				region = regions[0] // 默认取第一个有权限的地区
			} else {
				c.JSON(http.StatusOK, []Product{})
				return
			}
		} else {
			allowed := false
			for _, r := range regions {
				if r == region {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问该地区产品"})
				return
			}
		}
	}

	rows, err := db.Query(`
		SELECT p.model_no, p.brand, p.commodities, p.descriptions, p.unit,
		p.updated_at, p.updated_by,
		pp.export_ref_price, pp.import_ref_cost
		FROM products p
		LEFT JOIN product_prices pp ON p.model_no = pp.model_no AND pp.region = ?
		ORDER BY p.model_no
	`, region)
	if err != nil {
		logger.Error("handleProductList 查询失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ModelNo, &p.Brand, &p.Commodities, &p.Descriptions, &p.Unit,
			&p.UpdatedAt, &p.UpdatedBy, &p.ExportRefPrice, &p.ImportRefCost)
		result = append(result, p)
	}
	if result == nil {
		result = []Product{}
	}
	c.JSON(http.StatusOK, result)
}

// 新增产品
func handleProductCreate(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	canEdit := c.GetBool("can_edit_products")

	if userRole != "admin" && !canEdit {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var req struct {
		Product
		Region string `json:"region"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// 插入 products 表
	_, err := db.Exec(`
		INSERT INTO products (model_no, brand, commodities, descriptions, unit, updated_at, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.ModelNo, req.Brand, req.Commodities, req.Descriptions, req.Unit, now, username)
	if err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "型号已存在"})
			return
		}
		logger.Error("handleProductCreate 插入 products 失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 插入 product_prices 表（如果有地区和价格信息）
	if req.Region != "" {
		_, err = db.Exec(`
			INSERT INTO product_prices (model_no, region, export_ref_price, import_ref_cost, updated_at, updated_by)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			export_ref_price = VALUES(export_ref_price),
			import_ref_cost = VALUES(import_ref_cost),
			updated_at = VALUES(updated_at),
			updated_by = VALUES(updated_by)
		`, req.ModelNo, req.Region, req.ExportRefPrice, req.ImportRefCost, now, username)
		if err != nil {
			logger.Error("handleProductCreate 插入 product_prices 失败", "err", err)
		}
	}

	logger.Info("handleProductCreate 创建成功", "model_no", req.ModelNo)
	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// 更新产品（model_no 不可改）
func handleProductUpdate(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	canEdit := c.GetBool("can_edit_products")
	modelNo := c.Param("model_no")

	if userRole != "admin" && !canEdit {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var req struct {
		Product
		Region string `json:"region"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// 更新 products 基础属性
	_, err := db.Exec(`
		UPDATE products SET brand = ?, commodities = ?, descriptions = ?, unit = ?,
		updated_at = ?, updated_by = ? WHERE model_no = ?
	`, req.Brand, req.Commodities, req.Descriptions, req.Unit, now, username, modelNo)
	if err != nil {
		logger.Error("handleProductUpdate 更新 products 失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// upsert product_prices
	if req.Region != "" {
		_, err = db.Exec(`
			INSERT INTO product_prices (model_no, region, export_ref_price, import_ref_cost, updated_at, updated_by)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			export_ref_price = VALUES(export_ref_price),
			import_ref_cost = VALUES(import_ref_cost),
			updated_at = VALUES(updated_at),
			updated_by = VALUES(updated_by)
		`, modelNo, req.Region, req.ExportRefPrice, req.ImportRefCost, now, username)
		if err != nil {
			logger.Error("handleProductUpdate 更新 product_prices 失败", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	logger.Info("handleProductUpdate 更新成功", "model_no", modelNo)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 删除产品（先查 invoice_items 引用）
func handleProductDelete(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以删除"})
		return
	}

	modelNo := c.Param("model_no")
	force := c.Query("force") == "true" // 前端确认后带上 force=true

	// 查询 invoice_items 引用数量
	var count int
	db.QueryRow("SELECT COUNT(*) FROM invoice_items WHERE model_no = ?", modelNo).Scan(&count)

	if count > 0 && force != true {
		// 有引用且未强制确认，返回引用数量让前端弹窗提示
		c.JSON(http.StatusConflict, gin.H{
			"error":       "该型号已在表单中使用",
			"ref_count":   count,
			"need_confirm": true,
		})
		return
	}

	// 删除 product_prices 和 products（invoice_items 保留）
	db.Exec("DELETE FROM product_prices WHERE model_no = ?", modelNo)
	_, err := db.Exec("DELETE FROM products WHERE model_no = ?", modelNo)
	if err != nil {
		logger.Error("handleProductDelete 删除失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("handleProductDelete 删除成功", "model_no", modelNo)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 一键复制地区价格（跳过已存在的记录）
func handleProductCopy(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以复制"})
		return
	}

	var req struct {
		FromRegion string `json:"from_region"`
		ToRegion   string `json:"to_region"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.FromRegion == "" || req.ToRegion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择源地区和目标地区"})
		return
	}
	if req.FromRegion == req.ToRegion {
		c.JSON(http.StatusBadRequest, gin.H{"error": "源地区和目标地区不能相同"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// INSERT IGNORE：已存在的记录跳过，不覆盖
	result, err := db.Exec(`
		INSERT IGNORE INTO product_prices (model_no, region, export_ref_price, import_ref_cost, updated_at, updated_by)
		SELECT model_no, ?, export_ref_price, import_ref_cost, ?, ?
		FROM product_prices WHERE region = ?
	`, req.ToRegion, now, username, req.FromRegion)
	if err != nil {
		logger.Error("handleProductCopy 复制失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	logger.Info("handleProductCopy 复制成功", "from", req.FromRegion, "to", req.ToRegion, "rows", rows)
	c.JSON(http.StatusOK, gin.H{
		"message": "复制成功",
		"copied":  rows,
	})
}

// 获取所有可用地区列表（从 product_prices 和 user_regions 合并）
func handleRegionList(c *gin.Context) {
	userRole := c.GetString("user_role")
	userID := c.GetInt("user_id")

	var regions []string

	if userRole == "admin" {
		// 管理员看全部地区
		rows, err := db.Query("SELECT DISTINCT region FROM user_regions ORDER BY region")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var r string
			rows.Scan(&r)
			regions = append(regions, r)
		}
	} else {
		var err error
		regions, err = getUserRegions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if regions == nil {
		regions = []string{}
	}
	c.JSON(http.StatusOK, regions)
}

// 检查 product_prices 里某个地区是否已有某个 model_no（供前端自动补全用）
func handleProductDetail(c *gin.Context) {
	modelNo := c.Param("model_no")

	var p Product
	err := db.QueryRow(`
		SELECT model_no, brand, commodities, descriptions, unit, updated_at, updated_by
		FROM products WHERE model_no = ?
	`, modelNo).Scan(&p.ModelNo, &p.Brand, &p.Commodities, &p.Descriptions, &p.Unit,
		&p.UpdatedAt, &p.UpdatedBy)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}
