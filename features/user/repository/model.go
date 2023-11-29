package repository

import (
	"mime/multipart"
	"os"
	"raihpeduli/config"
	"raihpeduli/features/user"
	"raihpeduli/helpers"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/wneessen/go-mail"
	"gorm.io/gorm"
)

type model struct {
	db         *gorm.DB
	connection *redis.Client
	clStorage  helpers.CloudStorageInterface
}

func New(db *gorm.DB, rdClient *redis.Client, clStorage helpers.CloudStorageInterface) user.Repository {
	return &model{
		db:         db,
		connection: rdClient,
		clStorage:  clStorage,
	}
}

func (mdl *model) Paginate(page, size int) []user.User {
	var users []user.User

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&users)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return users
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

func (mdl *model) ValidateVerification(verificationKey string) string {
	email, statusCMD := mdl.connection.Get(verificationKey).Result()
	if statusCMD != nil {
		logrus.Error(statusCMD.Error())
		return ""
	}

	_, err := mdl.connection.Del(verificationKey).Result()
	if err != nil {
		return ""
	}
	return email
}

func (mdl *model) InsertUser(newUser *user.User) (*user.User, error) {
	result := mdl.db.Table("users").Create(newUser)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return newUser, nil
}

func (mdl *model) SelectByID(userID int) *user.User {
	var user user.User
	result := mdl.db.First(&user, userID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &user
}

func (mdl *model) SelectByEmail(email string) (*user.User, error) {
	var user user.User
	result := mdl.db.Table("users").Where("email = ?", email).First(&user)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (mdl *model) UpdateUser(user user.User) int64 {
	result := mdl.db.Updates(&user)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(userID int) int64 {
	result := mdl.db.Delete(&user.User{}, userID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
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

func (mdl *model) GetTotalData() int64 {
	var totalData int64

	result := mdl.db.Table("users").Where("deleted_at IS NULL").Count(&totalData)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return totalData
}

func (mdl *model) UploadFile(file multipart.File, oldFilename string) (string, error) {
	var config = config.LoadCloudStorageConfig()
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/users/")
	var objectName string

	if file == nil {
		return oldFilename, nil
	}

	if oldFilename != "" {
		objectName = oldFilename[urlLength:]

		if objectName == "default" {
			objectName = ""
		} else if err := mdl.clStorage.DeleteFile(objectName); err != nil {
			return "", err
		}
	}
	objectName = uuid.New().String()

	if err := mdl.clStorage.UploadFile(file, objectName); err != nil {
		return "", err
	}

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/users/" + objectName, nil
}
