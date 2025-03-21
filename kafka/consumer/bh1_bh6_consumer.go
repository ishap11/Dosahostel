package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type BigBoysConsumer struct {
	Consumer sarama.Consumer
	Topics   []string
}

func NewBigBoysConsumer(broker string, topics []string) (*BigBoysConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create the consumer
	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Printf("Error creating consumer: %v", err)
		return nil, fmt.Errorf("error creating consumer: %v", err)
	}

	log.Printf("Kafka consumer created successfully, subscribing to topics: %v", topics)

	// Return a NorthConsumer instance with the list of topics
	return &BigBoysConsumer{Consumer: consumer, Topics: topics}, nil
}

func (nc *BigBoysConsumer) Listen() {
	defer func() {
		if err := nc.Consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v\n", err)
		}
	}()

	// Create a consumer for each topic, listening to partition 0 for each one
	for _, topic := range nc.Topics {
		partitionConsumer, err := nc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to start consumer for partition 0 of topic %s: %v", topic, err)
		}
		defer partitionConsumer.Close()

		// Log that the consumer has started listening
		log.Printf("Consumer is now listening to topic: %s from partition 0\n", topic)

		// Start a goroutine for each topic to consume messages concurrently
		go nc.consumeMessages(partitionConsumer)
	}

	// Block indefinitely (since we're listening to multiple topics concurrently)
	select {}
}

func (nc *BigBoysConsumer) consumeMessages(partitionConsumer sarama.PartitionConsumer) {
	for msg := range partitionConsumer.Messages() {
		// Log received message and metadata
		log.Printf("Received message from topic %s: %s\n", msg.Topic, string(msg.Value))
		log.Printf("Message metadata - Partition: %d, Offset: %d\n", msg.Partition, msg.Offset)

		// Process the message (e.g., save data to the database or trigger further actions)
		if err := processMessage(msg.Topic, msg); err != nil {
			log.Printf("Error processing message: %v", err)
		} else {
			log.Printf("Message processed successfully")
		}
	}
}

func processMessage(topic string, msg *sarama.ConsumerMessage) error {
	log.Printf("Processing message: %s \n", string(msg.Value))
	switch topic {
	case "hospital_admin":
		log.Printf("Processing hospital_admin message: %s", string(msg.Value))

		var admin database.HospitalAdmin
		if err := json.Unmarshal(msg.Value, &admin); err != nil {
			log.Printf("Failed to unmarshal hospital_admin message: %v", err)
			return err
		}
		if err := database.NorthDB.Create(&admin); err != nil {
			log.Printf("Failed to save hospital_admin data: %v", err.Error)
			return fmt.Errorf("Failed to write to the DB", err.Error, err)
		}
	case "hospital_registration":
		log.Printf("Processing hospital_registration message: %s", string(msg.Value))

		var hospital database.Hospitals
		if err := json.Unmarshal(msg.Value, &hospital); err != nil {
			log.Printf("Error unmarshalling hospital data: %v", err)
			return err

		}
		hospital.Username = fmt.Sprintf("DEL%d", hospital.HospitalId)

		if err := database.NorthDB.Create(&hospital).Error; err != nil {
			log.Printf("Error creating hospital in database: %v", err)
			return fmt.Errorf(err.Error(), err)
		}
	default:
		// Handle any other topics or log an error if the topic is not recognized
		log.Printf("Received message from unknown topic: %s", topic)
		// Add your default logic here
	}
	return nil
}
