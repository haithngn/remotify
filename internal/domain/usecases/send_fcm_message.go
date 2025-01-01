package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"net/http"
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/dto"
	"notification-deployer/internal/domain/values"
)

type SendFCMMessageParams struct {
	DeviceToken string
	DeviceType  values.FCMDeviceType
	PayloadData string
	ToSave      bool
}

func SendFCM(params SendFCMMessageParams, db *gorm.DB) string {
	if params.ToSave == true {
		if _, err := repositories.SaveFCMMessage(dto.FCMMessageDTO{
			DeviceToken: params.DeviceToken,
			DeviceType:  params.DeviceType,
			PayloadData: params.PayloadData,
		}, db); err != nil {
			log.Error("Failed to save FCM message: ", err)
		}
	}

	setting, err := repositories.FindSetting(values.SettingTypeServiceAccount, db)
	if err != nil {
		log.Error(err)
		return values.Failed(errors.New("FCM service account not found"), http.StatusInternalServerError)
	}

	opts := []option.ClientOption{option.WithCredentialsFile(setting.Value)}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Error("Failed to initialize app: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	// Obtain a messaging.Client from the App.
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Error("Failed to get messaging client: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	// See documentation on defining a message payload.
	var message messaging.Message
	err = json.Unmarshal([]byte(params.PayloadData), &message)

	if err != nil {
		log.Error("Failed to unmarshal message: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}
	message.Token = params.DeviceToken

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, &message)
	if err != nil {
		log.Error("Failed to send message: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}
	// Response is a message ID string.

	return values.Succeed(response)
}
