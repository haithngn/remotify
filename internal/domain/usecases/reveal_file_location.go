package usecases

import (
	"fmt"
	"notification-deployer/internal/domain/values"
	"os/exec"
	"path/filepath"
)

func RevealFileLocation(file string) string {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("open %s", filepath.Dir(file)))
	_, err := cmd.Output()

	if err != nil {
		return fmt.Sprintf(`{
				"status_code": 500,
				"body": {
					"message": "Failed to reveal file location",
				}
		}`)
	}

	return values.Succeed("File location revealed")
}
