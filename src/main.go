package main

import (
	"fmt"
	"net/http"

	"go-url-shortener/dal"
	"go-url-shortener/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db := dal.MustGetDb()
	router := mux.NewRouter()

	baseUrl := "localhost:8080"

	router.HandleFunc(
		"/status", handlers.StatusHandler(db)).Methods(http.MethodGet)

	router.HandleFunc(
		"/", handlers.IndexHandler()).Methods(http.MethodGet)

	router.HandleFunc(
		"/", handlers.CreateShortenedURLHandler(db)).Methods(http.MethodPost)

	router.Handle(
		"/success", &handlers.SuccessHandler{BaseUrl: baseUrl}).Methods(http.MethodGet)

	router.HandleFunc(
		"/{urlId}", handlers.RedirectHandler(db)).Methods(http.MethodGet)

	err := http.ListenAndServe(baseUrl, router)
	if err != nil {
		panic(fmt.Sprintf("canot start server because: %s", err.Error()))
	}
}
