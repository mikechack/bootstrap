package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sqbu.com/MediaFusion/bootstrap/models"
)

type TokenResponse struct {
	Access_token             string `json:"access_token,omitempty"`
	Expires_in               int    `json:"expires_in,omitempty"`
	Refresh_token            string `json:"refresh_token,omitempty"`
	Refresh_token_expires_in int    `json:"refresh_token_expires_in,omitempty"`
}

type bearerTokenResponse struct {
	BearerToken string `json:"BearerToken,omitempty"`
}

type machineAccountRequest struct {
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	AdminUser bool   `json:"adminUser,omitempty"`
}

var clientId = "C71e2f13edd03a6307b9591f529345a90447d83814b6db35c26c18fc81044da2e"
var clientSecret = "be96c991c972f84696263449e4cdade3a12b4acec52996b3a529666f5a8b1237"

func GetBearerTokenForMachineAccountStoredSession() (btoken string, err error) {
	session := models.GetSession()

	defer func() {
		models.ReturnSession(session)
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}

			err, _ = r.(error)

			if str, ok := r.(string); ok {
				err = errors.New(str)
			}
			log.Printf("GetBearerTokenForMachineAccountStoredSession - %s\n", err)
		}
	}()

	btoken, err = GetBearerTokenForMachineAccount(session.Ma_name, session.Ma_password, session.Ma_organization, true)
	return btoken, err

}

func GetBearerTokenForMachineAccount(name, password, organization string, admin bool) (btoken string, err error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}

			err, _ = r.(error)

			if str, ok := r.(string); ok {
				err = errors.New(str)
				return
			}
			log.Printf("GetBearerTokenForMachineAccount - %s\n", err)
			return
		}
	}()

	log.Printf("GetBearerTokenForMachineAccount - name = %s\n", name)
	log.Printf("GetBearerTokenForMachineAccount - password = %s\n", password)

	var creds = machineAccountRequest{Name: name, Password: password, AdminUser: true}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://idbroker.webex.com/idb/token/"+organization+"/v1/actions/GetBearerToken/invoke", nil)
	req.Header.Set("Content-Type", "application/json")
	buf, err := json.Marshal(creds)
	req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	if err != nil {
		log.Printf("Json error %v\n", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Bearer token request error")
		return btoken, err
	}

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Printf("Bearer token request error \n%s\n", string(body))
		return btoken, errors.New(res.Status)
	}

	var token = bearerTokenResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &token); err != nil {
		log.Fatal(err)
	}

	//log.Printf("Here is the token %v", m.Access_token)
	res.Body.Close()

	return token.BearerToken, err
}

func GetTokenForMachineAccount(bearer string) (t string, err error) {
	session := models.GetSession()
	defer func() {
		models.ReturnSession(session)
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}

			err, _ = r.(error)

			if str, ok := r.(string); ok {
				err = errors.New(str)
				return
			}
			log.Printf("GetTokenForMachineAccount - %s\n", err)
			return
		}
	}()

	//var creds = oauthCredentials{clientId, clientSecret}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://idbroker.webex.com/idb/oauth2/v1/access_token", nil)
	//req.Header.Set("Authorization", creds.getAuthorization())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	v := url.Values{}
	v.Add("grant_type", "urn:ietf:params:oauth:grant-type:saml2-bearer")
	v.Add("assertion", bearer)
	v.Add("client_id", clientId)
	v.Add("client_secret", clientSecret)
	v.Add("scope", "squared-fusion-mgmt:management squared-fusion-media:device_connect Identity:SCIM Identity:Organization")
	//req.Body = nopCloser{bytes.NewBufferString(v.Encode())}
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(v.Encode())))

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Token request error")
		return t, err
	}

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Printf("Token request error \n%s\n", string(body))
		return t, errors.New(res.Status)
	}

	var token = TokenResponse{}
	body, err := ioutil.ReadAll(res.Body)

	log.Printf("Token for Machine Response Code      %s\n", res.Status)
	log.Printf("Token for Machine Response        %s\n", body)
	if err = json.Unmarshal(body, &token); err != nil {
		log.Printf("Token request error")
		return t, errors.New(res.Status)
	}

	log.Printf("Token for Machine Here is the token  %s\n", token.Access_token)
	res.Body.Close()

	session.Access_token = token.Access_token
	session.Refresh_token = token.Refresh_token
	session.Expires_in = token.Expires_in
	session.Refresh_expires_in = token.Refresh_token_expires_in

	return token.Access_token, err
}
