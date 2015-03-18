package controllers

import (
	"encoding/json"
	"github.com/mikechack/bootstrap/models"
	"io/ioutil"
	"log"
	"net/http"
)

type accountRequest struct {
	AccessToken string `json:"accessToken"`
	SessionId   string `json:"sessionId"`
}

func GetMachineAccount_controller(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	var ma models.MachineAccount
	//var ma = models.MachineAccount{}

	log.Printf("GetMachineAccount URI %s\n", r.Method)

	if r.Method == "POST" {

		buf, _ := ioutil.ReadAll(r.Body)
		req := accountRequest{}
		err := json.Unmarshal(buf[0:], &req)
		if err != nil {
			resp = models.ResponseWrapperError{}.JsonEncode(500, "GetMachineAccountFailed - "+err.Error())
			http.Error(w, string(resp), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		ma, err = models.GetMachineAccount(req.AccessToken, req.SessionId)
		if err != nil {
			resp = models.ResponseWrapperError{}.JsonEncode(500, "GetMachineAccountFailed - "+err.Error())
			http.Error(w, string(resp), 500)
			return
		} else {
			w.WriteHeader(200)
			resp = models.ResponseWrapper{}.JsonEncode(200, "all is well", ma)
		}
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
