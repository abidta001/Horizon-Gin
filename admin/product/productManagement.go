package product

import (
	"net/http"

	db "admin/DB"
	"admin/models"

	"github.com/gin-gonic/gin"
)

func ViewProducts(c *gin.Context) {
	var products []models.Product
	result := db.Db.Order("product_id ASC").Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for i := range products {
		if products[i].Quantity == 0 {
			products[i].Status = 2 // Out of stock
			db.Db.Model(&products[i]).Update("status", products[i].Status)
		} else {
			products[i].Status = 1 // Available
		}
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No products listed"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func AddProducts(c *gin.Context) {
	var products models.Product

	if err := c.ShouldBind(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if products.Price < 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Price cannot be a negative value"})
		return
	}

	// Determine initial status
	products.Status = 1 // Assume Available by default
	if products.Quantity == 0 {
		products.Status = 2 // Out of stock
	}

	if err := db.Db.Create(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var product models.Product

	if err := db.Db.Where("deleted_at IS NULL").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		ProductName string  `json:"product_name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImgURL      string  `json:"img_url"`
		Status      int     `json:"status"` // Changed to int
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := models.Product{
		ProductName: input.ProductName,
		Description: input.Description,
		Price:       input.Price,
		ImgURL:      input.ImgURL,
		Status:      input.Status,
	}

	if err := db.Db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product_name": product.ProductName})
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product

	if err := db.Db.Where("product_id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := db.Db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func UpdateProductStock(c *gin.Context) {
	var product models.Product
	productID := c.Param("id")

	if err := db.Db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Quantity = input.Quantity

	db.Db.Save(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Stock updated"})
}
