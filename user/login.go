package user

import (
	"net/http"

	db "admin/DB"
	"admin/middleware"
	"admin/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input models.LoginInput
	var user models.User

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Db.Where("email=?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	if user.Status == "Blocked" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User has been blocked by the Admin"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}
	token, err := middleware.CreateToken("user", user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Error Generating jwt"})
	}
	c.Header("Authorization", "Bearer"+token)
	c.JSON(http.StatusOK, gin.H{"message": "Login successfull", "token": token})
}
