package models

import (
	"bytes"
	"encoding/json"
	"github.com/mikechack/bootstrap/utility"
	"log"
	"net/url"
)

type redirectUriT struct {
	Uri string `json:"uri"`
}

func (a redirectUriT) ToJSON() []byte {
	buf := make([]byte, 1500)
	buf, _ = json.Marshal(a)
	return buf
}

type finalRedirectUriT struct {
	Redirect_uri  string `json:"redirect_uri"`
	Session_id    string `json:"session_id"`
	Shared_secret string `json:"shared_secret"`
	Box_name      string `json:"box_name"`
}

func (a finalRedirectUriT) ToJSON() []byte {
	buf := make([]byte, 1500)
	buf, _ = json.Marshal(a)
	return buf
}

func (a *finalRedirectUriT) create(scheme string, ipaddress string, path string) {
	redir := url.URL{}
	redir.Scheme = scheme
	redir.Host = ipaddress
	redir.Path = path
	a.Redirect_uri = redir.String()
	a.Session_id = utility.GetRandomString(48)
	a.Shared_secret = utility.GetRandomString(32)

	session := GetSession()
	log.Printf("previous session id = %s\n", session.session_id)
	session.session_id = a.Session_id
	session.shared_secret = a.Shared_secret
	session.Session_state = SessionStateInit
	ReturnSession(session)
}

func DecryptToken(cryptToken string) string {
	token, _ := utility.DecryptAesCBC([]byte(finalRedirect.Shared_secret), cryptToken)
	return token
}

var finalRedirect = finalRedirectUriT{}

var genurl = url.URL{}

func createRawQuery() string {
	var buffer bytes.Buffer

	buffer.WriteString("response_type=token")
	buffer.WriteString("&")
	buffer.WriteString("client_id=C71e2f13edd03a6307b9591f529345a90447d83814b6db35c26c18fc81044da2e")
	buffer.WriteString("&")
	buffer.WriteString("redirect_uri=https%3A%2F%2Fhercules.ladidadi.org%2Ffuse_redirect")
	buffer.WriteString("&")
	buffer.WriteString("scope=Identity%3ASCIM%20Identity%3AOrganization%20squared-fusion-mgmt%3Amanagement")
	buffer.WriteString("&")
	buffer.WriteString("state=")

	return buffer.String()
}

func GetRedirectUri(scheme string, ipaddress, path string) ([]byte, string) {

	log.Printf("URI Scheme    %s\n", scheme)
	log.Printf("URI ipaddress %s\n", ipaddress)
	log.Printf("URI path      %s\n", path)

	genurl.Scheme = "https"
	genurl.Host = "idbroker.webex.com"
	genurl.Path = "idb/oauth2/v1/authorize"
	genurl.RawQuery = createRawQuery()

	finalRedirect.create(scheme, ipaddress, path)

	_redir := redirectUriT{}
	_redir.Uri = genurl.String()
	_redir.Uri += utility.EncryptPKCS1v15(finalRedirect.ToJSON())

	log.Printf("\n\nRedirect URI String = %s\n\n", _redir.Uri)

	return ResponseWrapper{}.JsonEncode(200, "all is well", _redir), _redir.Uri
}
