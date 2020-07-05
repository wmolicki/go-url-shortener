package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"go-url-shortener/dal"
	"go-url-shortener/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db := dal.MustGetDb()
	router := mux.NewRouter()

	indexTemplate := template.Must(template.ParseFiles("static/index.html"))
	baseUrl := "localhost:8080"

	router.HandleFunc(
		"/status", handlers.StatusHandler(db)).Methods(http.MethodGet)

	router.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			// will serve the main page
			indexTemplate.Execute(w, nil)
		}).Methods(http.MethodGet)

	router.Handle(
		"/success", &handlers.SuccessHandler{BaseUrl: baseUrl}).Methods(http.MethodGet)

	router.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			// will save posted url as new
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%v", err)))
			}

			url := r.FormValue("url")
			if url == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("You have to provide nonempty url"))
			}
			if strings.HasPrefix(url, "https://") == false {
				url = "https://" + url
			}

			userAgent := r.Header.Get("User-Agent")
			shortened := dal.InsertShortenedURL(db, url, userAgent)

			http.Redirect(w, r, fmt.Sprintf("/success?s=%s", shortened), http.StatusSeeOther)

		}).Methods(http.MethodPost)

	router.HandleFunc(
		"/{urlId}", func(w http.ResponseWriter, r *http.Request) {
			// will redirect to the shortened url, or 404 if not registered
			vars := mux.Vars(r)
			urlId := vars["urlId"]

			redirectURL, exists := dal.GetOriginalURL(db, urlId)

			if exists == false {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)

		}).Methods(http.MethodGet)

	err := http.ListenAndServe(baseUrl, router)
	if err != nil {
		panic(fmt.Sprintf("canot start server because: %s", err.Error()))
	}
}
