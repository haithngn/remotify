package usecases

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

func SetAppleDeveloperTeamID(teamID string, db *gorm.DB) string {
	return updateSetting(values.SettingTypeAppleDeveloperTeamID, teamID, db)
}
