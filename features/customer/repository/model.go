package repository

import (
	"os"
	"raihpeduli/features/customer"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/wneessen/go-mail"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) customer.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []customer.Customer {
	var customers []customer.Customer

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&customers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return customers
}

func (mdl *model) InsertCustomer(newCustomer *customer.Customer) (*customer.Customer, error) {
	result := mdl.db.Table("customers").Create(newCustomer)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return newCustomer, nil
}

func (mdl *model) InsertUser(newUser *customer.User) (*customer.User, error) {
	result := mdl.db.Table("users").Create(newUser)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return newUser, nil
}

func (mdl *model) InsertOTP(otp *customer.OTP) error {
	result := mdl.db.Table("otps").Create(otp)

	if result.Error != nil {
		log.Error(result.Error)
		return result.Error
	}

	return nil
}

func (mdl *model) SelectByID(customerID int) *customer.Customer {
	var customer customer.Customer
	result := mdl.db.First(&customer, customerID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &customer
}

func (mdl *model) SelectByEmail(email string) (*customer.User, error) {
	var user customer.User
	result := mdl.db.Table("users").Where("email = ?", email).First(&user)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (mdl *model) SelectOTP(userID int, otp string) (*customer.OTP, error) {
	var OTP customer.OTP
	result := mdl.db.Table("otps").Where("user_id = ? AND otp = ?", userID, otp).First(&OTP)
	if result == nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return &OTP, nil
}

func (mdl *model) UpdateCustomer(customer customer.Customer) int64 {
	result := mdl.db.Save(&customer)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) UpdateUser(user customer.User) error {
	result := mdl.db.Table("users").Save(&user)

	if result.Error != nil {
		log.Error(result.Error)
		return result.Error
	}

	return nil
}

func (mdl *model) DeleteByID(customerID int) int64 {
	result := mdl.db.Delete(&customer.Customer{}, customerID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}

func (mdl *model) DeleteOTP(otp customer.OTP) error {
	result := mdl.db.Delete(&otp)

	if result.Error != nil {
		log.Error(result.Error)
		return result.Error
	}

	return nil
}

func (mdl *model) SendOTPByEmail(email string, otp string) error {
	secret_user := os.Getenv("SMTP_USER")
	secret_pass := os.Getenv("SMTP_PASS")
	secret_port := os.Getenv("SMTP_PORT")

	convPort, err := strconv.Atoi(secret_port)
	if err != nil {
		return err
	}

	m := mail.NewMsg()
	if err := m.From(secret_user); err != nil {
		return err
	}
	if err := m.To(email); err != nil {
		return err
	}
	m.Subject("Verifikasi Email - Raih Peduli")
	m.SetBodyString(mail.TypeTextPlain, "Kode OTP anda adalah : "+otp)

	c, err := mail.NewClient("smtp.gmail.com", mail.WithPort(convPort), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(secret_user), mail.WithPassword(secret_pass))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
