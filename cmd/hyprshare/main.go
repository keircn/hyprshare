package main

import (
	"fmt"
	"os"

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
		err = upload.ToHost(screenshotPath, opts.Host, opts.Debug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "\nThe screenshot was still saved to: %s\n", screenshotPath)
			if opts.Debug {
				fmt.Fprintln(os.Stderr, "\nTip: Check that hostman works when used directly:")
				fmt.Fprintf(os.Stderr, "  hostman upload %s\n", screenshotPath)
			}
			os.Exit(1)
		}
	}
}
