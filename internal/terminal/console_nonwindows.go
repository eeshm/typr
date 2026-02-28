//go:build !windows

package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type KeyType int

const (
	KeyRune KeyType = iota
	KeyBackspace
	KeyCtrlC
)

type KeyEvent struct {
	Type KeyType
	Rune rune
}

type Console struct {
	stdinFd  int
	oldState *term.State
}

func New() *Console {
	return &Console{stdinFd: int(os.Stdin.Fd())}
}

func (c *Console) Init() error {
	state, err := term.MakeRaw(c.stdinFd)
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	c.oldState = state
	return nil
}

func (c *Console) Restore() error {
	if c.oldState == nil {
		return nil
	}
	if err := term.Restore(c.stdinFd, c.oldState); err != nil {
		return fmt.Errorf("failed restoring terminal state: %w", err)
	}
	return nil
}

func (c *Console) ReadKey() (KeyEvent, error) {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return KeyEvent{}, err
		}
		if n == 0 {
			continue
		}
		b := buf[0]
		switch b {
		case 3:
			return KeyEvent{Type: KeyCtrlC}, nil
		case 8, 127:
			return KeyEvent{Type: KeyBackspace}, nil
		case 10, 13:
			continue
		default:
			if b >= 32 && b <= 126 {
				return KeyEvent{Type: KeyRune, Rune: rune(b)}, nil
			}
		}
	}
}

func ClearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func MoveHome() {
	fmt.Print("\x1b[H")
}

func EraseDisplay() {
	fmt.Print("\x1b[J")
}
