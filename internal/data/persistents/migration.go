package persistents

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"notification-deployer/internal/domain/entities"
	"notification-deployer/internal/domain/values"
)

func Migration(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entities.APNSMessage{},
		&entities.FCMMessage{},
		&entities.Setting{},
	)

	if err != nil {
		return err
	}

	history := entities.Setting{
		Type:  values.SettingTypeHistoryLength,
		Value: "100",
	}

	apns_mode := entities.Setting{
		Type:  values.SettingTypeAPNSType,
		Value: string(values.APNSTypeTokenBase),
	}

	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "type"}},
		DoNothing: true,
	}).Save(&history).Save(&apns_mode).Error

	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) == false {
		return err
	}

	return err
}
