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
    
	router.With(handlers.JWTMiddleware).Get("/admin", handlers.AdminPage)
	router.With(handlers.JWTMiddleware).Get("/admin/categories", handlers.CategoriesPage)
	router.With(handlers.JWTMiddleware).Get("/admin/products", handlers.ProductsPage)

	router.With(handlers.JWTMiddleware).Get("/admin/products/view_all", handlers.ViewAllProducts)
	router.With(handlers.JWTMiddleware).Get("/admin/products/create", handlers.CreateProductPage)
	router.With(handlers.JWTMiddleware).Post("/admin/products/create", handlers.CreateProduct)
	router.With(handlers.JWTMiddleware).Get("/admin/products/remove", handlers.RemoveProductPage)
	router.With(handlers.JWTMiddleware).Post("/admin/products/remove", handlers.RemoveProduct)
	
	router.With(handlers.JWTMiddleware).Get("/admin/categories/view_all", handlers.ViewAllCategories)
	router.With(handlers.JWTMiddleware).Get("/admin/categories/create", handlers.CreateCategoryPage)
	router.With(handlers.JWTMiddleware).Post("/admin/categories/create", handlers.CreateCategory)
	router.With(handlers.JWTMiddleware).Get("/admin/categories/remove", handlers.RemoveCategoryPage)
	router.With(handlers.JWTMiddleware).Post("/admin/categories/remove", handlers.RemoveCategory)


	fmt.Println("Admin server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
