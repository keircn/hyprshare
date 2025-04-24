package main

import (
	"fmt"
	"os"

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

	screenshotPath, err := screenshot.Capture(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if opts.ClipboardOnly {
		fmt.Println("Screenshot copied to clipboard")
		return
	}

	fmt.Printf("Screenshot saved to: %s\n", screenshotPath)

	if opts.Upload {
		err = upload.ToHost(screenshotPath, opts.Host, opts.FallbackHost, opts.Debug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "\nThe screenshot was still saved to: %s\n", screenshotPath)
			if opts.Debug {
				fmt.Fprintln(os.Stderr, "\nTip: Verify hostman configuration with 'hostman config'")
				fmt.Fprintln(os.Stderr, "You may need to update your API key or try a different host.")
			}
			os.Exit(1)
		}
	}
}
