package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.POST("/api/login", handleLogin)

	api := r.Group("/api", authMiddleware())
	{
		api.GET("/me", handleMe)

		// 出口表单
		api.GET("/invoices/export", handleExportList)
		api.GET("/invoices/export/:inv_no", handleExportDetail)
		api.POST("/invoices/export", handleExportCreate)
		api.PUT("/invoices/export/:inv_no", handleExportUpdate)
		api.DELETE("/invoices/export/:inv_no", handleExportDelete)
	}
}
func handleMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id":           c.GetInt("user_id"),
		"username":          c.GetString("username"),
		"user_role":         c.GetString("user_role"),
		"can_edit_products": c.GetBool("can_edit_products"),
	})
}

func handleExportList(c *gin.Context) {
	userRole := c.GetString("user_role")
	invNo := c.Query("inv_no")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if userRole != "admin" && userRole != "export_input" && userRole != "export_review" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	query := `SELECT id, inv_no, invoice_date, contract_no, 
			  shipping_line, reviewed, reviewer, created_at 
			  FROM invoices WHERE invoice_type = 'export'`
	args := []interface{}{}

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

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type InvoiceSummary struct {
		ID           int     `json:"id"`
		InvNo        string  `json:"inv_no"`
		InvoiceDate  *string `json:"invoice_date"`
		ContractNo   *string `json:"contract_no"`
		ShippingLine *string `json:"shipping_line"`
		Reviewed     bool    `json:"reviewed"`
		Reviewer     *string `json:"reviewer"`
		CreatedAt    string  `json:"created_at"`
	}

	var result []InvoiceSummary
	for rows.Next() {
		var inv InvoiceSummary
		rows.Scan(&inv.ID, &inv.InvNo, &inv.InvoiceDate, &inv.ContractNo,
			&inv.ShippingLine, &inv.Reviewed, &inv.Reviewer, &inv.CreatedAt)
		result = append(result, inv)
	}
	if result == nil {
		result = []InvoiceSummary{}
	}
	c.JSON(http.StatusOK, result)
}
