package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"
	kafkamanager "github.com/adityjoshi/Dosahostel/kafka/manager"
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

func PostComplaintKafka(c *gin.Context) {
	km, exists := c.Get("km")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KafkaManager not found"})
		return
	}

	kafkaManager, ok := km.(*kafkamanager.KafkaManager)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid KafkaManager"})
		return
	}

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

	complaints, err := json.Marshal(complaint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal complaint admin data to JSON"})
		return
	}

	var errKafka error

	switch blockID {
	case "BH1", "BH6":
		errKafka = kafkaManager.ComplaintRegistration(blockID, "complaint_registration", string(complaints))
	case "BH2", "BH3", "BH4", "BH5":
		errKafka = kafkaManager.ComplaintRegistration(blockID, "complaint_registration", string(complaints))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid region: %s", blockID)})
		return
	}

	if errKafka != nil {
		log.Printf("Failed to send hospital registration data to Kafka: %v", errKafka)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Kafka"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Complaint created successfully", "region": complaint.HostelName})
}
