package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"codex-brainstorming-telegram/internal/virtualcodex"
)

func main() {
	os.Exit(run(context.Background(), os.Stdout, os.Stderr, os.Args[1:]))
}

func run(parent context.Context, stdout io.Writer, stderr io.Writer, args []string) int {
	fs := flag.NewFlagSet("virtual-codex", flag.ContinueOnError)
	fs.SetOutput(stderr)

	input := fs.String("input", "", "input text for virtual codex")
	timeout := fs.Duration("timeout", 5*time.Minute, "overall timeout")
	delay := fs.Duration("delay", 0, "simulated processing delay")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *timeout <= 0 {
		fmt.Fprintln(stderr, "timeout must be greater than 0")
		return 2
	}

	finalInput := strings.TrimSpace(*input)
	if finalInput == "" && fs.NArg() > 0 {
		finalInput = strings.TrimSpace(strings.Join(fs.Args(), " "))
	}
	if finalInput == "" {
		fmt.Fprintln(stderr, "input is required (use --input or positional text)")
		return 2
	}

	ctx, cancel := context.WithTimeout(parent, *timeout)
	defer cancel()

	engine := virtualcodex.NewEngine(virtualcodex.Config{ProcessingDelay: *delay})
	resp, err := engine.Respond(ctx, finalInput)
	if err != nil {
		fmt.Fprintf(stderr, "virtual-codex error: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, resp)
	return 0
}
