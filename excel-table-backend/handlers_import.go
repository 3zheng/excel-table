package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/3zheng/logger"
	"github.com/gin-gonic/gin"
)

// 进口明细行（含进口特有字段）
type ImportInvoiceItem struct {
	ID              int      `json:"id"`
	InvNo           string   `json:"inv_no"`
	InvoiceType     string   `json:"invoice_type"`
	ItemNo          *int     `json:"item_no"`
	Brand           *string  `json:"brand"`
	Commodities     *string  `json:"commodities"`
	ModelNo         *string  `json:"model_no"`
	CartonQty       *int     `json:"carton_qty"`
	UnitQty         *int     `json:"unit_qty"`
	Unit            *string  `json:"unit"`
	Descriptions    *string  `json:"descriptions"`
	ExportUnitPrice *float64 `json:"export_unit_price"`
	TotalAmount     *float64 `json:"total_amount"`
	GrossWeight     *float64 `json:"gross_weight"`
	NetWeight       *float64 `json:"net_weight"`
	Volume          *float64 `json:"volume"`
	ShareRate       *float64 `json:"share_rate"`
	SeaFreight      *float64 `json:"sea_freight"`
	CifPrice        *float64 `json:"cif_price"`
	CifTotal        *float64 `json:"cif_total"`
	CifBsTotal      *float64 `json:"cif_bs_total"`
	ImportDutyVat   *float64 `json:"import_duty_vat"`
	Transportation  *float64 `json:"transportation"`
	OthersCharge    *float64 `json:"others_charge"`
	TotalCost       *float64 `json:"total_cost"`
	ImportUnitCost  *float64 `json:"import_unit_cost"`
	SystemPrice     *float64 `json:"system_price"`
	ExchangeRate    *float64 `json:"exchange_rate"`
	PriceAlert      string   `json:"price_alert"`
}

