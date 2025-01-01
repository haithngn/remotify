package entities

import (
	"encoding/json"
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

type JSON json.RawMessage

type FCMMessage struct {
	gorm.Model
	DeviceToken string               `gorm:"not null"`
	DeviceType  values.FCMDeviceType `gorm:"not null"`
	Payload     JSON                 `gorm:"not null"`
	Note        string               `gorm:"type:text"`
}

func (receiver *FCMMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if receiver.Note == "" {
		receiver.Note = receiver.DeviceToken
	}
	return
}
