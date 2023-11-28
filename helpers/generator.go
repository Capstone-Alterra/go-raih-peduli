package helpers

import (
	"math/rand"
	"time"
)

type generator struct{}

func NewGenerator() GeneratorInterface {
	return &generator{}
}

func (g generator) GenerateRandomOTP() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	const number = "0123456789"

	otp := make([]byte, 6)
	for i := range otp {
		otp[i] = number[r.Intn(len(number))]
	}

	return string(otp)
}
func (g generator) GenerateRandomID() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000000
	max := 9999999
	return int(min + rand.Intn(max-min+1))
}
