package main

import (
	"context"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

var kafkaClient *kgo.Client

func initKafka() {
	var err error
	seeds := []string{"k3:9092"}
	kafkaClient, err = kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %s", err)
	} else {
		log.Println("Connected to Kafka successfully.")
	}

	// Check if Kafka is reachable by listing the existing topics
	admin := kadm.NewClient(kafkaClient)
	defer admin.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List all topics to ensure connectivity
	topics, err := admin.ListTopics(ctx)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	} else {
		log.Printf("Kafka topics: %v", topics.Names())
	}

	// Check if the topic exists
	topicName := "email-topic"
	if _, exists := topics[topicName]; !exists {
		log.Printf("Topic %s does not exist. Creating it now.", topicName)
		err := createTopicIfNotExists(admin, topicName, 1, 1)
		if err != nil {
			log.Fatalf("Failed to create topic %s: %v", topicName, err)
		}
	} else {
		log.Printf("Topic %s already exists.", topicName)
	}
}

func createTopicIfNotExists(admin *kadm.Client, topic string, partitions int32, replicationFactor int16) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pass empty config map and an empty list of configs for simplicity
	_, err := admin.CreateTopic(ctx, topic, partitions, replicationFactor, nil)
	if err != nil {
		log.Printf("Error creating topic %s: %v\n", topic, err)
		return err
	}

	log.Printf("Topic %s created successfully.\n", topic)
	return nil
}
