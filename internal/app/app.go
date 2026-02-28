package app

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"terminal-wpm/internal/content"
	"terminal-wpm/internal/engine"
	"terminal-wpm/internal/terminal"
)

type Config struct {
	Mode      string
	TimeLimit time.Duration
}

func Run(cfg Config) error {
	text, err := content.RandomText(cfg.Mode)
	if err != nil {
		return err
	}

	cons := terminal.New()
	if err := cons.Init(); err != nil {
		return err
	}
	defer func() {
		_ = cons.Restore()
		terminal.ShowCursor()
		fmt.Print("\n")
	}()

	terminal.HideCursor()
	terminal.ClearScreen()

	sess := engine.NewSession(text, cfg.TimeLimit)
	keyCh := make(chan terminal.KeyEvent, 32)
	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		for {
			key, readErr := cons.ReadKey()
			if readErr != nil {
				errCh <- readErr
				return
			}
			keyCh <- key
		}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	cancelled := false
	timedOut := false
	ended := false

	renderNow := func(now time.Time) {
		terminal.MoveHome()
		terminal.EraseDisplay()
		metrics := sess.Snapshot(now, timedOut, cancelled)
		remaining := cfg.TimeLimit - sess.Elapsed(now)
		renderLive(cfg.Mode, text, sess.Input(), metrics, remaining, cfg.TimeLimit > 0)
	}

	renderNow(time.Now())

	for !ended {
		select {
		case <-ticker.C:
			now := time.Now()
			if sess.IsTimedOut(now) {
				timedOut = true
				ended = true
			}
			renderNow(now)
		case sig := <-sigCh:
			if sig != nil {
				cancelled = true
				ended = true
			}
		case readErr := <-errCh:
			if errors.Is(readErr, os.ErrClosed) {
				ended = true
				break
			}
			return fmt.Errorf("input error: %w", readErr)
		case key := <-keyCh:
			now := time.Now()
			switch key.Type {
			case terminal.KeyCtrlC:
				cancelled = true
				ended = true
			case terminal.KeyBackspace:
				sess.Backspace()
			case terminal.KeyRune:
				sess.ApplyRune(key.Rune, now)
			}

			if sess.IsCompleted() {
				ended = true
			}
			renderNow(now)
		}
	}

	terminal.ClearScreen()
	final := sess.Snapshot(time.Now(), timedOut, cancelled)
	renderSummary(final)
	return nil
}
