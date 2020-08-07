package handlers

import (
	"go-url-shortener/dal"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RedirectHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// will redirect to the shortened url, or 404 if not registered
		vars := mux.Vars(r)
		urlId := vars["urlId"]

		redirectURL, exists := dal.GetOriginalURL(db, urlId)

		if exists == false {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}
