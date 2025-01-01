package usecases

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

func SetAppleDeveloperKeyID(id string, db *gorm.DB) string {
	return updateSetting(values.SettingTypeAppleDeveloperKeyID, id, db)
}
