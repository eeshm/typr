package history

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Record stores the result of a single typing test.
type Record struct {
	Date         time.Time `json:"date"`
	Mode         string    `json:"mode"`
	WordCount    int       `json:"word_count"`
	WPM          float64   `json:"wpm"`
	RawWPM       float64   `json:"raw_wpm"`
	Accuracy     float64   `json:"accuracy"`
	Errors       int       `json:"errors"`
	TimeTaken    float64   `json:"time_taken_sec"`
	Completed    bool      `json:"completed"`
	Tier         string    `json:"tier"`
}

const maxRecords = 50

// historyPath returns the path to the JSON file inside the user's config dir.
func historyPath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cfgDir, "terminal-wpm")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

// Load reads all saved records from disk.
func Load() ([]Record, error) {
	path, err := historyPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // no history yet
		}
		return nil, err
	}

	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Save appends a record and persists to disk, keeping only the last maxRecords.
func Save(r Record) error {
	records, err := Load()
	if err != nil {
		// If history is corrupted, start fresh rather than blocking the user.
		records = nil
	}

	records = append(records, r)

	// Trim to last N records.
	if len(records) > maxRecords {
		records = records[len(records)-maxRecords:]
	}

	path, err := historyPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Recent returns the last n records (most recent last).
func Recent(n int) []Record {
	records, err := Load()
	if err != nil || len(records) == 0 {
		return nil
	}
	if n > len(records) {
		n = len(records)
	}
	return records[len(records)-n:]
}
