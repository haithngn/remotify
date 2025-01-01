package usecases

import (
	"context"
	"errors"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/token"
	"net/http"
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/dto"
	"notification-deployer/internal/domain/values"
	"time"

	"github.com/sideshow/apns2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type APNSMessageParams struct {
	DeviceToken string
	BundleID    string
	Payload     string
	APNSID      *string
	CollapseID  *string
	ExpiredAt   *string
	Priority    values.APNSPriority
	PushType    values.APNSPushType
	ToSave      bool
}

func SendAPNS(params APNSMessageParams, db *gorm.DB) string {
	dtoMsg := dto.APNSMessageDTO{
		DeviceToken: params.DeviceToken,
		BundleID:    params.BundleID,
		Payload:     params.Payload,
		APNSID:      params.APNSID,
		CollapseID:  params.CollapseID,
		ExpiredAt:   params.ExpiredAt,
		Priority:    params.Priority,
		PushType:    params.PushType,
	}

	// ToSave APNS message into sending history
	if params.ToSave == true {
		if _, err := repositories.SaveAPNSMessage(dtoMsg, db); err != nil {
			log.Error("Failed to save APNS message: ", err)
		}
	}

	apnsInUseSetting, err := repositories.FindSetting(values.SettingTypeAPNSType, db)
	if err != nil {
		log.Error("Failed to get APNS setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	switch apnsInUseSetting.Value {
	case string(values.APNSTypeTokenBase):
		return sendAPNSMessageInJWTBase(params, db)
	case string(values.APNSTypeLegacy):
		return sendAPNSMessageInLegacy(params, db)
	default:
		return values.Failed(errors.New("Unknown APNS type: "+apnsInUseSetting.Value), http.StatusInternalServerError)
	}
}

func sendAPNSMessageInLegacy(params APNSMessageParams, db *gorm.DB) string {
	expiredTime := time.Time{}
	if params.ExpiredAt != nil {
		t, err := time.Parse("2006-01-02T15:04", *params.ExpiredAt)
		if err != nil {
			return values.Failed(errors.New("Invalid expired time : "+*params.ExpiredAt), http.StatusBadRequest)
		}
		expiredTime = t
	}

	priority := apns2.PriorityHigh
	switch params.Priority {
	case values.APNSPriorityImmediately:
		priority = apns2.PriorityHigh
	case values.APNSPriorityThrottled:
		priority = apns2.PriorityLow
	}

	pushType := apns2.PushTypeAlert
	switch params.PushType {

	case values.APNSPushTypeAlert:
		pushType = apns2.PushTypeAlert
	case values.APNSPushTypeBackground:
		pushType = apns2.PushTypeBackground
	case values.APNSPushTypeLocation:
		pushType = apns2.PushTypeLocation
	case values.APNSPushTypeVoIP:
		pushType = apns2.PushTypeVOIP
	case values.APNSPushTypeComplication:
		pushType = apns2.PushTypeComplication
	case values.APNSPushTypeFileProvider:
		pushType = apns2.PushTypeFileProvider
	case values.APNSPushTypeMDM:
		pushType = apns2.PushTypeMDM
	}

	notification := &apns2.Notification{
		DeviceToken: params.DeviceToken,
		Topic:       params.BundleID,
		Payload:     []byte(params.Payload),
		Priority:    priority,
		PushType:    pushType,
		Expiration:  expiredTime,
	}

	if params.APNSID != nil {
		notification.ApnsID = *params.APNSID
	}
	if params.CollapseID != nil {
		notification.CollapseID = *params.CollapseID
	}

	apnsEnvSetting, err := repositories.FindSetting(values.SettingTypeAPNSLegacyEnvironment, db)
	if err != nil {
		log.Error("Failed to get apnsEnvSetting setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	notification.PushType = apns2.EPushType(apnsEnvSetting.Value)

	certSetting, err := repositories.FindSetting(values.SettingTypeAppleCertificateFilePath, db)
	if err != nil {
		log.Error("Failed to get APNS setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	decryptPasswordSetting, err := repositories.FindSetting(values.SettingTypeAppleCertificateDecryptPassword, db)
	if err != nil {
		log.Error("Failed to get apnsEnvSetting setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	cert, err := certificate.FromP12File(certSetting.Value, decryptPasswordSetting.Value)
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	var client *apns2.Client

	//Sandbox or Production
	if apnsEnvSetting.Value == string(values.APNSEnvironmentSandbox) {
		client = apns2.NewClient(cert).Development()
	} else {
		client = apns2.NewClient(cert).Production()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()
	res, err := client.PushWithContext(ctx, notification)

	if err != nil {
		log.Error("Failed to send APNS message: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	if res.Sent() == false {
		log.Error("Failed to send APNS message: ", res, " push type: ", pushType)
		return values.Failed(errors.New("Failed to send APNS message: "+string(res.Reason)), res.StatusCode)
	}

	return values.Succeed(res)
}

func sendAPNSMessageInJWTBase(params APNSMessageParams, db *gorm.DB) string {
	var expiredTime *time.Time
	if params.ExpiredAt != nil {
		t, err := time.Parse("2006-01-02T15:04", *params.ExpiredAt)
		if err != nil {
			return values.Failed(errors.New("Invalid expired time : "+*params.ExpiredAt), http.StatusBadRequest)
		}
		expiredTime = &t
	}

	priority := apns2.PriorityHigh
	switch params.Priority {
	case values.APNSPriorityImmediately:
		priority = apns2.PriorityHigh
	case values.APNSPriorityThrottled:
		priority = apns2.PriorityLow
	}

	pushType := apns2.PushTypeAlert
	switch params.PushType {
	case values.APNSPushTypeAlert:
		pushType = apns2.PushTypeAlert
	case values.APNSPushTypeBackground:
		pushType = apns2.PushTypeBackground
	case values.APNSPushTypeLocation:
		pushType = apns2.PushTypeLocation
	case values.APNSPushTypeVoIP:
		pushType = apns2.PushTypeVOIP
	case values.APNSPushTypeComplication:
		pushType = apns2.PushTypeComplication
	case values.APNSPushTypeFileProvider:
		pushType = apns2.PushTypeFileProvider
	case values.APNSPushTypeMDM:
		pushType = apns2.PushTypeMDM
	}

	notification := &apns2.Notification{
		DeviceToken: params.DeviceToken,
		Topic:       params.BundleID,
		Payload:     []byte(params.Payload),
		Priority:    priority,
		PushType:    pushType,
	}

	if expiredTime != nil {
		notification.Expiration = *expiredTime
	}

	if params.APNSID != nil {
		notification.ApnsID = *params.APNSID
	}
	if params.CollapseID != nil {
		notification.CollapseID = *params.CollapseID
	}

	jwtSetting, err := repositories.FindSetting(values.SettingTypeAPNSJWTFilePath, db)
	if err != nil {
		log.Error("Failed to get JWT setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	keyIDSetting, err := repositories.FindSetting(values.SettingTypeAppleDeveloperKeyID, db)
	if err != nil {
		log.Error("Failed to get Key ID setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	teamIDSetting, err := repositories.FindSetting(values.SettingTypeAppleDeveloperTeamID, db)
	if err != nil {
		log.Error("Failed to get Team ID setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	authKey, err := token.AuthKeyFromFile(jwtSetting.Value)
	if err != nil {
		log.Error("Failed to create JWT client: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	apnsToken := &token.Token{
		AuthKey: authKey,
		KeyID:   keyIDSetting.Value,
		TeamID:  teamIDSetting.Value,
	}

	jwtEnvSetting, err := repositories.FindSetting(values.SettingTypeJWTAPNSEnvironment, db)
	if err != nil {
		log.Error("Failed to get JWT environment setting: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	var client *apns2.Client
	if jwtEnvSetting.Value == string(values.APNSEnvironmentSandbox) {
		client = apns2.NewTokenClient(apnsToken).Development()
	} else {
		client = apns2.NewTokenClient(apnsToken).Production()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	res, err := client.PushWithContext(ctx, notification)

	if err != nil {
		log.Error("Failed to send: ", err)
		return values.Failed(err, http.StatusInternalServerError)
	}

	if res.Sent() == false {
		log.Error("APNS failed: ", res.Reason, " pushType: ", pushType, " env type: ", jwtEnvSetting.Value)
		return values.Failed(errors.New(res.Reason), res.StatusCode)
	}

	return values.Succeed(res)
}
