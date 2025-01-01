package usecases

import (
	"notification-deployer/internal/data/repositories"
	"notification-deployer/internal/domain/values"

	"gorm.io/gorm"
)

type appSettings struct {
	ServiceAccount  string `json:"service_account" binding:"required"`
	P8              string `json:"p8" binding:"required"`
	TeamId          string `json:"team_id" binding:"required"`
	KeyId           string `json:"key_id" binding:"required"`
	P12             string `json:"p12" binding:"required"`
	DecryptPassword string `json:"decrypt_password" binding:"required"`
	APNSMode        string `json:"apns_mode" binding:"required"`
	JWTAPNSMode     string `json:"jwt_apns_mode" binding:"required"`
	APNSType        string `json:"apns_type" binding:"required"`
	ThemeMode       string `json:"theme_mode" binding:"required"`
}

func GetSettings(db *gorm.DB) string {
	settings := appSettings{
		ServiceAccount:  getSettingValueOrEmpty(values.SettingTypeServiceAccount, db),
		P8:              getSettingValueOrEmpty(values.SettingTypeAPNSJWTFilePath, db),
		TeamId:          getSettingValueOrEmpty(values.SettingTypeAppleDeveloperTeamID, db),
		KeyId:           getSettingValueOrEmpty(values.SettingTypeAppleDeveloperKeyID, db),
		P12:             getSettingValueOrEmpty(values.SettingTypeAppleCertificateFilePath, db),
		DecryptPassword: getSettingValueOrEmpty(values.SettingTypeAppleCertificateDecryptPassword, db),
		APNSMode:        getSettingValueOrEmpty(values.SettingTypeAPNSLegacyEnvironment, db),
		JWTAPNSMode:     getSettingValueOrEmpty(values.SettingTypeJWTAPNSEnvironment, db),
		APNSType:        getSettingValueOrEmpty(values.SettingTypeAPNSType, db),
		ThemeMode:       getSettingValueOrEmpty(values.SettingTypeThemeMode, db),
	}
	return values.Succeed(settings)
}

func getSettingValueOrEmpty(settingType values.SettingType, db *gorm.DB) string {
	serviceAccountSetting, err := repositories.FindSetting(settingType, db)

	if err != nil {
		return ""
	}

	return serviceAccountSetting.Value
}

func getSettingValueInBooleanOrFalse(settingType values.SettingType, db *gorm.DB) bool {
	serviceAccountSetting, err := repositories.FindSetting(settingType, db)

	if err != nil {
		return false
	}

	return serviceAccountSetting.Value == "true"
}
