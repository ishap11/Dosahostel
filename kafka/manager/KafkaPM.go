package manager

import (
	"fmt"
	"log"

	"github.com/adityjoshi/Dosahostel/kafka/producer"
)

type KafkaManager struct {
	BigBoysHostel *producer.BigBoysHostel
	BoysHostel    *producer.BoysHostel
}

func NewKafkaManager(BigBoysBroker, BoysBroker []string) (*KafkaManager, error) {
	BigBoysHostel, err := producer.BigBoysProducer(BigBoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing North producer: %w", err)
	}

	BoysHostel, err := producer.BoysProducer(BoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing North producer: %w", err)
	}

	// Return the KafkaManager instance with both producers
	return &KafkaManager{
		BigBoysHostel: BigBoysHostel,
		BoysHostel:    BoysHostel,
	}, nil
}

func (km *KafkaManager) ComplaintRegistration(hostel, topic, messageString string) error {
	var err error

	switch hostel {
	case "BH1", "BH6":
		log.Printf("Sending message to North region, topic: %s", topic)
		err = km.BigBoysHostel.SendMessage(topic, messageString)
	case "BH2", "BH3", "BH4", "BH5":

		log.Printf("Sending message to South region, topic: %s", topic)
		err = km.BoysHostel.SendMessageBoys(topic, messageString)
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
