package repository

import (
	"raihpeduli/config"
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
	config     *config.SMTPConfig
}

func New(db *gorm.DB, rdClient *redis.Client, config *config.SMTPConfig) auth.Repository {
	return &model{
		db:         db,
		connection: rdClient,
		config:     config,
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
	user := mdl.config.SMTP_USER
	password := mdl.config.SMTP_PASS
	port := mdl.config.SMTP_PORT

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
	m.Subject("Verifikasi Email - Raih Peduli")
	m.SetBodyString(mail.TypeTextPlain, "Kode OTP anda adalah : "+otp)

	c, err := mail.NewClient("smtp.gmail.com", mail.WithPort(convPort), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(user), mail.WithPassword(password))
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
