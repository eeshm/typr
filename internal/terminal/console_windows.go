//go:build windows

package terminal

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
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
	stdinFd    int
	stdoutFd   int
	oldState   *term.State
	oldInMode  uint32
	oldOutMode uint32
}

func New() *Console {
	return &Console{
		stdinFd:  int(os.Stdin.Fd()),
		stdoutFd: int(os.Stdout.Fd()),
	}
}

func (c *Console) Init() error {
	if err := c.enableVirtualTerminal(); err != nil {
		return err
	}

	state, err := term.MakeRaw(c.stdinFd)
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	c.oldState = state
	return nil
}

func (c *Console) Restore() error {
	if c.oldState != nil {
		if err := term.Restore(c.stdinFd, c.oldState); err != nil {
			return fmt.Errorf("failed restoring terminal state: %w", err)
		}
	}
	if c.oldInMode != 0 {
		hIn := windows.Handle(c.stdinFd)
		_ = windows.SetConsoleMode(hIn, c.oldInMode)
	}
	if c.oldOutMode != 0 {
		hOut := windows.Handle(c.stdoutFd)
		_ = windows.SetConsoleMode(hOut, c.oldOutMode)
	}
	return nil
}

func (c *Console) enableVirtualTerminal() error {
	hOut := windows.Handle(c.stdoutFd)
	var outMode uint32
	if err := windows.GetConsoleMode(hOut, &outMode); err == nil {
		c.oldOutMode = outMode
		newOutMode := outMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		if err := windows.SetConsoleMode(hOut, newOutMode); err != nil {
			return fmt.Errorf("failed to enable VT output mode: %w", err)
		}
	}

	hIn := windows.Handle(c.stdinFd)
	var inMode uint32
	if err := windows.GetConsoleMode(hIn, &inMode); err == nil {
		c.oldInMode = inMode
		newInMode := inMode | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
		_ = windows.SetConsoleMode(hIn, newInMode)
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
		case 13, 10:
			continue
		case 0, 224:
			seq := make([]byte, 1)
			_, _ = os.Stdin.Read(seq)
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
