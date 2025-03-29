package mq

import (
	// "diploma/services/courier/pkg/storage"
	"diploma/services/courier/pkg/models"
	"diploma/services/courier/pkg/storage"
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
		"group.id":           "my-group-3",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})

	topics := []string{"distribution", "order"}
	Consumer.SubscribeTopics(topics, nil)
}

func HandleMessages() {
	for {
		msg, err := Consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Received message: %s\n", string(msg.Key))
			err = ParseMessageAndProduce(msg)
			if err != nil {
				log.Printf("Failed to parse message: %s\n", string(msg.Value))
			}
		} else {
			log.Println("Consumer error:", err)
			break
		}
	}
}

func ProduceMessage(msg models.OrderMessage, key string) error {
	topic := "courier"

	log.Print(msg)

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

func ProduceState(msg models.Courier, key string) error {
	topic := "courier"

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
	case "Order sent to couriers":
		err = storage.InsertOrder(orderMessage)
		if err != nil {
			return err
		}
	case "Order collected":
		time.Sleep(5 * time.Second)
		err = storage.UpdateOrderStatus(orderMessage.OrderId, "order collected")
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}