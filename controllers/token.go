package controllers

import (
	"log"
	"net/http"
	"sqbu.com/MediaFusion/bootstrap/models"
	"sqbu.com/MediaFusion/bootstrap/oauth"
	_ "sqbu.com/MediaFusion/bootstrap/rabbitmq"
)

func Token(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		crypttoken := r.FormValue("access_token")
		cryptId := r.FormValue("session_id")
		log.Printf("Token  crypted  - %s\n", crypttoken)
		log.Printf("Session crypted - %s\n", cryptId)
		models.GetMachineAccount(crypttoken, cryptId)
		bearerToken, _ := oauth.GetBearerTokenForMachineAccountStoredSession()
		log.Printf("Bearer Token \n%s\n", bearerToken)
		log.Printf("About to get Real Token \n")
		token, _ := oauth.GetTokenForMachineAccount(bearerToken)
		log.Printf("Real Token  %s\n", token)
		//err := rabbitmq.ConnectAmqp()
		//log.Printf("Rabbit MQ logon error if any %s\n", err)
		/*
			log.Printf("Session valid   - %v\n", models.ValidateSessionId(cryptId, crypttoken))
			sId = oauth.DecryptAesCBC(key, sId)
			log.Printf("Decrypted Session = %s\n", sId)
			tempToken = oauth.DecryptAesCBC(key, token)
			log.Printf("Decrypted Token   = %s\n", tempToken)
			ma := getMachineAccount(tempToken, sId)
			machine = ma
			bearerToken = oauth.GetBearerTokenForMachineAccount(machine)
			token = oauth.GetTokenForMachineAccount(bearerToken)
			fmsRegisterDevice(token)
			//http.Redirect(w, r, "https://int-admin.wbx2.com/#/login", 302)
		*/
		w.WriteHeader(200)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
