package repository

import (
	"context"
	"raihpeduli/config"
	"raihpeduli/features/auth"
	"raihpeduli/helpers"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type model struct {
	db         *gorm.DB
	connection *redis.Client
	config     *config.SMTPConfig
	collection *mongo.Collection
}

func New(db *gorm.DB, rdClient *redis.Client, config *config.SMTPConfig, collection *mongo.Collection) auth.Repository {
	return &model{
		db:         db,
		connection: rdClient,
		config:     config,
		collection: collection,
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

func (mdl *model) SendOTPByEmail(fullname, email, otp, status string) error {
	err := helpers.EmailService(fullname, email, otp, status)
	if err != nil {
		return err
	}
	return nil
}

func (mdl *model) InsertToken(userID int, fcmToken string) error {
	filter := bson.M{"user_id": userID}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	set := bson.M{"$set": bson.M{"user_id": userID, "device_token": fcmToken, "timestamp": time.Now().UTC()}}

	if result := mdl.collection.FindOneAndUpdate(context.Background(), filter, set, opts); result.Err() != nil {
		logrus.Error(result.Err())
	}

	return nil
}