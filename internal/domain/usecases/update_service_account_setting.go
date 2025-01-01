package usecases

import (
	"notification-deployer/internal/domain/values"

	"gorm.io/gorm"
)

func UpdateServiceAccount(f string, db *gorm.DB) string {
	return updateSetting(values.SettingTypeServiceAccount, f, db)
}
