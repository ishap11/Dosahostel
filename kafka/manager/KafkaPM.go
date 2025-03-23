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

func NewKafkaManager(BigBoysBroker, BoysBroker []string) (*KafkaManager, error) {
	NorthProducer, err := producer.NewNorthProducer(BigBoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing North producer: %w", err)
	}

	SouthProducer, err := producer.NewSouthProducer(BoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing North producer: %w", err)
	}

	// Return the KafkaManager instance with both producers
	return &KafkaManager{
		NorthProducer: NorthProducer,
		SouthProducer: SouthProducer,
	}, nil
}

func (km *KafkaManager) ComplaintRegistration(hostel, topic, messageString string) error {
	var err error

	switch hostel {
	case "north":
		log.Printf("Sending message to North region, topic: %s", topic)
		err = km.NorthProducer.SendMessage(topic, messageString)
	case "south":

		log.Printf("Sending message to South region, topic: %s", topic)
		err = km.SouthProducer.SouthMessage(topic, messageString)
	default:

		return fmt.Errorf("invalid hostel: %s", hostel)
	}
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s in %s region: %w", topic, hostel, err)
	}

	// Log successful message sending
	log.Printf("Message successfully sent to topic %s in %s region", topic, hostel)
	return nil
}
