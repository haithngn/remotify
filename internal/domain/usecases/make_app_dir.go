package usecases

import (
	"fmt"
	"os"
)

// MakeAppDirIfNotExist creates a hidden directory in the user's home directory
func MakeAppDirIfNotExist() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	hiddenDirPath := home + "/.remotify"

	// Check if the hidden directory exists
	if _, err := os.Stat(hiddenDirPath); !os.IsNotExist(err) {
		return hiddenDirPath, nil
	}

	err = os.Mkdir(hiddenDirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create hidden directory: %w", err)
	}

	return hiddenDirPath, nil
}
