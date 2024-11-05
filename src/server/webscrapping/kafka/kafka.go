package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Global variable for Kafka producer
var producer sarama.SyncProducer
var admin sarama.ClusterAdmin

// InitKafkaProducer initializes a Kafka producer
func InitKafkaProducer(brokers []string) error {
	config := sarama.NewConfig()
	// Require acknowledgment from all brokers for message to be considered successful
	config.Producer.RequiredAcks = sarama.WaitForAll
	// Set max retries for sending messages
	config.Producer.Retry.Max = 5
	// Ensure successes are returned after message is sent
	config.Producer.Return.Successes = true
	// Set max message size to 25MB
	config.Producer.MaxMessageBytes = 104857600

	var err error
	// Initialize the Kafka producer
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}

	// Initialize admin to manage topics
	admin, err = sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return err
	}

	// Create the topic "scrapped-data"
	err = CreateKafkaTopic("scrapped-data", 1, 1)
	if err != nil {
		log.Println("Error creating topic:", err)
		return err
	}

	return nil
}

// CreateKafkaTopic creates a new topic in Kafka
func CreateKafkaTopic(topic string, numPartitions int32, replicationFactor int16) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	// Check if the topic already exists
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	if _, exists := topics[topic]; exists {
		log.Printf("Topic %s already exists.\n", topic)
		return nil
	}

	// Create the topic
	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		return err
	}

	log.Printf("Topic %s created successfully.\n", topic)
	return nil
}

// SendMessageToKafka sends a simple message to Kafka
func SendMessageToKafka(topic, message string) {
	if producer == nil {
		log.Println("Error: Kafka producer is not initialized.")
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("Error sending message to Kafka:", err)
		return
	}

	log.Printf("Message sent to Kafka. Partition: %d, Offset: %d\n", partition, offset)
}

// SendEncodedFileToKafka sends a base64-encoded file to Kafka
func SendEncodedFileToKafka(topic, filePath string, encodedData string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encodedData),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("Error sending file to Kafka:", err)
		return
	}

	log.Printf("Encoded file sent to Kafka. Partition: %d, Offset: %d\n", partition, offset)
}

// CloseKafkaProducer closes the Kafka producer connection
func CloseKafkaProducer() {
	if producer != nil {
		producer.Close()
	}

	if admin != nil {
		admin.Close()
	}
}
