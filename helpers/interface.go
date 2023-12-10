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
	RefereshJWT(refreshToken *jwt.Token) map[string]any
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
	CreateTransactionBank(IDTransaction string, PaymentType string, Amount int64) (string, string, error)
	CreateTransactionGopay(IDTransaction string, PaymentType string, Amount int64) (string, string, error)
	CreateTransactionQris(IDTransaction string, PaymentType string, Amount int64) (string, string, error)
	TransactionStatus(transactionStatusResp *coreapi.TransactionStatusResponse) transaction.Status
	CheckTransactionStatus(IDTransaction string) (string, error)
	MappingPaymentName(paymentType string) string
}

type NotificationInterface interface {
	SendNotifications(tokens string, userID string, message string) error
}

type CloudStorageInterface interface {
	UploadFile(file multipart.File, object string) error
	DeleteFile(object string) error
}

type ConverterInterface interface {
	Convert(target any, value any) error
}

type OpenAIInterface interface {
	GetAppInformation(question string, qnaList map[string]string) (string, error)
	GetNewsContent(prompt string) (string, error)
}
