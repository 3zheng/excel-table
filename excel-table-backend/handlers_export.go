package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/3zheng/logger"
	"github.com/gin-gonic/gin"
)

type Invoice struct {
	ID               int           `json:"id"`
	InvNo            string        `json:"inv_no"`
	InvoiceType      string        `json:"invoice_type"`
	Region           *string       `json:"region"`
	InvoiceDate      *DateOnly     `json:"invoice_date"`
	ContractNo       *string       `json:"contract_no"`
	ShippingLine     *string       `json:"shipping_line"`
	BlNo             *string       `json:"bl_no"`
	ContainerNo      *string       `json:"container_no"`
	PortOfLoading    *string       `json:"port_of_loading"`
	PortOfDischarge  *string       `json:"port_of_discharge"`
	FinalDestination *string       `json:"final_destination"`
	BuyerName        *string       `json:"buyer_name"`
	BuyerAddress     *string       `json:"buyer_address"`
	BuyerTel         *string       `json:"buyer_tel"`
	SellerName       *string       `json:"seller_name"`
	PaymentTerm      *string       `json:"payment_term"`
	FobTotal         *float64      `json:"fob_total"`
	Reviewed         bool          `json:"reviewed"`
	Reviewer         *string       `json:"reviewer"`
	ReviewComment    *string       `json:"review_comment"`
	CreatedAt        string        `json:"created_at"`
	UpdatedAt        string        `json:"updated_at"`
	Items            []InvoiceItem `json:"items"`
}

type InvoiceItem struct {
	ID              int      `json:"id"`
	InvNo           string   `json:"inv_no"`
	InvoiceType     string   `json:"invoice_type"`
	ItemNo          *int     `json:"item_no"`
	Brand           *string  `json:"brand"`
	Commodities     *string  `json:"commodities"`
	ModelNo         *string  `json:"model_no"`
	CartonQty       *int     `json:"carton_qty"`
	QtyPerCarton    *int     `json:"qty_per_carton"` // 新增
	UnitQty         *int     `json:"unit_qty"`
	Unit            *string  `json:"unit"`
	Descriptions    *string  `json:"descriptions"`
	ExportUnitPrice *float64 `json:"export_unit_price"`
	TotalAmount     *float64 `json:"total_amount"`
	GrossWeight     *float64 `json:"gross_weight"`
	NetWeight       *float64 `json:"net_weight"`
	Volume          *float64 `json:"volume"`
	PriceAlert      string   `json:"price_alert"`
}

