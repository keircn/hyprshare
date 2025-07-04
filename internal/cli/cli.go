package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Mode          []string
	OutputFolder  string
	Filename      string
	Freeze        bool
	Debug         bool
	Silent        bool
	Raw           bool
	NotifTimeout  int
	ClipboardOnly bool
	Command       []string
	Upload        bool
	Host          string
	UseAt         bool
	UseFlameshot  bool
}

func ParseOptions() (Options, error) {
	var opts Options
	var modesStr string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return opts, fmt.Errorf("failed to get home directory: %v", err)
	}
	defaultOutput := filepath.Join(homeDir, "Pictures", "Screenshots")

	flag.StringVar(&modesStr, "mode", "", "Screenshot mode (output, window, region, active, or OUTPUT_NAME)")
	flag.StringVar(&modesStr, "m", "", "Screenshot mode (shorthand)")
	flag.StringVar(&opts.OutputFolder, "output-folder", defaultOutput, "Directory to save screenshots")
	flag.StringVar(&opts.OutputFolder, "o", defaultOutput, "Directory to save screenshots (shorthand)")
	flag.StringVar(&opts.Filename, "filename", "", "Custom filename for screenshot")
	flag.StringVar(&opts.Filename, "f", "", "Custom filename for screenshot (shorthand)")
	flag.BoolVar(&opts.Freeze, "freeze", false, "Freeze screen on initialization")
	flag.BoolVar(&opts.Freeze, "z", false, "Freeze screen on initialization (shorthand)")
	flag.BoolVar(&opts.Debug, "debug", false, "Print debug information")
	flag.BoolVar(&opts.Debug, "d", false, "Print debug information (shorthand)")
	flag.BoolVar(&opts.Silent, "silent", false, "Don't send notification when screenshot is saved")
	flag.BoolVar(&opts.Silent, "s", false, "Don't send notification when screenshot is saved (shorthand)")
	flag.BoolVar(&opts.Raw, "raw", false, "Output raw image data to stdout")
	flag.BoolVar(&opts.Raw, "r", false, "Output raw image data to stdout (shorthand)")
	flag.IntVar(&opts.NotifTimeout, "notif-timeout", 5000, "Notification timeout in milliseconds")
	flag.IntVar(&opts.NotifTimeout, "t", 5000, "Notification timeout in milliseconds (shorthand)")
	flag.BoolVar(&opts.ClipboardOnly, "clipboard-only", false, "Copy screenshot to clipboard only")
	flag.BoolVar(&opts.Upload, "upload", true, "Upload screenshot after capturing")
	flag.BoolVar(&opts.Upload, "u", true, "Upload screenshot (shorthand)")
	flag.StringVar(&opts.Host, "host", "default", "Image host to use (default: system default host)")
	flag.BoolVar(&opts.UseAt, "use-at", false, "Use 'at' binary instead of 'hostman' for uploads")
	flag.BoolVar(&opts.UseAt, "a", false, "Use 'at' binary instead of 'hostman' for uploads (shorthand)")
	flag.BoolVar(&opts.UseFlameshot, "use-flameshot", false, "Use flameshot instead of hyprshot for screenshots")
	flag.BoolVar(&opts.UseFlameshot, "F", false, "Use flameshot instead of hyprshot for screenshots (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "HyprShare: Screenshot and upload utility for Hyprland\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -m <mode> [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Modes:\n")
		fmt.Fprintf(os.Stderr, "  output    - full monitor screenshot\n")
		fmt.Fprintf(os.Stderr, "  window    - screenshot of a window\n")
		fmt.Fprintf(os.Stderr, "  region    - screenshot of selected area\n")
		fmt.Fprintf(os.Stderr, "  active    - active window/output (use with -m)\n\n")

		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  hyprshare -m window                      # Window screenshot with default host\n")
		fmt.Fprintf(os.Stderr, "  hyprshare -m region --host imgur         # Region screenshot to imgur\n")
		fmt.Fprintf(os.Stderr, "  hyprshare -m window --u=false            # Window screenshot without upload\n")
		fmt.Fprintf(os.Stderr, "  hyprshare -m region -d                   # With debug output\n")
		fmt.Fprintf(os.Stderr, "  hyprshare -m region --use-at             # Use 'at' binary for upload\n")
		fmt.Fprintf(os.Stderr, "  hyprshare -m region --use-flameshot      # Use flameshot for screenshot\n")
	}

	commandIndex := -1
	for i, arg := range os.Args {
		if arg == "--" {
			commandIndex = i
			break
		}
	}

	var args []string
	if commandIndex != -1 {
		args = os.Args[1:commandIndex]
		opts.Command = os.Args[commandIndex+1:]
	} else {
		args = os.Args[1:]
	}

	err = flag.CommandLine.Parse(args)
	if err != nil {
		return opts, err
	}

	if modesStr != "" {
		opts.Mode = strings.Split(modesStr, ",")
	}

	if len(opts.Mode) == 0 {
		flag.Usage()
		return opts, fmt.Errorf("no screenshot mode specified, use -m/--mode")
	}

	return opts, nil
}
