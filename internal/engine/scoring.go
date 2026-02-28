package engine

import "time"

// CalculateNetWPM returns WPM based on correctly typed characters only.
// This matches how real typing test sites (monkeytype, typing.com) work:
// Net WPM = (correct characters / 5) / minutes
func CalculateNetWPM(correctChars int, elapsed time.Duration) float64 {
	if correctChars <= 0 || elapsed <= 0 {
		return 0
	}
	minutes := elapsed.Minutes()
	if minutes <= 0 {
		return 0
	}
	return (float64(correctChars) / 5.0) / minutes
}

// CalculateRawWPM returns WPM based on all keystrokes (including errors).
func CalculateRawWPM(totalChars int, elapsed time.Duration) float64 {
	if totalChars <= 0 || elapsed <= 0 {
		return 0
	}
	minutes := elapsed.Minutes()
	if minutes <= 0 {
		return 0
	}
	return (float64(totalChars) / 5.0) / minutes
}

// CalculateAccuracy returns percentage of correct keystrokes.
func CalculateAccuracy(correct, total int) float64 {
	if total <= 0 {
		return 100
	}
	return (float64(correct) / float64(total)) * 100
}

// CountCorrectWords counts fully correct space-delimited words
// by comparing input against target rune-by-rune.
func CountCorrectWords(target, input []rune) (correctWords, totalWords int) {
	wordStart := 0
	for i := 0; i <= len(target); i++ {
		// word boundary: space or end of target
		atBoundary := i == len(target) || target[i] == ' '
		if !atBoundary {
			continue
		}

		// only count words the user has reached
		if wordStart >= len(input) {
			break
		}

		totalWords++

		// check if every char in this word was typed correctly
		wordCorrect := true
		for j := wordStart; j < i; j++ {
			if j >= len(input) || input[j] != target[j] {
				wordCorrect = false
				break
			}
		}
		// also check the space after the word was typed (if applicable)
		if i < len(target) && i < len(input) && input[i] != ' ' {
			wordCorrect = false
		}
		if wordCorrect {
			correctWords++
		}

		wordStart = i + 1
	}
	return
}
