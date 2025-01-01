package usecases

import (
	"gorm.io/gorm"
	"net/http"
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/values"
)

func GetRecentAPNSMessages(db *gorm.DB) string {
	messages, err := repositories.FindAPNSMessages(values.RecentItemsLimit, db)
	if err != nil {
		return values.Failed(err, http.StatusInternalServerError)
	}

	return values.Succeed(messages)
}
