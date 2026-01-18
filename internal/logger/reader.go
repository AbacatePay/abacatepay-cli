package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LogEntry struct {
	ID         string `json:"id"`
	Event      string `json:"event"`
	Time       string `json:"time"`
	Level      string `json:"level"`
	Msg        string `json:"msg"`
	Timestamp  string `json:"timestamp,omitempty"`
	URL        string `json:"url,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	SizeBytes  int    `json:"size_bytes,omitempty"`
	RawMessage string `json:"raw_message,omitempty"`
	Error      string `json:"error,omitempty"`
}

type ReadOptions struct {
	Limit      int
	TypeFilter string
}

func GetLogFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}

	return filepath.Join(homeDir, ".abacatepay", "logs", "transactions.log"), nil
}

func ReadTransactionLogs(opts ReadOptions) ([]LogEntry, error) {
	logPath, err := GetLogFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if opts.TypeFilter != "" && entry.Msg != opts.TypeFilter {
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	if opts.Limit > 0 && len(entries) > opts.Limit {
		entries = entries[len(entries)-opts.Limit:]
	}

	return entries, nil
}

func FindLogEntryByID(id string) (*LogEntry, error) {
	logPath, err := GetLogFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if entry.ID == id {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("event with ID %s not found in local logs", id)
}
