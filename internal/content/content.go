package content

import (
	"errors"
	"math/rand"
	"time"
)

var quotes = []string{
	"Small daily improvements over time lead to stunning results.",
	"Simplicity is the soul of efficiency in both code and life.",
	"Focus on consistency first; speed follows naturally.",
	"Ship one clear thing well before chasing many things poorly.",
}

var codeSnippets = []string{
	"for i := 0; i < len(items); i++ { if items[i]%2 == 0 { total += items[i] } }",
	"const greet = (name) => { return `Hello, ${name}!`; }; console.log(greet('dev'));",
	"if err != nil { return fmt.Errorf(\"save failed: %w\", err) }",
	"function debounce(fn, ms){ let t; return (...args)=>{ clearTimeout(t); t=setTimeout(()=>fn(...args), ms); }; }",
}

func RandomText(mode string) (string, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	switch mode {
	case "quote":
		return quotes[rng.Intn(len(quotes))], nil
	case "code":
		return codeSnippets[rng.Intn(len(codeSnippets))], nil
	default:
		return "", errors.New("unsupported mode")
	}
}
