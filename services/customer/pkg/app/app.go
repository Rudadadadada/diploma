package app

import (
	"diploma/services/customer/pkg/handlers"
	"diploma/services/customer/pkg/mq"
	"diploma/services/customer/pkg/redis"
	"diploma/services/customer/pkg/storage"
	"fmt"
	"log"
	"net/http"

	// "github.com/confluentinc/confluent-kafka-go/kafka"
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
	storage.PostgresCfg.GetConfig("customer")
	err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	
	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("front/static/"))))
	
	router.With(handlers.JWTMiddleware).Get("/customer", handlers.CustomerPage)
	router.With(handlers.JWTMiddleware).Get("/customer/logout", handlers.Logout)
	router.With(handlers.JWTMiddleware).Post("/customer/make_order", handlers.MakeOrder)
	router.With(handlers.JWTMiddleware).Post("/customer/decline_order", handlers.DeclineOrder)
	router.With(handlers.JWTMiddleware).Get("/customer/view_orders", handlers.ViewOrdersPage)
	router.With(handlers.JWTMiddleware).Get("/customer/view_orders/view_order_items", handlers.ViewOrderItems)

	router.With(handlers.JWTMiddleware).Get("/customer/select_category", handlers.SelectCategoryPage)
	router.With(handlers.JWTMiddleware).Get("/customer/select_products/by_category", handlers.SelectProductsByCategoryIdPage)
	
	router.With(handlers.JWTMiddleware).Get("/customer/bucket", handlers.BucketPage)
	router.With(handlers.JWTMiddleware).Post("/customer/insert_into_bucket", handlers.InsertIntoBucket)
	router.With(handlers.JWTMiddleware).Post("/customer/bucket/remove_item_from_bucket", handlers.RemoveItemFromBucket)

	router.With(handlers.JWTMiddleware).Get("/order/get_statuses", handlers.GetOrderStatuses)

	fmt.Println("Customer server is running on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func Launch() {
	go kafkaLaunch()
	go serverLaunch()

	select {}
}