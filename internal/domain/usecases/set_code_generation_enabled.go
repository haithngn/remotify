package usecases

import (
	"notification-deployer/internal/domain/values"

	"gorm.io/gorm"
)

func SetCodeGenerationEnabled(enabled bool, db *gorm.DB) string {
	value := "false"
	if enabled {
		value = "true"
	}

	return updateSetting(values.SettingTypeCodeGenerationEnabled, value, db)
}
