package engine

import "time"

func CalculateWPM(totalTyped int, elapsed time.Duration) float64 {
	if totalTyped <= 0 || elapsed <= 0 {
		return 0
	}
	minutes := elapsed.Minutes()
	if minutes <= 0 {
		return 0
	}
	return (float64(totalTyped) / 5.0) / minutes
}

func CalculateAccuracy(correct, total int) float64 {
	if total <= 0 {
		return 100
	}
	return (float64(correct) / float64(total)) * 100
}
