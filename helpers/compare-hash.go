package helpers

import "golang.org/x/crypto/bcrypt"

func CompareHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}