package content

import (
	"errors"
	"math/rand/v2"
	"strings"
)

var quoteWords = []string{
	"focus", "habit", "steady", "practice", "moment", "clarity", "rhythm", "progress", "simple", "calm",
	"result", "daily", "effort", "learn", "repeat", "strong", "patience", "consistency", "skill", "growth",
	"speed", "accuracy", "discipline", "improve", "craft", "quality", "mindset", "builder", "future", "master",
}

var codeWords = []string{
	"func", "return", "error", "nil", "range", "slice", "map", "struct", "interface", "goroutine",
	"channel", "mutex", "context", "package", "module", "compile", "testing", "pointer", "method", "receiver",
	"import", "deploy", "commit", "branch", "refactor", "backend", "frontend", "async", "buffer", "runtime",
}

func RandomText(mode string, wordCount int) (string, error) {
	if wordCount <= 0 {
		return "", errors.New("word count must be greater than zero")
	}

	words := make([]string, 0, wordCount)

	pool, err := wordPool(mode)
	if err != nil {
		return "", err
	}

	for range wordCount {
		words = append(words, pool[rand.IntN(len(pool))])
	}

	return strings.Join(words, " "), nil
}

func wordPool(mode string) ([]string, error) {
	switch mode {
	case "quote":
		return quoteWords, nil
	case "code":
		return codeWords, nil
	default:
		return nil, errors.New("unsupported mode")
	}
}
