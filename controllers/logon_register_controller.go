package controllers

import (
	"encoding/json"
	"github.com/mikechack/bootstrap/models"
	"github.com/mikechack/bootstrap/oauth"
	"github.com/mikechack/bootstrap/rabbitmq"
	"io/ioutil"
	"log"
	"net/http"
)

/*
type logonRegisterResponse struct {
	Scheme    string `json:"scheme"`
	Ipaddress string `json:"ipaddress"`
	Path      string `json:"path"`
}
*/

func doLogonRegister(name, password, organization string, serialNumber string, admin bool) (err error) {
	var bearerToken string
	var token string
	bearerToken, err = oauth.GetBearerTokenForMachineAccount(name, password, organization, admin)
	if err != nil {
		return err
	}
	log.Printf("Logon Register Bearer Token Successful\n")
	bearerToken = bearerToken
	token, err = oauth.GetTokenForMachineAccount(bearerToken)
	if err != nil {
		return err
	}
	log.Printf("Real Token  %s\n", token)
	session := models.GetSession()
	session.Session_state = models.SessionStateOauthToken
	models.ReturnSession(session)

	err = rabbitmq.ConnectAmqp()
	if err != nil {
		log.Printf("Rabbit MQ logon error %s\n", err)
		return err
	}

	dev := models.DeviceInfo{
		SerialNumber:    serialNumber,
		DeviceType:      "vTS",
		DeviceId:        "12341234",
		Ipaddress:       "10.10.10.20",
		SoftwareVersion: "1.0",
		OsVersion:       "CentOS-7",
	}

	models.AddDevice(dev)

	rabbitmq.SendRegJoin(serialNumber)

	return nil
}

func LogonRegister_controller(w http.ResponseWriter, r *http.Request) {

	log.Printf("LogonRegister URI %s\n", r.Method)

	if r.Method == "POST" {

		buf, _ := ioutil.ReadAll(r.Body)

		w.Header().Set("Content-Type", "application/json")
		ma := models.MachineAccount{}
		err := json.Unmarshal(buf[0:], &ma)
		var resp []byte
		if err != nil {
			resp = models.ResponseWrapperError{}.JsonEncode(400, "Logon Request Failed - JSON Parse error - "+err.Error())
			http.Error(w, string(resp), 400)
			return
		} else {
			models.CreateSessionFromCache(ma)
			err := doLogonRegister(ma.Username, ma.Password, ma.Organization_id, ma.SerialNumber, true)
			if err != nil {
				resp = models.ResponseWrapperError{}.JsonEncode(400, "Logon Request Failed - GetBearerToken - "+err.Error())
				http.Error(w, string(resp), 400)
				return
			}
			w.WriteHeader(200)
			resp = models.ResponseWrapper{}.JsonEncode(200, "Logon - Register Successful", nil)
		}
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
