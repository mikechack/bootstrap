package models

import (
	"encoding/json"
)

type responseStatus struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseWrapper struct {
	Status   responseStatus `json:"status,omitempty"`
	Response interface{}    `json:"response,omitempty"`
}

type ResponseWrapperError struct {
	Status responseStatus `json:"status,omitempty"`
}

func (wrapper ResponseWrapper) JsonEncode(code int, message string, response interface{}) []byte {
	wrapper.Status = responseStatus{Code: code, Message: message}
	if response != nil {
		wrapper.Response = response
	}

	buf, _ := json.MarshalIndent(wrapper, "", "   ")
	return buf
}

func (wrapper ResponseWrapperError) JsonEncode(code int, message string) []byte {
	wrapper.Status = responseStatus{Code: code, Message: message}

	buf, _ := json.MarshalIndent(wrapper, "", "   ")
	return buf
}
