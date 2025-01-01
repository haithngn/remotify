package repositories

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/dto"
	"notification-deployer/internal/domain/entities"
)

func SaveAPNSMessage(messageDTO dto.APNSMessageDTO, db *gorm.DB) (dto.APNSMessageDTO, error) {
	msg := entities.APNSMessage{
		DeviceToken: messageDTO.DeviceToken,
		BundleID:    messageDTO.BundleID,
		Payload:     messageDTO.Payload,
		APNSID:      messageDTO.APNSID,
		CollapseID:  messageDTO.CollapseID,
		ExpiredAt:   messageDTO.ExpiredAt,
		Priority:    messageDTO.Priority,
		PushType:    messageDTO.PushType,
	}

	err := db.Create(&msg).Error
	if err != nil {
		return messageDTO, err
	}

	return dto.APNSMessageDTO{
		ID:          msg.ID,
		DeviceToken: msg.DeviceToken,
		BundleID:    msg.BundleID,
		Payload:     msg.Payload,
	}, nil
}

func FindAPNSMessages(top int, db *gorm.DB) ([]dto.APNSMessageDTO, error) {
	var messages []entities.APNSMessage

	err := db.Order("updated_at desc").
		Limit(top).
		Find(&messages).
		Error

	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return []dto.APNSMessageDTO{}, nil
	}

	return dto.APNSMessageDTOs(messages), nil
}

func RemoveAPNSMessage(id uint, db *gorm.DB) error {
	return db.Where("id = ?", id).Delete(&entities.APNSMessage{}).Error
}

func SaveAPNSMessageNote(note string, id uint, db *gorm.DB) error {
	return db.Model(&entities.APNSMessage{}).
		Where("id = ?", id).
		Update("note", note).
		Error
}
