package repository

import (
	"raihpeduli/features/user"

	"raihpeduli/features/fundraise"
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

func (mdl *model) Paginate(page, size int, keyword string) []transaction.Transaction {
	var transactions []transaction.Transaction

	offset := (page - 1) * size
	searching := "%" + keyword + "%"

	result := mdl.db.Preload("User").
		Table("transactions").
		Joins("JOIN users ON transactions.user_id = users.id").
		Where("users.fullname LIKE ?", searching).
		Offset(offset).Limit(size).
		Find(&transactions)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return transactions
}

func (mdl *model) CountByID(fundraiseID int) (int64, error) {
	var count int64

	if err := mdl.db.Model(&fundraise.Fundraise{}).Where("id = ?", fundraiseID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (mdl *model) SelectUserByID(userID int) *user.User {
	var user user.User
	result := mdl.db.First(&user, userID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &user
}

func (mdl *model) GetTotalData(keyword string) int64 {
	var totalData int64
	searching := "%" + keyword + "%"

	result := mdl.db.Model(&transaction.Transaction{}).
		Preload("User", "fullname LIKE ?", searching).
		Joins("JOIN users ON transactions.user_id = users.id").
		Where("users.fullname LIKE ?", searching).
		Count(&totalData)

	//result := mdl.db.Model(&transaction.Transaction{}).Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) GetTotalDataByUser(userID int, keyword string) int64 {
	var totalData int64
	searching := "%" + keyword + "%"

	result := mdl.db.Model(&transaction.Transaction{}).
		Preload("User", "fullname LIKE ?", searching).
		Joins("JOIN users ON transactions.user_id = users.id").
		Where("users.fullname LIKE ?", searching).
		Count(&totalData)

	// result := mdl.db.Model(&transaction.Transaction{}).Where("user_id = ?", userID).Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) PaginateUser(page, size, userID int) []transaction.Transaction {
	var transactions []transaction.Transaction

	offset := (page - 1) * size

	result := mdl.db.Preload("User").Where("user_id = ?", userID).Offset(offset).Limit(size).Find(&transactions)

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
	result := mdl.db.Preload("User").First(&transaction, transactionID)

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
