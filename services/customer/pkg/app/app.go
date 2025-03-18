package app

import (
	// "diploma/services/customer/pkg/storage"
	"diploma/services/customer/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func Launch() {
	// storage.PostgresCfg.GetConfig("customer")
	// err := storage.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("front/static/"))))
	
	router.Get("/customer", handlers.CustomerPage)

	// http.HandleFunc("/customer", handlers.CreateProduct)
	// http.HandleFunc("/customer", handlers.EditProduct)
	// http.HandleFunc("/products/view_all", handlers.ViewAllProducts)
	
	// http.HandleFunc("/categories/create", handlers.CreateCategory)
	// http.HandleFunc("/categories/edit", handlers.EditCategory)

	fmt.Println("Server is running on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
