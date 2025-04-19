package mq

import (
	"diploma/services/customer/pkg/models"
	"diploma/services/customer/pkg/storage"
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
		"group.id":           "my-group-1",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})

	topics := []string{"order", "admin"}
	Consumer.SubscribeTopics(topics, nil)
}

func HandleMessages() {
	for {
		msg, err := Consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Received message: %s\n", string(msg.Key))
			err = ParseMessage(msg)
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
	topic := "customer"

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

func ParseMessage(msg *kafka.Message) error {
	var orderMessage models.OrderMessage
	var syncDatabasesMessage models.SyncDatabasesMessage
	key := string(msg.Key)

	if key != "Sync databases" {
		err := json.Unmarshal(msg.Value, &orderMessage)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(msg.Value, &syncDatabasesMessage)
		if err != nil {
			return err
		}
	}

	switch key {
	case "Sync databases":
		err := storage.SyncDatabases(syncDatabasesMessage)
		if err != nil {
			return err
		}
	case "Waiting for courier", "No couriers", "Order distributed", "Order taken from shop", "Order delivered", "Declined by courier":
		storage.UpdateStatus(orderMessage.OrderId, orderMessage.Status)
	case "Order collected":
		ptrChanged, err := storage.GetChangesAndUpdate(orderMessage.OrderItems, orderMessage.OrderId, int(orderMessage.TotalCost))
		if err != nil {
			return err
		}

		var changed bool
		if ptrChanged != nil {
			changed = *ptrChanged
		}

		if !changed {
			storage.UpdateStatus(orderMessage.OrderId, "order collected")
		} else {
			ptrIsEmpty, err := storage.CheckOrderIsEmpty(orderMessage.OrderId)
			if err != nil {
				return err
			}

			var isEmpty bool
			if ptrIsEmpty != nil {
				isEmpty = *ptrIsEmpty
			}

			if !isEmpty {
				storage.UpdateStatus(orderMessage.OrderId, "order collected with some changes")
			} else {
				storage.UpdateStatus(orderMessage.OrderId, "declined because no products left")
			}
		}
	}

	return nil
}
