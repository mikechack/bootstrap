package models

import (
	"log"
	"testing"
)

func TestState(t *testing.T) {
	state := GetState()
	if state == false {
		log.Printf("state = %v", state)
		//t.("State was not in the correct state: ", state)
	}
}
