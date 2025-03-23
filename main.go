package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"

	kafkamanager "github.com/adityjoshi/Dosahostel/kafka/manager"
	"github.com/adityjoshi/Dosahostel/routes"
	"github.com/gin-gonic/gin"
)

var km *kafkamanager.KafkaManager

func main() {

	db.InitDB()
	fmt.Print("jai shree ram \n")

	router := gin.Default()
	setupRouter(router)
	router.GET("/PING", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PONG",
		})
	})

	northBrokers := []string{"kafka-broker:9092"}
	southBrokers := []string{"kafka-broker:9092"}

	var err error
	km, err = kafkamanager.NewKafkaManager(northBrokers, southBrokers)
	if err != nil {
		log.Fatal("Failed to initialize Kafka Manager:", err)
	}

	regions := []string{"north", "south"}
	for _, region := range regions {
		go func(r string) {
			log.Printf("Starting Kafka consumer for region: %s\n", r)
			kafkamanager.StartConsumer(r)
		}(region)
	}
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
