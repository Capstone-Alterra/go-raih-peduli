package helpers

import (
	"context"
	"raihpeduli/utils"

	"firebase.google.com/go/v4/messaging"
	"github.com/sirupsen/logrus"
)

type notificationService struct{}

func NewNotificationService() NotificationInterface {
	return &notificationService{}
}

func (ns *notificationService) SendNotifications(tokens string, title string, message string) error {
	fcmClient := utils.FirebaseInit()

	_, err := fcmClient.Send(context.Background(), &messaging.Message{
		Token: tokens,
		Data: map[string]string{
			"Title":   title,
			"Message": message,
		},
	})
	logrus.Info(tokens)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return err
}
