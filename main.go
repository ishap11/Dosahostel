package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"

	"github.com/adityjoshi/Dosahostel/kafka/consumer"
	kafkamanager "github.com/adityjoshi/Dosahostel/kafka/manager"
	"github.com/adityjoshi/Dosahostel/routes"
	"github.com/gin-gonic/gin"
)

var km *kafkamanager.KafkaManager

func main() {
	// Initialize database
	db.InitDB()
	fmt.Print("jai shree ram \n")

	// Initialize router
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("KafkaManager", km) // Setting KafkaManager into the context
		c.Next()
	})

	// Set up PING endpoint
	router.GET("/PING", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PONG",
		})
	})

	// Kafka setup
	northBrokers := []string{"kafka:9092"}
	southBrokers := []string{"kafka:9092"}

	var err error
	km, err = kafkamanager.NewKafkaManager(northBrokers, southBrokers)
	if err != nil {
		log.Fatal("Failed to initialize Kafka Manager:", err)
	}

	log.Printf("KafkaManager initialized successfully: %v", km)

	if km == nil {
		log.Fatal("KafkaManager is not initialized")
	}
	setupRouter(router)
	// Start Kafka consumers
	regions := []string{"north", "south"}
	for _, region := range regions {
		go func(r string) {
			log.Printf("Starting Kafka consumer for region: %s\n", r)
			consumer.StartConsumer(r)
		}(region)
	}

	// Run server
	server := &http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	log.Println("Server is running at :8001...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupRouter(router *gin.Engine) {
	routes.StudentRoutes(router, km)
}
