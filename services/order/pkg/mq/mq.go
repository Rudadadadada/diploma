package mq

import (
	"diploma/services/order/pkg/models"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

var Producer *kafka.Producer
var Consumer *kafka.Consumer

func New() {
	Producer, _ = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":      "localhost:9092",
		"queue.buffering.max.ms": "1000",
	})

	Consumer, _ = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "my-group-2",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})

	topics := []string{"customer", "distribution", "admin", "courier"}
	Consumer.SubscribeTopics(topics, nil)
}

func HandleMessages() {
	for {
		msg, err := Consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Received message: %s\n", string(msg.Key))
			err = ParseMessageAndProduce(msg)
			if err != nil {
				log.Printf("Failed to parse and produce on message: %s\n", string(msg.Value))
			}
		} else {
			log.Println("Consumer error:", err)
			break
		}
	}
}

func ProduceMessage(msg models.OrderMessage, key string) error {
	topic := "order"

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonMsg,
		Key:            []byte(key),
	}, nil)

	if err != nil {
		return err
	}

	return nil
}

func ParseMessageAndProduce(msg *kafka.Message) error {
	var orderMessage models.OrderMessage
	key := string(msg.Key)

	err := json.Unmarshal(msg.Value, &orderMessage)
	if err != nil {
		return err
	}

	switch key {
	case "Made order":
		orderMessage.Status = "waiting for courier"
		ProduceMessage(orderMessage, "Waiting for courier")
	case "No couriers":
		orderMessage.Status = "canceled because no couriers"
		ProduceMessage(orderMessage, "No couriers")
	case "Order distributed":
		orderMessage.Status = "preparing"
		ProduceMessage(orderMessage, "Order distributed")
	case "Order collected":
		ProduceMessage(orderMessage, "Order collected")
	case "Order taken from shop":
		orderMessage.Status = "order taken from shop"
		ProduceMessage(orderMessage, "Order taken from shop")
	case "Delivered":
		orderMessage.Status = "order delivered"
		ProduceMessage(orderMessage, "Order delivered")
	case "Order declined":
		orderMessage.Status = "order declined"
		ProduceMessage(orderMessage, "Order declined")
	case "Declined by courier":
		orderMessage.Status = "declined by courier"
		ProduceMessage(orderMessage, "Declined by courier")
	}

	return nil
}
