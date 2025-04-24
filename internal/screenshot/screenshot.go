// internal/screenshot/screenshot.go
package screenshot

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

func Capture(screenshotType, outputDir string) (string, error) {
    timestamp := time.Now().Format("20060102-150405")
    filename := fmt.Sprintf("screenshot-%s-%s.png", screenshotType, timestamp)
    outputPath := filepath.Join(outputDir, filename)

    cmd := exec.Command("hyprshot", "-m", screenshotType, "-o", outputPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("hyprshot command failed: %v", err)
    }

    if _, err := os.Stat(outputPath); os.IsNotExist(err) {
        return "", fmt.Errorf("screenshot file was not created at: %s", outputPath)
    }

    return outputPath, nil
}
