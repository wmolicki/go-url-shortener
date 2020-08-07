package handlers

import (
	"fmt"
	"go-url-shortener/dal"
	"html/template"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

var indexTemplate = template.Must(template.ParseFiles("static/index.html"))

func CreateShortenedURLHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

	}
}

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// will serve the main page
		indexTemplate.Execute(w, nil)
	}
}
