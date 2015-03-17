package models

import (
	_ "fmt"
)

type SessionStateResponse struct {
	SessionState string `json:"sessionState"`
}

func GetStateResponse() SessionStateResponse {
	session := GetSession()
	ReturnSession(session)
	return SessionStateResponse{SessionState: session.Session_state.String()}
}
