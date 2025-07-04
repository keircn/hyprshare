package screenshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/keircn/hyprshare/internal/cli"
)

func Capture(opts cli.Options) (string, error) {
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

	if !opts.ClipboardOnly {
		if err := os.MkdirAll(opts.OutputFolder, 0755); err != nil {
			return "", fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	var toolPath string
	var err error

	if opts.UseFlameshot {
		err = captureWithFlameshot(outputPath, opts)
		if err != nil {
			return "", err
		}
	} else {
		toolPath, err = exec.LookPath("hyprshot")
		if err != nil {
			if opts.Debug {
				fmt.Printf("Debug: hyprshot not found, trying flameshot: %v\n", err)
			}
			err = captureWithFlameshot(outputPath, opts)
			if err != nil {
				return "", fmt.Errorf("neither hyprshot nor flameshot found in PATH")
			}
		} else {
			err = captureWithHyprshot(toolPath, outputPath, opts)
			if err != nil {
				return "", err
			}
		}
	}

	time.Sleep(200 * time.Millisecond)

	if opts.ClipboardOnly {
		return "", nil
	}

	info, statErr := os.Stat(outputPath)
	if statErr != nil {
		if opts.Debug {
			fmt.Printf("Debug: file check error at %s: %v\n", outputPath, statErr)
		}

		files, _ := filepath.Glob(filepath.Join(opts.OutputFolder, "hyprshare-*"))
		if len(files) > 0 {
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

func captureWithHyprshot(hypshotPath, outputPath string, opts cli.Options) error {
	args := []string{}
	for _, mode := range opts.Mode {
		args = append(args, "-m", mode)
	}

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

	if opts.Debug {
		fmt.Printf("Command: %s %s\n", hypshotPath, strings.Join(args, " "))
	}

	cmd := exec.Command(hypshotPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func captureWithFlameshot(outputPath string, opts cli.Options) error {
	flameshotPath, err := exec.LookPath("flameshot")
	if err != nil {
		return fmt.Errorf("flameshot not found in PATH: %v", err)
	}

	if opts.Debug {
		fmt.Println("Using flameshot for screenshot capture")
	}

	var flameShotMode string
	if len(opts.Mode) > 0 {
		switch opts.Mode[0] {
		case "output", "screen":
			flameShotMode = "screen"
		case "window":
			return fmt.Errorf("flameshot does not support window mode directly, use region mode instead")
		case "region":
			flameShotMode = "gui"
		case "active":
			flameShotMode = "screen"
		default:
			flameShotMode = "gui"
		}
	} else {
		flameShotMode = "gui"
	}

	args := []string{flameShotMode}

	if opts.ClipboardOnly {
		args = append(args, "--clipboard")
	} else {
		args = append(args, "--path", outputPath)
	}

	if opts.NotifTimeout > 0 && flameShotMode == "gui" {
		args = append(args, "--delay", strconv.Itoa(opts.NotifTimeout))
	}

	if opts.Debug {
		fmt.Printf("Command: %s %s\n", flameshotPath, strings.Join(args, " "))
	}

	cmd := exec.Command(flameshotPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
