package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Invoice struct {
	ID               int           `json:"id"`
	InvNo            string        `json:"inv_no"`
	InvoiceType      string        `json:"invoice_type"`
	Region           *string       `json:"region"`
	InvoiceDate      *string       `json:"invoice_date"`
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
	ID           int      `json:"id"`
	InvNo        string   `json:"inv_no"`
	InvoiceType  string   `json:"invoice_type"`
	ItemNo       *int     `json:"item_no"`
	Brand        *string  `json:"brand"`
	Commodities  *string  `json:"commodities"`
	ModelNo      *string  `json:"model_no"`
	CartonQty    *int     `json:"carton_qty"`
	UnitQty      *int     `json:"unit_qty"`
	Unit         *string  `json:"unit"`
	Descriptions *string  `json:"descriptions"`
	UnitPrice    *float64 `json:"unit_price"`
	TotalAmount  *float64 `json:"total_amount"`
	GrossWeight  *float64 `json:"gross_weight"`
	NetWeight    *float64 `json:"net_weight"`
	Volume       *float64 `json:"volume"`
}

// 获取出口表单详情
func handleExportDetail(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "admin" && userRole != "export_input" && userRole != "export_review" {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "表单不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取明细行
	rows, err := db.Query(`
		SELECT id, inv_no, invoice_type, item_no, brand, commodities, model_no,
		carton_qty, unit_qty, unit, descriptions, unit_price, total_amount,
		gross_weight, net_weight, volume
		FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'
		ORDER BY item_no
	`, invNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	inv.Items = []InvoiceItem{}
	for rows.Next() {
		var item InvoiceItem
		rows.Scan(&item.ID, &item.InvNo, &item.InvoiceType, &item.ItemNo,
			&item.Brand, &item.Commodities, &item.ModelNo, &item.CartonQty,
			&item.UnitQty, &item.Unit, &item.Descriptions, &item.UnitPrice,
			&item.TotalAmount, &item.GrossWeight, &item.NetWeight, &item.Volume)
		inv.Items = append(inv.Items, item)
	}

	c.JSON(http.StatusOK, inv)
}

// 新建出口表单
func handleExportCreate(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "admin" && userRole != "export_input" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var inv Invoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// 插入主表
	result, err := db.Exec(`
		INSERT INTO invoices (inv_no, invoice_type, invoice_date, contract_no,
		shipping_line, bl_no, container_no, port_of_loading, port_of_discharge,
		final_destination, buyer_name, buyer_address, buyer_tel, seller_name,
		payment_term, fob_total, created_at, updated_at)
		VALUES (?, 'export', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, inv.InvNo, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
		inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
		inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
		inv.PaymentTerm, inv.FobTotal, now, now,
	)
	if err != nil {
		// 检查是否是唯一键冲突
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "INV. NO 已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = result

	// 插入明细行
	for i, item := range inv.Items {
		itemNo := i + 1
		db.Exec(`
			INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
			model_no, carton_qty, unit_qty, unit, descriptions, unit_price, total_amount,
			gross_weight, net_weight, volume)
			VALUES (?, 'export', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, inv.InvNo, itemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.UnitQty, item.Unit, item.Descriptions,
			item.UnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume)
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "inv_no": inv.InvNo})
}

// 更新出口表单
func handleExportUpdate(c *gin.Context) {
	userRole := c.GetString("user_role")
	username := c.GetString("username")
	invNo := c.Param("inv_no")

	if userRole != "admin" && userRole != "export_input" && userRole != "export_review" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var inv Invoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	if userRole == "export_review" {
		fmt.Println("走了审核员分支")
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
		c.JSON(http.StatusOK, gin.H{"message": "审核更新成功"})
		return
	}

	if userRole == "export_input" {
		fmt.Println("走了录入员分支")
		// 录入员只能修改表单内容字段，不能碰审核字段
		db.Exec(`
			UPDATE invoices SET invoice_date = ?, contract_no = ?, shipping_line = ?,
			bl_no = ?, container_no = ?, port_of_loading = ?, port_of_discharge = ?,
			final_destination = ?, buyer_name = ?, buyer_address = ?, buyer_tel = ?,
			seller_name = ?, payment_term = ?, fob_total = ?, updated_at = ?
			WHERE inv_no = ? AND invoice_type = 'export'
		`, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
			inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
			inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
			inv.PaymentTerm, inv.FobTotal, now, invNo)
	} else {
		fmt.Println("走了管理员分支")
		reviewComment := ""
		if inv.ReviewComment != nil {
			reviewComment = *inv.ReviewComment
		}
		fmt.Printf("reviewed=%v, review_comment=%s\n", inv.Reviewed, reviewComment)
		// 管理员可以修改所有字段包括审核字段
		var reviewer interface{}
		if inv.Reviewed {
			reviewer = username
		} else {
			reviewer = nil
		}
		result, err := db.Exec(`
    	UPDATE invoices SET invoice_date = ?, contract_no = ?, shipping_line = ?,
   		bl_no = ?, container_no = ?, port_of_loading = ?, port_of_discharge = ?,
    	final_destination = ?, buyer_name = ?, buyer_address = ?, buyer_tel = ?,
    	seller_name = ?, payment_term = ?, fob_total = ?,
    	reviewed = ?, reviewer = ?, review_comment = ?, updated_at = ?
    	WHERE inv_no = ? AND invoice_type = 'export'
		`, inv.InvoiceDate, inv.ContractNo, inv.ShippingLine, inv.BlNo,
			inv.ContainerNo, inv.PortOfLoading, inv.PortOfDischarge, inv.FinalDestination,
			inv.BuyerName, inv.BuyerAddress, inv.BuyerTel, inv.SellerName,
			inv.PaymentTerm, inv.FobTotal,
			inv.Reviewed, reviewer, inv.ReviewComment, now, invNo)
		if err != nil {
			fmt.Println("UPDATE error:", err)
		} else {
			rows, _ := result.RowsAffected()
			fmt.Println("UPDATE rows affected:", rows)
		}
	}

	// 删除旧明细行，重新插入
	db.Exec("DELETE FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'", invNo)
	for i, item := range inv.Items {
		itemNo := i + 1
		db.Exec(`
			INSERT INTO invoice_items (inv_no, invoice_type, item_no, brand, commodities,
			model_no, carton_qty, unit_qty, unit, descriptions, unit_price, total_amount,
			gross_weight, net_weight, volume)
			VALUES (?, 'export', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, invNo, itemNo, item.Brand, item.Commodities, item.ModelNo,
			item.CartonQty, item.UnitQty, item.Unit, item.Descriptions,
			item.UnitPrice, item.TotalAmount, item.GrossWeight, item.NetWeight, item.Volume)
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 删除出口表单
func handleExportDelete(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以删除"})
		return
	}

	invNo := c.Param("inv_no")

	db.Exec("DELETE FROM invoice_items WHERE inv_no = ? AND invoice_type = 'export'", invNo)
	db.Exec("DELETE FROM invoices WHERE inv_no = ? AND invoice_type = 'export'", invNo)

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 检查是否是唯一键冲突错误
func isDuplicateError(err error) bool {
	return err != nil && len(err.Error()) > 0 &&
		(contains(err.Error(), "Duplicate entry") || contains(err.Error(), "UNIQUE constraint failed"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
