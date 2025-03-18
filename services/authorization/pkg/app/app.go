package app

import (
	"diploma/services/authorization/pkg/storage"
	"diploma/services/authorization/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func Launch() {
	storage.PostgresCfg.GetConfig("authorization")
	err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("front/static"))))

	router.Get("/registration/customer", handlers.CustomerRegistrationPage)
	router.Post("/registration/customer", handlers.CustomerRegistration)
	router.Get("/registration/courier", handlers.CourierRegistrationPage)
	router.Post("/registration/courier", handlers.CourierRegistration)

	router.Get("/authorization/customer", handlers.CustomerAuthorizationPage)
	router.Post("/authorization/customer", handlers.CustomerAuthorization)
	router.Get("/authorization/courier", handlers.CourierAuthorizationPage)
	router.Post("/authorization/courier", handlers.CourierAuthorization)
	router.Get("/authorization/admin", handlers.AdminAuthorizationPage)
	router.Post("/authorization/admin", handlers.AdminAuthorization)


	fmt.Println("Authorization server is running on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
