package app

import (
	"diploma/db_configs"
	"diploma/services/admin/pkg/handlers"
	"diploma/services/admin/pkg/storage"
	"fmt"
	"log"
	"net/http"
)

func LaunchAdminService(db_config *db_configs.DBConfig) {
	err := storage.New(db_config)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/products/create", handlers.CreateProduct)
	http.HandleFunc("/products/edit", handlers.EditProduct)
	http.HandleFunc("/products/view_all", handlers.ViewAllProducts)
	
	http.HandleFunc("/categories/create", handlers.CreateCategory)
	http.HandleFunc("/categories/edit", handlers.EditCategory)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
