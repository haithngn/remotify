package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"notification-deployer/internal/domain/entities"
	"notification-deployer/internal/domain/values"
)

func FindSetting(settingType values.SettingType, db *gorm.DB) (entities.Setting, error) {
	var setting entities.Setting
	err := db.Where("type = ?", settingType).First(&setting).Error

	return setting, err
}

func SaveSetting(settingType values.SettingType, value string, db *gorm.DB) error {
	setting := entities.Setting{
		Type:  settingType,
		Value: value,
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Save(&setting).Error
}