func handleExportList(c *gin.Context) {
	logger.Info("进入handleExportList")
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	invNo := c.Query("inv_no")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	region := c.Query("region")

	if userRole != "admin" && userRole != "export_input" && userRole != "export_review" {
		logger.Error("无用户权限", "username", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	query := `SELECT id, inv_no, region, invoice_date, contract_no, 
			  shipping_line, reviewed, reviewer, created_at 
			  FROM invoices WHERE invoice_type = 'export'`
	args := []interface{}{}

	if invNo != "" {
		query += " AND inv_no LIKE ?"
		args = append(args, "%"+invNo+"%")
	}
	if region != "" {
		query += " AND region = ?"
		args = append(args, region)
	}
	if startDate != "" {
		query += " AND invoice_date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND invoice_date <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY created_at DESC"
	info := fmt.Sprintf("ExportList query SQL is %s", query)
	logger.Info(info)

	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Error("handleExportList StatusInternalServerError")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type InvoiceSummary struct {
		ID           int       `json:"id"`
		InvNo        string    `json:"inv_no"`
		Region       *string   `json:"region"`
		InvoiceDate  *DateOnly `json:"invoice_date"`
		ContractNo   *string   `json:"contract_no"`
		ShippingLine *string   `json:"shipping_line"`
		Reviewed     bool      `json:"reviewed"`
		Reviewer     *string   `json:"reviewer"`
		CreatedAt    string    `json:"created_at"`
	}

	var result []InvoiceSummary
	for rows.Next() {
		var inv InvoiceSummary
		rows.Scan(&inv.ID, &inv.InvNo, &inv.Region, &inv.InvoiceDate, &inv.ContractNo,
			&inv.ShippingLine, &inv.Reviewed, &inv.Reviewer, &inv.CreatedAt)
		result = append(result, inv)
	}
	if result == nil {
		result = []InvoiceSummary{}
	}
	c.JSON(http.StatusOK, result)
}

// 获取出口表单详情
func handleExportDetail(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	if userRole != "admin" && userRole != "export_input" && userRole != "export_review" {
		logger.Error("handleExportDetail 无用户权限", "username", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	invNo := c.Param("inv_no")

	var inv Invoice
	err := db.QueryRow(`
		SELECT id, inv_no, invoice_type, region, invoice_date, contract_no,
		shipping_line, bl_no, container_no, port_of_loading, port_of_discharge,
		final_destination, buyer_name, buyer_address, buyer_tel, seller_name,
		payment_term, fob_total, reviewed, reviewer, review_comment,
		created_at, updated_at
		FROM invoices WHERE inv_no = ? AND invoice_type = 'export'
	`, invNo).Scan(
		&inv.ID, &inv.InvNo, &inv.InvoiceType, &inv.Region, &inv.InvoiceDate,
		&inv.ContractNo, &inv.ShippingLine, &inv.BlNo, &inv.ContainerNo,
		&inv.PortOfLoading, &inv.PortOfDischarge, &inv.FinalDestination,
		&inv.BuyerName, &inv.BuyerAddress, &inv.BuyerTel, &inv.SellerName,
		&inv.PaymentTerm, &inv.FobTotal, &inv.Reviewed, &inv.Reviewer,
		&inv.ReviewComment, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		logger.Error("handleExportDetail表单不存在", "invNo", invNo)
		c.JSON(http.StatusNotFound, gin.H{"error": "表单不存在"})
		return
	}
	if err != nil {
		logger.Error("handleExportDetail StatusInternalServerError", "invNo", invNo, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取明细行
	rows, err := db.Query(`
		SELECT id, inv_no, invoice_type, item_no, brand, commodities, model_no,
		carton_qty, qty_per_carton, unit_qty, unit, descriptions, export_unit_price, total_amount,
		gross_weight, net_weight, volume, price_alert
		FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'
		ORDER BY item_no
	`, invNo)
	if err != nil {
		logger.Error("handleExportDetail StatusInternalServerError", "invNo", invNo, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	inv.Items = []InvoiceItem{}
	for rows.Next() {
		var item InvoiceItem
		rows.Scan(&item.ID, &item.InvNo, &item.InvoiceType, &item.ItemNo,
			&item.Brand, &item.Commodities, &item.ModelNo, &item.CartonQty,
			&item.QtyPerCarton, &item.UnitQty, &item.Unit, &item.Descriptions,
			&item.ExportUnitPrice, &item.TotalAmount, &item.GrossWeight,
			&item.NetWeight, &item.Volume, &item.PriceAlert)
		inv.Items = append(inv.Items, item)
	}

	c.JSON(http.StatusOK, inv)
}

// 新建出口表单
func handleExportCreate(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	if userRole != "admin" && userRole != "export_input" {
		logger.Error("handleExportCreate 无用户权限", "username", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var inv Invoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		logger.Error("handleExportCreate StatusBadRequest", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// 插入主表
	result, err := db.Exec(`
		INSERT INTO invoices (inv_no, invoice_type, region, invoice_date, contract_no,
		shipping_line, bl_no, container_no, port_of_loading, port_of_discharge,
		final_destination, buyer_name, buyer_address, buyer_tel, seller_name,
		payment_term, fob_total, created_at, updated_at)
		VALUES (?, 'export', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, inv.InvNo, inv.Region, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
		inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
		inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
		inv.PaymentTerm, inv.FobTotal, now, now,
	)
	if err != nil {
		// 检查是否是唯一键冲突
		if isDuplicateError(err) {
			logger.Error("handleExportCreate INV. NO 已存在", "err", err)
			c.JSON(http.StatusConflict, gin.H{"error": "INV. NO 已存在"})
			return
		}
		logger.Error("handleExportCreate StatusInternalServerError", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = result

	// 插入明细行
	for i, item := range inv.Items {
		itemNo := i + 1
		priceAlert := calcExportPriceAlert(item.ModelNo, inv.Region, item.ExportUnitPrice)
		db.Exec(`
			INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
			model_no, carton_qty, qty_per_carton, unit_qty, unit, descriptions,
			export_unit_price, total_amount, gross_weight, net_weight, volume, price_alert)
			VALUES (?, 'export', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, inv.InvNo, itemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.QtyPerCarton, item.UnitQty, item.Unit, item.Descriptions,
			item.ExportUnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume,
			priceAlert)
	}
	logger.Info("handleExportCreate创建成功", "inv_no", inv.InvNo)
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "inv_no": inv.InvNo})
}

// 更新出口表单
func handleExportUpdate(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	invNo := c.Param("inv_no")

	if userRole != "admin" && userRole != "export_input" && userRole != "export_review" {
		logger.Error("handleExportUpdate 无用户权限", "username", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var inv Invoice
	/*
		bodyBytes, _ := c.GetRawData()
		logger.Info("收到的原始请求体", "body", string(bodyBytes))
		// GetRawData 会清空 c.Request.Body，需要重新填充回去，否则 ShouldBindJSON 读不到数据
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	*/
	if err := c.ShouldBindJSON(&inv); err != nil {
		logger.Error("handleExportUpdate StatusBadRequest", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if userRole == "export_review" {
		logger.Info("走了审核员分支")
		// 审核员只能修改审核字段
		var reviewer interface{}
		if inv.Reviewed {
			reviewer = username
		} else {
			reviewer = nil
		}
		db.Exec(`UPDATE invoices SET reviewed = ?, reviewer = ?, review_comment = ?, updated_at = ?
			WHERE inv_no = ? AND invoice_type = 'export'`,
			inv.Reviewed, reviewer, inv.ReviewComment, now, invNo)

		// 审核通过时自动创建进口表单
		if inv.Reviewed {
			if err := createImportFromExport(invNo); err != nil {
				logger.Error("创建进口表单失败", "err", err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "审核更新成功"})
		return
	}

	if userRole == "export_input" {
		logger.Info("走了录入员分支", "export_input", username)
		// 录入员只能修改表单内容字段，不能碰审核字段
		db.Exec(`
			UPDATE invoices SET region = ?, invoice_date = ?, contract_no = ?, shipping_line = ?,
			bl_no = ?, container_no = ?, port_of_loading = ?, port_of_discharge = ?,
			final_destination = ?, buyer_name = ?, buyer_address = ?, buyer_tel = ?,
			seller_name = ?, payment_term = ?, fob_total = ?, updated_at = ?
			WHERE inv_no = ? AND invoice_type = 'export'
		`, inv.Region, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
			inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
			inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
			inv.PaymentTerm, inv.FobTotal, now, invNo)
	} else {
		logger.Info("走了管理员分支")
		reviewComment := ""
		if inv.ReviewComment != nil {
			reviewComment = *inv.ReviewComment
		}
		info := fmt.Sprintf("reviewed=%v, review_comment=%s\n", inv.Reviewed, reviewComment)
		logger.Info(info)
		// 管理员可以修改所有字段包括审核字段
		var reviewer interface{}
		if inv.Reviewed {
			reviewer = username
		} else {
			reviewer = nil
		}
		result, err := db.Exec(`
    	UPDATE invoices SET region = ?, invoice_date = ?, contract_no = ?, shipping_line = ?,
   		bl_no = ?, container_no = ?, port_of_loading = ?, port_of_discharge = ?,
    	final_destination = ?, buyer_name = ?, buyer_address = ?, buyer_tel = ?,
    	seller_name = ?, payment_term = ?, fob_total = ?,
    	reviewed = ?, reviewer = ?, review_comment = ?, updated_at = ?
    	WHERE inv_no = ? AND invoice_type = 'export'
		`, inv.Region, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
			inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
			inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
			inv.PaymentTerm, inv.FobTotal,
			inv.Reviewed, reviewer, inv.ReviewComment, now, invNo)
		if err != nil {
			logger.Error("UPDATE失败", "err", err)
		} else {
			rows, _ := result.RowsAffected()
			logger.Info("UPDATE成功", "rows", rows)
		}

		// 审核通过时自动创建进口表单
		if inv.Reviewed {
			if err := createImportFromExport(invNo); err != nil {
				logger.Error("创建进口表单失败", "err", err)
			}
		}
	}

	// 删除旧明细行，重新插入。使用事务来处理避免数据删除后没有成功插入导致数据丢失
	tx, err := db.Begin()
	if err != nil {
		logger.Error("invoice_items事务创建失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback() // 如果没有走到 tx.Commit()，会自动回滚

	// 删除旧明细行
	if _, err := tx.Exec("DELETE FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'", invNo); err != nil {
		logger.Error("删除旧明细行失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除旧明细失败: " + err.Error()})
		return
	}

	// 重新插入
	for i, item := range inv.Items {
		itemNo := i + 1
		priceAlert := calcExportPriceAlert(item.ModelNo, inv.Region, item.ExportUnitPrice)
		_, err := tx.Exec(`
		INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
		model_no, carton_qty, qty_per_carton, unit_qty, unit, descriptions,
		export_unit_price, total_amount, gross_weight, net_weight, volume, price_alert)
		VALUES (?, 'export', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, inv.InvNo, itemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.QtyPerCarton, item.UnitQty, item.Unit, item.Descriptions,
			item.ExportUnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume,
			priceAlert)
		if err != nil {
			logger.Error("插入明细行失败", "err", err, "item_no", itemNo)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "插入明细行失败: " + err.Error()})
			return // defer tx.Rollback() 会自动回滚，之前的 DELETE 也会撤销
		}
	}

	// 全部成功才提交
	if err := tx.Commit(); err != nil {
		logger.Error("事务提交失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 删除出口表单
func handleExportDelete(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	if userRole != "admin" {
		logger.Error("handleExportDelete 只有管理员可以删除", "username", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以删除"})
		return
	}

	invNo := c.Param("inv_no")

	db.Exec("DELETE FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'", invNo)
	db.Exec("DELETE FROM invoices WHERE inv_no = ? AND invoice_type = 'export'", invNo)

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 计算出口价格预警等级
// 出口：export_unit_price vs export_ref_price（按地区）
func calcExportPriceAlert(modelNo *string, region *string, unitPrice *float64) string {
	if modelNo == nil || region == nil || unitPrice == nil || *unitPrice == 0 {
		return "normal"
	}

	var refPrice float64
	err := db.QueryRow(
		"SELECT export_ref_price FROM product_prices WHERE model_no = ? AND region = ?",
		*modelNo, *region,
	).Scan(&refPrice)
	if err != nil || refPrice == 0 {
		return "normal"
	}

	diff := (*unitPrice - refPrice) / refPrice
	if diff < 0 {
		diff = -diff
	}
	if diff >= 0.5 {
		return "red"
	}
	if diff >= 0.1 {
		return "yellow"
	}
	return "normal"
}

// 检查是否是唯一键冲突错误
func isDuplicateError(err error) bool {
	return err != nil && len(err.Error()) > 0 &&
		(contains(err.Error(), "Duplicate entry") || contains(err.Error(), "UNIQUE constraint failed"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
