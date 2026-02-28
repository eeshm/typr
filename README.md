# terminal-wpm

A production-style terminal typing speed test CLI in Go, built to help you learn practical Go architecture while shipping a real tool.

## Features
- Random text at start (`quote` or `code` mode)
- Starts timing on first typed character
- Real-time key capture (no Enter needed)
- Backspace support
- Live colored feedback:
  - Green = correct
  - Red = incorrect
- Live stats:
  - WPM
  - Accuracy
  - Errors
- End conditions:
  - Text completed
  - Optional time limit reached (`--time`)
- Graceful Ctrl+C handling

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
- `internal/app` - app loop + rendering
- `internal/engine` - typing session state + scoring
- `internal/content` - random quote/code text provider
- `internal/terminal` - raw key input + terminal control

## Learn-Go notes
This project demonstrates:
- Package-based design with `internal/`
- State management via a `Session` type
- Time-based calculations with `time.Duration`
- Signal handling (`os/signal`) for graceful shutdown
- Windows terminal mode handling for real-time input

## Validate
```powershell
go test ./...
go vet ./...
```