// 进口表单完整结构（复用 Invoice 主表结构，items 用进口明细）
type ImportInvoice struct {
	ID               int                 `json:"id"`
	InvNo            string              `json:"inv_no"`
	InvoiceType      string              `json:"invoice_type"`
	Region           *string             `json:"region"`
	InvoiceDate      *DateOnly           `json:"invoice_date"`
	ContractNo       *string             `json:"contract_no"`
	ShippingLine     *string             `json:"shipping_line"`
	BlNo             *string             `json:"bl_no"`
	ContainerNo      *string             `json:"container_no"`
	PortOfLoading    *string             `json:"port_of_loading"`
	PortOfDischarge  *string             `json:"port_of_discharge"`
	FinalDestination *string             `json:"final_destination"`
	BuyerName        *string             `json:"buyer_name"`
	BuyerAddress     *string             `json:"buyer_address"`
	BuyerTel         *string             `json:"buyer_tel"`
	SellerName       *string             `json:"seller_name"`
	PaymentTerm      *string             `json:"payment_term"`
	FobTotal         *float64            `json:"fob_total"`
	Reviewed         bool                `json:"reviewed"`
	Reviewer         *string             `json:"reviewer"`
	ReviewComment    *string             `json:"review_comment"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
	Items            []ImportInvoiceItem `json:"items"`
}

// 获取当前用户有权访问的地区列表
func getUserRegions(userID int) ([]string, error) {
	rows, err := db.Query("SELECT region FROM user_regions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []string
	for rows.Next() {
		var r string
		rows.Scan(&r)
		regions = append(regions, r)
	}
	return regions, nil
}

// 构建地区过滤的 IN 子句
func buildRegionFilter(regions []string) (string, []interface{}) {
	if len(regions) == 0 {
		return "AND 1=0", nil // 没有地区权限，返回空结果
	}
	clause := "AND region IN ("
	args := []interface{}{}
	for i, r := range regions {
		if i > 0 {
			clause += ","
		}
		clause += "?"
		args = append(args, r)
	}
	clause += ")"
	return clause, args
}

// 进口表单列表
func handleImportList(c *gin.Context) {
	userRole := c.GetString("user_role")
	userID := c.GetInt("user_id")

	if userRole != "admin" && userRole != "import_input" && userRole != "import_review" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		logger.Error("StatusForbidden 403错误, handleImportList无用户权限", "userRole", userRole)
		return
	}

	invNo := c.Query("inv_no")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	regionFilter := c.Query("region")

	query := `SELECT id, inv_no, region, invoice_date, contract_no,
			  shipping_line, reviewed, reviewer, created_at
			  FROM invoices WHERE invoice_type = 'import'`
	args := []interface{}{}

	// 非管理员按地区权限过滤
	if userRole != "admin" {
		regions, err := getUserRegions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		clause, regionArgs := buildRegionFilter(regions)
		query += " " + clause
		args = append(args, regionArgs...)
	}

	if invNo != "" {
		query += " AND inv_no LIKE ?"
		args = append(args, "%"+invNo+"%")
	}
	if startDate != "" {
		query += " AND invoice_date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND invoice_date <= ?"
		args = append(args, endDate)
	}
	if regionFilter != "" {
		query += " AND region = ?"
		args = append(args, regionFilter)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error("StatusInternalServerError 500错误:", "err", err)
		return
	}
	defer rows.Close()

	type ImportSummary struct {
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

	var result []ImportSummary
	for rows.Next() {
		var inv ImportSummary
		rows.Scan(&inv.ID, &inv.InvNo, &inv.Region, &inv.InvoiceDate,
			&inv.ContractNo, &inv.ShippingLine, &inv.Reviewed, &inv.Reviewer, &inv.CreatedAt)
		result = append(result, inv)
	}
	if result == nil {
		result = []ImportSummary{}
	}
	c.JSON(http.StatusOK, result)
}

// 进口表单详情
func handleImportDetail(c *gin.Context) {
	if debugMode {
		fmt.Println("开始查询import detail")
	}
	userRole := c.GetString("user_role")
	userID := c.GetInt("user_id")
	invNo := c.Param("inv_no")

	if userRole != "admin" && userRole != "import_input" && userRole != "import_review" {
		logger.Error("handleImportDetail无用户权限", "userID", userID, "userRole", userRole)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var inv ImportInvoice
	err := db.QueryRow(`
		SELECT id, inv_no, invoice_type, region, invoice_date, contract_no,
		shipping_line, bl_no, container_no, port_of_loading, port_of_discharge,
		final_destination, buyer_name, buyer_address, buyer_tel, seller_name,
		payment_term, fob_total, reviewed, reviewer, review_comment,
		created_at, updated_at
		FROM invoices WHERE inv_no = ? AND invoice_type = 'import'
	`, invNo).Scan(
		&inv.ID, &inv.InvNo, &inv.InvoiceType, &inv.Region, &inv.InvoiceDate,
		&inv.ContractNo, &inv.ShippingLine, &inv.BlNo, &inv.ContainerNo,
		&inv.PortOfLoading, &inv.PortOfDischarge, &inv.FinalDestination,
		&inv.BuyerName, &inv.BuyerAddress, &inv.BuyerTel, &inv.SellerName,
		&inv.PaymentTerm, &inv.FobTotal, &inv.Reviewed, &inv.Reviewer,
		&inv.ReviewComment, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		logger.Error("表单不存在", "inv_no", invNo)
		c.JSON(http.StatusNotFound, gin.H{"error": "表单不存在"})
		return
	}
	if err != nil {
		logger.Error("StatusInternalServerError", "inv_no", invNo)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 非管理员校验地区权限
	if userRole != "admin" && inv.Region != nil {
		regions, err := getUserRegions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allowed := false
		for _, r := range regions {
			if r == *inv.Region {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问该地区表单"})
			return
		}
	}

	// 获取进口明细行（含所有进口字段）
	rows, err := db.Query(`
		SELECT id, inv_no, invoice_type, item_no, brand, commodities, model_no,
		carton_qty, unit_qty, unit, descriptions, export_unit_price, total_amount,
		gross_weight, net_weight, volume, share_rate, sea_freight, cif_price,
		cif_total, cif_bs_total, import_duty_vat, transportation, others_charge,
		total_cost, import_unit_cost, system_price, exchange_rate, price_alert
		FROM invoice_items WHERE inv_no = ? AND invoice_type = 'import'
		ORDER BY item_no
	`, invNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	inv.Items = []ImportInvoiceItem{}
	for rows.Next() {
		var item ImportInvoiceItem
		rows.Scan(
			&item.ID, &item.InvNo, &item.InvoiceType, &item.ItemNo,
			&item.Brand, &item.Commodities, &item.ModelNo, &item.CartonQty,
			&item.UnitQty, &item.Unit, &item.Descriptions, &item.ExportUnitPrice,
			&item.TotalAmount, &item.GrossWeight, &item.NetWeight, &item.Volume,
			&item.ShareRate, &item.SeaFreight, &item.CifPrice, &item.CifTotal,
			&item.CifBsTotal, &item.ImportDutyVat, &item.Transportation,
			&item.OthersCharge, &item.TotalCost, &item.ImportUnitCost,
			&item.SystemPrice, &item.ExchangeRate, &item.PriceAlert,
		)
		inv.Items = append(inv.Items, item)
	}
	if debugMode {
		jInv, _ := json.MarshalIndent(inv, "", "  ")
		fmt.Println("查询import detail成功, ImportInvoice:", string(jInv))
	}
	c.JSON(http.StatusOK, inv)
}

// 更新进口表单（进口录入员填写进口字段，审核员只改审核字段）
func handleImportUpdate(c *gin.Context) {
	userRole := c.GetString("user_role")
	userID := c.GetInt("user_id")
	username := c.GetString("username")
	invNo := c.Param("inv_no")

	if userRole != "admin" && userRole != "import_input" && userRole != "import_review" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	// 先查出该表单的地区，做权限校验
	var region sql.NullString
	err := db.QueryRow("SELECT region FROM invoices WHERE inv_no = ? AND invoice_type = 'import'", invNo).Scan(&region)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "表单不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 非管理员校验地区权限
	if userRole != "admin" && region.Valid {
		regions, err := getUserRegions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allowed := false
		for _, r := range regions {
			if r == region.String {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作该地区表单"})
			return
		}
	}

	var inv ImportInvoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if userRole == "import_review" {
		// 审核员只改审核字段
		var reviewer interface{}
		if inv.Reviewed {
			reviewer = username
		} else {
			reviewer = nil
		}
		db.Exec(`UPDATE invoices SET reviewed = ?, reviewer = ?, review_comment = ?, updated_at = ?
			WHERE inv_no = ? AND invoice_type = 'import'`,
			inv.Reviewed, reviewer, inv.ReviewComment, now, invNo)
		c.JSON(http.StatusOK, gin.H{"message": "审核更新成功"})
		return
	}

	// import_input 或 admin：更新主表（region 不可改，由出口表单复制而来）
	db.Exec(`
		UPDATE invoices SET invoice_date = ?, contract_no = ?, shipping_line = ?,
		bl_no = ?, container_no = ?, port_of_loading = ?, port_of_discharge = ?,
		final_destination = ?, buyer_name = ?, buyer_address = ?, buyer_tel = ?,
		seller_name = ?, payment_term = ?, fob_total = ?, updated_at = ?
		WHERE inv_no = ? AND invoice_type = 'import'
	`, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
		inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
		inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
		inv.PaymentTerm, inv.FobTotal, now, invNo)

	// 删除旧明细，重新插入（含价格预警计算）
	db.Exec("DELETE FROM invoice_items WHERE inv_no = ? AND invoice_type = 'import'", invNo)
	for i, item := range inv.Items {
		itemNo := i + 1
		priceAlert := calcPriceAlert(item.ModelNo, region.String, item.ImportUnitCost)
		db.Exec(`
			INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
			model_no, carton_qty, unit_qty, unit, descriptions, export_unit_price, total_amount,
			gross_weight, net_weight, volume, share_rate, sea_freight, cif_price,
			cif_total, cif_bs_total, import_duty_vat, transportation, others_charge,
			total_cost, import_unit_cost, system_price, exchange_rate, price_alert)
			VALUES (?, 'import', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, invNo, itemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.UnitQty, item.Unit, item.Descriptions,
			item.ExportUnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume,
			item.ShareRate, item.SeaFreight, item.CifPrice, item.CifTotal, item.CifBsTotal,
			item.ImportDutyVat, item.Transportation, item.OthersCharge,
			item.TotalCost, item.ImportUnitCost, item.SystemPrice, item.ExchangeRate, priceAlert)
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 计算价格预警等级
// 进口：import_unit_cost vs import_ref_cost
func calcPriceAlert(modelNo *string, region string, unitCost *float64) string {
	if modelNo == nil || unitCost == nil || *unitCost == 0 {
		return "normal"
	}

	var refCost float64
	err := db.QueryRow(
		"SELECT import_ref_cost FROM product_prices WHERE model_no = ? AND region = ?",
		*modelNo, region,
	).Scan(&refCost)
	if err != nil || refCost == 0 {
		return "normal"
	}

	diff := (*unitCost - refCost) / refCost
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

// 新建进口表单（直接创建，不依赖出口表单）
func handleImportCreate(c *gin.Context) {
	userRole := c.GetString("user_role")
	userID := c.GetInt("user_id")

	if userRole != "admin" && userRole != "import_input" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var inv ImportInvoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 非管理员校验地区权限
	if userRole != "admin" && inv.Region != nil {
		regions, err := getUserRegions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allowed := false
		for _, r := range regions {
			if r == *inv.Region {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作该地区表单"})
			return
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Exec(`
		INSERT INTO invoices (inv_no, invoice_type, region, invoice_date, contract_no,
		shipping_line, bl_no, container_no, port_of_loading, port_of_discharge,
		final_destination, buyer_name, buyer_address, buyer_tel, seller_name,
		payment_term, fob_total, created_at, updated_at)
		VALUES (?, 'import', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, inv.InvNo, inv.Region, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine,
		inv.BlNo, inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge,
		inv.FinalDestination, inv.BuyerName, inv.BuyerAddress, inv.BuyerTel,
		inv.SellerName, inv.PaymentTerm, inv.FobTotal, now, now)
	if err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "INV. NO 已存在"})
			return
		}
		logger.Error("handleImportCreate 创建失败", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 插入明细行
	for i, item := range inv.Items {
		itemNo := i + 1
		regionStr := ""
		if inv.Region != nil {
			regionStr = *inv.Region
		}
		priceAlert := calcPriceAlert(item.ModelNo, regionStr, item.ImportUnitCost)
		db.Exec(`
			INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
			model_no, carton_qty, unit_qty, unit, descriptions, export_unit_price, total_amount,
			gross_weight, net_weight, volume, share_rate, sea_freight, cif_price,
			cif_total, cif_bs_total, import_duty_vat, transportation, others_charge,
			total_cost, import_unit_cost, system_price, exchange_rate, price_alert)
			VALUES (?, 'import', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, inv.InvNo, itemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.UnitQty, item.Unit, item.Descriptions,
			item.ExportUnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume,
			item.ShareRate, item.SeaFreight, item.CifPrice, item.CifTotal, item.CifBsTotal,
			item.ImportDutyVat, item.Transportation, item.OthersCharge,
			item.TotalCost, item.ImportUnitCost, item.SystemPrice, item.ExchangeRate, priceAlert)
	}

	logger.Info("handleImportCreate 创建成功", "inv_no", inv.InvNo)
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "inv_no": inv.InvNo})
}

// 删除进口表单（仅管理员，独立于出口表单）
func handleImportDelete(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以删除"})
		return
	}

	invNo := c.Param("inv_no")

	db.Exec("DELETE FROM invoice_items WHERE inv_no = ? AND invoice_type = 'import'", invNo)
	db.Exec("DELETE FROM invoices WHERE inv_no = ? AND invoice_type = 'import'", invNo)

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
func createImportFromExport(invNo string) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// 检查进口表单是否已存在
	var count int
	db.QueryRow("SELECT COUNT(*) FROM invoices WHERE inv_no = ? AND invoice_type = 'import'", invNo).Scan(&count)
	if count > 0 {
		return nil // 已存在，不处理
	}

	// 查出口主表
	var inv Invoice
	err := db.QueryRow(`
		SELECT inv_no, region, invoice_date, contract_no, shipping_line, bl_no,
		container_no, port_of_loading, port_of_discharge, final_destination,
		buyer_name, buyer_address, buyer_tel, seller_name, payment_term, fob_total
		FROM invoices WHERE inv_no = ? AND invoice_type = 'export'
	`, invNo).Scan(
		&inv.InvNo, &inv.Region, &inv.InvoiceDate, &inv.ContractNo,
		&inv.ShippingLine, &inv.BlNo, &inv.ContainerNo, &inv.PortOfLoading,
		&inv.PortOfDischarge, &inv.FinalDestination, &inv.BuyerName,
		&inv.BuyerAddress, &inv.BuyerTel, &inv.SellerName, &inv.PaymentTerm, &inv.FobTotal,
	)
	if err != nil {
		return err
	}

	// 插入进口主表
	_, err = db.Exec(`
		INSERT INTO invoices (inv_no, invoice_type, region, invoice_date, contract_no,
		shipping_line, bl_no, container_no, port_of_loading, port_of_discharge,
		final_destination, buyer_name, buyer_address, buyer_tel, seller_name,
		payment_term, fob_total, created_at, updated_at)
		VALUES (?, 'import', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, inv.InvNo, inv.Region, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine,
		inv.BlNo, inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge,
		inv.FinalDestination, inv.BuyerName, inv.BuyerAddress, inv.BuyerTel,
		inv.SellerName, inv.PaymentTerm, inv.FobTotal, now, now)
	if err != nil {
		return err
	}

	// 复制出口明细行（只复制共用字段，进口特有字段留空）
	rows, err := db.Query(`
		SELECT item_no, brand, commodities, model_no, carton_qty, unit_qty,
		unit, descriptions, export_unit_price, total_amount, gross_weight, net_weight, volume
		FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'
		ORDER BY item_no
	`, invNo)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item InvoiceItem
		rows.Scan(&item.ItemNo, &item.Brand, &item.Commodities, &item.ModelNo,
			&item.CartonQty, &item.UnitQty, &item.Unit, &item.Descriptions,
			&item.ExportUnitPrice, &item.TotalAmount, &item.GrossWeight, &item.NetWeight, &item.Volume)
		db.Exec(`
			INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
			model_no, carton_qty, unit_qty, unit, descriptions, export_unit_price, total_amount,
			gross_weight, net_weight, volume)
			VALUES (?, 'import', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, invNo, item.ItemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.UnitQty, item.Unit, item.Descriptions,
			item.ExportUnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume)
	}

	return nil
}
