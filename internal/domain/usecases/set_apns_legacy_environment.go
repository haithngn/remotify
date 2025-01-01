package usecases

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

func SetLegacyAPNSEnvironment(mode values.APNSEnvironment, db *gorm.DB) string {
	return updateSetting(values.SettingTypeAPNSLegacyEnvironment, string(mode), db)
}

func SetJWTAPNSEnvironment(mode values.APNSEnvironment, db *gorm.DB) string {
	return updateSetting(values.SettingTypeJWTAPNSEnvironment, string(mode), db)
}
