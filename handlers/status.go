package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type StatusResponse struct {
	DBStatus string
}

func StatusHandler(db *sqlx.DB) http.HandlerFunc {  // type alias, same as func(http.ResponseWriter, *http.Request)
	return func(w http.ResponseWriter, r *http.Request) {
		var dbOk string
		err := db.Ping()

		if err != nil {
			dbOk = fmt.Sprintf("false, err: %s", err)
		} else {
			dbOk = "ok"
		}

		s := StatusResponse{DBStatus: dbOk}

		jsonData, err := json.Marshal(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error: %v", err)))
		} else {
			w.Write(jsonData)
		}

	}
}
