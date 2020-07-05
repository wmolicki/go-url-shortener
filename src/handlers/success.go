package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

var successTemplate = template.Must(template.ParseFiles("static/success.html"))

type SuccessPageData struct {
	BaseUrl string
	Key     string
}

type SuccessHandler struct {
	BaseUrl string
}

func (h *SuccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	key := r.FormValue("s")

	data := SuccessPageData{BaseUrl: h.BaseUrl, Key: key}
	//data := map[string]interface{}{"BaseUrl": h.BaseUrl, "Key": key}

	err = successTemplate.Execute(w, data)
	if err != nil {
		panic(fmt.Sprintf("Template render failure: %v", err))
	}
}
