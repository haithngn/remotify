package usecases

import (
	"errors"
	"net/http"
	"notification-deployer/internal/domain/values"

	"gorm.io/gorm"
)

// UpdateAPNSCertificate opens a file dialog for selecting a .p12 file and returns the selected file path.
//
// Parameters:
// - ctx: the context.Context for the function.
//
// Returns:
// - string: the file path of the selected .p12 file, or a values.Failed struct if an error occurs.
func UpdateAPNSCertificate(f string, db *gorm.DB) string {
	if !validateFile(f) {
		return values.Failed(errors.New("invalid file"), http.StatusBadRequest)
	}

	return updateSetting(values.SettingTypeAppleCertificateFilePath, f, db)
}
