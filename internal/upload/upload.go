package upload

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ToHost(screenshotPath, host string, debug bool) error {
	hostmanPath, err := exec.LookPath("hostman")
	if err != nil {
		return fmt.Errorf("hostman not found in PATH: %v", err)
	}

	var args []string
	if host == "anonhost" || host == "default" || host == "" {
		args = []string{"upload", screenshotPath}
		if debug {
			fmt.Println("Uploading using default host (anonhost)...")
		}
	} else {
		args = []string{"upload", screenshotPath, "--host", host}
		if debug {
			fmt.Printf("Uploading to %s...\n", host)
		}
	}

	if debug {
		fmt.Printf("Command: %s %s\n", hostmanPath, strings.Join(args, " "))
	}

	cmd := exec.Command(hostmanPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("upload failed: exit status %d", exitErr.ExitCode())
		}
		return fmt.Errorf("upload failed: %v", err)
	}

	return nil
}
