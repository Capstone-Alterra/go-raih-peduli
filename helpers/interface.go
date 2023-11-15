package helpers

import (
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
)

type ValidationInterface interface {
	ValidateRequest(request any) []string
}

type JWTInterface interface {
	GenerateJWT(userID string, roleID string) map[string]any
	GenerateToken(userID string, roleID string) string
	GenerateTokenResetPassword(userID string, roleID string) string
	ExtractToken(token *jwt.Token) any
	ValidateToken(token string, secret string) (*jwt.Token, error)
}

type HashInterface interface {
	HashPassword(password string) string
	CompareHash(password, hashed string) bool
}

type GeneratorInterface interface {
	GenerateRandomOTP() string
}

type CloudStorageInterface interface {
	UploadFile(file multipart.File, object string) error
}

type ConverterInterface interface {
	Convert(target any, value any) error
}
