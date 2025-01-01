package usecases

import (
	"gorm.io/gorm"
	"net/http"
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/values"
)

func SaveAPNSMessageNote(note string, messageID uint, db *gorm.DB) string {
	err := repositories.SaveAPNSMessageNote(note, messageID, db)
	if err != nil {
		return values.Failed(err, http.StatusInternalServerError)
	}

	return values.Succeed(note)
}
