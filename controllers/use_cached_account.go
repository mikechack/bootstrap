package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sqbu.com/MediaFusion/bootstrap/models"
	"sqbu.com/MediaFusion/bootstrap/oauth"
	"sqbu.com/MediaFusion/bootstrap/rabbitmq"
)

func UseCachedAccount(w http.ResponseWriter, r *http.Request) {
	var token string
	log.Printf("Http Use Cached Account - Method %s", r.Method)
	r.ParseForm()
	fname := r.FormValue("filename")

	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		resp := models.ResponseWrapperError{}.JsonEncode(400, "Use cached account - file read error - "+err.Error())
		log.Printf(string(resp))
		http.Error(w, string(resp), 400)
		return
	}

	var ma = models.MachineAccount{}
	if err = json.Unmarshal(buf, &ma); err != nil {
		resp := models.ResponseWrapperError{}.JsonEncode(400, "Use cached account- JSON Parse error - "+err.Error())
		http.Error(w, string(resp), 400)
		return
	}

	log.Printf("User Name = %s\n", ma.Username)
	log.Printf("Password  = %s\n", ma.Password)

	models.CreateSessionFromCache(ma)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	bearerToken, _ := oauth.GetBearerTokenForMachineAccountStoredSession()
	log.Printf("Bearer Token \n%s\n", bearerToken)
	token, err = oauth.GetTokenForMachineAccount(bearerToken)
	log.Printf("Real Token  %s\n", token)
	err = rabbitmq.ConnectAmqp()
	log.Printf("Rabbit MQ logon error if any %s\n", err)

	//resp, _ := json.Marshal(ma)
	resp := models.ResponseWrapper{}.JsonEncode(200, "all is well", ma)
	w.Write(resp)

}
