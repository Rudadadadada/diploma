package handlers

import (
	"log"
	"net/http"
	"html/template"
)

type SuccessData struct {
    Title string
    Path  string
}

func CustomerPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/customer/customer.html")
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