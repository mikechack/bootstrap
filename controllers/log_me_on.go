package controllers

import (
	"log"
	"net/http"
	"sqbu.com/MediaFusion/bootstrap/models"
)

func LogMeOn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Http LogMeOn - Method %s", r.Method)
	r.ParseForm()
	scheme := r.FormValue("scheme")
	ipaddress := r.FormValue("ipaddress")
	path := r.FormValue("path")
	_, uri := models.GetRedirectUri(scheme, ipaddress, path)
	log.Printf("Log Me On redirect URI = %s\n", uri)
	http.Redirect(w, r, uri, 302)

}
