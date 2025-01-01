package dto

import (
	"notification-deployer/internal/domain/entities"
	"notification-deployer/internal/domain/values"
)

type FCMMessageDTO struct {
	ID          uint                 `json:"id" binding:"required"`
	DeviceToken string               `json:"device_token" binding:"required"`
	DeviceType  values.FCMDeviceType `json:"device_type" binding:"required"`
	Title       string               `json:"title" binding:"required"`
	Message     string               `json:"message" binding:"required"`
	ImageURL    string               `json:"image_url" binding:"required"`
	PayloadData string               `json:"payload_data" binding:"required"`
	Config      string               `json:"config" binding:"required"`
	Note        string               `json:"note" binding:"required"`
	CreatedAt   string               `json:"created_at" binding:"required"`
}

func FCMMessageDTOs(messages []entities.FCMMessage) []FCMMessageDTO {
	var result []FCMMessageDTO
	for _, msg := range messages {
		result = append(result, FCMMessageDTO{
			ID:          msg.ID,
			DeviceToken: msg.DeviceToken,
			DeviceType:  msg.DeviceType,
			PayloadData: string(msg.Payload),
			Note:        msg.Note,
			CreatedAt:   msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result
}
