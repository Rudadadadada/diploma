package app

import (
	"diploma/services/admin/pkg/storage"
	"diploma/services/admin/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func Launch() {
	storage.PostgresCfg.GetConfig("admin")
	err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("front/static"))))
    
	router.Get("/admin", handlers.AdminPage)
	router.Get("/admin/categories", handlers.CategoriesPage)
	router.Get("/admin/products", handlers.ProductsPage)

	router.Get("/admin/products/view_all", handlers.ViewAllProducts)
	router.Get("/admin/products/create", handlers.CreateProductPage)
	router.Post("/admin/products/create", handlers.CreateProduct)
	router.Get("/admin/products/remove", handlers.RemoveProductPage)
	router.Post("/admin/products/remove", handlers.RemoveProduct)
	
	router.Get("/admin/categories/view_all", handlers.ViewAllCategories)
	router.Get("/admin/categories/create", handlers.CreateCategoryPage)
	router.Post("/admin/categories/create", handlers.CreateCategory)
	router.Get("/admin/categories/remove", handlers.RemoveCategoryPage)
	router.Post("/admin/categories/remove", handlers.RemoveCategory)


	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
