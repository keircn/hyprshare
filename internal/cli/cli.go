// internal/cli/cli.go
package cli

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

type Options struct {
    Type      string
    Host      string
    OutputDir string
}

func ParseOptions() (Options, error) {
    var opts Options

    flag.StringVar(&opts.Host, "host", "e-z", "Image host to use (default: e-z)")
    flag.StringVar(&opts.OutputDir, "output-dir", "", "Directory to save screenshots (default: ~/Pictures/Screenshots)")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [options] <fullscreen|window|region>\n\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
    }

    flag.Parse()

    if flag.NArg() != 1 {
        return opts, fmt.Errorf("exactly one screenshot type must be specified: fullscreen, window, or region")
    }

    opts.Type = flag.Arg(0)
    if opts.Type != "fullscreen" && opts.Type != "window" && opts.Type != "region" {
        return opts, fmt.Errorf("invalid screenshot type: %s (must be fullscreen, window, or region)", opts.Type)
    }

    if opts.OutputDir == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return opts, fmt.Errorf("failed to get user home directory: %v", err)
        }
        opts.OutputDir = filepath.Join(homeDir, "Pictures", "Screenshots")
    }

    if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
        return opts, fmt.Errorf("failed to create output directory: %v", err)
    }

    return opts, nil
}
