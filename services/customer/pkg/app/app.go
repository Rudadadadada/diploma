package app

import (
	"diploma/services/customer/pkg/redis"
	"diploma/services/customer/pkg/storage"
	"diploma/services/customer/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func Launch() {
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
	router.With(handlers.JWTMiddleware).Get("/customer/view_orders", handlers.ViewOrdersPage)
	router.With(handlers.JWTMiddleware).Get("/customer/view_orders/view_order_items", handlers.ViewOrderItems)

	router.With(handlers.JWTMiddleware).Get("/customer/select_category", handlers.SelectCategoryPage)
	router.With(handlers.JWTMiddleware).Get("/customer/select_products/by_category", handlers.SelectProductsByCategoryIdPage)
	
	router.With(handlers.JWTMiddleware).Get("/customer/bucket", handlers.BucketPage)
	router.With(handlers.JWTMiddleware).Post("/customer/insert_into_bucket", handlers.InsertIntoBucket)
	router.With(handlers.JWTMiddleware).Post("/customer/bucket/remove_item_from_bucket", handlers.RemoveItemFromBucket)

	fmt.Println("Customer server is running on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
