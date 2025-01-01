package usecases

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

func SetAPNSType(mode values.APNSType, db *gorm.DB) string {
	return updateSetting(values.SettingTypeAPNSType, string(mode), db)
}
