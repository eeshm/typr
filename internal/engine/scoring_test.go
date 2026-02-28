package engine

import (
	"math"
	"testing"
	"time"
)

func TestCalculateWPM(t *testing.T) {
	wpm := CalculateWPM(250, 2*time.Minute)
	if math.Abs(wpm-25.0) > 0.001 {
		t.Fatalf("expected 25.0, got %.3f", wpm)
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
