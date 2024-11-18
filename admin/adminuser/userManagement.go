package adminuser

import (
	"net/http"

	db "admin/DB"
	"admin/models"
	"admin/models/responsemodels"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	var users []responsemodels.User

	if err := db.Db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Users": users})
}

func BlockUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := db.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Status == "Blocked" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User aldready blocked"})
		return
	}
	user.Status = "Blocked"
	if err := db.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

func UnblockUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	if err := db.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Status == "Available" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User aldready active"})
		return
	}
	user.Status = "Available"
	if err := db.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Unblocked successfully"})
}
