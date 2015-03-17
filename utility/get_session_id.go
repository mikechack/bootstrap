package utility

import (
	"log"
	"math/rand"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func GetRandomString(length int) string {

	alphabet := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	b := make([]byte, length)

	for i := 0; i < len(b); i++ {
		b[i] = alphabet[random(0, len(alphabet)-1)]
	}

	s := string(b[:])
	log.Printf("Generate Random String %s", s)
	return s
}
