package app

import (
	"diploma/services/courier/pkg/redis"
	"diploma/services/courier/pkg/handlers"
	"diploma/services/courier/pkg/storage"
	"diploma/services/courier/pkg/mq"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func kafkaLaunch() {
	mq.New()
	
	for i := 0; i < 5; i++ {
		go mq.HandleMessages()
	}

	select {}
}

func serverLaunch() {
	redis.New()
	storage.PostgresCfg.GetConfig("courier")
	err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	
	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("front/static/"))))
	
	router.With(handlers.JWTMiddleware).Get("/courier", handlers.CourierPage)
	router.With(handlers.JWTMiddleware).Get("/courier/logout", handlers.Logout)

	router.With(handlers.JWTMiddleware).Get("/courier/set_state", handlers.SetState)
	router.With(handlers.JWTMiddleware).Get("/courier/get_state", handlers.GetState)

	fmt.Println("Courier server is running on port 8083")
	if err := http.ListenAndServe(":8083", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func Launch() {
	go kafkaLaunch()
	go serverLaunch()

	select {}
}