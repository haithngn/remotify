package dto

import (
	"notification-deployer/internal/domain/entities"
	"notification-deployer/internal/domain/values"
)

type APNSMessageDTO struct {
	ID          uint                `json:"id" binding:"required"`
	DeviceToken string              `json:"device_token" binding:"required"`
	BundleID    string              `json:"bundle_id" binding:"required"`
	Payload     string              `json:"payload" binding:"required"`
	APNSID      *string             `json:"apns_id"`
	CollapseID  *string             `json:"collapse_id"`
	ExpiredAt   *string             `json:"expired_at"`
	Priority    values.APNSPriority `json:"priority"`
	PushType    values.APNSPushType `json:"push_type"`
	Note        string              `json:"note" binding:"required"`
	CreatedAt   string              `json:"created_at"`
}

func APNSMessageDTOs(messages []entities.APNSMessage) []APNSMessageDTO {
	var result []APNSMessageDTO
	for _, msg := range messages {
		result = append(result, APNSMessageDTO{
			ID:          msg.ID,
			DeviceToken: msg.DeviceToken,
			BundleID:    msg.BundleID,
			Payload:     msg.Payload,
			APNSID:      msg.APNSID,
			CollapseID:  msg.CollapseID,
			ExpiredAt:   msg.ExpiredAt,
			Priority:    msg.Priority,
			PushType:    msg.PushType,
			Note:        msg.Note,
			CreatedAt:   msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result
}
