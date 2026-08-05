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

		// 进口表单
		api.GET("/invoices/import", handleImportList)
		api.GET("/invoices/import/:inv_no", handleImportDetail)
		api.POST("/invoices/import", handleImportCreate)
		api.PUT("/invoices/import/:inv_no", handleImportUpdate)
		api.DELETE("/invoices/import/:inv_no", handleImportDelete)

		// 产品属性
		api.GET("/products", handleProductList)
		api.GET("/products/:model_no", handleProductDetail)
		api.POST("/products", handleProductCreate)
		api.PUT("/products/:model_no", handleProductUpdate)
		api.DELETE("/products/:model_no", handleProductDelete)
		api.POST("/products/copy-prices", handleProductCopy)
		api.GET("/regions", handleRegionList)

		// 用户管理（仅管理员）
		api.GET("/users", handleUserList)
		api.POST("/users", handleUserCreate)
		api.PUT("/users/:id", handleUserUpdate)
		api.DELETE("/users/:id", handleUserDelete)
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
