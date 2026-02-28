package content

import (
	"errors"
	"math/rand/v2"
	"strings"
)

var quoteWords = []string{
	// common English words (matches monkeytype / typeracer pools)
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
	"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
	"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
	"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
	"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	"great", "between", "need", "large", "often", "hand", "high", "place", "find", "here",
	"thing", "many", "still", "long", "made", "before", "world", "life", "right", "old",
	"same", "tell", "does", "set", "three", "group", "under", "let", "end", "move",
	"try", "point", "city", "home", "small", "found", "own", "part", "off", "much",
	"while", "name", "should", "school", "every", "keep", "never", "last", "read", "run",
	"each", "left", "start", "house", "turn", "state", "play", "live", "near", "head",
	"open", "add", "next", "change", "began", "seem", "help", "talk", "where", "side",
	"been", "may", "call", "might", "stop", "must", "put", "thought", "went", "line",
	"walk", "ask", "door", "close", "feel", "plan", "sure", "build", "face", "light",
	"love", "stand", "bring", "hard", "begin", "air", "kind", "mean", "leave", "story",
}

var codeWords = []string{
	// programming keywords and concepts
	"func", "return", "error", "nil", "range", "slice", "map", "struct", "interface", "goroutine",
	"channel", "mutex", "context", "package", "module", "compile", "testing", "pointer", "method", "receiver",
	"import", "deploy", "commit", "branch", "refactor", "backend", "frontend", "async", "buffer", "runtime",
	"const", "break", "case", "continue", "default", "defer", "else", "for", "goto", "select",
	"switch", "type", "var", "string", "int", "bool", "float", "byte", "rune", "append",
	"make", "new", "close", "delete", "copy", "panic", "recover", "print", "println", "len",
	"cap", "true", "false", "iota", "init", "main", "config", "server", "client", "request",
	"response", "handler", "router", "middleware", "database", "query", "schema", "table", "index", "cache",
	"token", "parse", "format", "encode", "decode", "marshal", "unmarshal", "serialize", "validate", "filter",
	"sort", "search", "hash", "encrypt", "decrypt", "compress", "extract", "stream", "socket", "listen",
	"connect", "send", "receive", "publish", "subscribe", "queue", "stack", "tree", "graph", "node",
	"edge", "loop", "array", "list", "linked", "binary", "linear", "recursive", "iterate", "traverse",
	"insert", "update", "merge", "split", "batch", "process", "thread", "spawn", "kill", "signal",
	"event", "callback", "promise", "await", "yield", "throw", "catch", "finally", "abstract", "static",
	"public", "private", "class", "object", "inherit", "extend", "implement", "override", "template", "generic",
	"docker", "container", "image", "volume", "network", "proxy", "load", "balance", "scale", "monitor",
	"debug", "trace", "profile", "benchmark", "assert", "mock", "stub", "fixture", "factory", "builder",
	"pattern", "design", "model", "view", "controller", "service", "layer", "module", "plugin", "library",
	"framework", "sandbox", "staging", "release", "version", "upgrade", "migrate", "backup", "restore", "rollback",
	"webhook", "endpoint", "payload", "header", "status", "timeout", "retry", "fallback", "circuit", "breaker",
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
