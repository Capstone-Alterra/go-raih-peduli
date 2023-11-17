package repository

import (
	"raihpeduli/features/transaction"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) transaction.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []transaction.Transaction {
	var transactions []transaction.Transaction

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&transactions)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return transactions
}

func (mdl *model) Insert(newTransaction transaction.Transaction) int64 {
	result := mdl.db.Create(&newTransaction)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newTransaction.ID)
}

func (mdl *model) SelectByID(transactionID int) *transaction.Transaction {
	var transaction transaction.Transaction
	result := mdl.db.First(&transaction, transactionID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &transaction
}

func (mdl *model) Update(transaction transaction.Transaction) int64 {
	result := mdl.db.Updates(&transaction)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(transactionID int) int64 {
	result := mdl.db.Delete(&transaction.Transaction{}, transactionID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}
