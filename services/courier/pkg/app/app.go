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
	router.With(handlers.JWTMiddleware).Get("/courier/view_orders", handlers.ViewOrdersPage)
	router.With(handlers.JWTMiddleware).Get("/courier/view_orders/view_order_items", handlers.ViewOrderItemsPage)
	router.With(handlers.JWTMiddleware).Get("/courier/take_order", handlers.TakeOrder)
	router.With(handlers.JWTMiddleware).Get("/courier/in_progress", handlers.InProgressPage)

	router.With(handlers.JWTMiddleware).Get("/order/get_status", handlers.GetOrderStatus)
	router.With(handlers.JWTMiddleware).Get("/order/declined", handlers.DeclinedPage)
	router.With(handlers.JWTMiddleware).Get("/order/not_yet", handlers.NotYetPage)
	router.With(handlers.JWTMiddleware).Post("/courier/take_order_from_shop", handlers.TakeOrderFromShop)

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