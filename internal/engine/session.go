package engine

import "time"

type Session struct {
	target       []rune
	input        []rune
	cursor       int
	started      bool
	startTime    time.Time
	endTime      time.Time
	timeLimit    time.Duration
	totalTyped   int
	correctTyped int
	errors       int
}

func NewSession(target string, timeLimit time.Duration) *Session {
	return &Session{
		target:    []rune(target),
		timeLimit: timeLimit,
	}
}

func (s *Session) Target() []rune {
	return s.target
}

func (s *Session) Input() []rune {
	return s.input
}

func (s *Session) Cursor() int {
	return s.cursor
}

func (s *Session) Started() bool {
	return s.started
}

func (s *Session) StartTime() time.Time {
	return s.startTime
}

func (s *Session) TimeLimit() time.Duration {
	return s.timeLimit
}

// ApplyRune types a character and returns true if it was correct.
func (s *Session) ApplyRune(ch rune, now time.Time) bool {
	if s.IsCompleted() {
		return false
	}
	if !s.started {
		s.started = true
		s.startTime = now
	}
	if s.cursor >= len(s.target) {
		return false
	}

	expected := s.target[s.cursor]
	s.totalTyped++
	correct := ch == expected
	if correct {
		s.correctTyped++
	} else {
		s.errors++
	}

	if s.cursor < len(s.input) {
		s.input[s.cursor] = ch
	} else {
		s.input = append(s.input, ch)
	}
	s.cursor++
	if s.IsCompleted() {
		s.endTime = now
	}
	return correct
}

func (s *Session) Backspace() {
	if s.cursor == 0 {
		return
	}
	s.cursor--

	// Undo the scoring for the character we're erasing.
	// This prevents inflated totalTyped/errors when the user corrects mistakes.
	if s.cursor < len(s.input) {
		typed := s.input[s.cursor]
		expected := s.target[s.cursor]
		s.totalTyped--
		if typed == expected {
			s.correctTyped--
		} else {
			s.errors--
		}
	}

	s.input = s.input[:s.cursor]
	if !s.endTime.IsZero() {
		s.endTime = time.Time{}
	}
}

func (s *Session) IsCompleted() bool {
	return s.cursor >= len(s.target)
}

func (s *Session) IsTimedOut(now time.Time) bool {
	if !s.started || s.timeLimit <= 0 {
		return false
	}
	return now.Sub(s.startTime) >= s.timeLimit
}

func (s *Session) Elapsed(now time.Time) time.Duration {
	if !s.started {
		return 0
	}
	if !s.endTime.IsZero() {
		return s.endTime.Sub(s.startTime)
	}
	elapsed := now.Sub(s.startTime)
	if elapsed < 0 {
		return 0
	}
	return elapsed
}

func (s *Session) Snapshot(now time.Time, timedOut, cancelled bool) Metrics {
	elapsed := s.Elapsed(now)
	correctWords, totalWords := CountCorrectWords(s.target, s.input)
	return Metrics{
		WPM:          CalculateNetWPM(s.correctTyped, elapsed),
		RawWPM:       CalculateRawWPM(s.totalTyped, elapsed),
		Accuracy:     CalculateAccuracy(s.correctTyped, s.totalTyped),
		Errors:       s.errors,
		TotalTyped:   s.totalTyped,
		Correct:      s.correctTyped,
		CorrectWords: correctWords,
		TotalWords:   totalWords,
		TimeTaken:    elapsed,
		Completed:    s.IsCompleted(),
		TimedOut:     timedOut,
		Cancelled:    cancelled,
	}
}
