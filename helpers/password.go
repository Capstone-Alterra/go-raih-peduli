package helpers

import "golang.org/x/crypto/bcrypt"

type HashInterface interface {
	HashPassword(password string) string
	CompareHash(password, hashed string) bool
}

type hash struct{}

func NewHash() HashInterface {
	return &hash{}
}

func (h hash) HashPassword(password string) string {
	result, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(result)
}

func (h hash) CompareHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
