package usecases

import (
	"errors"
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/values"
	"os"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func updateSetting(settingType values.SettingType, value string, db *gorm.DB) string {
	err := repositories.SaveSetting(settingType, value, db)
	if err != nil {
		log.Error("Failed to save setting: ", settingType, " error ", err)
	}

	return values.Succeed(value)
}

func validateFile(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
