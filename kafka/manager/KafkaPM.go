package manager

import (
	"fmt"
	"log"

	"github.com/adityjoshi/Dosahostel/kafka/producer"
)

type KafkaManager struct {
	NorthProducer *producer.NorthProducer
	SouthProducer *producer.SouthProducer
}

// NewKafkaManager initializes the Kafka producers and returns a KafkaManager instance.
func NewKafkaManager(BigBoysBroker, BoysBroker []string) (*KafkaManager, error) {
	// Initialize the NorthProducer
	NorthProducer, err := producer.NewNorthProducer(BigBoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing North producer: %w", err)
	}

	// Initialize the SouthProducer
	SouthProducer, err := producer.NewSouthProducer(BoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing South producer: %w", err) // Fixed error message to reflect SouthProducer
	}

	// Return the KafkaManager instance with both producers
	return &KafkaManager{
		NorthProducer: NorthProducer,
		SouthProducer: SouthProducer,
	}, nil
}

// ComplaintRegistration sends a message to the appropriate Kafka topic based on the hostel region.
func (km *KafkaManager) ComplaintRegistration(hostel, topic, messageString string) error {
	var err error

	// Ensure that producers are initialized
	if km.NorthProducer == nil || km.SouthProducer == nil {
		return fmt.Errorf("Kafka producers are not properly initialized")
	}

	// Send message based on the hostel region
	switch hostel {
	case "north":
		if km.NorthProducer == nil {
			return fmt.Errorf("North producer is not initialized")
		}
		log.Printf("Sending message to North region, topic: %s", topic)
		err = km.NorthProducer.SendMessage(topic, messageString)

	case "south":
		if km.SouthProducer == nil {
			return fmt.Errorf("South producer is not initialized")
		}
		log.Printf("Sending message to South region, topic: %s", topic)
		err = km.SouthProducer.SouthMessage(topic, messageString)

	default:
		return fmt.Errorf("invalid hostel: %s", hostel)
	}

	// Check for errors in message sending
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s in %s region: %w", topic, hostel, err)
	}

	// Log successful message sending
	log.Printf("Message successfully sent to topic %s in %s region", topic, hostel)
	return nil
}
