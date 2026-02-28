package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	panelWidth = 72
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	hintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	currentStyle = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("16")).Background(lipgloss.Color("229"))
	endCursor    = lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("229")).Render(" ")
	remainStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	statsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	textStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(1, 1)

	finalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("69")).
			Padding(1, 3)

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

func (m model) viewLive() string {
	metrics := m.session.Snapshot(m.now, false, false)
	elapsed := m.session.Elapsed(m.now)

	statsRows := []string{
		fmt.Sprintf("WPM: %.1f", metrics.WPM),
		fmt.Sprintf("Accuracy: %.1f%%", metrics.Accuracy),
		fmt.Sprintf("Elapsed: %s", formatDuration(elapsed)),
		fmt.Sprintf("Errors: %d", metrics.Errors),
	}
	if m.cfg.TimeLimit > 0 {
		remaining := m.cfg.TimeLimit - elapsed
		if remaining < 0 {
			remaining = 0
		}
		statsRows = append(statsRows, fmt.Sprintf("Time Left: %s", formatDuration(remaining)))
	}

	header := titleStyle.Render("Terminal WPM") + "\n" +
		hintStyle.Render(fmt.Sprintf("Mode: %s  •  Words: %d  •  Start typing to begin timer", m.cfg.Mode, m.cfg.WordCount))

	typedText := renderTarget(m.target, m.session.Input())
	main := textStyle.Width(panelWidth).Render(typedText)
	stats := statsStyle.Width(panelWidth).Render(strings.Join(statsRows, "\n"))
	footer := hintStyle.Render("Backspace to correct • Ctrl+C to stop")

	return lipgloss.JoinVertical(lipgloss.Left, header, "", main, "", stats, "", footer)
}

func (m model) viewSummary() string {
	metrics := m.final
	if metrics.TimeTaken == 0 {
		metrics = m.session.Snapshot(m.now, m.timedOut, m.cancelled)
	}

	resultLabel := "Text completed"
	if metrics.TimedOut {
		resultLabel = "Time limit reached"
	}
	if metrics.Cancelled {
		resultLabel = "Stopped by user"
	}

	body := strings.Join([]string{
		titleStyle.Render("Typing Test Results"),
		"",
		fmt.Sprintf("Final WPM: %.1f", metrics.WPM),
		fmt.Sprintf("Accuracy: %.1f%%", metrics.Accuracy),
		fmt.Sprintf("Total errors: %d", metrics.Errors),
		fmt.Sprintf("Time taken: %s", formatDuration(metrics.TimeTaken)),
		fmt.Sprintf("Tier: %s", performanceTier(metrics.WPM)),
		fmt.Sprintf("Result: %s", resultLabel),
		"",
		hintStyle.Render("Press Enter, q, or Esc to exit"),
	}, "\n")

	boxed := finalStyle.Render(body)
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxed)
	}
	return boxed
}

func renderTarget(target string, input []rune) string {
	targetRunes := []rune(target)
	var builder strings.Builder
	cursor := len(input)

	for i, r := range targetRunes {
		if i < len(input) {
			if input[i] == r {
				builder.WriteString(correctStyle.Render(string(r)))
			} else {
				builder.WriteString(wrongStyle.Render(string(r)))
			}
		} else if i == cursor {
			if r == ' ' {
				builder.WriteString(currentStyle.Render("·"))
			} else {
				builder.WriteString(currentStyle.Render(string(r)))
			}
		} else {
			builder.WriteString(remainStyle.Render(string(r)))
		}
	}

	if cursor >= len(targetRunes) {
		builder.WriteString(endCursor)
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

func performanceTier(wpm float64) string {
	switch {
	case wpm < 30:
		return "Beginner"
	case wpm < 50:
		return "Average"
	case wpm < 80:
		return "Fast"
	default:
		return "Elite"
	}
}
