package upload

import (
	"fmt"
	"os"
	"os/exec"
)

func ToHost(screenshotPath, host string) error {
	fmt.Printf("Uploading screenshot to %s...\n", host)
	
	cmd := exec.Command("hostman", "upload", screenshotPath, "--host", host)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hostman upload failed: %v", err)
	}

	return nil
}