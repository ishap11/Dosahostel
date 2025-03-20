package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/adityjoshi/Dosahostel/database"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()
	fmt.Print("jai shree ram \n")

	router := gin.Default()
	router.GET("/PING", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PONG",
		})
	})

	server := &http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	log.Println("Server is running at :8001...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
