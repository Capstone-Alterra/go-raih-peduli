package helpers

import (
	"math/rand"
	"time"
)

func GenerateRandomOTP() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	const number = "0123456789"

	otp := make([]byte, 6)
	for i := range otp {
		otp[i] = number[r.Intn(len(number))]
	}

	return string(otp)
}
