package models

import (
	"crypto/tls"
	"github.com/mikechack/bootstrap/certs"
	"github.com/mikechack/bootstrap/utility"
	"log"
)

type SessionState int

const (
	SessionStateIdle SessionState = 1 + iota
	SessionStateInit
	SessionStateMachineAccount
	SessionStateOauthToken
)

var states = [...]string{
	"SessionStateIdle",
	"SessionStateInit",
	"SessionStateMachineAccount",
	"SessionStateOauthToken",
}

func (s SessionState) String() string { return states[s-1] }

var cert tls.Certificate

type ConnectorSession struct {
	session_id         string
	shared_secret      string
	temp_token         string
	Access_token       string
	Refresh_token      string
	Expires_in         int
	Refresh_expires_in int
	Ma_name            string
	Ma_password        string
	Ma_account         string
	Ma_organization    string
	Session_state      SessionState
}

var sessionChan = make(chan ConnectorSession, 1)

func CreateSessionFromCache(ma MachineAccount) {
	session := GetSession()

	session.Ma_account = ma.Account_id
	session.Ma_name = ma.Username
	session.Ma_organization = ma.Organization_id
	session.Ma_password = ma.Password

	ReturnSession(session)

}

func GetSession() ConnectorSession {

	var session ConnectorSession

	select {
	case session = <-sessionChan:
		return session
	}
}

func ReturnSession(session ConnectorSession) {

	sessionChan <- session

}

func ValidateSessionId(cryptSessionId string, crypttoken string) bool {
	session := GetSession()

	log.Printf("previous session id = %s\n", session.session_id)
	sId, _ := utility.DecryptAesCBC([]byte(session.shared_secret), cryptSessionId)
	session.temp_token, _ = utility.DecryptAesCBC([]byte(session.shared_secret), crypttoken)

	ReturnSession(session)
	return sId == session.session_id
}

func init() {

	sessionChan <- ConnectorSession{Session_state: SessionStateIdle}

	var err error
	//cert, err = tls.LoadX509KeyPair("/Users/mchack/Documents/GitRepositories/wsrouter/src/certs/client0.crt", "/Users/mchack/Documents/GitRepositories/wsrouter/src/certs/client0.key")
	clientCert, _ := certs.Asset("certs/client0.crt")
	clientKey, _ := certs.Asset("certs/client0.key")
	cert, err = tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

}
