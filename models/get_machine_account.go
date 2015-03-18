package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/mikechack/bootstrap/utility"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

type MachineAccount struct {
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	Location        string `json:"location,omitempty"`
	Organization_id string `json:"organization_id,omitempty"`
	Account_id      string `json:"account_id,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
}

type jsonSessionId struct {
	Session_id     string `json:"session_id,omitempty"`
	Connector_type string `json:"connector_type,omitempty"`
}

func saveMachineAccount(fname string, machine MachineAccount) {
	buf, _ := json.Marshal(machine)
	ioutil.WriteFile(fname, buf, os.ModePerm)
}

func GetMachineAccount(cryptToken, cryptSessionId string) (ma MachineAccount, err error) {
	session := GetSession()

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}

			err, _ = r.(error)

			if str, ok := r.(string); ok {
				err = errors.New(str)
			}
			log.Printf("GetMachineAccount - %s\n", err)
			ReturnSession(session)
		}
	}()

	if session.Session_state != SessionStateInit {
		panic(errors.New("Incorrect state, expected SessionStateInit - got " + session.Session_state.String()))
	}

	sId, err := utility.DecryptAesCBC([]byte(session.shared_secret), cryptSessionId)
	if err != nil {
		panic(err)
	}

	session.temp_token, _ = utility.DecryptAesCBC([]byte(session.shared_secret), cryptToken)
	log.Printf("Machine account session check %v\n", sId == session.session_id)

	sessionid := jsonSessionId{Session_id: sId, Connector_type: "dmc_management_connector"}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("POST", "https://hercules.ladidadi.org/v1/machine_accounts", nil)
	//req, _ := http.NewRequest("POST", "https://hercules.hitest.huron-dev.com/v1/machine_accounts", nil)
	req.Header.Set("Authorization", "Bearer "+session.temp_token)
	req.Header.Set("Content-Type", "application/json")
	buf, err := json.Marshal(sessionid)
	log.Printf("Hercules-GetMA - json - %s\n", buf)
	req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	if err != nil {
		log.Printf("Json error %v\n", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Hercules-GetMA - Status - %s\n", res.Status)

	ma = MachineAccount{}
	body, err := ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &ma); err != nil {
		log.Fatal(err)
	}

	session.Ma_password = ma.Password
	session.Ma_name = ma.Username
	session.Ma_organization = ma.Organization_id

	re := regexp.MustCompile("([^/]+)$")
	id := re.Find([]byte(ma.Location))
	ma.Account_id = string(id)
	session.Ma_account = ma.Account_id
	saveMachineAccount("./machine."+string(id)+".conf", ma)

	log.Printf("Hercules-GetMA - Name         - %s\n", session.Ma_name)
	log.Printf("Hercules-GetMA - Password     - %s\n", session.Ma_password)
	log.Printf("Hercules-GetMA - Organization - %s\n", session.Ma_organization)
	log.Printf("Hercules-GetMA - Accountid    - %s\n", session.Ma_account)

	res.Body.Close()

	session.Session_state = SessionStateMachineAccount

	ReturnSession(session)
	return ma, err
}
