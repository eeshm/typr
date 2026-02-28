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

var quoteWordsExtra = []string{
	"able", "account", "across", "action", "actually", "address", "administration", "admit", "adult", "affect",
	"against", "age", "agency", "agent", "ago", "agree", "agreement", "ahead", "allow", "almost",
	"alone", "along", "already", "although", "always", "among", "amount", "analysis", "animal", "another",
	"answer", "anyone", "anything", "appear", "apply", "approach", "area", "argue", "around", "arrive",
	"article", "artist", "assume", "attack", "attention", "attorney", "audience", "author", "available", "avoid",
	"away", "baby", "bank", "base", "beat", "beautiful", "become", "bed", "believe", "benefit",
	"best", "better", "beyond", "big", "bill", "billion", "black", "blood", "blue", "board",
	"body", "book", "born", "both", "box", "boy", "break", "brother", "budget", "build",
	"building", "business", "buy", "camera", "campaign", "cancer", "candidate", "capital", "car", "card",
	"care", "career", "carry", "cause", "cell", "center", "central", "century", "certain", "certainly",
	"chair", "challenge", "chance", "character", "charge", "check", "child", "choice", "choose", "church",
	"citizen", "civil", "claim", "class", "clear", "clearly", "coach", "cold", "collection", "college",
	"color", "commercial", "common", "community", "company", "compare", "computer", "concern", "condition", "conference",
	"consider", "consumer", "contain", "continue", "control", "cost", "country", "couple", "course", "court",
	"cover", "create", "crime", "cultural", "culture", "cup", "current", "customer", "cut", "dark",
	"data", "daughter", "dead", "deal", "debate", "decade", "decide", "decision", "deep", "defense",
	"degree", "democrat", "democratic", "describe", "design", "despite", "detail", "determine", "develop", "development",
	"difference", "different", "difficult", "direction", "director", "discover", "discuss", "discussion", "disease", "doctor",
	"dog", "down", "draw", "dream", "drive", "drop", "drug", "during", "east", "easy",
	"economic", "economy", "education", "effect", "eight", "either", "election", "employee", "energy", "enjoy",
	"enough", "enter", "entire", "environment", "especially", "establish", "evening", "event", "ever", "evidence",
	"exactly", "example", "executive", "exist", "expect", "experience", "expert", "explain", "eye", "fact",
	"factor", "fail", "fall", "family", "far", "fast", "father", "fear", "federal", "free",
	"federal", "field", "fight", "figure", "fill", "film", "final", "financial", "fine", "fire",
	"fish", "five", "floor", "focus", "follow", "food", "foot", "force", "foreign", "form",
	"former", "forward", "four", "friend", "full", "fund", "game", "garden", "general", "generation",
	"girl", "glass", "goal", "government", "green", "ground", "grow", "growth", "guess", "gun",
	"guy", "hair", "half", "happen", "happy", "heart", "heat", "heavy", "herself", "history",
	"hit", "hold", "hospital", "hot", "hotel", "hour", "however", "huge", "human", "hundred",
	"husband", "idea", "identify", "image", "imagine", "impact", "important", "improve", "include", "including",
	"increase", "indeed", "indicate", "individual", "industry", "information", "inside", "instead", "institution", "interest",
	"interesting", "international", "interview", "investment", "involve", "issue", "item", "itself", "job", "join",
	"judge", "jump", "key", "kid", "kitchen", "knowledge", "land", "language", "laugh", "law",
	"lawyer", "leader", "learn", "least", "legal", "less", "letter", "level", "likely", "limit",
	"list", "listen", "little", "local", "lose", "loss", "lot", "low", "machine", "magazine",
	"main", "maintain", "major", "majority", "manage", "management", "manager", "market", "marriage", "material",
	"matter", "maybe", "measure", "media", "medical", "meeting", "member", "memory", "message", "method",
	"middle", "military", "million", "mind", "minute", "mission", "modern", "moment", "money", "month",
	"more", "morning", "mother", "mouth", "movie", "music", "myself", "nation", "national", "natural",
	"nature", "necessary", "network", "news", "newspaper", "night", "north", "note", "nothing", "notice",
	"occur", "offer", "office", "officer", "official", "once", "operation", "opportunity", "option", "order",
	"organization", "others", "outside", "owner", "page", "pain", "painting", "paper", "parent", "particular",
	"particularly", "partner", "party", "pass", "past", "patient", "pattern", "peace", "perform", "performance",
	"perhaps", "period", "person", "personal", "phone", "physical", "pick", "picture", "piece", "place",
	"plan", "plant", "player", "police", "policy", "political", "politics", "poor", "popular", "population",
	"position", "positive", "possible", "power", "practice", "prepare", "present", "president", "pressure", "pretty",
	"prevent", "price", "private", "probably", "problem", "process", "produce", "product", "production", "professional",
	"professor", "program", "project", "property", "protect", "prove", "provide", "public", "pull", "purpose",
	"push", "quality", "quickly", "quite", "race", "radio", "raise", "range", "rate", "rather",
	"reach", "real", "reality", "realize", "really", "reason", "receive", "recent", "recently", "recognize",
	"record", "red", "reduce", "reflect", "region", "relate", "relationship", "religious", "remain", "remember",
	"remove", "report", "represent", "require", "research", "resource", "respond", "response", "responsibility", "rest",
	"result", "reveal", "rich", "risk", "road", "rock", "role", "room", "rule", "safe",
	"save", "scene", "science", "scientist", "score", "season", "seat", "second", "section", "security",
	"seek", "sell", "senior", "sense", "series", "serious", "serve", "service", "seven", "several",
	"sex", "sexual", "shake", "share", "shoot", "short", "shot", "show", "sign", "significant",
	"similar", "simply", "since", "single", "sister", "site", "situation", "six", "size", "skill",
	"skin", "social", "society", "soldier", "somebody", "someone", "something", "sometimes", "son", "song",
	"soon", "source", "south", "southern", "space", "speak", "special", "specific", "speech", "spend",
	"sport", "spring", "staff", "stage", "standard", "star", "statement", "station", "stay", "step",
	"stock", "stop", "store", "strategy", "street", "strong", "student", "study", "stuff", "style",
	"subject", "success", "successful", "such", "suddenly", "suffer", "suggest", "summer", "support", "surface",
	"system", "table", "task", "tax", "teach", "teacher", "team", "technology", "television", "tend",
	"term", "test", "thank", "their", "themselves", "theory", "therefore", "third", "those", "though",
	"thousand", "threat", "through", "throughout", "throw", "today", "together", "tonight", "total", "toward",
	"town", "trade", "traditional", "training", "travel", "treat", "treatment", "tree", "trial", "trip",
	"trouble", "true", "truth", "type", "understand", "unit", "until", "upon", "usually", "value",
	"various", "very", "victim", "view", "violence", "visit", "voice", "vote", "wait", "wall",
	"watch", "water", "weapon", "wear", "week", "weight", "west", "western", "whatever", "whether",
	"white", "whole", "whose", "wife", "wind", "window", "wish", "within", "without", "woman",
	"wonder", "word", "worker", "writing", "wrong", "yard", "yeah", "young", "yourself", "zero",
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

var codeWordsExtra = []string{
	"actor", "adapter", "aggregate", "algorithm", "allocator", "analytics", "annotation", "api", "argument", "artifact",
	"assembler", "authentication", "authorization", "autoscale", "availability", "backoff", "bandwidth", "bastion", "bluegreen", "bootstrap",
	"boundary", "bucket", "buildkite", "bypass", "bytecode", "checkpoint", "checksum", "choreography", "cipher", "cleanup",
	"cluster", "coalesce", "codec", "collation", "command", "commitlint", "compatibility", "compiler", "concurrency", "configmap",
	"connection", "consistency", "constraint", "consumer", "containerd", "contention", "contract", "coordinator", "coroutine", "coverage",
	"crashloop", "credential", "cronjob", "cursor", "customresource", "datasource", "deadletter", "declarative", "decorator", "degradation",
	"dependency", "deserializer", "deterministic", "devcontainer", "diagnostic", "diff", "digest", "directive", "dispatcher", "distributed",
	"dns", "domain", "downtime", "driver", "durability", "ebpf", "edgecache", "egress", "elasticsearch", "endpointslice",
	"engine", "env", "eslint", "eventual", "exception", "executor", "featureflag", "federation", "filesystem", "firewall",
	"flake", "fluentd", "forwarder", "gateway", "garbagecollector", "gitops", "graphql", "grpc", "handshake", "healthcheck",
	"heap", "helm", "hpa", "http", "https", "hypervisor", "idempotent", "immutable", "indexer", "ingress",
	"injector", "instance", "instrumentation", "integration", "interceptor", "io", "isolation", "iterator", "javascript", "json",
	"jwt", "kafka", "kernel", "keyspace", "kubernetes", "latency", "leaderboard", "leak", "lifecycle", "linter",
	"loadbalancer", "lockfree", "logging", "lookup", "manifest", "matmul", "memcache", "memory", "messagebus", "metadata",
	"microservice", "middlewarestack", "migration", "minikube", "mirroring", "mockserver", "modular", "monorepo", "multitenant", "namespace",
	"negotiation", "nodelocal", "normalization", "nosql", "notifier", "oauth", "observability", "opcode", "operator", "optimizer",
	"orchestration", "overload", "packet", "parallel", "parameter", "partition", "passthrough", "patch", "pipeline", "platform",
	"pointermap", "policyengine", "pool", "port", "postmortem", "pragma", "precommit", "prefetch", "priority", "probe",
	"profiler", "protocol", "protobuf", "provisioner", "pseudocode", "publisher", "pullrequest", "purge", "pvc", "qos",
	"quorum", "rateLimit", "rbac", "readiness", "reconcile", "redis", "redrive", "reference", "registry", "reindex",
	"replica", "replication", "repository", "requestid", "resilience", "resolver", "resourcequota", "retention", "reusability", "rollout",
	"roundtrip", "routetable", "saga", "sanitizer", "scheduler", "schemafirst", "scope", "scripting", "secret", "segfault",
	"semantic", "serviceaccount", "session", "shard", "sidecar", "signature", "singleton", "snapshot", "snippet", "softdelete",
	"sourcecode", "spec", "spill", "splaytree", "sql", "stateful", "stateless", "stdout", "stdin", "storageclass",
	"streaming", "strictmode", "subnet", "supervisor", "suspension", "swagger", "symlink", "sync", "syscall", "telemetry",
	"terraform", "testcase", "throttle", "timeseries", "tls", "tokenizer", "topology", "tracing", "transaction", "transport",
	"ttl", "tuning", "typecheck", "udp", "ulid", "unified", "unixsocket", "upstream", "url", "uuid",
	"validator", "vector", "verifier", "virtualization", "vm", "vpc", "waf", "websocket", "workqueue", "workload",
	"yaml", "zonal", "zeroDowntime", "zipkin", "authn", "authz", "ci", "cd", "cli", "sdk",
	"checksummer", "debouncer", "memoization", "backpressure", "throughput", "fanout", "fanin", "coordinator", "watcher", "loader",
	"serializer", "deserializer", "marshaller", "demarshaller", "adapterpattern", "statepattern", "eventstore", "snapshotter", "commandbus", "querybus",
}

var quotePool = uniqueWords(append(append([]string{}, quoteWords...), quoteWordsExtra...))
var codePool = uniqueWords(append(append([]string{}, codeWords...), codeWordsExtra...))

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
		return quotePool, nil
	case "code":
		return codePool, nil
	default:
		return nil, errors.New("unsupported mode")
	}
}

func uniqueWords(words []string) []string {
	seen := make(map[string]struct{}, len(words))
	result := make([]string, 0, len(words))
	for _, word := range words {
		if _, exists := seen[word]; exists {
			continue
		}
		seen[word] = struct{}{}
		result = append(result, word)
	}
	return result
}
