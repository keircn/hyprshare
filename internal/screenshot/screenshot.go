package screenshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/q4ow/hyprshare/internal/cli"
)

func Capture(opts cli.Options) (string, error) {
	// Generate our own unique filepath
	timestamp := time.Now().Format("20060102-150405")
	uniqueID := fmt.Sprintf("%d", time.Now().UnixNano()%10000)
	filename := fmt.Sprintf("hyprshare-%s-%s.png", timestamp, uniqueID)

	var outputPath string
	if opts.ClipboardOnly {
		outputPath = filepath.Join(os.TempDir(), filename)
	} else if opts.Filename != "" {
		outputPath = filepath.Join(opts.OutputFolder, opts.Filename)
	} else {
		outputPath = filepath.Join(opts.OutputFolder, filename)
	}

	// Make sure output directory exists
	if !opts.ClipboardOnly {
		if err := os.MkdirAll(opts.OutputFolder, 0755); err != nil {
			return "", fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Find hyprshot in PATH
	hypshotPath, err := exec.LookPath("hyprshot")
	if err != nil {
		return "", fmt.Errorf("hyprshot not found in PATH: %v", err)
	}

	// Prepare arguments for hyprshot
	args := []string{}
	for _, mode := range opts.Mode {
		args = append(args, "-m", mode)
	}

	// Always specify our output path
	args = append(args, "-o", filepath.Dir(outputPath))
	args = append(args, "-f", filepath.Base(outputPath))

	if opts.Freeze {
		args = append(args, "-z")
	}
	if opts.Debug {
		args = append(args, "-d")
	}
	if opts.Silent {
		args = append(args, "-s")
	}
	if opts.Raw {
		args = append(args, "-r")
	}
	args = append(args, "-t", strconv.Itoa(opts.NotifTimeout))

	if opts.ClipboardOnly {
		args = append(args, "--clipboard-only")
	}

	if len(opts.Command) > 0 {
		args = append(args, "--")
		args = append(args, opts.Command...)
	}

	// Print command for debugging
	if opts.Debug {
		fmt.Printf("Command: %s %s\n", hypshotPath, strings.Join(args, " "))
	}

	// Run hyprshot
	cmd := exec.Command(hypshotPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // Ignore error since hyprshot might exit with non-zero code

	// Give hyprshot a moment to finish writing the file
	time.Sleep(200 * time.Millisecond)

	// For clipboard-only mode, we're done
	if opts.ClipboardOnly {
		return "", nil
	}

	// For regular mode, check if the file exists
	info, statErr := os.Stat(outputPath)
	if statErr != nil {
		if opts.Debug {
			fmt.Printf("Debug: file check error at %s: %v\n", outputPath, statErr)
		}

		// Try looking in the directory to see what's there
		files, _ := filepath.Glob(filepath.Join(opts.OutputFolder, "hyprshare-*"))
		if len(files) > 0 {
			// Sort files by modification time to find the most recent
			var newestFile string
			var newestTime time.Time

			for _, file := range files {
				fileInfo, err := os.Stat(file)
				if err == nil && (newestFile == "" || fileInfo.ModTime().After(newestTime)) {
					newestTime = fileInfo.ModTime()
					newestFile = file
				}
			}

			if newestFile != "" {
				if opts.Debug {
					fmt.Printf("Debug: found alternative file %s (modified %s)\n",
						newestFile, newestTime.Format(time.RFC3339))
				}
				return newestFile, nil
			}
		}

		return "", fmt.Errorf("screenshot not found at %s", outputPath)
	}

	if opts.Debug {
		fmt.Printf("Debug: found file at %s (size: %d bytes)\n", outputPath, info.Size())
	}
	return outputPath, nil
}
