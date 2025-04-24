package screenshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/q4ow/hyprshare/internal/cli"
)

func Capture(opts cli.Options) (string, error) {
	var outputPath string
	if !opts.ClipboardOnly {
		if opts.Filename == "" {
			timestamp := time.Now().Format("20060102-150405")
			filename := fmt.Sprintf("screenshot-%s.png", timestamp)
			outputPath = filepath.Join(opts.OutputFolder, filename)
		} else {
			outputPath = filepath.Join(opts.OutputFolder, opts.Filename)
		}
	}

	args := []string{}
	for _, mode := range opts.Mode {
		args = append(args, "-m", mode)
	}

	if !opts.ClipboardOnly {
		args = append(args, "-o", opts.OutputFolder)
		if opts.Filename != "" {
			args = append(args, "-f", opts.Filename)
		}
	}

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

	cmd := exec.Command("hyprshot", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("hyprshot command failed: %v", err)
	}

	if !opts.ClipboardOnly && opts.Filename == "" {
		return outputPath, nil
	} else if !opts.ClipboardOnly {
		return outputPath, nil
	}

	return "", nil
}