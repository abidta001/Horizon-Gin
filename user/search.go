package user

import (
	"net/http"

	db "admin/DB"
	"admin/models"
	"github.com/gin-gonic/gin"
)

func SearchProducts(c *gin.Context) {
	query := c.Query("query")
	sort := c.Query("sort")
	order := c.Query("order")
	categoryID := c.Query("categoryID")

	var products []models.Product
	db := db.Db.Model(&models.Product{}).Where("product_name ILIKE ?", "%"+query+"%")

	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	switch sort {
	case "popularity":
		db = db.Order("popularity " + order)
	case "price":
		db = db.Order("price " + order)
	case "new_arrivals":
		db = db.Order("created_at " + order)
	case "featured":
		db = db.Order("featured " + order)
	case "name":
		if order == "asc" {
			db = db.Order("product_name ASC")
		} else {
			db = db.Order("product_name DESC")
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort parameter"})
		return
	}

	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
