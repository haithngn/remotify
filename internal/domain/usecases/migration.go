package usecases

import (
	"gorm.io/gorm"
	"notification-deployer/internal/data/persistents"
)

func Migration(db *gorm.DB) error {
	return persistents.Migration(db)
}
