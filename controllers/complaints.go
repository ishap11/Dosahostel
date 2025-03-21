package controllers

import (
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"
	"github.com/adityjoshi/Dosahostel/models"
	"github.com/adityjoshi/Dosahostel/utils"
	"github.com/gin-gonic/gin"
)

func PostComplaint(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Decode token to extract student details
	claims, err := utils.DecodeStudentJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userClaims, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token structure"})
		return
	}

	studentID, ok := userClaims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
		return
	}

	regNo, ok := userClaims["reg_no"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid reg_no in token"})
		return
	}

	blockID, ok := userClaims["block_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid block_id in token"})
		return
	}

	var complaintReq struct {
		Complaint   string `json:"complaint"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&complaintReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	database, err := db.GetDB(blockID)
	if err != nil || database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error or invalid hostel"})
		return
	}
	var student models.Student
	if err := database.Where("reg_no = ?", regNo).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	complaint := models.Complaint{
		RegNo:          regNo,
		ComplaintType:  models.ComplaintType(complaintReq.Complaint),
		Description:    complaintReq.Description,
		StudentID:      uint(studentID),
		HostelName:     blockID,
		Room:           student.Room,
		ContactDetails: student.ContactDetails,
	}

	if err := database.Create(&complaint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User complaint registered successfully",
	})
}
