package engine

import (
	"testing"
	"time"
)

func TestSessionTypingAndBackspace(t *testing.T) {
	start := time.Now()
	s := NewSession("abc", 0)

	s.ApplyRune('a', start)
	s.ApplyRune('x', start.Add(100*time.Millisecond)) // wrong
	if s.Cursor() != 2 {
		t.Fatalf("expected cursor 2, got %d", s.Cursor())
	}

	// Backspace should undo the 'x' error from counts
	s.Backspace()
	if s.Cursor() != 1 {
		t.Fatalf("expected cursor 1 after backspace, got %d", s.Cursor())
	}

	s.ApplyRune('b', start.Add(200*time.Millisecond))
	s.ApplyRune('c', start.Add(300*time.Millisecond))

	if !s.IsCompleted() {
		t.Fatal("expected session to be completed")
	}

	m := s.Snapshot(start.Add(400*time.Millisecond), false, false)
	// After backspace undo: typed a, (x undone), b, c = 3 total
	if m.TotalTyped != 3 {
		t.Fatalf("expected total typed 3, got %d", m.TotalTyped)
	}
	// The 'x' error was undone by backspace, all 3 remaining are correct
	if m.Errors != 0 {
		t.Fatalf("expected 0 errors (backspace undid the mistake), got %d", m.Errors)
	}
	if m.Correct != 3 {
		t.Fatalf("expected 3 correct, got %d", m.Correct)
	}
}

func TestSessionTimeoutStartsOnFirstKey(t *testing.T) {
	now := time.Now()
	s := NewSession("abc", 2*time.Second)

	if s.IsTimedOut(now.Add(5 * time.Second)) {
		t.Fatal("should not timeout before first key")
	}

	s.ApplyRune('a', now)
	if !s.IsTimedOut(now.Add(3 * time.Second)) {
		t.Fatal("expected timeout after first key and limit")
	}
}
