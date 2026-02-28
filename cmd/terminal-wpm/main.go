package main

import (
	"fmt"
	"os"

	"terminal-wpm/internal/app"
)

func main() {
	cfg := app.Config{
		Mode: "quote",
	}

	if err := app.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
