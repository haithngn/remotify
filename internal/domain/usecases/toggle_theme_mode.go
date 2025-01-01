package usecases

import (
	"notification-deployer/internal/domain/values"

	"gorm.io/gorm"
)

func ToggleThemeMode(theme values.ThemeMode, db *gorm.DB) string {
	return updateSetting(values.SettingTypeThemeMode, string(theme), db)
}
