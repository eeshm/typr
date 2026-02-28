# terminal-wpm

A production-style terminal typing speed test CLI in Go.

## Features
- Structured TUI loop powered by Bubble Tea (smooth in-place updates)
- Styling and color rendering via Lip Gloss
- Random word test generated at start (`quote` or `code` word bank)
- Startup prompt to choose `30` or `60` words each run
- Starts timing on first typed character
- Real-time key capture (no Enter needed)
- Backspace support
- Live colored feedback:
  - Green = correct
  - Red = incorrect
  - Highlighted cursor block + underline = current target character
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

No flags are required. Launch it directly and start typing.

You will be prompted each run:

```text
Choose word count [30/60]:
```

## WPM & Accuracy formula
- `WPM = (total characters typed / 5) / minutes`
- `Accuracy = correct characters / total characters * 100`

## Project layout
- `cmd/terminal-wpm/main.go` - CLI entrypoint (runs directly, no required flags)
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

## Publish (GitHub Releases)
This repo is configured with GoReleaser + GitHub Actions.

1. Commit and push all changes.
2. Create a semantic version tag:

```powershell
git tag v1.0.0
git push origin v1.0.0
```

3. GitHub Actions runs `.github/workflows/release.yml` and publishes binaries for:
  - Windows (`amd64`, `arm64`)
  - Linux (`amd64`, `arm64`)
  - macOS (`amd64`, `arm64`)

Release artifacts are uploaded automatically to the GitHub Release for that tag.

## Publish to npm
An npm wrapper package is provided at `npm/terminal-wpm`.

### One-time setup
- Create an npm automation token and add it as GitHub secret: `NPM_TOKEN`
- Ensure your GitHub release exists (tag like `v1.0.0`)

### Automated publish
When a GitHub Release is published, `.github/workflows/npm-publish.yml` publishes the wrapper package.

### Package users install
```bash
npm i -g @eeshm/typr
```

or:

```bash
npx @eeshm/typr
```
