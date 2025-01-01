package entities

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

type APNSMessage struct {
	gorm.Model
	DeviceToken string `json:"device_token" binding:"required"`
	BundleID    string `json:"bundle_id" binding:"required"`
	Payload     string `json:"payload" binding:"required"`
	APNSID      *string
	CollapseID  *string
	ExpiredAt   *string
	Priority    values.APNSPriority `gorm:"default:10"`
	PushType    values.APNSPushType `gorm:"default:alert"`
	Note        string              `gorm:"type:text"`
}

func (receiver *APNSMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if receiver.Note == "" {
		receiver.Note = receiver.BundleID
	}
	return
}
