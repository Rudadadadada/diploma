package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func CourierPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/courier/courier.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}