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

func PostInventory(c *gin.Context) {
	// Extract token from header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Decode token to extract admin details
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

	adminID, ok := userClaims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
		return
	}
	region, ok := userClaims["region"].(string) // Extracting region from JWT
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid region in token"})
		return
	}

	database, err := db.GetDB(region)
	if err != nil || database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var admin models.Users
	if err := database.Where("id = ?", uint(adminID)).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Bind JSON request body
	var inventoryReq struct {
		ProductName string `json:"product_name"`
		GSTNumber   string `json:"gst_number"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&inventoryReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create new inventory record
	inventory := models.Inventory{
		AdminID:      uint(adminID),
		BusinessName: admin.BusinessName,
		GSTNumber:    inventoryReq.GSTNumber,
		ProductName:  inventoryReq.ProductName,
		Quantity:     inventoryReq.Quantity,
	}

	if err := database.Create(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory item registered successfully",
	})
}

/*
 */
func PostComplaintKafka(c *gin.Context) {
	// Check if KafkaManager is available in the context
	km, exists := c.Get("km")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KafkaManager not found"})
		return
	}

	// Check if the type of KafkaManager is correct
	kafkaManager, ok := km.(*kafkamanager.KafkaManager)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid KafkaManager"})
		return
	}

	// Check if Authorization token is present
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

	// Check if the token contains valid user claims
	userClaims, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token structure"})
		return
	}

	// Check if user_id is present and valid in the token
	adminID, ok := userClaims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
		return
	}

	// Check if region is available in the token
	region, ok := userClaims["region"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid region in token"})
		return
	}

	// Struct for inventory request validation
	var inventoryReq struct {
		ProductName string `json:"product_name"`
		GSTNumber   string `json:"gst_number"`
		Quantity    int    `json:"quantity"`
	}

	// Bind and validate the incoming JSON body
	if err := c.ShouldBindJSON(&inventoryReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Fetch the database for the given region
	database, err := db.GetDB(region)
	if err != nil || database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error or invalid hostel"})
		return
	}

	// Fetch the admin information from the database
	var admin models.Users
	if err := database.Where("id = ?", uint(adminID)).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Prepare the inventory data
	Inventory := models.Inventory{
		AdminID:      uint(adminID),
		BusinessName: admin.BusinessName,
		GSTNumber:    inventoryReq.GSTNumber,
		ProductName:  inventoryReq.ProductName,
		Quantity:     inventoryReq.Quantity,
	}

	// Marshal the inventory data into JSON
	complaints, err := json.Marshal(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal complaint admin data to JSON"})
		return
	}

	// Declare a variable to capture Kafka errors
	var errKafka error

	// Send the complaint data to Kafka based on region
	switch region {
	case "north":
		errKafka = kafkaManager.ComplaintRegistration(region, "inventory", string(complaints))
	case "south":
		errKafka = kafkaManager.ComplaintRegistration(region, "inventory", string(complaints))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid region: %s", region)})
		return
	}

	// Check if sending data to Kafka failed
	if errKafka != nil {
		log.Printf("Failed to send hospital registration data to Kafka: %v", errKafka)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Kafka"})
		return
	}

	// Send a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Complaint created successfully", "region": Inventory.BusinessName})
}
