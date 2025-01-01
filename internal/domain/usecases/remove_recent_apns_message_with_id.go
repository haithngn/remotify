package usecases

import (
	"gorm.io/gorm"
	"net/http"
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/values"
)

func RemoveRecentAPNSMessage(id uint, db *gorm.DB) interface{} {
	err := repositories.RemoveAPNSMessage(id, db)

	if err != nil {
		return values.Failed(err, http.StatusInternalServerError)
	}

	return values.Succeed(nil)
}
