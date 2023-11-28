package transaction

import (
	"raihpeduli/features/transaction/dtos"
	"raihpeduli/features/user"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, keyword string) []Transaction
	CountByID(fundraiseID int) (int64, error)
	PaginateUser(page, size, userID int) []Transaction
	SelectUserByID(userID int) *user.User
	Insert(newTransaction Transaction) int64
	SelectByID(transactionID int) *Transaction
	Update(transaction Transaction) int64
	DeleteByID(transactionID int) int64
	GetTotalData(keyword string) int64
	GetTotalDataByUser(userID int, keyword string) int64
}

type Usecase interface {
	FindAll(page, size, roleID, userID int, keyword string) ([]dtos.ResTransaction, int64)
	FindByID(transactionID int) *dtos.ResTransaction
	Create(userID int, newTransaction dtos.InputTransaction) (*dtos.ResTransaction, error, []string)
	Modify(transactionData dtos.InputTransaction, transactionID int) bool
	Remove(transactionID int) bool
	Notifications(notificationPayload map[string]any) error
}

type Handler interface {
	GetTransactions() echo.HandlerFunc
	TransactionDetails() echo.HandlerFunc
	CreateTransaction() echo.HandlerFunc
	UpdateTransaction() echo.HandlerFunc
	DeleteTransaction() echo.HandlerFunc
	Notifications() echo.HandlerFunc
}
