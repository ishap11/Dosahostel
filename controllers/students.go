package controllers

import (
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"
	"github.com/adityjoshi/Dosahostel/models"
	"github.com/adityjoshi/Dosahostel/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func BusinessAdminReg(c *gin.Context) {
	var admin models.Users

	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database, err := db.GetDB(admin.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid region"})
		return
	}
	var existingUser models.Users
	if err := database.Where("email =?", admin.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	admin.Password = string(hashedPassword)
	admin.User_type = "Admin"

	userStudent := models.Users{
		Full_Name:     admin.Full_Name,
		GenderInfo:    models.Gender(admin.GenderInfo),
		ContactNumber: admin.ContactNumber,
		Email:         admin.Email,
		Password:      admin.Password,
		BusinessName:  admin.BusinessName,
		Region:        admin.Region,
		User_type:     models.Buyer,
	}
	if err := database.Create(&userStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}

func StudentLogin(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Region   string `json:"region"`
	}
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database, err := db.GetDB(loginRequest.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid region"})
		return
	}

	var admin models.Users
	if err := database.Where("email = ?", loginRequest.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	jwtToken, err := utils.GenerateStudentJWT(int(admin.ID), admin.Region, string(admin.User_type))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	utils.GenerateAndSendOTP(admin.Email)

	c.JSON(http.StatusOK, gin.H{"Status": "Login successful", "token": jwtToken})
}
