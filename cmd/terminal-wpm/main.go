package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"terminal-wpm/internal/app"
)

func main() {
	var mode string
	var seconds int

	flag.StringVar(&mode, "mode", "quote", "typing mode: quote or code")
	flag.IntVar(&seconds, "time", 0, "optional time limit in seconds")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--mode quote|code] [--time seconds]\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "A real-time terminal typing speed test.")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if mode != "quote" && mode != "code" {
		fmt.Fprintln(os.Stderr, "invalid mode: use --mode quote or --mode code")
		os.Exit(2)
	}
	if seconds < 0 {
		fmt.Fprintln(os.Stderr, "invalid --time: must be >= 0")
		os.Exit(2)
	}

	wordCount, err := promptWordCount()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read word count:", err)
		os.Exit(1)
	}

	cfg := app.Config{
		Mode:      mode,
		TimeLimit: time.Duration(seconds) * time.Second,
		WordCount: wordCount,
	}

	if err := app.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func promptWordCount() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Choose word count [30/60]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		switch strings.TrimSpace(input) {
		case "30":
			return 30, nil
		case "60":
			return 60, nil
		default:
			fmt.Println("Please enter only 30 or 60.")
		}
	}
}
