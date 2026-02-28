package engine

import (
	"testing"
	"time"
)

func TestSessionTypingAndBackspace(t *testing.T) {
	start := time.Now()
	s := NewSession("abc", 0)

	s.ApplyRune('a', start)
	s.ApplyRune('x', start.Add(100*time.Millisecond))
	if s.Cursor() != 2 {
		t.Fatalf("expected cursor 2, got %d", s.Cursor())
	}

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
	if m.TotalTyped != 4 {
		t.Fatalf("expected total typed 4, got %d", m.TotalTyped)
	}
	if m.Errors != 1 {
		t.Fatalf("expected 1 error, got %d", m.Errors)
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
