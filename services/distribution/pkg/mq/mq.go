package mq

import (
	"diploma/services/distribution/pkg/storage"
	"time"
	// "diploma/services/customer/pkg/models"
	// "encoding/json"
	"diploma/services/distribution/pkg/models"
	"encoding/json"
	"log"

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
		"group.id":           "my-group-4",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})

	topics := []string{"courier", "order"}
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
	topic := "distribution"

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
	var courierMessage models.CourierState

	key := string(msg.Key)

	if key == "Courier state" {
		err := json.Unmarshal(msg.Value, &courierMessage)
		if err != nil {
			return err
		}

		storage.AddCourier(courierMessage.CourierId)
		storage.SetState(courierMessage.CourierId, courierMessage.State)
	} else {
		err := json.Unmarshal(msg.Value, &orderMessage)
		if err != nil {
			return err
		}

		switch key {
		case "Waiting for courier":
			time.Sleep(5 * time.Second)
			courierStates, err := storage.GetActiveCouriers()
			if err != nil {
				return err
			}
	
			if len(courierStates) == 0 {
				orderMessage.Status = "canceled because no couriers"
				ProduceMessage(orderMessage, "No couriers")
			} else {
				orderMessage.Status = "order sent to couriers"
				ProduceMessage(orderMessage, "Order sent to couriers")
			}
		case "Order taken":
			ProduceMessage(orderMessage, "Order distributed")
		}
	}

	return nil
}