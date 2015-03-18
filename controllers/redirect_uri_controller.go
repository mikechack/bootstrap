package controllers

import (
	"encoding/json"
	"github.com/mikechack/bootstrap/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type redirectUri struct {
	Scheme    string `json:"scheme"`
	Ipaddress string `json:"ipaddress"`
	Path      string `json:"path"`
}

func RedirectUri_controller(w http.ResponseWriter, r *http.Request) {

	log.Printf("Redirect URI %s\n", r.Method)

	if r.Method == "POST" {

		buf, _ := ioutil.ReadAll(r.Body)

		w.Header().Set("Content-Type", "application/json")
		uri := redirectUri{}
		err := json.Unmarshal(buf[0:], &uri)
		var resp []byte
		if err != nil {
			resp = models.ResponseWrapperError{}.JsonEncode(400, "Get Redirect URI Failed - "+err.Error())
			http.Error(w, string(resp), 400)
			return
		} else {
			w.WriteHeader(200)
			resp, _ = models.GetRedirectUri(uri.Scheme, uri.Ipaddress, uri.Path)
			resp = []byte(strings.Replace(string(resp), "\\u0026", "&", -1))
		}
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
