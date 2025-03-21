package manager

import (
	"fmt"

	"github.com/adityjoshi/Dosahostel/kafka/producer"
)

type KafkaManager struct {
	BigBoysHostel *producer.BigBoysHostel
}

func NewKafkaManager(BigBoysBroker, BoysBroker []string) (*KafkaManager, error) {
	BigBoysHostel, err := producer.BigBoysProducer(BigBoysBroker)
	if err != nil {
		return nil, fmt.Errorf("error initializing North producer: %w", err)
	}

	// Return the KafkaManager instance with both producers
	return &KafkaManager{
		BigBoysHostel: BigBoysHostel,
	}, nil
}
