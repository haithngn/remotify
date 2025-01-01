package usecases

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

func SetAppleCertificateDecryptPassword(id string, db *gorm.DB) string {
	return updateSetting(values.SettingTypeAppleCertificateDecryptPassword, id, db)
}
