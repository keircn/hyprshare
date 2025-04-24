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
	flag.BoolVar(&opts.ClipboardOnly, "clipboard-only", false, "Copy screenshot to clipboard and don't save image")
	flag.BoolVar(&opts.Upload, "upload", true, "Upload screenshot after capturing")
	flag.BoolVar(&opts.Upload, "u", true, "Upload screenshot after capturing (shorthand)")
	flag.StringVar(&opts.Host, "host", "e-z", "Image host to use with hostman")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] -- [command]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A wrapper around hyprshot that uploads screenshots to an image host\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nModes:\n")
		fmt.Fprintf(os.Stderr, "  output        take screenshot of an entire monitor\n")
		fmt.Fprintf(os.Stderr, "  window        take screenshot of an open window\n")
		fmt.Fprintf(os.Stderr, "  region        take screenshot of selected region\n")
		fmt.Fprintf(os.Stderr, "  active        take screenshot of active window|output\n")
		fmt.Fprintf(os.Stderr, "  OUTPUT_NAME   take screenshot of output with OUTPUT_NAME\n")
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
		os.Exit(1)
	}

	return opts, nil
}
