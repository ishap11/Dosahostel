package manager

// consumer/consumerFile.go

import (
	"fmt"
	"log"

	"github.com/adityjoshi/Dosahostel/kafka/consumer"
)

func StartConsumer(region string) {
	kafkaBroker := "kafka-broker:9092"
	var topic = []string{
		"inventory",

		// Add other topics as necessary
	}
	switch region {
	case "north":
		//topic := "hospital_admin"
		northConsumer, err := consumer.NewNorthConsumer(kafkaBroker, topic)
		if err != nil {
			log.Fatalf("Failed to create north consumer: %v", err)
		}
		northConsumer.Listen()

	case "south":
		// Similar setup for the south region consumer
		fmt.Println("Starting south consumer...")

	case "east":
		// Similar setup for the east region consumer
		fmt.Println("Starting east consumer...")

	case "west":
		// Similar setup for the west region consumer
		fmt.Println("Starting west consumer...")

	default:
		fmt.Println("Unknown region:", region)
	}
}
