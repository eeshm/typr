package app

import (
	"fmt"
	"strings"
	"time"

	"terminal-wpm/internal/engine"
)

const (
	colorReset = "\x1b[0m"
	colorGreen = "\x1b[32m"
	colorRed   = "\x1b[31m"
	colorCyan  = "\x1b[36m"
	colorGray  = "\x1b[90m"
)

func renderLive(mode, target string, input []rune, metrics engine.Metrics, remaining time.Duration, hasLimit bool) {
	fmt.Printf("%sTerminal WPM Test%s\n", colorCyan, colorReset)
	fmt.Printf("Mode: %s\n", mode)
	if hasLimit {
		if remaining < 0 {
			remaining = 0
		}
		fmt.Printf("Time left: %s\n", formatDuration(remaining))
	}
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("Type the text below:")
	fmt.Println(renderTarget(target, input))
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("WPM: %.1f | Accuracy: %.1f%% | Errors: %d\n", metrics.WPM, metrics.Accuracy, metrics.Errors)
	fmt.Println("Backspace to correct. Ctrl+C to stop.")
}

func renderSummary(metrics engine.Metrics) {
	fmt.Printf("%sTest Complete%s\n", colorCyan, colorReset)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Final WPM: %.1f\n", metrics.WPM)
	fmt.Printf("Accuracy: %.1f%%\n", metrics.Accuracy)
	fmt.Printf("Time taken: %s\n", formatDuration(metrics.TimeTaken))
	fmt.Printf("Total errors: %d\n", metrics.Errors)
	fmt.Printf("Total typed: %d\n", metrics.TotalTyped)
	if metrics.TimedOut {
		fmt.Println("Result: Time limit reached")
	} else if metrics.Cancelled {
		fmt.Println("Result: Stopped by user")
	} else if metrics.Completed {
		fmt.Println("Result: Text completed")
	}
}

func renderTarget(target string, input []rune) string {
	targetRunes := []rune(target)
	var builder strings.Builder

	for i, r := range targetRunes {
		if i < len(input) {
			if input[i] == r {
				builder.WriteString(colorGreen)
				builder.WriteRune(r)
				builder.WriteString(colorReset)
			} else {
				builder.WriteString(colorRed)
				builder.WriteRune(r)
				builder.WriteString(colorReset)
			}
		} else {
			builder.WriteString(colorGray)
			builder.WriteRune(r)
			builder.WriteString(colorReset)
		}
	}

	return builder.String()
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Round(time.Second).Seconds())
	minutes := total / 60
	seconds := total % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
