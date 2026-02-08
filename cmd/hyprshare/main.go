package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/keircn/hyprshare/internal/cli"
	"github.com/keircn/hyprshare/internal/screenshot"
	"github.com/keircn/hyprshare/internal/upload"
)

func main() {
	opts, err := cli.ParseOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	screenshotPath, err := screenshot.Capture(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if opts.ClipboardOnly {
		fmt.Println("Screenshot copied to clipboard")
		return
	}

	if opts.Upload {
		err = upload.ToHost(screenshotPath, opts.Host, opts.UseAt, opts.Debug)
		if err != nil {
			if !opts.Silent {
				_ = exec.Command("notify-send", "HyprShare", fmt.Sprintf("Upload failed: %v", err)).Run()
			}
			fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "\nThe screenshot was still saved to: %s\n", screenshotPath)
			if opts.Debug {
				binaryName := "hostman"
				if opts.UseAt {
					binaryName = "at"
				}
				fmt.Fprintf(os.Stderr, "\nTip: Check that %s works when used directly:\n", binaryName)
				fmt.Fprintf(os.Stderr, "  %s upload %s\n", binaryName, screenshotPath)
			}
			os.Exit(1)
		} else if !opts.Silent {
			_ = exec.Command("notify-send", "HyprShare", "Screenshot uploaded successfully").Run()
		}
	}
}