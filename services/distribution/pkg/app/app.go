package app

import (
	"diploma/services/distribution/pkg/storage"
	"diploma/services/distribution/pkg/mq"
	"log"
)

func Launch() {
	mq.New()
	storage.PostgresCfg.GetConfig("distribution")
	err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		go mq.HandleMessages()
	}

	select {}
}