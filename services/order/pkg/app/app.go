package app

import (
	"diploma/services/order/pkg/mq"
)

func Launch() {
	mq.New()
	
	for i := 0; i < 5; i++ {
		go mq.HandleMessages()
	}

	select {}
}