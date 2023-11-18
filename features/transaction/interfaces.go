package transaction

import (
	"raihpeduli/features/transaction/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Transaction
	PaginateUser(page, size, userID int) []Transaction
	Insert(newTransaction Transaction) int64
	SelectByID(transactionID int) *Transaction
	Update(transaction Transaction) int64
	DeleteByID(transactionID int) int64
	GetTotalData() int64
	GetTotalDataByUser(userID int) int64
}

type Usecase interface {
	FindAll(page, size, roleID, userID int) ([]dtos.ResTransaction, int64)
	FindByID(transactionID int) *dtos.ResTransaction
	Create(userID int, newTransaction dtos.InputTransaction) (*dtos.ResTransaction, error)
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
