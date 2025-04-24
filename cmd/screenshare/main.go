package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/q4ow/hyprshare/internal/cli"
    "github.com/q4ow/hyprshare/internal/screenshot"
    "github.com/q4ow/hyprshare/internal/upload"
)

func main() {
	opts, err := cli.ParseOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !opts.ClipboardOnly {
		if err := os.MkdirAll(opts.OutputFolder, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
			os.Exit(1)
		}
	}

	screenshotPath, err := screenshot.Capture(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to capture screenshot: %v\n", err)
		os.Exit(1)
	}

	if opts.ClipboardOnly {
		fmt.Println("Screenshot copied to clipboard")
		return
	}

	fmt.Printf("Screenshot saved to: %s\n", screenshotPath)

	if opts.Upload {
		err = upload.ToHost(screenshotPath, opts.Host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to upload screenshot: %v\n", err)
			os.Exit(1)
		}
	}
}