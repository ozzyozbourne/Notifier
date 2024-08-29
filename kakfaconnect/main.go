package main

import (
	"context"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"time"
)

func main() {
	brokersAddress := []string{"localhost:9192"}

	// Kafka client configuration
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokersAddress...),
		kgo.WithLogger(kgo.BasicLogger(log.Default().Writer(), kgo.LogLevelInfo, nil)),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client -> \n%v\n", err)
	}
	defer client.Close()

	// Create a Kafka admin client
	adminClient := kadm.NewClient(client)
	defer adminClient.Close()

	// Create a topic
	topic := "send-email"
	printTopicList(adminClient)
	createTopic(adminClient, topic)
	printTopicList(adminClient)
	sendMessage(client, topic, "test_key", "A test message from golang")
}

func printTopicList(adminClient *kadm.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topics, err := adminClient.ListTopics(ctx)
	if err != nil {
		log.Fatalf("Unable to fetch list of topics \n%v\n", err)
	}
	if len(topics) == 0 {
		log.Printf("No available topic\n")
	} else {
		log.Printf("List of fetched topics ->\n")
		for topic := range topics {
			log.Printf("%s\n", topic)
		}
	}
}

func createTopic(adminClient *kadm.Client, topicName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topics, err := adminClient.ListTopics(ctx)
	if err != nil {
		log.Fatalf("Unable to fetch list of topics \n%v\n", err)
	}
	if _, exists := topics[topicName]; !exists {
		_, err = adminClient.CreateTopics(ctx, 1, 1, nil, topicName)
		if err != nil {
			log.Fatalf("Failed to create topic %s\n%v\v", topicName, err)
		}
		log.Printf("Topic %s created successfully\n", topicName)
	} else {
		log.Printf("Topic %s already exists\n", topicName)
	}
}

func sendMessage(client *kgo.Client, topic, key, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(message),
	}
	err := client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		log.Fatalf("Failed to send message -> %s to topic %s\n", message, topic)
	} else {
		log.Printf("Message %s sent to topic %s successfully\n", message, topic)
	}
}
