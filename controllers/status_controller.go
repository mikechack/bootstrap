package controllers

import (
	"github.com/mikechack/bootstrap/models"
	"log"
	"net/http"
)

func Status_controller(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	log.Printf("DMC-Agent status %s\n", r.Method)
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		buf := models.ResponseWrapper{}.JsonEncode(200, "all is well", models.GetStateResponse())
		w.Write(buf)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
