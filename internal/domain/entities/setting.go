package entities

import (
	"gorm.io/gorm"
	"notification-deployer/internal/domain/values"
)

type Setting struct {
	gorm.Model
	Type  values.SettingType `gorm:"not null;uniqueIndex"`
	Value string             `gorm:"not null"`
}
