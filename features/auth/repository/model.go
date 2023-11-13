package repository

import (
	"os"
	"raihpeduli/features/auth"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/wneessen/go-mail"
	"gorm.io/gorm"
)

type model struct {
	db         *gorm.DB
	connection *redis.Client
}

func New(db *gorm.DB, rdClient *redis.Client) auth.Repository {
	return &model{
		db:         db,
		connection: rdClient,
	}
}

func (mdl *model) Login(email string) (*auth.User, error) {
	var user auth.User
	result := mdl.db.Table("users").Where("email = ? AND is_verified = ?", email, 1).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		log.Error(result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (mdl *model) SelectByEmail(email string) (*auth.User, error) {
	var user auth.User
	result := mdl.db.Table("users").Where("email = ?", email).First(&user)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (mdl *model) GetNameAdmin(id int) (string, error) {
	var fullname string

	result := mdl.db.Table("admins").Where("user_id = ?", id).Select("fullname").Scan(&fullname)
	if result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return fullname, nil
}

func (mdl *model) GetNameCustomer(id int) (string, error) {
	var fullname string

	result := mdl.db.Table("customers").Where("user_id = ?", id).Select("fullname").Scan(&fullname)
	if result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return fullname, nil
}

func (mdl *model) Register(newUser *auth.User) (*auth.User, error) {
	result := mdl.db.Table("users").Create(newUser)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return newUser, nil
}

func (mdl *model) InsertVerification(email string, verificationKey string) error {
	statusCMD := mdl.connection.Set(verificationKey, email, time.Minute*10)
	if statusCMD.Err() != nil {
		logrus.Error(statusCMD.Err())
		return statusCMD.Err()
	}

	res, err := statusCMD.Result()
	if err != nil {
		logrus.Error(err.Error())
	}

	logrus.Info("OTP Inserted to Redis ", res)

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

	query := mdl.InsertVerification(email, otp)
	if query != nil {
		return query
	}

	return nil
}
