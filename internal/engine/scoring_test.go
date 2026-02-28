package engine

import (
	"math"
	"testing"
	"time"
)

func TestNetWPM(t *testing.T) {
	// 200 correct chars in 2 min => (200/5)/2 = 20 WPM
	wpm := CalculateNetWPM(200, 2*time.Minute)
	if math.Abs(wpm-20.0) > 0.001 {
		t.Fatalf("expected 20.0, got %.3f", wpm)
	}
}

func TestRawWPM(t *testing.T) {
	// 250 total chars in 2 min => (250/5)/2 = 25 WPM
	wpm := CalculateRawWPM(250, 2*time.Minute)
	if math.Abs(wpm-25.0) > 0.001 {
		t.Fatalf("expected 25.0, got %.3f", wpm)
	}
}

func TestNetWPMIgnoresErrors(t *testing.T) {
	// 0 correct chars => 0 WPM regardless of time
	wpm := CalculateNetWPM(0, 1*time.Minute)
	if wpm != 0 {
		t.Fatalf("expected 0, got %.3f", wpm)
	}
}

func TestCalculateAccuracy(t *testing.T) {
	acc := CalculateAccuracy(45, 50)
	if math.Abs(acc-90.0) > 0.001 {
		t.Fatalf("expected 90.0, got %.3f", acc)
	}
}

func TestCalculateAccuracyZeroTotal(t *testing.T) {
	acc := CalculateAccuracy(0, 0)
	if acc != 100 {
		t.Fatalf("expected 100, got %.3f", acc)
	}
}

func TestCountCorrectWords(t *testing.T) {
	target := []rune("hello world test")
	input := []rune("hello wxrld test")
	correct, total := CountCorrectWords(target, input)
	if total != 3 {
		t.Fatalf("expected 3 total words, got %d", total)
	}
	if correct != 2 {
		t.Fatalf("expected 2 correct words (hello, test), got %d", correct)
	}
}

func TestCountCorrectWordsPartialInput(t *testing.T) {
	target := []rune("one two three")
	input := []rune("one tw")
	correct, total := CountCorrectWords(target, input)
	// "one" is complete+correct, "tw" is partial (counts as attempted but incomplete)
	if correct != 1 {
		t.Fatalf("expected 1 correct word, got %d", correct)
	}
	if total != 2 {
		t.Fatalf("expected 2 total words attempted, got %d", total)
	}
}
