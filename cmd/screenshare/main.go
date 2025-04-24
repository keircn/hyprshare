package main

import (
    "fmt"
    "os"

    "github.com/q4ow/screenshare/internal/cli"
    "github.com/q4ow/screenshare/internal/screenshot"
    "github.com/q4ow/screenshare/internal/upload"
)

func main() {
    opts, err := cli.ParseOptions()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing options: %v\n", err)
        os.Exit(1)
    }

    screenshotPath, err := screenshot.Capture(opts.Type, opts.OutputDir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to capture screenshot: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Screenshot saved to: %s\n", screenshotPath)

    err = upload.ToHost(screenshotPath, opts.Host)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to upload screenshot: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Screenshot uploaded successfully!")
}
