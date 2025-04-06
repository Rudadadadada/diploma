package mq

import (
	"diploma/services/admin/pkg/models"
	"diploma/services/admin/pkg/storage"
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
		"group.id":           "my-group-5",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})

	topics := []string{"order"}
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
	topic := "admin"

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

func ProduceSyncMessage(msg models.SyncDatabasesMessage) error {
	topic := "admin"

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonMsg,
		Key:            []byte("Sync databases"),
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
	case "Order distributed":
		actualProductsState, err := storage.GetActualState(orderMessage.OrderItems)
		if err != nil {
			return err
		}

		var toUpdate []models.Product
		var newTotalCost float32

		for i := 0; i < len(actualProductsState); i++ {
			id := actualProductsState[i].Id

			amount := int(orderMessage.OrderItems[i].Amount)
			actualAmount := int(actualProductsState[i].Amount)
			cost := actualProductsState[i].Cost

			var tmp models.Product
			tmp.Id = id

			if actualAmount-amount >= 0 {
				tmp.Amount = uint(actualAmount - amount)
				newTotalCost += float32(amount) * cost
			} else {
				tmp.Amount = 0
				newTotalCost += float32(actualAmount) * cost
				
				orderMessage.OrderItems[i].Amount = uint(actualAmount)
				orderMessage.OrderItems[i].TotalCost = float32(actualAmount) * cost
			}

			toUpdate = append(toUpdate, tmp)
		}

		orderMessage.TotalCost = newTotalCost

		err = storage.UpadteProducts(toUpdate)
		if err != nil {
			return err
		}

		syncDatabasesMessage, err := storage.SyncDatabases()
		if err != nil {
			return err
		}

		err = ProduceSyncMessage(*syncDatabasesMessage)
		if err != nil {
			return err
		}

		ProduceMessage(orderMessage, "Order collected")
	case "Order declined", "Declined by courier":
		actualProductsState, err := storage.GetActualState(orderMessage.OrderItems)
		if err != nil {
			return err
		}

		var toUpdate []models.Product
		for i := 0; i < len(actualProductsState); i++ {
			id := actualProductsState[i].Id

			amount := int(orderMessage.OrderItems[i].Amount)
			actualAmount := int(actualProductsState[i].Amount)

			var tmp models.Product
			tmp.Id = id

			tmp.Amount = uint(actualAmount + amount)

			toUpdate = append(toUpdate, tmp)
		}

		err = storage.UpadteProducts(toUpdate)
		if err != nil {
			return err
		}

		syncDatabasesMessage, err := storage.SyncDatabases()
		if err != nil {
			return err
		}

		err = ProduceSyncMessage(*syncDatabasesMessage)
		if err != nil {
			return err
		}
	}

	return nil
}
