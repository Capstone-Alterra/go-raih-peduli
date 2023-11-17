package helpers

import (
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
	"github.com/midtrans/midtrans-go/coreapi"

	"raihpeduli/features/transaction"
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
	GenerateRandomID() int
}

type MidtransInterface interface {
	CreateTransactionBank(IDTransaction string, PaymentType string, Amount int64) (string, error)
	CreateTransactionGopay(IDTransaction string, PaymentType string, Amount int64) (string, error)
	TransactionStatus(transactionStatusResp *coreapi.TransactionStatusResponse) transaction.Status
}

type CloudStorageInterface interface {
	UploadFile(file multipart.File, object string) error
}

type ConverterInterface interface {
	Convert(target any, value any) error
}
