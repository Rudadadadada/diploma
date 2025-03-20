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

	router.With(handlers.JWTMiddleware).Get("/customer/select_category", handlers.SelectCategoryPage)
	router.With(handlers.JWTMiddleware).Get("/customer/select_products/by_category", handlers.SelectProductsByCategoryId)

	// http.HandleFunc("/customer", handlers.CreateProduct)
	// http.HandleFunc("/customer", handlers.EditProduct)
	// http.HandleFunc("/products/view_all", handlers.ViewAllProducts)
	
	// http.HandleFunc("/categories/create", handlers.CreateCategory)
	// http.HandleFunc("/categories/edit", handlers.EditCategory)

	fmt.Println("Customer server is running on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
