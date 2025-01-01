package repositories

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/dto"
	"notification-deployer/internal/domain/entities"
)

func SaveFCMMessage(message dto.FCMMessageDTO, db *gorm.DB) (dto.FCMMessageDTO, error) {
	msg := entities.FCMMessage{
		DeviceToken: message.DeviceToken,
		DeviceType:  message.DeviceType,
		Payload:     entities.JSON(message.PayloadData),
	}

	err := db.Create(&msg).Error
	if err != nil {
		return dto.FCMMessageDTO{}, err
	}

	return dto.FCMMessageDTO{
		ID:          msg.ID,
		DeviceToken: msg.DeviceToken,
		DeviceType:  msg.DeviceType,
		PayloadData: string(msg.Payload),
	}, nil
}

func FindFCMMessages(top int, db *gorm.DB) ([]dto.FCMMessageDTO, error) {
	var messages []entities.FCMMessage

	err := db.
		Order("created_at desc").
		Limit(top).
		Find(&messages).
		Error

	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return []dto.FCMMessageDTO{}, nil
	}

	return dto.FCMMessageDTOs(messages), nil
}

func SaveFCMMessageNote(note string, id uint, db *gorm.DB) error {
	return db.Model(&entities.FCMMessage{}).
		Where("id = ?", id).
		Update("note", note).
		Error
}

func RemoveFCMMessage(id uint, db *gorm.DB) error {
	return db.Delete(&entities.FCMMessage{}, id).Error
}
