package upload

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ToHost(screenshotPath, primaryHost, fallbackHost string, debug bool) error {
	err := uploadToSingleHost(screenshotPath, primaryHost, debug)
	if err == nil {
		return nil
	}

	if fallbackHost != "" {
		fmt.Printf("Primary host failed, trying fallback host: %s\n", fallbackHost)
		return uploadToSingleHost(screenshotPath, fallbackHost, debug)
	}

	return err
}

func uploadToSingleHost(screenshotPath, host string, debug bool) error {
	fmt.Printf("Uploading to %s...\n", host)

	hostmanPath, err := exec.LookPath("hostman")
	if err != nil {
		return fmt.Errorf("hostman not found in PATH: %v", err)
	}

	args := []string{"upload", screenshotPath, "--host", host}

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
