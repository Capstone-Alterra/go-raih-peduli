package helpers

import (
	"context"
	"raihpeduli/utils"

	"firebase.google.com/go/v4/messaging"
)

type notificationService struct{}

func NewNotificationService() NotificationInterface {
	return &notificationService{}
}

func (ns *notificationService) SendNotifications(tokens string, userID string, message string) error {
	fcmClient := utils.FirebaseInit()

	_, err := fcmClient.Send(context.Background(), &messaging.Message{
		Token: tokens,
		Data: map[string]string{
			message: message,
		},
	})
	if err != nil {
		return err
	}

	return err
}
