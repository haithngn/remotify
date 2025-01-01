package usecases

import (
	"errors"
	"net/http"
	"notification-deployer/internal/domain/values"

	"gorm.io/gorm"
)

func UpdateAPNSToken(f string, db *gorm.DB) string {
	if validateFile(f) == false {
		return values.Failed(errors.New("invalid file"), http.StatusBadRequest)
	}

	return updateSetting(values.SettingTypeAPNSJWTFilePath, f, db)
}
