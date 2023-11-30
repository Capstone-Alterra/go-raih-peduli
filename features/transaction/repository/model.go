package repository

import (
	"fmt"
	"raihpeduli/config"
	"raihpeduli/features/user"
	"strconv"

	"raihpeduli/features/fundraise"
	"raihpeduli/features/transaction"

	"github.com/labstack/gommon/log"
	"github.com/wneessen/go-mail"
	"gorm.io/gorm"
)

type model struct {
	db     *gorm.DB
	config *config.SMTPConfig
}

func New(db *gorm.DB, config *config.SMTPConfig) transaction.Repository {
	return &model{
		db:     db,
		config: config,
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

func (mdl *model) GetFundraiseByID(fundraiseID int) (*fundraise.Fundraise, error) {
	var result fundraise.Fundraise

	if err := mdl.db.Where("id = ?", fundraiseID).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
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

func (mdl *model) SendPaymentConfirmation(email string, amount int, idFundraise int, paymentType string) error {
	user := mdl.config.SMTP_USER
	password := mdl.config.SMTP_PASS
	port := mdl.config.SMTP_PORT

	fundraise, err := mdl.GetFundraiseByID(idFundraise)
	if err != nil {
		return err
	}

	convPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	m := mail.NewMsg()
	if err := m.From(user); err != nil {
		return err
	}
	if err := m.To(email); err != nil {
		return err
	}
	m.Subject("Konfirmasi Pembayaran - Raih Peduli")

	// HTML body
	body := fmt.Sprintf(`
        <html>
            <head>
                <style>
                    /* Add your CSS styles here */
                    body {
                        font-family: Arial, sans-serif;
                    }
                    .confirmation-container {
                        max-width: 600px;
                        margin: 0 auto;
                        padding: 20px;
                        border: 1px solid #ccc;
                        border-radius: 5px;
                    }
                </style>
            </head>
            <body>
                <div class="confirmation-container">
                    <h2>Terima kasih, Orang Baik!</h2>
                    <h3>Detail Penggalangan Dana</h3>
                    <p>Judul: %s</p>
                    <p>Jumlah Donasi: Rp. %d</p>
                    <p>Metode Pembayaran: %s</p>
                </div>
            </body>
        </html>
    `, fundraise.Title, amount, paymentType)

	m.SetBodyString(mail.TypeTextHTML, body)

	c, err := mail.NewClient("smtp.gmail.com", mail.WithPort(convPort), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(user), mail.WithPassword(password))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
