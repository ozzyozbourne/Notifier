package main

import (
	"github.com/IBM/sarama"
	"log"
)

func main() {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll

	brokersAddress := []string{"localhost:9192"}

	//create a sync producer
	producer, err := sarama.NewSyncProducer(brokersAddress, config)
	if err != nil {
		log.Fatalf("Failed to create producer -> \n%v\n", err)
	}
	defer producer.Close()

	//create a admin client
	admin, err := sarama.NewClusterAdmin(brokersAddress, config)
	if err != nil {
		log.Fatalf("Failed to create admin client -> \n%v\n", err)
	}
	defer admin.Close()
	//create topic
	topic := "send-email"
	printTopicList(admin)
	createTopic(admin, topic)
	printTopicList(admin)
	sendMessage(producer, topic, "test_key", "A test message from golang")

}

func printTopicList(admin sarama.ClusterAdmin) {
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Unable to fetch list of topics \n%v\n", err)
	}
	if topics != nil && len(topics) == 0 {
		log.Printf("No available topic\n")
	} else {
		log.Printf("List of fetched topics ->\n%v\n", topics)
	}
}

func createTopic(admin sarama.ClusterAdmin, topicName string) {
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Unable to fetch list of topics \n%v\n", err)
	}
	if _, exists := topics[topicName]; !exists {
		topicDetail := sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
		err = admin.CreateTopic(topicName, &topicDetail, false)
		if err != nil {
			log.Fatalf("Failed to create topic %s\n%v\v", topicName, err)
		}
		log.Printf("Topic %s created successfully\n", topicName)
	} else {
		log.Printf("Topic %s already exists\n", topicName)
	}
}

func sendMessage(producer sarama.SyncProducer, topic, key, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message -> %s to topic %s\n", message, topic)
	} else {
		log.Printf("Message %s send to topic %s successfully\n", message, topic)
		log.Printf("Partition value -> %v\n", partition)
		log.Printf("Offset value -> %v\n", offset)
	}
}
