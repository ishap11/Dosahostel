package controllers

import (
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"
	"github.com/adityjoshi/Dosahostel/models"
	"github.com/adityjoshi/Dosahostel/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func StudentRegistration(c *gin.Context) {
	var student models.Student

	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database, err := db.GetDB(student.HostelName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid region"})
		return
	}
	var existingUser models.Student
	if err := database.Where("email =?", student.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	student.Password = string(hashedPassword)
	student.UserType = "student"

	userStudent := models.Student{
		FullName:       student.FullName,
		Email:          student.Email,
		ContactDetails: student.ContactDetails,
		RegNo:          student.RegNo,
		Room:           student.Room,
		Password:       student.Password,
		HostelName:     student.HostelName,
		UserType:       student.UserType,
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
		Email      string `json:"email"`
		Password   string `json:"password"`
		HostelName string `json:"hostel_name"`
	}
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database, err := db.GetDB(loginRequest.HostelName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid region"})
		return
	}

	var student models.Student
	if err := database.Where("email = ?", loginRequest.Email).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	jwtToken, err := utils.GenerateStudentJWT(int(student.StudentID), student.HostelName, student.RegNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Login successful", "token": jwtToken})
}
