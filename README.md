# terminal-wpm

A production-style terminal typing speed test CLI in Go, built to help you learn practical Go architecture while shipping a real tool.

## Features
- Structured TUI loop powered by Bubble Tea (smooth in-place updates)
- Styling and color rendering via Lip Gloss
- Random text at start (`quote` or `code` mode)
- Starts timing on first typed character
- Real-time key capture (no Enter needed)
- Backspace support
- Live colored feedback:
  - Green = correct
  - Red = incorrect
  - Underlined = current target character
  - Dim gray = remaining characters
- Live stats:
  - WPM
  - Accuracy
  - Elapsed time
  - Errors
- End conditions:
  - Text completed
  - Optional time limit reached (`--time`)
- Graceful Ctrl+C handling
- Final centered results screen with performance tier:
  - `<30` Beginner
  - `30-50` Average
  - `50-80` Fast
  - `80+` Elite

## Install Go (Windows)
Choose one:

- `winget install --id GoLang.Go -e`
- Or install MSI from https://go.dev/dl/

Verify:

```powershell
go version
```

## Build & Run
From the project root:

```powershell
go mod tidy
go run ./cmd/terminal-wpm
```

### CLI flags
- `--mode quote` (default)
- `--mode code`
- `--time <seconds>` (optional time-limited run)
- `--help`

Examples:

```powershell
go run ./cmd/terminal-wpm --mode quote
go run ./cmd/terminal-wpm --mode code --time 60
```

## WPM & Accuracy formula
- `WPM = (total characters typed / 5) / minutes`
- `Accuracy = correct characters / total characters * 100`

## Project layout
- `cmd/terminal-wpm/main.go` - CLI entrypoint and flag parsing
- `internal/app` - Bubble Tea model/update/view + Lip Gloss rendering
- `internal/engine` - typing session state + scoring
- `internal/content` - random quote/code text provider
- `internal/terminal` - legacy terminal helpers (kept for compatibility)

## Learn-Go notes
This project demonstrates:
- Package-based design with `internal/`
- State management via a `Session` type
- Time-based calculations with `time.Duration`
- Event-driven terminal UI with Bubble Tea
- Reusable styling with Lip Gloss

## Validate
```powershell
go test ./...
go vet ./...
```
