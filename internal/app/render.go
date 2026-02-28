package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"terminal-wpm/internal/history"
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

	selectedStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("63")).Padding(0, 2)
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Padding(0, 2)

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

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 3)

	historyStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("244")).
			Padding(0, 2)

	historyDimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

func (m model) viewMenu() string {
	var rows []string
	rows = append(rows, titleStyle.Render("Terminal WPM"))
	rows = append(rows, "")
	rows = append(rows, "Choose word count:")
	rows = append(rows, "")

	for i, opt := range wordOptions {
		if i == m.menuIdx {
			rows = append(rows, selectedStyle.Render("▸ "+opt.label))
		} else {
			rows = append(rows, unselectedStyle.Render("  "+opt.label))
		}
	}

	rows = append(rows, "")
	rows = append(rows, hintStyle.Render("↑/↓ to move • Enter to start • Ctrl+C to quit"))

	box := menuStyle.Render(strings.Join(rows, "\n"))
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return box
}

func (m model) viewLive() string {
	metrics := m.session.Snapshot(m.now, false, false)
	elapsed := m.session.Elapsed(m.now)

	statsRows := []string{
		fmt.Sprintf("WPM: %.1f", metrics.WPM),
		fmt.Sprintf("Raw WPM: %.1f", metrics.RawWPM),
		fmt.Sprintf("Accuracy: %.1f%%", metrics.Accuracy),
		fmt.Sprintf("Words: %d/%d correct", metrics.CorrectWords, metrics.TotalWords),
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

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", main, "", stats, "", footer)
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

func (m model) viewSummary() string {
	metrics := m.final

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
		fmt.Sprintf("WPM: %.1f", metrics.WPM),
		fmt.Sprintf("Raw WPM: %.1f", metrics.RawWPM),
		fmt.Sprintf("Accuracy: %.1f%%", metrics.Accuracy),
		fmt.Sprintf("Correct Words: %d / %d", metrics.CorrectWords, metrics.TotalWords),
		fmt.Sprintf("Total errors: %d", metrics.Errors),
		fmt.Sprintf("Time taken: %s", formatDuration(metrics.TimeTaken)),
		fmt.Sprintf("Tier: %s", performanceTier(metrics.WPM)),
		fmt.Sprintf("Result: %s", resultLabel),
		"",
		hintStyle.Render("Press Enter, q, or Esc to exit"),
	}, "\n")

	boxed := finalStyle.Render(body)

	// Append recent history below the result box.
	historyBox := renderHistory(m.history)

	scrollHint := hintStyle.Render("↑/↓ to scroll")
	combined := lipgloss.JoinVertical(lipgloss.Center, boxed, "", historyBox, "", scrollHint)

	// Apply vertical scrolling when content exceeds terminal height.
	if m.height > 0 {
		lines := strings.Split(combined, "\n")
		totalLines := len(lines)

		// Clamp scroll so we don't scroll past the end.
		maxScroll := totalLines - m.height
		if maxScroll < 0 {
			maxScroll = 0
		}
		offset := m.scrollY
		if offset > maxScroll {
			offset = maxScroll
		}

		if maxScroll == 0 {
			// Content fits — center normally.
			if m.width > 0 {
				return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, combined)
			}
			return combined
		}

		// Content overflows — show windowed slice, no vertical re-centering.
		end := offset + m.height
		if end > totalLines {
			end = totalLines
		}
		visible := strings.Join(lines[offset:end], "\n")

		if m.width > 0 {
			// Only center horizontally, fill full height so Bubble Tea doesn't jump.
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, visible)
		}
		return visible
	}
	return combined
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

func renderHistory(records []history.Record) string {
	if len(records) == 0 {
		return historyStyle.Render(historyDimStyle.Render("No previous sessions yet."))
	}

	var rows []string
	rows = append(rows, hintStyle.Render("Recent Sessions"))
	rows = append(rows, historyDimStyle.Render(fmt.Sprintf("%-12s %6s %6s %7s %s", "Date", "WPM", "Raw", "Acc", "Tier")))

	for _, r := range records {
		dateStr := r.Date.Format("Jan 02 15:04")
		line := fmt.Sprintf("%-12s %6.1f %6.1f %6.1f%% %s", dateStr, r.WPM, r.RawWPM, r.Accuracy, r.Tier)
		rows = append(rows, historyDimStyle.Render(line))
	}

	return historyStyle.Render(strings.Join(rows, "\n"))
}
